package kube

import (
	"bt/internal/globals"
	"bytes"
	"encoding/json"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
	"log"
	"os"
	"os/exec"
	"strconv"
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
	storedTokenReference := "env://BOUNDARY_TOKEN"
	storedToken := os.Getenv("BOUNDARY_TOKEN")

	if storedToken == "" {

		storedToken, err = globals.GetStoredToken()
		if err != nil {
			log.Printf("fallo al pillar el token: %s", err.Error())
			return
		}

		storedTokenReference = "file://" + globals.BtTemporaryDir + "/BOUNDARY_TOKEN"
	}

	// We need a target to connect to
	if len(args) != 1 {
		log.Print("we need a target baby")
		return
	}

	//
	boundaryArgs := []string{"targets", "authorize-session", "-id=" + args[0], "-token=" + storedTokenReference, "-format=json"}
	consoleCommand := exec.Command("boundary", boundaryArgs...)
	consoleCommand.Stdout = &consoleStdout
	consoleCommand.Stderr = &consoleStderr

	err = consoleCommand.Run()
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

	//
	boundaryArgs = []string{"connect", "-authz-token=" + targetSessionToken, "-token=" + storedTokenReference, "-format=json"}

	consoleCommand = exec.Command("boundary", boundaryArgs...)

	sessionFileName := targetSessionToken[:10]
	consoleCommand.Stdout, _ = os.OpenFile(globals.BtTemporaryDir+"/"+sessionFileName+".out", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0700)
	consoleCommand.Stderr, _ = os.OpenFile(globals.BtTemporaryDir+"/"+sessionFileName+".err", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0700)

	consoleCommand.SysProcAttr = &syscall.SysProcAttr{
		Setsid: true,
	}

	err = consoleCommand.Start()
	if err != nil {
		log.Printf("Error ejecutando el comando: %v", err)
		return
	}

	connectSessionStdoutEmpty := true
	var connectSessionStdoutRaw []byte
	for loop := 0; loop <= 10 && connectSessionStdoutEmpty == true; loop++ {

		connectSessionStdoutRaw, err = os.ReadFile(globals.BtTemporaryDir + "/" + sessionFileName + ".out")
		if err != nil {
			log.Printf("pepito err 0: %s", err.Error())

			// TODO
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

	//
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
		log.Printf("pepito err: %s", err.Error())
		// TODO
		return
	}

	err = os.WriteFile(globals.BtTemporaryDir+"/"+connectSessionStdout.SessionId+".yaml", kubeconfigContent, 0700)
	if err != nil {
		log.Printf("pepito err 2: %s", err.Error())
		// TODO
		return
	}

	// Craftear el kubeconfig
	log.Printf("Execute your kubectl commands using auto-generated kubeconfig like this: kubectl --kubeconfig=%s.yaml get pods", globals.BtTemporaryDir+"/"+connectSessionStdout.SessionId)
}
