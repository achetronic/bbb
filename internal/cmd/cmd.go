package cmd

import (
	"bt/internal/cmd/auth"
	"bt/internal/cmd/connect"
	"bt/internal/cmd/version"

	"github.com/spf13/cobra"
)

const (
	descriptionShort = `TODO` // TODO

	// descriptionLong TODO
	descriptionLong = `TODO` // TODO
)

func NewRootCommand(name string) *cobra.Command {
	c := &cobra.Command{
		Use:   name,
		Short: descriptionShort,
		Long:  descriptionLong,
	}

	c.AddCommand(
		version.NewCommand(),
		auth.NewCommand(),
		connect.NewCommand(),
	)

	return c
}
