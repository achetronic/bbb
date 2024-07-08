package list

import (
	"bt/internal/globals"
	"fmt"
	"github.com/spf13/cobra"
	"log"
	"strconv"
	"strings"
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

	//
	storedTokenReference, err := globals.GetStoredTokenReference()
	if err != nil {
		log.Fatalf("fallo al pillar el token: %s", err.Error())
	}

	// Retrieve and classify the scopes by scope
	scopes, err := GetScopes(storedTokenReference)
	if err != nil {
		// TODO
		return
	}

	scopesByScope := GetScopesByScope(scopes)
	if _, globalScopeFound := scopesByScope["global"]; !globalScopeFound {
		log.Fatal("No hay scopes en tu H.Boundary")
	}

	// Craft a map with abbreviation to improve UX and its related scope ID
	projectAbbreviationToScopeMap := AbbreviationToScopeMapT{}

	if len(args) == 0 {
		fmt.Println(strings.ReplaceAll(ListOrganizationsCommandHeader, "\t", ""))
	}

	if len(scopesByScope["global"]) == 0 {
		fmt.Println(strings.ReplaceAll(ListCommandEmpty, "\t", ""))
		return
	}

	// Iterate over Global scope looking for Organizations
	for _, organization := range scopesByScope["global"] {

		organizationTableHeader := fmt.Sprintf("%s: %s", organization.Name, organization.Description)
		organizationTableContent := [][]string{
			{"Abbreviation", "Project Name", "Description"},
		}

		// Iterate over Organizations looking for Projects
		for _, project := range scopesByScope[organization.Id] {

			projectAbbreviationToScopeMap[GenerateAbbreviation(project.Name)] = project.Id

			organizationTableContent = append(organizationTableContent, []string{
				fmt.Sprintf(Cyan+Bold+"%s"+Reset, GenerateAbbreviation(project.Name)),
				project.Name,
				project.Description,
			})
		}

		// Show the organization data when no specific project is selected.
		// Projects for this organization will appear later "inside"
		if len(args) == 0 {
			PrintTable(organizationTableHeader, organizationTableContent)
		}
	}

	// We need a project to list its targets from this point, honey
	if len(args) != 1 {
		return
	}

	fmt.Println(strings.ReplaceAll(ListProjectsCommandHeader, "\t", ""))

	// Look for the targets for desired project
	targets, err := GetScopeTargets(projectAbbreviationToScopeMap[args[0]], storedTokenReference)
	if err != nil {
		// TODO
		return
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
			fmt.Sprintf(Cyan+Bold+"%s"+Reset, target.Id),
		})
	}

	if len(projectTableContent) < 2 {
		fmt.Println(strings.ReplaceAll(ListCommandEmpty, "\t", ""))
		return
	}
	PrintTable(projectTableHeader, projectTableContent)
}
