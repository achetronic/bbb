package globals

const (
	//
	BoundaryCliNotFoundErrorMessage = `
	{Red}Boundary CLI is not detected on your system.

	{Magenta}This is a CLI that improve a lot of aspects of the original CLI,
	but it works on top of it, so sadly, it needs it... *sniff, sniff* {Cyan}(yet)

	{White}It is possible to remediate this issue installing it from:

	{White}Pretty website: {Reset}{Cyan}https://developer.hashicorp.com/boundary/install
	{White}Experts website: {Reset}{Cyan}https://releases.hashicorp.com/boundary
	`

	//
	BoundaryCliNotFoundSuggestionErrorMessage = `
	{Magenta}May be I can improve your day. Let me think... 
	Is it your desire to have a direct URL to download it?

	{Bold}{White}Detected OS: {Reset}{Cyan}%s
	{Bold}{White}Detected Architecture: {Reset}{Cyan}%s

	{White}Direct link: {Reset}{Cyan}%s
	`

	//
	UnexpectedErrorMessage = `
	{Red}There was an unexpected error.

	{White}What happened under the hood:
	{Italic}%s`

	// TokenRetrievalErrorMessage message is thrown when BOUNDARY_TOKEN can not be retrieved neither from env nor filesystem
	TokenRetrievalErrorMessage = `
	{Red}Impossible to get the authentication token from your system.

	{White}It is possible to remediate this performing auth step again:
	{Bold}{White}console ~$ {Green}bbb auth`
)
