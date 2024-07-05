package globals

import (
	"errors"
	"os"
)

var BtTemporaryDir = os.TempDir() + "/bt" // TODO

// GetStoredToken TODO
func GetStoredToken() (token string, err error) {
	fileContent, err := os.ReadFile(BtTemporaryDir + "/BOUNDARY_TOKEN")
	token = string(fileContent)
	if token == "" {
		err = errors.New("no token found")
	}
	return token, err
}

// StoreToken stores the token into a temporary file
func StoreToken(token string) (err error) {

	//
	err = os.MkdirAll(BtTemporaryDir, 0700)
	if err != nil {
		return err
	}

	//
	err = os.WriteFile(BtTemporaryDir+"/BOUNDARY_TOKEN", []byte(token), 0700)
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

		storedTokenReference = "file://" + BtTemporaryDir + "/BOUNDARY_TOKEN"
	}

	return storedTokenReference, err
}
