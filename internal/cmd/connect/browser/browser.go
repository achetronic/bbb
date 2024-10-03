package browser

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/spf13/cobra"

	"bbb/internal/boundary"
	"bbb/internal/fancy"
	"bbb/internal/globals"

	"net/http"
	"net/http/httputil"
	"net/url"
)

const (
	descriptionShort = `Create a connection to a Browser target`
	descriptionLong  = `
	Create a connection to a Browser target.
	It authorizes a session, and opens a Browser using it`

	//
	webserverPortRangeMin = 10000
	webserverPortRangeMax = 11000
)

var (
	webserverPort int = GetFreeRandomPort(webserverPortRangeMin, webserverPortRangeMax)

	//
	targetProxyAddressPattern = "%s://127.0.0.1:%d"
)

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "browser",
		DisableFlagsInUseLine: true,
		Short:                 descriptionShort,
		Long:                  strings.ReplaceAll(descriptionLong, "\t", ""),

		Run: RunCommand,
	}

	cmd.Flags().IntVar(&webserverPort, "port", webserverPort, "Port for the local webserver where browser will connect")

	return cmd
}

// RunCommand TODO
// Ref: https://pkg.go.dev/github.com/spf13/pflag#StringSlice
func RunCommand(cmd *cobra.Command, args []string) {

	var err error
	var consoleStderr bytes.Buffer
	var consoleStdout bytes.Buffer

	//
	storedTokenReference, err := globals.GetStoredTokenReference()
	if err != nil {
		fancy.Fatalf(globals.TokenRetrievalErrorMessage)
	}

	// We need a target to connect to
	if len(args) != 1 {
		fancy.Fatalf(CommandArgsNoTargetErrorMessage)
	}

	var browserCli string
	switch runtime.GOOS {
	case "darwin":
		browserCli = "open"
	case "linux":
		browserCli = "xdg-open"
	default:
		fancy.Fatalf(globals.UnexpectedErrorMessage,
			"Usupported SO, we don't know which command to use to open the browser")
	}

	_, err = exec.LookPath(browserCli)
	if err != nil {
		fancy.Fatalf(BrowserCliNotFoundErrorMessage)
	}

	////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
	// 1. Ask H.Boundary for an authorized session
	// This request will provide a session ID and brokered credentials associated to the target
	// (AuthorizeSession & Connect) are performed in separated steps to check type of target before connecting
	_, err = boundary.GetTargetAuthorizedSession(storedTokenReference, args[0], &consoleStdout, &consoleStderr)
	if err != nil {
		// Brutally fail when there is no output or error to handle anything
		if len(consoleStderr.Bytes()) == 0 && len(consoleStdout.Bytes()) == 0 {
			fancy.Fatalf(AuthorizeSessionErrorMessage, err.Error(), consoleStderr.String())
		}

		// Forward stderr to stdout for later processing
		consoleStdout = consoleStderr
	}

	//
	var response boundary.AuthorizeSessionResponseT
	err = json.Unmarshal(consoleStdout.Bytes(), &response)
	if err != nil {
		fancy.Fatalf(globals.UnexpectedErrorMessage, "Failed converting JSON object into Struct: "+err.Error())
	}

	// On user failures, just inform the user
	if response.StatusCode >= 400 && response.StatusCode < 500 {
		fancy.Fatalf(AuthorizeSessionUserErrorMessage, consoleStdout.String())
	}

	// Check brokered credentials to guess which type of auth to use as Browser target
	// If some credentials are brokered, password is mandatory for all types of authentication. Username is checked later
	credentialsIndex := -1
	var authenticationMethod string
	for credentialIndex, credential := range response.Item.Credentials {
		if credential.Credential.Password != "" {
			credentialsIndex = credentialIndex
		}
	}

	if credentialsIndex == -1 {
		fancy.Printf(TargetWithNoCredentials)
	}

	//
	targetSessionToken := response.Item.AuthorizationToken

	// Extract host of the target in Boundary for later usage.
	// Remember some proxies use this to route
	targetSessionUrl, err := url.Parse(response.Item.Endpoint)
	if err != nil {
		fancy.Fatalf(globals.UnexpectedErrorMessage,
			"Failed parsing session URL. You have to configure a valid URL in Boundary: "+err.Error())
	}

	//
	targetSessionHost := strings.Split(targetSessionUrl.Host, ":")
	if len(targetSessionHost) != 2 {
		fancy.Fatalf(globals.UnexpectedErrorMessage,
			"Failed parsing session Host. Session URL must have <address>:<port> format: "+err.Error())
	}

	////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
	// 2. Create a TCP connection to the target with authorized session previously created
	// User commands will be performed over this connection
	sessionFileName := targetSessionToken[:10]
	connectCommand, err := boundary.GetSessionConnection(storedTokenReference, targetSessionToken)
	if err != nil {
		fancy.Fatalf(globals.UnexpectedErrorMessage,
			"Failed executing 'boundary connect' command: "+err.Error()+"\nCommand stderr: "+consoleStderr.String())
	}

	//
	stdoutFile := globals.BbbTemporaryDir + "/" + sessionFileName + ".out"
	connectSessionStdoutRaw, err := globals.GetFileContents(stdoutFile, true)
	if err != nil {
		fancy.Fatalf(globals.UnexpectedErrorMessage, err.Error())
	}

	//
	var connectSessionStdout boundary.ConnectSessionStdoutT
	err = json.Unmarshal(connectSessionStdoutRaw, &connectSessionStdout)
	if err != nil {
		fancy.Fatalf(globals.UnexpectedErrorMessage, "Failed converting JSON object into Struct: "+err.Error())
	}

	////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
	// 3. Create a webserver to inject Authorization Header from Username and Password retrieved from H.Boundary
	// It is needed so xdg-open/open can not inject headers by itself

	// Check the protocol used by the brokered target
	targetProxyAddress := fmt.Sprintf(targetProxyAddressPattern, "http", connectSessionStdout.Port)
	_, err = http.Get(targetProxyAddress)
	if err != nil {
		targetProxyAddress = fmt.Sprintf(targetProxyAddressPattern, "https", connectSessionStdout.Port)
	}

	targetProxyAddressObj, err := url.Parse(targetProxyAddress)
	if err != nil {
		fancy.Fatalf(globals.UnexpectedErrorMessage,
			"Failed parsing target URL: "+err.Error())
	}

	// If there are credentials, assume bearer authentication using password as token by default.
	// Assume basic auth when username is also present
	var authHeader string = ""

	if credentialsIndex != -1 {
		var headerAuthUsername string

		authenticationMethod = "bearer"
		if response.Item.Credentials[credentialsIndex].Credential.Username != "" {
			authenticationMethod = "basic"
			headerAuthUsername = response.Item.Credentials[credentialsIndex].Credential.Username
		}

		authHeader, err = GetAuthHeaderValue(
			headerAuthUsername,
			response.Item.Credentials[credentialsIndex].Credential.Password,
			authenticationMethod)
		if err != nil {
			fancy.Fatalf(globals.UnexpectedErrorMessage,
				"Failed crafting authorization header: "+err.Error())
		}
	}

	// Create and start the webserver
	webserverAddress := fmt.Sprintf("127.0.0.1:%d", webserverPort)
	webserver := httputil.NewSingleHostReverseProxy(targetProxyAddressObj)
	webserver.Transport = &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		r.Header.Set("Authorization", authHeader)
		r.Host = targetSessionHost[0]
		webserver.ServeHTTP(w, r)
	})

	go func() {
		err := http.ListenAndServe(webserverAddress, nil)
		if err != nil {
			fancy.Fatalf(globals.UnexpectedErrorMessage, "Error creating local webserver: "+err.Error())
		}
	}()

	////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
	// 4. Open browser to the webserver created in step 3
	// We use xdg-open or open command for Linux and MacOS systems respectively
	sourceWebserverAddress := "http://" + webserverAddress
	browserCommand := exec.Command(browserCli, sourceWebserverAddress)
	err = browserCommand.Run()

	browserCommand.Stdin = os.Stdin
	browserCommand.Stdout = os.Stdout
	browserCommand.Stderr = &consoleStderr

	if err != nil {
		fancy.Fatalf(globals.UnexpectedErrorMessage,
			"Failed executing 'xdg-open' command: "+err.Error()+"\nCommand stderr: "+consoleStderr.String())
	}

	// Capture some OS signals to close the process gracefully before closing
	WaitSignalAfter(func() {
		durationStringFromNow, err := globals.GetDurationStringFromNow(connectSessionStdout.Expiration)
		if err != nil {
			fancy.Fatalf(globals.UnexpectedErrorMessage, "Error getting session duration: "+err.Error())
		}

		fancy.Printf(ConnectionSuccessfulMessage,
			connectSessionStdout.SessionId,
			durationStringFromNow,
			sourceWebserverAddress)
	})

	// Clean up the connection
	err = connectCommand.Process.Kill()
	if err != nil {
		fancy.Fatalf(globals.UnexpectedErrorMessage,
			"\nFailed killing background connection to H.Boundary: %v\n", err)
	}
}
