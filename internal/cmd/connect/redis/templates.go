package redis

const (
	RedisCliNotFoundErrorMessage = `
	{Red}The package redis-cli is not detected on your system.

	{Magenta}If you want to open a connection directly to a Redis target,
	you need to have redis-cli installed on your system.

	{White}It is possible to remediate this issue installing it from:

	{White}For Linux users: {Reset}{Cyan}Use your package manager and install 'redis'
	{White}For MacOS users: {Reset}{Cyan}Use your package manager and install 'redis'

	{White}Install redis-cli and try again if you want ;)`

	// CommandArgsNoTargetErrorMessage is the message thrown when the user is trying to establish a connection
	// without defining a target
	CommandArgsNoTargetErrorMessage = `
	{Red}Impossible to get target ID from arguments.

	{White}To connect to a target, you need to connect using {Cyan}Target ID {White}as follows:
	{Bold}{White}console ~$ {Green}bbb connect redis {Cyan}{ttcp_example}`

	// AuthorizeSessionErrorMessage message is thrown when there is an error different from 4xx on authorize-session
	// command execution
	AuthorizeSessionErrorMessage = `
	{Red}Error executing 'authorize-session' command.

	{White}Command execution returned:
	{Italic}%s{Reset}

	{White}H.Boundary command under the hood returned:
	{Italic}%s`

	// AuthorizeSessionUserErrorMessage message is thrown when there is a 4xx error on authorize-session command execution
	AuthorizeSessionUserErrorMessage = `
	{Red}There is something wrong with your target when authorizing the session

	{Magenta}Review following points:
	* Your H.Boundary token is still valid
	* You have properly written the target id
	* You have enough permissions to authorize a session against to that target

	{White}Under the hood:
	{Italic}%s`

	//
	TargetWithNoCredentials = `
	{Magenta}Selected target has no credentials associated.

	Anyways we will try to connect to your target with no credentials, but it may fail.`

	// ConnectionSuccessfulMessage represents the message thrown when everything finished as expected
	// and shows how to use recently created (connection + webserver) to the user
	ConnectionSuccessfulMessage = `
	{Magenta}Some reminders: {Reset}

	* Your H.Boundary session ID is: {Bold}{Cyan}%s {Reset}
	* Your session will expire in: {Bold}{Cyan}%s {Reset}


	Press following combination to kill the connection once you don't need it:
	{Bold}{White}console ~$ {Green}Cntrl + C {Reset}


	{Magenta}Your redis CLI will automatically open pointing to desired URL.
	If redis-cli does not open automatically, probably you don't have the package installed
	Just install it and try again{Reset}
	{Bold}{White}Redis listening on: {Green}%s`
)
