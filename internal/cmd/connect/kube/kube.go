package kube

import (
	"bytes"
	"encoding/json"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"

	"bt/internal/fancy"
	"bt/internal/globals"
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

	// 1. Ask H.Boundary for an authorized session
	// This request will provide a session ID and brokered credentials associated to the target
	boundaryArgs := []string{"targets", "authorize-session", "-id=" + args[0], "-token=" + storedTokenReference, "-format=json"}
	authorizeSessionCommand := exec.Command("boundary", boundaryArgs...)
	authorizeSessionCommand.Stdout = &consoleStdout
	authorizeSessionCommand.Stderr = &consoleStderr

	err = authorizeSessionCommand.Run()
	if err != nil {
		// Brutally fail when there is no output or error to handle anything
		if len(consoleStderr.Bytes()) == 0 && len(consoleStdout.Bytes()) == 0 {
			fancy.Fatalf(AuthorizeSessionErrorMessage, err.Error(), consoleStderr.String())
		}

		// Forward stderr to stdout for later processing
		consoleStdout = consoleStderr
	}

	//
	var response AuthorizeSessionResponseT
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

	// 2. Create a TCP connection to the target with authorized session previously created
	// User commands will be performed over this connection
	boundaryArgs = []string{"connect", "-authz-token=" + targetSessionToken, "-token=" + storedTokenReference, "-format=json"}
	connectCommand := exec.Command("boundary", boundaryArgs...)

	sessionFileName := targetSessionToken[:10]
	connectCommand.Stdout, _ = os.OpenFile(globals.BtTemporaryDir+"/"+sessionFileName+".out", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0700)
	connectCommand.Stderr, _ = os.OpenFile(globals.BtTemporaryDir+"/"+sessionFileName+".err", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0700)

	connectCommand.SysProcAttr = &syscall.SysProcAttr{
		Setsid: true,
	}

	err = connectCommand.Start()
	if err != nil {
		fancy.Fatalf(globals.UnexpectedErrorMessage,
			"Failed executing 'boundary connect' command: "+err.Error()+"\nCommand stderr: "+consoleStderr.String())
	}

	connectSessionStdoutEmpty := true
	var connectSessionStdoutRaw []byte
	for loop := 0; loop <= 10 && connectSessionStdoutEmpty == true; loop++ {

		stdoutFile := globals.BtTemporaryDir + "/" + sessionFileName + ".out"

		connectSessionStdoutRaw, err = os.ReadFile(stdoutFile)
		if err != nil {
			fancy.Fatalf(globals.UnexpectedErrorMessage, "Failed reading file '"+stdoutFile+"': "+err.Error())
		}

		if len(connectSessionStdoutRaw) > 0 {
			connectSessionStdoutEmpty = false
		}

		time.Sleep(500 * time.Millisecond)
	}

	if connectSessionStdoutEmpty {
		fancy.Fatalf(globals.UnexpectedErrorMessage, "There is no content on 'connect' stdout command execution")
	}

	//
	var connectSessionStdout ConnectSessionStdoutT
	err = json.Unmarshal(connectSessionStdoutRaw, &connectSessionStdout)
	if err != nil {
		fancy.Fatalf(globals.UnexpectedErrorMessage, "Failed converting JSON object into Struct: "+err.Error())
	}

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

	err = os.WriteFile(globals.BtTemporaryDir+"/"+connectSessionStdout.SessionId+".yaml", kubeconfigContent, 0700)
	if err != nil {
		fancy.Fatalf(globals.UnexpectedErrorMessage, "Failed writing kubeconfig YAML in temporary directory: "+err.Error())
	}

	// 4. Show final message to the user
	durationStringFromNow, err := GetDurationStringFromNow(connectSessionStdout.Expiration)
	if err != nil {
		fancy.Fatalf(globals.UnexpectedErrorMessage, "Error getting session duration: "+err.Error())
	}

	fancy.Printf(ConnectionSuccessfulMessage,
		connectSessionStdout.SessionId,
		durationStringFromNow,
		"kill -INT "+strconv.Itoa(authorizeSessionCommand.Process.Pid),
		"pkill -f '^boundary connect'",
		"kubectl --kubeconfig="+globals.BtTemporaryDir+"/"+connectSessionStdout.SessionId+".yaml get pods")
}
