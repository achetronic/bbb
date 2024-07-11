package globals

import (
	"errors"
	"fmt"
	"os"
	"time"
)

var BbbTemporaryDir = os.TempDir() + "/bbb"

// GetStoredToken TODO
func GetStoredToken() (token string, err error) {
	fileContent, err := os.ReadFile(BbbTemporaryDir + "/BOUNDARY_TOKEN")
	token = string(fileContent)
	if token == "" {
		err = errors.New("no token found")
	}
	return token, err
}

// StoreToken stores the token into a temporary file
func StoreToken(token string) (err error) {

	//
	err = os.MkdirAll(BbbTemporaryDir, 0700)
	if err != nil {
		return err
	}

	//
	err = os.WriteFile(BbbTemporaryDir+"/BOUNDARY_TOKEN", []byte(token), 0700)
	return err
}

// GetStoredTokenReference TODO
func GetStoredTokenReference() (storedTokenReference string, err error) {

	storedTokenReference = "env://BOUNDARY_TOKEN"
	storedToken := os.Getenv("BOUNDARY_TOKEN")

	if storedToken == "" {

		storedToken, err = GetStoredToken()
		if err != nil {
			return storedTokenReference, err
		}

		storedTokenReference = "file://" + BbbTemporaryDir + "/BOUNDARY_TOKEN"
	}

	return storedTokenReference, err
}

// GetFileContents TODO
func GetFileContents(filePath string, waitBetweenChecks bool) (content []byte, err error) {

	fileEmpty := true
	var fileContentRaw []byte

	for loop := 0; loop <= 10 && fileEmpty == true; loop++ {
		fileContentRaw, err = os.ReadFile(filePath)
		if err != nil {
			err = errors.New("error reading file '" + filePath + "': " + err.Error())
			return content, err
		}

		if len(fileContentRaw) > 0 {
			fileEmpty = false
		}

		if waitBetweenChecks {
			time.Sleep(500 * time.Millisecond)
		}
	}

	if fileEmpty {
		err = errors.New("there is no content in '" + filePath + "' after several reading retries")
	}

	return fileContentRaw, err
}

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
