package list

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os/exec"
)

// GetScopes TODO
func GetScopes(storedTokenReference string) (scopes []ScopeT, err error) {
	var consoleStderr bytes.Buffer
	var consoleStdout bytes.Buffer

	boundaryArgs := []string{"scopes", "list", "-recursive", "-format=json", "-token=" + storedTokenReference}
	consoleCommand := exec.Command("boundary", boundaryArgs...)
	consoleCommand.Stdout = &consoleStdout
	consoleCommand.Stderr = &consoleStderr

	err = consoleCommand.Run()
	if err != nil {
		return scopes, fmt.Errorf("failed executing command: %v; %s", err, consoleStderr.String())
	}

	// Extract scopes from stdout
	var scopesResponse ListScopesResponseT
	err = json.Unmarshal(consoleStdout.Bytes(), &scopesResponse)
	if err != nil {
		err = fmt.Errorf("failed decoding json output from stdout: %v", err)
	}

	return scopesResponse.Items, err
}

// GetScopesByScope TODO
func GetScopesByScope(scopes []ScopeT) (result map[string][]ScopeT) {

	result = make(map[string][]ScopeT)

	for _, scope := range scopes {
		result[scope.ScopeId] = append(result[scope.ScopeId], scope)
	}

	return result
}

// GetScopeTargets TODO
func GetScopeTargets(scopeId string, storedTokenReference string) (targets []TargetT, err error) {

	var consoleStderr bytes.Buffer
	var consoleStdout bytes.Buffer

	// Look for the targets for desired project
	boundaryArgs := []string{"targets", "list", "-scope-id=" + scopeId, "-format=json", "-recursive", "-token=" + storedTokenReference}
	consoleCommand := exec.Command("boundary", boundaryArgs...)

	consoleStderr.Reset()
	consoleStdout.Reset()

	consoleCommand.Stdout = &consoleStdout
	consoleCommand.Stderr = &consoleStderr

	err = consoleCommand.Run()
	if err != nil {
		return targets, fmt.Errorf("failed executing command: %v; %s", err, consoleStderr.String())
	}

	// Extract targets from stdout
	var targetsResponse ListTargetsResponseT
	err = json.Unmarshal(consoleStdout.Bytes(), &targetsResponse)
	if err != nil {
		err = fmt.Errorf("failed decoding json output from stdout: %v", err)
	}

	return targetsResponse.Items, err
}
