package kube

const (

	// CommandArgsNoTargetErrorMessage is the message thrown when the user is trying to establish a connection
	// without defining a target
	CommandArgsNoTargetErrorMessage = `
	{Red}Impossible to get target ID from arguments.

	{White}To connect to a target, you need to connect using {Cyan}Target ID {White}as follows:
	{Bold}{White}console ~$ {Green}bt connect kube {Cyan}{ttcp_example}`

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

	// NotKubeTargetErrorMessage message is thrown when user is trying to use 'kube' connection type on non-kube targets
	NotKubeTargetErrorMessage = `
	{Red}Selected target is not configured as Kubernetes target.
	
	{Magenta}Contact your H.Boundary administrators and may the force be with you`

	// ConnectionSuccessfulMessage represents the message thrown when everything finished as expected
	// and shows how to use recently created (connection + kubeconfig) to the user
	ConnectionSuccessfulMessage = `
	{Magenta}Some reminders: {Reset}

	* Your H.Boundary session ID is: {Bold}{Cyan}%s {Reset}
	* Your session will expire in: {Bold}{Cyan}%s {Reset}


	Execute the following command to kill the connection once you don't need it: 	
	{Bold}{White}console ~$ {Green}%s {Reset}

	Execute the following command to kill ALL the connections at once:
	{Bold}{White}console ~$ {Green}%s {Reset}


	{Magenta}You are ready to query your Kubernetes cluster using command as follows: {Reset}
	{Bold}{White}console ~$ {Green}%s`
)
