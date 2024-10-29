package redis

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"

	"bbb/internal/boundary"
	"bbb/internal/fancy"
	"bbb/internal/globals"

	"net/url"
)

const (
	descriptionShort = `Create a connection to a Redis target`
	descriptionLong  = `
	Create a connection to a Redis target.
	It authorizes a session if needed, and open a redis-cli using it`
)

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "redis",
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
		fancy.Fatalf(CommandArgsNoTargetErrorMessage)
	}

	// Redis-cli package used to open the connection
	var redisCli string = "redis-cli"

	// Check if redis-cli is present in the system, if not, exit with error
	_, err = exec.LookPath(redisCli)
	if err != nil {
		fancy.Fatalf(RedisCliNotFoundErrorMessage)
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

	if len(response.Item.Credentials) == 0 {
		fancy.Printf(TargetWithNoCredentials)
	}

	//
	targetSessionToken := response.Item.AuthorizationToken

	// Extract host of the target in Boundary for later usage.
	// Remember some proxies use this to route
	targetSessionUrl, err := url.Parse(response.Item.Endpoint)
	if err != nil {
		fancy.Fatalf(globals.UnexpectedErrorMessage,
			"Failed parsing session URL. You have to configure a valid URL in Boundary: "+err.Error())
	}

	//
	targetSessionHost := strings.Split(targetSessionUrl.Host, ":")
	if len(targetSessionHost) != 2 {
		fancy.Fatalf(globals.UnexpectedErrorMessage,
			"Failed parsing session Host. Session URL must have <address>:<port> format: "+err.Error())
	}

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
	// 3. Open redis-cli to the socket opened by the connection
	// We use redis-cli or just return the connection to the user

	// Redis URL to connect to
	redisUrl := fmt.Sprintf("redis://%s:%d", connectSessionStdout.Address, connectSessionStdout.Port)

	// Redis-cli arguments for authentication if needed
	redisCliArgs := []string{"-u", redisUrl}

	// Use password/access key to the cli command
	if response.Item.Credentials[0].Credential.Password != "" {
		redisCliArgs = append(redisCliArgs,
			"--pass", response.Item.Credentials[0].Credential.Password)
	}

	// Boundary does not allow to define password and not an username
	// As workaround, we set in boundary username and password to the same value
	// rescued from Vault
	if response.Item.Credentials[0].Credential.Username != "" &&
		response.Item.Credentials[0].Credential.Username != response.Item.Credentials[0].Credential.Password {
		redisCliArgs = append(redisCliArgs,
			"--user", response.Item.Credentials[0].Credential.Username)
	}

	// Set redis-cli command
	redisCommand := exec.Command(redisCli, redisCliArgs...)

	redisCommand.Stdin = os.Stdin
	redisCommand.Stdout = os.Stdout
	redisCommand.Stderr = &consoleStderr

	durationStringFromNow, err := globals.GetDurationStringFromNow(connectSessionStdout.Expiration)
	if err != nil {
		fancy.Fatalf(globals.UnexpectedErrorMessage, "Error getting session duration: "+err.Error())
	}

	fancy.Printf(ConnectionSuccessfulMessage,
		connectSessionStdout.SessionId,
		durationStringFromNow,
		redisUrl)

	err = redisCommand.Run()

	if err != nil {
		fancy.Fatalf(globals.UnexpectedErrorMessage,
			"Failed executing '"+redisCli+"' command: "+err.Error()+"\nCommand stderr: "+consoleStderr.String())
	}

	// Clean up the connection
	err = connectCommand.Process.Kill()
	if err != nil {
		fancy.Fatalf(globals.UnexpectedErrorMessage,
			"\nFailed killing background connection to H.Boundary: %v\n", err)
	}
}
