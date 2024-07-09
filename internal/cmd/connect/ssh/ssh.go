package ssh

import (
	"bytes"
	"log"
	"os"
	"os/exec"

	"github.com/spf13/cobra"

	"bt/internal/cmd/connect/kube"
	"bt/internal/fancy"
	"bt/internal/globals"
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
		fancy.Fatalf(globals.TokenRetrievalErrorMessage)
	}

	// We need a target to connect to
	if len(args) != 1 {
		fancy.Fatalf(kube.CommandArgsNoTargetErrorMessage)
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
