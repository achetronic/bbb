package direct

import (
	"bytes"
	"encoding/json"
	"strconv"
	"strings"

	"github.com/spf13/cobra"

	"bbb/internal/boundary"
	"bbb/internal/fancy"
	"bbb/internal/globals"

	"net/url"
)

const (
	descriptionShort = `Create a direct connection to a target`
	descriptionLong  = `
	Create a direct connection to a target.
	It authorizes a session, and give you freedom`
)

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "direct",
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
	if response.StatusCode >= 400 && response.StatusCode < 600 {
		fancy.Fatalf(AuthorizeSessionUserErrorMessage, consoleStdout.String())
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
			"Failed parsing session Host. Session URL must have <address>:<port> format")
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

	// Capture some OS signals to close the process gracefully before closing
	WaitSignalAfter(func() {
		durationStringFromNow, err := globals.GetDurationStringFromNow(connectSessionStdout.Expiration)
		if err != nil {
			fancy.Fatalf(globals.UnexpectedErrorMessage, "Error getting session duration: "+err.Error())
		}

		fancy.Printf(ConnectionSuccessfulMessage,
			connectSessionStdout.SessionId,
			durationStringFromNow,
			connectSessionStdout.Address+":"+strconv.Itoa(connectSessionStdout.Port))
	})

	// Clean up the connection
	err = connectCommand.Process.Kill()
	if err != nil {
		fancy.Fatalf(globals.UnexpectedErrorMessage,
			"Failed killing background connection to H.Boundary: "+err.Error())
	}
}
