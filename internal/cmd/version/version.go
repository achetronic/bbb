package version

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

const (
	descriptionShort = `Print the current version`
	descriptionLong  = `
	Version show the current bbb version client.`
)

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "version",
		DisableFlagsInUseLine: true,
		Short:                 descriptionShort,
		Long:                  strings.ReplaceAll(descriptionLong, "\t", ""),

		Run: RunCommand,
	}

	return cmd
}

func RunCommand(cmd *cobra.Command, args []string) {
	fmt.Print("version: 0.2.0\n")
}
