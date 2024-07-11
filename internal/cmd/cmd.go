package cmd

import (
	"github.com/spf13/cobra"
	"strings"

	"bt/internal/cmd/auth"
	"bt/internal/cmd/connect"
	"bt/internal/cmd/list"
	"bt/internal/cmd/version"
)

const (
	descriptionShort = `TODO` // TODO
	descriptionLong  = `TODO` // TODO

)

func NewRootCommand(name string) *cobra.Command {
	c := &cobra.Command{
		Use:   name,
		Short: descriptionShort,
		Long:  strings.ReplaceAll(descriptionLong, "\t", ""),
	}

	c.AddCommand(
		version.NewCommand(),
		auth.NewCommand(),
		list.NewCommand(),
		connect.NewCommand(),
	)

	return c
}
