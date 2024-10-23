package cmd

import (
	"strings"

	"github.com/spf13/cobra"

	"bbb/internal/cmd/auth"
	"bbb/internal/cmd/connect"
	"bbb/internal/cmd/list"
	"bbb/internal/cmd/upgrade"
	"bbb/internal/cmd/version"
)

const (
	descriptionShort = `A super UX friendly CLI to make daily connections through H.Boundary easy to do`
	descriptionLong  = `
	A super UX friendly CLI to make daily connections through H.Boundary easy to do.
	It covers common auth, targets listing, target connections by SSH, Kubernetes, etc.
	`
)

func NewRootCommand(name string) *cobra.Command {
	c := &cobra.Command{
		Use:   name,
		Short: descriptionShort,
		Long:  strings.ReplaceAll(descriptionLong, "\t", ""),
	}

	c.AddCommand(
		auth.NewCommand(),
		connect.NewCommand(),
		list.NewCommand(),
		upgrade.NewCommand(),
		version.NewCommand(),
	)

	return c
}
