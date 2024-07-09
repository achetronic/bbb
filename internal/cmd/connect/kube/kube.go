package kube

import (
	"bt/internal/globals"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"syscall"
	"time"
)

const (
	descriptionShort = `TODO` // TODO

	descriptionLong = `TODO` // TODO

)

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "kube",
		DisableFlagsInUseLine: true,
		Short:                 descriptionShort,
		Long:                  descriptionLong,

		Run: RunCommand,
	}

	return cmd
}

// RunCommand TODO
// Ref: https://pkg.go.dev/github.com/spf13/pflag#StringSlice
// Ref: https://support.hashicorp.com/hc/en-us/articles/21521422906131-BOUNDARY-Different-ways-to-connect-to-Kubernetes-cluster
// Ref: https://discuss.hashicorp.com/t/boundary-connect-usage-in-script/37658
func RunCommand(cmd *cobra.Command, args []string) {
	var err error
	var consoleStdout bytes.Buffer
	var consoleStderr bytes.Buffer

	//
	storedTokenReference, err := globals.GetStoredTokenReference()
	if err != nil {
		log.Fatalf("fallo al pillar el token: %s", err.Error())
	}

	// We need a target to connect to
	if len(args) != 1 {
		log.Fatal("we need a target baby")
	}

	// 1. Ask H.Boundary for an authorized session
	// This request will provide a session ID and brokered credentials associated to the target
	boundaryArgs := []string{"targets", "authorize-session", "-id=" + args[0], "-token=" + storedTokenReference, "-format=json"}
	authorizeSessionCommand := exec.Command("boundary", boundaryArgs...)
	authorizeSessionCommand.Stdout = &consoleStdout
	authorizeSessionCommand.Stderr = &consoleStderr

	err = authorizeSessionCommand.Run()
	if err != nil {
		log.Printf("failed executing command: %v; %s", err, consoleStderr.String())
		return
	}

	//
	var response AuthorizeSessionResponseT
	err = json.Unmarshal(consoleStdout.Bytes(), &response)
	if err != nil {
		// TODO
		return
	}

	//
	credentialsIndex := -1
	for credentialIndex, credential := range response.Item.Credentials {
		if credential.Secret.Decoded.ServiceAccountName != "" {
			credentialsIndex = credentialIndex
		}
	}

	if credentialsIndex == -1 {
		log.Fatal("Target is not configured as Kubernetes target. Quitting...")
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
		log.Printf("Error ejecutando el comando: %v", err)
		return
	}

	connectSessionStdoutEmpty := true
	var connectSessionStdoutRaw []byte
	for loop := 0; loop <= 10 && connectSessionStdoutEmpty == true; loop++ {

		connectSessionStdoutRaw, err = os.ReadFile(globals.BtTemporaryDir + "/" + sessionFileName + ".out")
		if err != nil {
			log.Printf("pepito err 0: %s", err.Error()) // TODO
			return
		}

		if len(connectSessionStdoutRaw) > 0 {
			connectSessionStdoutEmpty = false
		}

		time.Sleep(500 * time.Millisecond)
	}

	if connectSessionStdoutEmpty {
		log.Print("no hay contenido perra sata")
		return
	}

	//
	var connectSessionStdout ConnectSessionStdoutT
	err = json.Unmarshal(connectSessionStdoutRaw, &connectSessionStdout)
	if err != nil {
		log.Printf("pepito err 1: %s", err.Error())
		// TODO
		return
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
		log.Printf("pepito err: %s", err.Error()) // TODO
		return
	}

	err = os.WriteFile(globals.BtTemporaryDir+"/"+connectSessionStdout.SessionId+".yaml", kubeconfigContent, 0700)
	if err != nil {
		log.Printf("pepito err 2: %s", err.Error()) // TODO
		return
	}

	// 4. Show final message to the user
	durationStringFromNow, err := GetDurationStringFromNow(connectSessionStdout.Expiration)
	if err != nil {
		log.Printf("Error getting session duration: %s ", err.Error()) // TODO
		return
	}

	fmt.Printf(strings.ReplaceAll(ConnectKubeFinalMessageContent, "\t", ""),
		connectSessionStdout.SessionId,
		durationStringFromNow,
		"kill -INT "+strconv.Itoa(authorizeSessionCommand.Process.Pid),
		"pkill -f '^boundary connect'",
		"kubectl --kubeconfig="+globals.BtTemporaryDir+"/"+connectSessionStdout.SessionId+".yaml get pods",
	)
}
