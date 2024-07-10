package ssh

import (
	"bt/internal/boundary"
	"bt/internal/cmd/connect/kube"
	"bt/internal/fancy"
	"bt/internal/globals"
	"bytes"
	"encoding/json"
	"github.com/spf13/cobra"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const (
	descriptionShort = `Create a connection to an SSH target`

	descriptionLong = `
	Create a connection to an SSH target.
	It authorizes a session, and performs SSH connection using it`
)

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "ssh",
		DisableFlagsInUseLine: true,
		Short:                 descriptionShort,
		Long:                  strings.ReplaceAll(descriptionLong, "\t", ""),

		Run: RunCommand,
	}

	return cmd
}

// RunCommand TODO
// Ref: https://pkg.go.dev/github.com/spf13/pflag#StringSlice
func RunCommand(cmd *cobra.Command, args []string) {

	var err error
	var consoleStderr bytes.Buffer
	var consoleStdout bytes.Buffer

	//
	storedTokenReference, err := globals.GetStoredTokenReference()
	if err != nil {
		fancy.Fatalf(globals.TokenRetrievalErrorMessage)
	}

	// We need a target to connect to
	if len(args) != 1 {
		fancy.Fatalf(kube.CommandArgsNoTargetErrorMessage)
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

	// Check brokered credentials to guess whether requested target is configured as SSH target
	credentialsIndex := -1
	for credentialIndex, credential := range response.Item.Credentials {
		if credential.Credential.Username != "" && credential.Credential.PrivateKey != "" {
			credentialsIndex = credentialIndex
		}
	}

	if credentialsIndex == -1 {
		fancy.Fatalf(NotSshTargetErrorMessage)
	}

	//
	targetSessionToken := response.Item.AuthorizationToken

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
	stdoutFile := globals.BtTemporaryDir + "/" + sessionFileName + ".out"
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
	// 3. Establish SSH connection over TCP tunnel.
	// Done manually as H.Boundary CLI helpers are not so much reliable
	targetSessionSshUsername := response.Item.Credentials[credentialsIndex].Credential.Username
	targetSessionSshPrivateKey := response.Item.Credentials[credentialsIndex].Credential.PrivateKey

	// Write PrivateKey in a temporary file to be used by SSH binary
	temporaryPrivatekeyFile := globals.BtTemporaryDir + "/" + sessionFileName + ".pem"
	err = os.WriteFile(temporaryPrivatekeyFile, []byte(targetSessionSshPrivateKey+"\n"), 0600)
	if err != nil {
		log.Print(err.Error())
		return
	}

	// TODO
	_, err = globals.GetFileContents(temporaryPrivatekeyFile, true)
	if err != nil {
		fancy.Fatalf(globals.UnexpectedErrorMessage, err.Error())
	}

	// Finally, establish connection
	sshConnectionArgs := []string{
		"-o", "IdentitiesOnly=yes",
		"-o", "StrictHostKeyChecking=no",
		"-o", "ServerAliveInterval=30",
		"-p", strconv.Itoa(connectSessionStdout.Port),
		"-A", targetSessionSshUsername + "@127.0.0.1",
		"-i", temporaryPrivatekeyFile}

	sshCommand := exec.Command("ssh", sshConnectionArgs...)

	sshCommand.Stdin = os.Stdin
	sshCommand.Stdout = os.Stdout
	sshCommand.Stderr = &consoleStderr

	err = sshCommand.Run()
	if err != nil {
		fancy.Fatalf(globals.UnexpectedErrorMessage,
			"Failed executing 'ssh' command: "+err.Error()+"\nCommand stderr: "+consoleStderr.String())
	}

	// Kill the proxy to H.Boundary worker that is launched in the background
	// Clean associated privatekey PEM file
	err = connectCommand.Process.Kill()
	if err != nil {
		fancy.Fatalf(globals.UnexpectedErrorMessage,
			"Failed killing background connection to H.Boundary: "+err.Error())
	}

	err = os.Remove(temporaryPrivatekeyFile)
	if err != nil {
		fancy.Fatalf(globals.UnexpectedErrorMessage,
			"Failed removing SSH brokered credentials from temporary dir: "+err.Error())
	}
}
