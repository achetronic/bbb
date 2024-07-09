package connect

import (
	"github.com/spf13/cobra"

	"bt/internal/cmd/connect/kube"
	"bt/internal/cmd/connect/ssh"
)

const (
	descriptionShort = `TODO`

	descriptionLong = `TODO`

	//
)

func NewCommand() *cobra.Command {
	c := &cobra.Command{
		Use:   "connect",
		Short: descriptionShort,
		Long:  descriptionLong,
	}

	c.AddCommand(
		kube.NewCommand(),
		ssh.NewCommand(),
	)

	return c
}
