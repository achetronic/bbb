package ssh

import (
	"bytes"
	"encoding/json"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"

	"bt/internal/cmd/connect/kube"
	"bt/internal/fancy"
	"bt/internal/globals"
)

const (
	descriptionShort = `Create a connection to an SSH target`

	descriptionLong = `
	Create a connection to an SSH target.
	It authorizes a session, and performs SSH connection using it`
)

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "ssh",
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
		fancy.Fatalf(kube.CommandArgsNoTargetErrorMessage)
	}

	// 1. Ask H.Boundary for an authorized session
	// This request will provide a session ID and brokered credentials associated to the target
	// (AuthorizeSession & Connect) are performed in separated steps to check type of target before connecting
	boundaryArgs := []string{"targets", "authorize-session", "-id=" + args[0], "-token=" + storedTokenReference, "-format=json"}
	authorizeSessionCommand := exec.Command("boundary", boundaryArgs...)
	authorizeSessionCommand.Stdout = &consoleStdout
	authorizeSessionCommand.Stderr = &consoleStderr

	err = authorizeSessionCommand.Run()
	if err != nil {
		// Brutally fail when there is no output or error to handle anything
		if len(consoleStderr.Bytes()) == 0 && len(consoleStdout.Bytes()) == 0 {
			fancy.Fatalf(AuthorizeSessionErrorMessage, err.Error(), consoleStderr.String())
		}

		// Forward stderr to stdout for later processing
		consoleStdout = consoleStderr
	}

	//
	var response AuthorizeSessionResponseT
	err = json.Unmarshal(consoleStdout.Bytes(), &response)
	if err != nil {
		fancy.Fatalf(globals.UnexpectedErrorMessage, "Failed converting JSON object into Struct: "+err.Error())
	}

	// On user failures, just inform the user
	if response.StatusCode >= 400 && response.StatusCode < 500 {
		fancy.Fatalf(AuthorizeSessionUserErrorMessage, consoleStdout.String())
	}

	// Check whether the target and user's requested type match
	credentialsIndex := -1
	for credentialIndex, credential := range response.Item.Credentials {
		if credential.Credential.Username != "" &&
			(credential.Credential.PrivateKey != "" || credential.Credential.Password != "") {
			credentialsIndex = credentialIndex
		}
	}

	if credentialsIndex == -1 {
		fancy.Fatalf(NotSshTargetErrorMessage)
	}

	//
	targetSessionToken := response.Item.AuthorizationToken
	targetSessionSshUsername := response.Item.Credentials[credentialsIndex].Credential.Username

	// 2. Create an SSH connection to the target with authorized session previously created
	boundaryArgs = []string{"connect", "ssh", "-authz-token=" + targetSessionToken, "-token=" + storedTokenReference}
	if targetSessionSshUsername != "" {
		boundaryArgs = append(boundaryArgs, "-username="+targetSessionSshUsername)
	}

	consoleCommand := exec.Command("boundary", boundaryArgs...)
	consoleCommand.Stdin = os.Stdin
	consoleCommand.Stdout = os.Stdout
	consoleCommand.Stderr = &consoleStderr

	err = consoleCommand.Run()
	if err != nil {
		fancy.Fatalf(globals.UnexpectedErrorMessage, "Failed establishing SSH connectiong to the target: "+err.Error())
	}

}
