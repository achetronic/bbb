package browser

const (
	XdgOpenCliNotFoundErrorMessage = `
	{Red}xdg-open is not detected on your system.

	{Magenta}How do you expect to connect by SSH without having SSH CLI in your system? 
	I'm not an expert... but flying without wings is a bit... difficult *cough, cough*

	{White}It is possible to remediate this issue installing it from:

	{White}Pretty way: {Reset}{Cyan}Use your OS package manager
	{White}Experts website: {Reset}{Cyan}https://cloudflare.cdn.openbsd.org/pub/OpenBSD/OpenSSH/portable/
	`

	// CommandArgsNoTargetErrorMessage is the message thrown when the user is trying to establish a connection
	// without defining a target
	CommandArgsNoTargetErrorMessage = `
	{Red}Impossible to get target ID from arguments.

	{White}To connect to a target, you need to connect using {Cyan}Target ID {White}as follows:
	{Bold}{White}console ~$ {Green}bbb connect browser {Cyan}{ttcp_example} [--insecure]`

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
	* Username & Password or Username or Password
	
	Contact your H.Boundary administrators and may the force be with you`
)
