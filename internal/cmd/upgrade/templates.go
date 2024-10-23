package upgrade

const (

	// UpgradeConfirmationMessage represents the message thrown when ...
	UpgradeConfirmationMessage = `
	{White}You are going to {Bold}{Cyan}upgrade{Reset} BBB {Reset}

	{White}If you want to continue, please type {Bold}{Green}yes{Reset} and press {Bold}{Cyan}[ENTER]{Reset}`

	// UpgradeSuccessfulMessage represents the message thrown when everything finished as expected on 'upgrade' command
	UpgradeSuccessfulMessage = `
	{White}You have {Bold}{Cyan}successfully{Reset} upgraded BBB {Reset}


	{White}To know the version you are using now, execute following command: {Reset}
	{Bold}{White}console ~$ {Green}bbb version {Reset}

	💡 Wait, wait, wait... why should you do it when I can print it for you?
	✨ %s`
)
