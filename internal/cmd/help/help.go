package help

import (
	"github.com/spf13/cobra"
	"strings"
)

const (
	descriptionShort = `Help about any command`

	descriptionLong = `
	Help provides help for any command in the application.
	Simply type bt help [path to command] for full details.`
)

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "help [command] | STRING_TO_SEARCH",
		DisableFlagsInUseLine: true,
		Short:                 descriptionShort,
		Long:                  strings.ReplaceAll(descriptionLong, "\t", ""),

		Run: RunCommand,
	}

	return cmd
}

func RunCommand(cmd *cobra.Command, args []string) {
}
