package connect

import (
	"bt/internal/cmd/connect/kube"
	"bt/internal/cmd/connect/ssh"

	"github.com/spf13/cobra"
	"os"
)

const (
	descriptionShort = `TODO`

	descriptionLong = `TODO`

	//

	//
	LogLevelFlagErrorMessage                   = "impossible to get flag --log-level: %s"
	DisableTraceFlagErrorMessage               = "impossible to get flag --disable-trace: %s"
	GroupFlagErrorMessage                      = "impossible to get flag --group: %s"
	SyncTimeFlagErrorMessage                   = "impossible to get flag --sync-time: %s"
	UnableParseDurationErrorMessage            = "unable to parse duration: %s"
	UnableCreateGroupErrorMessage              = "unable to create group in Boundary: %s"
	UnableSetGroupMembersErrorMessage          = "unable to set group members: %s"
	EnvironmentVariableErrorMessage            = "environment variable not found"
	GsuiteCreateAdminErrorMessage              = "Unable to create new admin: %s"
	BoundaryAuthMethodNotSupportedErrorMessage = "boundary auth method not supported"
	BoundaryOidcIdFlagErrorMessage             = "impossible to get flag --boundary-oidc-id: %s"
	BoundaryScopeIdFlagErrorMessage            = "impossible to get flag --boundary-scope-id: %s"
	GetBoundaryPersonMapErrorMessage           = "Unable to get boundary persons: %s"
	SetupBoundaryObjErrorMessage               = "Fail to setup boundary object: %s"
	BoundaryOidcIdRequiredFlagErrorMessage     = "Mark boundary-oidc-id flag as required fail: %s"
	UnableBoundaryGetGroupsErrorMessage        = "Unable to get boundary groups: %s"
)

var (
	bAddressEnv            = os.ExpandEnv(os.Getenv("BOUNDARY_ADDR"))
	bAuthMethodPassIdEnv   = os.ExpandEnv(os.Getenv("BOUNDARY_AUTHMETHODPASS_ID"))
	bAuthMethodPassUserEnv = os.ExpandEnv(os.Getenv("BOUNDARY_AUTHMETHODPASS_USER"))
	bAuthMethodPassPassEnv = os.ExpandEnv(os.Getenv("BOUNDARY_AUTHMETHODPASS_PASS"))
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

	//
	//cmd.Flags().String("log-level", "info", "Verbosity level for logs")
	//cmd.Flags().Bool("disable-trace", true, "Disable showing traces in logs")
	//cmd.Flags().String("sync-time", "10m", "Waiting time between group synchronizations (in duration type)")
	//cmd.Flags().String("google-sa-credentials-path", "google.json", "Google ServiceAccount credentials JSON file path")
	//cmd.Flags().StringSlice("google-group", []string{}, "(Repeatable or comma-separated list) G.Workspace groups")
	//cmd.Flags().String("boundary-oidc-id", "amoidc_changeme", "Boundary oidc auth method ID to compare its users against G.Workspace")
	//cmd.Flags().String("boundary-scope-id", "global", "Boundary scope ID where the users and groups are synchronized")

	return c
}
