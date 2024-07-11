package connect

import (
	"github.com/spf13/cobra"
	"strings"

	"bt/internal/cmd/connect/kube"
	"bt/internal/cmd/connect/ssh"
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
	}

	c.AddCommand(
		kube.NewCommand(),
		ssh.NewCommand(),
	)

	return c
}
