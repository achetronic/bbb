package list

import (
	"bt/internal/globals"
	"bytes"
	"encoding/json"
	"github.com/spf13/cobra"
	"log"
	"os/exec"
)

const (
	descriptionShort = `TODO` // TODO

	descriptionLong = `TODO` // TODO

)

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "list",
		DisableFlagsInUseLine: true,
		Short:                 descriptionShort,
		Long:                  descriptionLong,

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
	_ = consoleStderr

	//
	storedTokenReference, err := globals.GetStoredTokenReference()
	if err != nil {
		log.Fatalf("fallo al pillar el token: %s", err.Error())
	}

	//
	boundaryArgs := []string{"scopes", "list", "-recursive", "-format=json", "-token=" + storedTokenReference}
	consoleCommand := exec.Command("boundary", boundaryArgs...)
	consoleCommand.Stdout = &consoleStdout
	consoleCommand.Stderr = &consoleStderr

	err = consoleCommand.Run()
	if err != nil {
		log.Fatalf("failed executing command: %v; %s", err, consoleStderr.String())
	}

	// Extract scopes from stdout
	var response ListScopesResponseT
	err = json.Unmarshal(consoleStdout.Bytes(), &response)
	if err != nil {
		// TODO
		return
	}

	// Retrieve and classify the scopes
	scopesByScope := GetScopesByScope(response)
	if _, globalScopeFound := scopesByScope["global"]; !globalScopeFound {
		log.Fatal("No hay scopes en tu H.Boundary")
	}

	//
	projectAbbreviationToScopeMap := AbbreviationToScopeMapT{}

	// Iterate over Global scope looking for Organizations
	for _, organization := range scopesByScope["global"] {

		// Show the organization data when no specific project is selected.
		// Projects for this organization will appear later "inside"
		if len(args) == 0 {
			log.Printf("%s (%s)\n", organization.Name, organization.Description)
		}

		// Iterate over Organizations looking for Projects
		for _, project := range scopesByScope[organization.Id] {

			projectAbbreviationToScopeMap[GenerateAbbreviation(project.Name)] = project.Id

			// List all the projects by organization when no specific one selected
			if len(args) == 0 {
				log.Printf("    %s => [%s] %s",
					GenerateAbbreviation(project.Name),
					project.Name,
					project.Description)
			}

		}
	}

	log.Print(projectAbbreviationToScopeMap)

	// We need a project to list its targets from this point, honey
	if len(args) != 1 {
		return
	}

	// Look for the targets for desired project
	boundaryArgs = []string{"targets", "list", "-scope-id=" + projectAbbreviationToScopeMap[args[0]], "-format=json", "-recursive", "-token=" + storedTokenReference}
	consoleCommand = exec.Command("boundary", boundaryArgs...)

	consoleStderr.Reset()
	consoleStdout.Reset()

	consoleCommand.Stdout = &consoleStdout
	consoleCommand.Stderr = &consoleStderr

	err = consoleCommand.Run()
	if err != nil {
		log.Fatalf("failed executing command: %v; %s", err, consoleStderr.String())
	}

	log.Print(consoleStdout.String())
}
