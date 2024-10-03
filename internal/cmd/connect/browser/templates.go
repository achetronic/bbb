package browser

const (
	BrowserCliNotFoundErrorMessage = `
	{Red}A browser opener CLI is not detected on your system.

	{Magenta}How do you expect to use a Browser connection without having a CLI in your system?
	C'mon everybody has xdg-open or open in their systems, are you using a toaster?

	{White}It is possible to remediate this issue installing it from:

	{White}For Linux users: {Reset}{Cyan}Use your package manager and install 'xdg-utils'
	{White}For MacOS users: {Reset}{Cyan}Nah I don't belive you, you already have it`

	// CommandArgsNoTargetErrorMessage is the message thrown when the user is trying to establish a connection
	// without defining a target
	CommandArgsNoTargetErrorMessage = `
	{Red}Impossible to get target ID from arguments.

	{White}To connect to a target, you need to connect using {Cyan}Target ID {White}as follows:
	{Bold}{White}console ~$ {Green}bbb connect browser {Cyan}{ttcp_example} [--insecure] [--port <port>]`

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

	// NotBrowserTargetErrorMessage message is thrown when user is trying to use 'browser' connection type on non-browser targets
	NotBrowserTargetErrorMessage = `
	{Red}Selected target is not configured as Browser target.

	{Magenta}To connect Browser its needed to have one of the following combinations:
	* Username & Password to use basic authentication
	* Password to use bearer authentication

	Contact your H.Boundary administrators and may the force be with you`

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


	{Magenta}Your browser will automatically open pointing to desired website.
	If browser does not open automatically, open it. It's all about your computer being lazy. {Reset}
	{Bold}{White}Website: {Green}%s`
)
