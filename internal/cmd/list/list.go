package list

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/spf13/cobra"

	"bbb/internal/fancy"
	"bbb/internal/globals"
)

const (
	descriptionShort = `List organizations' projects and their targets`
	descriptionLong  = `
	List organizations' projects and their targets.`
)

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "list",
		DisableFlagsInUseLine: true,
		Short:                 descriptionShort,
		Long:                  strings.ReplaceAll(descriptionLong, "\t", ""),

		PersistentPreRun: PreRunCommand,
		Run:              RunCommand,
	}

	return cmd
}

// PreRunCommand TODO
func PreRunCommand(cmd *cobra.Command, args []string) {
	err := globals.CheckEnv()
	if err != nil {
		fancy.Fatalf(globals.UnexpectedErrorMessage, err.Error())
	}
}

// RunCommand TODO
// Ref: https://pkg.go.dev/github.com/spf13/pflag#StringSlice
func RunCommand(cmd *cobra.Command, args []string) {

	var err error

	//
	storedTokenReference, err := globals.GetStoredTokenReference()
	if err != nil {
		fancy.Fatalf(globals.TokenRetrievalErrorMessage)
	}

	// 1. Retrieve and classify the scopes by scope
	scopes, err := GetScopes(storedTokenReference)
	if err != nil {
		fancy.Fatalf(globals.UnexpectedErrorMessage, "Failed getting scopes: "+err.Error())
	}

	scopesByScope := GetScopesByScope(scopes)
	if _, globalScopeFound := scopesByScope["global"]; !globalScopeFound {
		log.Fatal("No hay scopes en tu H.Boundary")
	}

	// 2. Craft a map with abbreviation to improve UX and its related scope ID
	projectAbbreviationToScopeMap := AbbreviationToScopeMapT{}

	if len(args) == 0 {
		fancy.Printf(ListOrganizationsCommandHeader)
	}

	if len(scopesByScope["global"]) == 0 {
		fancy.Fatalf(ListCommandEmpty)
	}

	// 3. Iterate over Global scope looking for Organizations. For each organization, this will
	// print a table with its scopes
	for _, organization := range scopesByScope["global"] {

		organizationTableHeader := fmt.Sprintf("%s: %s", organization.Name, organization.Description)
		organizationTableContent := [][]string{
			{"Project Name", "Description", "Abbreviation"},
		}

		// Iterate over Organizations looking for Projects
		for _, project := range scopesByScope[organization.Id] {

			projectAbbreviationToScopeMap[fancy.GenerateAbbreviation(project.Name)] = project.Id

			organizationTableContent = append(organizationTableContent, []string{
				project.Name,
				project.Description,
				fmt.Sprintf(fancy.Bold+fancy.Cyan+"%s"+fancy.Reset, fancy.GenerateAbbreviation(project.Name)),
			})
		}

		// Show the organization data when no specific project is selected.
		// Projects for this organization will appear later "inside"
		if len(args) == 0 {
			fancy.PrintTable(organizationTableHeader, organizationTableContent)
		}
	}

	// 4. Retrieve targets for passed scope. This will print a table with all the targets available
	// for authenticated user on that scope

	// We need a project to list its targets from this point, honey
	if len(args) != 1 {
		return
	}

	fancy.Printf(ListProjectsCommandHeader)

	// Look for the targets for desired project
	targets, err := GetScopeTargets(projectAbbreviationToScopeMap[args[0]], storedTokenReference)
	if err != nil {
		fancy.Fatalf(globals.UnexpectedErrorMessage, "Failed getting targets from scope '"+projectAbbreviationToScopeMap[args[0]]+"': "+err.Error())
	}

	// Print the table with the targets
	projectTableHeader := "Project: " + strings.ToLower(args[0])
	projectTableContent := [][]string{
		{"Name", "Address", "Port", "Target ID"},
	}

	for _, target := range targets {
		projectTableContent = append(projectTableContent, []string{
			target.Name,
			target.Address,
			strconv.Itoa(target.Attributes.DefaultPort),
			fmt.Sprintf(fancy.Cyan+fancy.Bold+"%s"+fancy.Reset, target.Id),
		})
	}

	if len(projectTableContent) < 2 {
		fancy.Fatalf(ListCommandEmpty)
	}
	fancy.PrintTable(projectTableHeader, projectTableContent)
}
