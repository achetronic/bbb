package auth

import (
	"bt/internal/globals"
	"bytes"
	"encoding/json"
	"github.com/spf13/cobra"
	"log"
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

		log.Printf("failed executing command: %v; %s", err, consoleStderr.String())
		return
	}

	//
	var response ResponseT
	log.Print(consoleStdout.String())
	err = json.Unmarshal(consoleStdout.Bytes(), &response)
	if err != nil {
		// TODO
		return
	}

	err = globals.StoreToken(response.Item.Attributes.Token)
	if err != nil {
		log.Print("mallll")
		return
	}

	log.Print("CP1")
}
