package list

const (
	// TODO
	UnexpectedErrorMessage = `
	{Red}There was an unexpected error.

	{White}What happened under the hood:
	{Italic}%s`

	// TokenRetrievalErrorMessage message is thrown when BOUNDARY_TOKEN can not be retrieved neither from env nor filesystem
	TokenRetrievalErrorMessage = `
	{Red}Impossible to get the authentication token from your system.

	{White}It is possible to remediate this performing auth step again:
	{Bold}{White}console ~$ {Green}bt auth`

	// TODO
	ListOrganizationsCommandHeader = `
	Following tables show the projects belonging to an organization.

	To list the targets inside, you need to list using {Bold}{Cyan}abbreviations{Reset} as follows:
	{Bold}{White}console ~$ {Green}bt list {Cyan}{abbreviation}`

	// TODO
	ListProjectsCommandHeader = `
	Following tables show the targets belonging to an project.

	To connect to a target, you need to connect using {Bold}{Cyan}Target ID{Reset} as follows: 
	
	{Bold}{White}console ~$ {Green}bt connect ssh {Cyan}{ttcp_example} {Reset}
	{Bold}{White}console ~$ {Green}bt connect kube {Cyan}{ttcp_example} {Reset}

	Remember to use {Bold}{Cyan}ssh{Reset} or {Bold}{Cyan}kube{Reset} subcommand
	depending on the target you are trying to connect to.`

	//
	ListCommandEmpty = `
	{Red}Dont you see any table? you might need some permissions.
	Contact your H.Boundary administrators and may the force be with you
	`
)
