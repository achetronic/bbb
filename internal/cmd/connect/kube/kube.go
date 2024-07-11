package kube

import (
	"bytes"
	"encoding/json"
	"gopkg.in/yaml.v3"
	"os"
	"strconv"
	"strings"

	"github.com/spf13/cobra"

	"bbb/internal/boundary"
	"bbb/internal/fancy"
	"bbb/internal/globals"
)

const (
	descriptionShort = `Create a connection to a Kubernetes target`
	descriptionLong  = `
	Create a connection to a Kubernetes target.
	It authorizes a session, maintains a TCP proxy connected to H.Boundary, and prepare kubectl to run commands`
)

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "kube",
		DisableFlagsInUseLine: true,
		Short:                 descriptionShort,
		Long:                  strings.ReplaceAll(descriptionLong, "\t", ""),

		Run: RunCommand,
	}

	return cmd
}

// RunCommand TODO
// Ref: https://support.hashicorp.com/hc/en-us/articles/21521422906131-BOUNDARY-Different-ways-to-connect-to-Kubernetes-cluster
// Ref: https://discuss.hashicorp.com/t/boundary-connect-usage-in-script/37658
func RunCommand(cmd *cobra.Command, args []string) {
	var err error
	var consoleStdout bytes.Buffer
	var consoleStderr bytes.Buffer

	//
	storedTokenReference, err := globals.GetStoredTokenReference()
	if err != nil {
		fancy.Fatalf(globals.TokenRetrievalErrorMessage)
	}

	// We need a target to connect to
	if len(args) != 1 {
		fancy.Fatalf(CommandArgsNoTargetErrorMessage)
	}

	////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
	// 1. Ask H.Boundary for an authorized session
	// This request will provide a session ID and brokered credentials associated to the target
	// (AuthorizeSession & Connect) are performed in separated steps to check type of target before connecting
	_, err = boundary.GetTargetAuthorizedSession(storedTokenReference, args[0], &consoleStdout, &consoleStderr)
	if err != nil {
		// Brutally fail when there is no output or error to handle anything
		if len(consoleStderr.Bytes()) == 0 && len(consoleStdout.Bytes()) == 0 {
			fancy.Fatalf(AuthorizeSessionErrorMessage, err.Error(), consoleStderr.String())
		}

		// Forward stderr to stdout for later processing
		consoleStdout = consoleStderr
	}

	//
	var response boundary.AuthorizeSessionResponseT
	err = json.Unmarshal(consoleStdout.Bytes(), &response)
	if err != nil {
		fancy.Fatalf(globals.UnexpectedErrorMessage, "Failed converting JSON object into Struct: "+err.Error())
	}

	// On user failures, just inform the user
	if response.StatusCode >= 400 && response.StatusCode < 500 {
		fancy.Fatalf(AuthorizeSessionUserErrorMessage, consoleStdout.String())
	}

	// Check whether the target and user's requested type match
	credentialsIndex := -1
	for credentialIndex, credential := range response.Item.Credentials {
		if credential.Secret.Decoded.ServiceAccountName != "" {
			credentialsIndex = credentialIndex
		}
	}

	if credentialsIndex == -1 {
		fancy.Fatalf(NotKubeTargetErrorMessage)
	}

	//
	targetSessionToken := response.Item.AuthorizationToken
	targetSessionKubernetesSaToken := response.Item.Credentials[credentialsIndex].Secret.Decoded.ServiceAccountToken

	////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
	// 2. Create a TCP connection to the target with authorized session previously created
	// User commands will be performed over this connection
	sessionFileName := targetSessionToken[:10]
	connectCommand, err := boundary.GetSessionConnection(storedTokenReference, targetSessionToken)
	if err != nil {
		fancy.Fatalf(globals.UnexpectedErrorMessage,
			"Failed executing 'boundary connect' command: "+err.Error()+"\nCommand stderr: "+consoleStderr.String())
	}

	//
	stdoutFile := globals.BbbTemporaryDir + "/" + sessionFileName + ".out"
	connectSessionStdoutRaw, err := globals.GetFileContents(stdoutFile, true)
	if err != nil {
		fancy.Fatalf(globals.UnexpectedErrorMessage, err.Error())
	}

	//
	var connectSessionStdout boundary.ConnectSessionStdoutT
	err = json.Unmarshal(connectSessionStdoutRaw, &connectSessionStdout)
	if err != nil {
		fancy.Fatalf(globals.UnexpectedErrorMessage, "Failed converting JSON object into Struct: "+err.Error())
	}

	////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
	// 3. Craft a temporary kubeconfig for this session in a temporary directory,
	// in a temporary heart, in a temp... oh wait, recursive comment detected
	kubeconfig := KubeconfigT{
		ApiVersion: "v1",
		Kind:       "Config",
		Clusters: []KubeconfigClustersT{
			{
				Name: "cluster_" + connectSessionStdout.SessionId,
				Cluster: KubeconfigClustersClusterT{
					Server:             "https://127.0.0.1:" + strconv.Itoa(connectSessionStdout.Port),
					InsecureSkipVerify: true,
				},
			},
		},
		Contexts: []KubeconfigContextT{
			{
				Name: "context_" + connectSessionStdout.SessionId,
				Context: KubeconfigContextContextT{
					Cluster: "cluster_" + connectSessionStdout.SessionId,
					User:    "user_" + connectSessionStdout.SessionId,
				},
			},
		},
		Users: []KubeconfigUsersT{
			{
				Name: "user_" + connectSessionStdout.SessionId,
				User: KubeconfigUsersUserT{
					Token: targetSessionKubernetesSaToken,
				},
			},
		},
		CurrentContext: "context_" + connectSessionStdout.SessionId,
	}

	kubeconfigContent, err := yaml.Marshal(kubeconfig)
	if err != nil {
		fancy.Fatalf(globals.UnexpectedErrorMessage, "Failed converting kubeconfig object into YAML: "+err.Error())
	}

	err = os.WriteFile(globals.BbbTemporaryDir+"/"+connectSessionStdout.SessionId+".yaml", kubeconfigContent, 0700)
	if err != nil {
		fancy.Fatalf(globals.UnexpectedErrorMessage, "Failed writing kubeconfig YAML in temporary directory: "+err.Error())
	}

	// 4. Show final message to the user
	durationStringFromNow, err := globals.GetDurationStringFromNow(connectSessionStdout.Expiration)
	if err != nil {
		fancy.Fatalf(globals.UnexpectedErrorMessage, "Error getting session duration: "+err.Error())
	}

	fancy.Printf(ConnectionSuccessfulMessage,
		connectSessionStdout.SessionId,
		durationStringFromNow,
		"kill -INT "+strconv.Itoa(connectCommand.Process.Pid),
		"pkill -f '^boundary connect'",
		"kubectl --kubeconfig="+globals.BbbTemporaryDir+"/"+connectSessionStdout.SessionId+".yaml get pods")
}
