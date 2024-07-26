package browser

import (
	"bytes"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
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

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "browser",
		DisableFlagsInUseLine: true,
		Short:                 descriptionShort,
		Long:                  strings.ReplaceAll(descriptionLong, "\t", ""),

		Run: RunCommand,
	}

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

	_, err = exec.LookPath("xdg-open")
	if err != nil {
		fancy.Fatalf(XdgOpenCliNotFoundErrorMessage)
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
	credentialsIndex := -1
	for credentialIndex, credential := range response.Item.Credentials {
		if credential.Credential.Username != "" || credential.Credential.Password != "" {
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
	// It is needed so xdg-open can not inject headers by itself

	// Retrieve target session credentials
	targetSessionBrowserUsername := response.Item.Credentials[credentialsIndex].Credential.Username
	targetSessionBrowserPassword := response.Item.Credentials[credentialsIndex].Credential.Password

	// Define the URL of the source proxy where the browser will be opened
	sourceProxyAddress := "127.0.0.1:10901"

	// Define the URL of the target to which the proxy_pass will be made
	targetProxyAddress, err := url.Parse(fmt.Sprintf("https://127.0.0.1:%d", connectSessionStdout.Port))
	if err != nil {
		fancy.Fatalf(globals.UnexpectedErrorMessage, "Failed parsing target URL: "+err.Error())
	}

	// Create the reverse proxy
	proxy := httputil.NewSingleHostReverseProxy(targetProxyAddress)
	proxy.Transport = &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	// Configure the webserver
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		// Authroization header creation in Basic format
		authHeader := "Basic " + base64.StdEncoding.EncodeToString([]byte(targetSessionBrowserUsername+":"+targetSessionBrowserPassword))
		r.Header.Set("Authorization", authHeader)

		// Redirect the request to the target server
		proxy.ServeHTTP(w, r)
	})

	// Start the webserver in a goroutine
	go func() {
		err := http.ListenAndServe(sourceProxyAddress, nil)
		if err != nil {
			fmt.Println("Error al iniciar el servidor web:", err)
		}
	}()

	////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
	// 4. Open browser to the webserver created in step 3
	// We use xdg-open or open command for Linux and MacOS systems respectively
	sourceWebserverAddress := "http://" + sourceProxyAddress
	browserCommand := exec.Command("xdg-open", sourceWebserverAddress)
	fmt.Println("Opening browser...", browserCommand)
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

	fmt.Println("Press Ctrl+C to exit...")

	// Wait for the signal
	<-c

	fmt.Println("\nExiting...")

	// Aquí podrías agregar código para limpiar o finalizar otros procesos si es necesario
	err = connectCommand.Process.Kill()
	if err != nil {
		fmt.Printf("Failed killing background connection to H.Boundary: %v\n", err)
	}

	fmt.Println("Cleaned up and exiting.")
}
