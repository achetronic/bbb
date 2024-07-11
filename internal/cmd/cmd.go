package cmd

import (
	"strings"

	"github.com/spf13/cobra"

	"bbb/internal/cmd/auth"
	"bbb/internal/cmd/connect"
	"bbb/internal/cmd/list"
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
		version.NewCommand(),
		auth.NewCommand(),
		list.NewCommand(),
		connect.NewCommand(),
	)

	return c
}
