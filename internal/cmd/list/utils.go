package list

import "bt/internal/fancy"

const (

	// TODO
	ListOrganizationsCommandHeader = `
	Following tables show the projects belonging to an organization.

	To list the targets inside, you need to list using ` + fancy.Cyan + fancy.Bold + ` 
	abbreviations ` + fancy.Reset + `as follows: 
	
	console ~$ ` + fancy.Green + fancy.Bold + `bt list ` + fancy.Cyan + fancy.Bold + `{abbreviation}` + fancy.Reset

	// TODO
	ListProjectsCommandHeader = `
	Following tables show the targets belonging to an project.

	To connect to a target, you need to connect using ` + fancy.Cyan + fancy.Bold + ` 
	Target ID ` + fancy.Reset + `as follows: 
	
	console ~$ ` + fancy.Green + fancy.Bold + `bt connect ssh ` + fancy.Cyan + fancy.Bold + `{ttcp_example}` + fancy.Reset + `
	console ~$ ` + fancy.Green + fancy.Bold + `bt connect kube ` + fancy.Cyan + fancy.Bold + `{ttcp_example}` + fancy.Reset + `

	Remember to use ` + fancy.Bold + `ssh` + fancy.Reset + ` or ` + fancy.Bold + `kube` + fancy.Reset + ` 
	subcommand depending on the target you are trying to connect to.
	`

	ListCommandEmpty = fancy.Magenta + fancy.Bold + `
	Dont you see any table? you might need some permissions.
	Contact your H.Boundary administrators and may the force be with you
	` + fancy.Reset
)
