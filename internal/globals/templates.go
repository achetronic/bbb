package globals

const (
	UnexpectedErrorMessage = `
	{Red}There was an unexpected error.

	{White}What happened under the hood:
	{Italic}%s`

	// TokenRetrievalErrorMessage message is thrown when BOUNDARY_TOKEN can not be retrieved neither from env nor filesystem
	TokenRetrievalErrorMessage = `
	{Red}Impossible to get the authentication token from your system.

	{White}It is possible to remediate this performing auth step again:
	{Bold}{White}console ~$ {Green}bt auth`
)
