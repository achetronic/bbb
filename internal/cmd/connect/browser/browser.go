package browser

import (
	"bytes"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"os/signal"
	"runtime"
	"strings"
	"syscall"

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
)

var (
	insecure bool
)

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "browser",
		DisableFlagsInUseLine: true,
		Short:                 descriptionShort,
		Long:                  strings.ReplaceAll(descriptionLong, "\t", ""),

		Run: RunCommand,
	}

	cmd.Flags().BoolVar(&insecure, "insecure", false, "Creates the local webserver without SSL/TLS")

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

	// Check brokered credentials to guess whether requested target is configured as Browser target
	// Also checks if the authentication header is with username and password or with bearer token
	credentialsIndex := -1
	var authenticationMethod string
	for credentialIndex, credential := range response.Item.Credentials {
		if credential.Credential.Password != "" {
			if credential.Credential.Username != "" {
				authenticationMethod = "basic"
			} else {
				authenticationMethod = "bearer"
			}
			credentialsIndex = credentialIndex
		}
	}

	if credentialsIndex == -1 {
		fancy.Fatalf(NotBrowserTargetErrorMessage)
	}

	//
	targetSessionToken := response.Item.AuthorizationToken

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

	// Retrieve target session credentials and creates authorization header
	var authHeader string
	if authenticationMethod == "basic" {
		authHeader = "Basic " + base64.StdEncoding.EncodeToString([]byte(
			response.Item.Credentials[credentialsIndex].Credential.Username+":"+response.Item.Credentials[credentialsIndex].Credential.Password))
	} else if authenticationMethod == "bearer" {
		authHeader = "Basic " + base64.StdEncoding.EncodeToString([]byte(
			response.Item.Credentials[credentialsIndex].Credential.Password))
	} else {
		fancy.Fatalf(globals.UnexpectedErrorMessage,
			"Unknown authentication method: "+authenticationMethod)
	}

	// Define the source proxy port randomnly between minPort and maxPort
	minPort := 10900
	maxPort := 11000
	sourcePort := rand.Intn(maxPort-minPort+1) + minPort

	// Define the URL of the source proxy where the browser will be opened
	sourceProxyAddress := fmt.Sprintf("127.0.0.1:%d", sourcePort)

	// Define the URL of the target to which the proxy_pass will be made
	targetProxyProtocol := "https"
	if insecure {
		targetProxyProtocol = "http"
	}
	targetProxyAddress, err := url.Parse(fmt.Sprintf("%s://127.0.0.1:%d", targetProxyProtocol, connectSessionStdout.Port))
	if err != nil {
		fancy.Fatalf(globals.UnexpectedErrorMessage,
			"Failed parsing target URL: "+err.Error())
	}

	// Create the reverse proxy
	proxy := httputil.NewSingleHostReverseProxy(targetProxyAddress)
	proxy.Transport = &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	// Configure the webserver
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		// Authroization header creation in Basic format
		r.Header.Set("Authorization", authHeader)

		// Redirect the request to the target server
		proxy.ServeHTTP(w, r)
	})

	// Start the webserver in a goroutine
	go func() {
		err := http.ListenAndServe(sourceProxyAddress, nil)
		if err != nil {
			fancy.Fatalf(globals.UnexpectedErrorMessage,
				"Error creating local webserver: "+err.Error())
		}
	}()

	////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
	// 4. Open browser to the webserver created in step 3
	// We use xdg-open or open command for Linux and MacOS systems respectively
	sourceWebserverAddress := "http://" + sourceProxyAddress
	browserCommand := exec.Command(browserCli, sourceWebserverAddress)
	fmt.Println("Opening browser: ", browserCommand)
	err = browserCommand.Run()

	browserCommand.Stdin = os.Stdin
	browserCommand.Stdout = os.Stdout
	browserCommand.Stderr = &consoleStderr

	if err != nil {
		fancy.Fatalf(globals.UnexpectedErrorMessage,
			"Failed executing 'xdg-open' command: "+err.Error()+"\nCommand stderr: "+consoleStderr.String())
	}

	// Capture the interrupt signal (Ctrl+C)
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	fmt.Println("Press Ctrl+C to close the connection...")

	// Wait for the signal
	<-c

	fmt.Printf("\nClosing connection %s...", sourceWebserverAddress)

	// Clean up the connection
	err = connectCommand.Process.Kill()
	if err != nil {
		fancy.Fatalf(globals.UnexpectedErrorMessage,
			"\nFailed killing background connection to H.Boundary: %v\n", err)
	}

	fmt.Println("\nCleaned up Boundary connection succesfully and exiting.")
}
