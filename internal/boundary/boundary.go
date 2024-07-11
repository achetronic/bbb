package boundary

import (
	"bytes"
	"os"
	"os/exec"
	"syscall"

	"bbb/internal/globals"
)

// GetTargetAuthorizedSession Ask H.Boundary for an authorized session
// This request will provide a session ID and brokered credentials associated to the target
func GetTargetAuthorizedSession(storedTokenReference, target string, consoleStdout, consoleStderr *bytes.Buffer) (command *exec.Cmd, err error) {

	boundaryArgs := []string{"targets", "authorize-session", "-id=" + target, "-token=" + storedTokenReference, "-format=json"}
	authorizeSessionCommand := exec.Command("boundary", boundaryArgs...)
	authorizeSessionCommand.Stdout = consoleStdout
	authorizeSessionCommand.Stderr = consoleStderr

	err = authorizeSessionCommand.Run()

	return command, err
}

// TODO
func GetSessionConnection(storedTokenReference, targetSessionToken string) (command *exec.Cmd, err error) {
	boundaryArgs := []string{"connect", "-authz-token=" + targetSessionToken, "-token=" + storedTokenReference, "-format=json"}
	connectCommand := exec.Command("boundary", boundaryArgs...)

	sessionFileName := targetSessionToken[:10]
	connectCommand.Stdout, _ = os.OpenFile(globals.BbbTemporaryDir+"/"+sessionFileName+".out", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0700)
	connectCommand.Stderr, _ = os.OpenFile(globals.BbbTemporaryDir+"/"+sessionFileName+".err", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0700)

	connectCommand.SysProcAttr = &syscall.SysProcAttr{
		Setsid: true,
	}

	err = connectCommand.Start()
	return connectCommand, err
}
