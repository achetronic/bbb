package ssh

import (
	"bt/internal/globals"
	"bytes"
	"github.com/spf13/cobra"
	"log"
	"os"
	"os/exec"
)

const (
	descriptionShort = `TODO` // TODO

	descriptionLong = `TODO` // TODO

)

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "ssh",
		DisableFlagsInUseLine: true,
		Short:                 descriptionShort,
		Long:                  descriptionLong,

		Run: RunCommand,
	}

	return cmd
}

// RunCommand TODO
// Ref: https://pkg.go.dev/github.com/spf13/pflag#StringSlice
func RunCommand(cmd *cobra.Command, args []string) {

	var err error
	var consoleStderr bytes.Buffer

	//
	storedTokenReference, err := globals.GetStoredTokenReference()
	if err != nil {
		log.Fatalf("fallo al pillar el token: %s", err.Error())
	}

	// We need a target to connect to
	if len(args) != 1 {
		log.Fatal("we need a target baby")
	}

	//
	boundaryArgs := []string{"connect", "ssh", "-target-id=" + args[0], "-token=" + storedTokenReference}
	consoleCommand := exec.Command("boundary", boundaryArgs...)
	consoleCommand.Stdin = os.Stdin
	consoleCommand.Stdout = os.Stdout
	consoleCommand.Stderr = &consoleStderr

	err = consoleCommand.Run()
	if err != nil {

		log.Printf("failed executing command: %v; %s", err, consoleStderr.String())
		return
	}

	log.Print(consoleCommand.Stdout)
}
