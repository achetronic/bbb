package connect

import (
	"strings"

	"github.com/spf13/cobra"

	"bbb/internal/cmd/connect/browser"
	"bbb/internal/cmd/connect/kube"
	"bbb/internal/cmd/connect/redis"
	"bbb/internal/cmd/connect/ssh"
	"bbb/internal/fancy"
	"bbb/internal/globals"
)

const (
	descriptionShort = `Create a connection to a target`
	descriptionLong  = `
	Create a connection to target.
	It authorizes a session, and performs a connection using it of a defined type`
)

func NewCommand() *cobra.Command {
	c := &cobra.Command{
		Use:   "connect",
		Short: descriptionShort,
		Long:  strings.ReplaceAll(descriptionLong, "\t", ""),

		PersistentPreRun: PreRunCommand,
	}

	c.AddCommand(
		kube.NewCommand(),
		ssh.NewCommand(),
		browser.NewCommand(),
		redis.NewCommand(),
	)

	return c
}

// PreRunCommand TODO
func PreRunCommand(cmd *cobra.Command, args []string) {
	err := globals.CheckEnv()
	if err != nil {
		fancy.Fatalf(globals.UnexpectedErrorMessage, err.Error())
	}
}
