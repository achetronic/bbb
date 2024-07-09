package auth

import (
	"bt/internal/fancy"
	"bt/internal/globals"
	"bytes"
	"encoding/json"
	"github.com/spf13/cobra"
	"os/exec"
)

const (
	descriptionShort = `TODO`

	descriptionLong = `TODO`
)

func NewCommand() *cobra.Command {
	c := &cobra.Command{
		Use:   "auth",
		Short: descriptionShort,
		Long:  descriptionLong,

		Run: RunCommand,
	}

	return c
}

// RunCommand TODO
// Ref: https://pkg.go.dev/github.com/spf13/pflag#StringSlice
func RunCommand(cmd *cobra.Command, args []string) {

	var consoleStdout bytes.Buffer
	var consoleStderr bytes.Buffer

	boundaryArgs := []string{"authenticate", "oidc", "-keyring-type=none", "-format=json"}
	consoleCommand := exec.Command("boundary", boundaryArgs...)
	consoleCommand.Stdout = &consoleStdout
	consoleCommand.Stderr = &consoleStderr

	err := consoleCommand.Run()
	if err != nil {

		// Brutally fail when there is no output or error to handle anything
		if len(consoleStderr.Bytes()) == 0 && len(consoleStdout.Bytes()) == 0 {
			fancy.Fatalf(AuthErrorMessage, err.Error(), consoleStderr.String())
		}

		// Forward stderr to stdout for later processing
		consoleStdout = consoleStderr
	}

	//
	var response ResponseT
	err = json.Unmarshal(consoleStdout.Bytes(), &response)
	if err != nil {
		fancy.Fatalf(globals.UnexpectedErrorMessage, "Failed converting JSON object into Struct: "+err.Error())
	}

	// On user failures, just inform the user
	if response.StatusCode >= 400 && response.StatusCode < 500 {
		fancy.Fatalf(AuthUserErrorMessage, consoleStdout.String())
	}

	err = globals.StoreToken(response.Item.Attributes.Token)
	if err != nil {
		fancy.Fatalf(globals.UnexpectedErrorMessage, "Failed to store H.Boundary token in your system: "+err.Error())
	}

	fancy.Printf(AuthSuccessfulMessage)
}
