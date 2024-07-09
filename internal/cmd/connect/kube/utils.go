package kube

import (
	"bt/internal/cmd/list"
	"errors"
	"fmt"
	"time"
)

const (
	// TODO
	ConnectKubeFinalMessageContent = `
	` + list.Magenta + `Some reminders: ` + list.Reset + `

	* Your H.Boundary session ID is: ` + list.Cyan + list.Bold + `%s` + list.Reset + `
	* Your session will expire in: ` + list.Cyan + list.Bold + `%s` + list.Reset + `


	Execute the following command to kill the connection once you don't need it: 	
	` + list.White + list.Bold + `console ~$ ` + list.Green + `%s` + list.Reset + `

	Execute the following command to kill ALL the connections at once:
	` + list.White + list.Bold + `console ~$ ` + list.Green + `%s` + list.Reset + `



	` + list.Magenta + `You are ready to query your Kubernetes cluster using command as follows:` + list.Reset + `

	` + list.White + list.Bold + `console ~$ ` + list.Green + `%s` + list.Reset
)

// GetDurationStringFromNow returns a string with a duration representation.
// That string is the lasting time between present moment and a date given as argument.
// If the date is in the past, it's returned as error
func GetDurationStringFromNow(date string) (duration string, err error) {

	dateParsed, err := time.Parse(time.RFC3339, date)
	if err != nil {
		return duration, err
	}

	//
	now := time.Now()
	if expiredDate := now.After(dateParsed); expiredDate {
		return duration, errors.New("date is in the past")
	}

	//
	durationRaw := dateParsed.Sub(now)
	duration = fmt.Sprintf("%dD %dH %dm %ds",
		int(durationRaw.Hours()/24),
		int(durationRaw.Hours())%24,
		int(durationRaw.Minutes())%60,
		int(durationRaw.Seconds())%60)

	return duration, err
}
