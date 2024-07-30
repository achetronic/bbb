package browser

import (
	"encoding/base64"
	"fmt"
	"math/rand/v2"
)

// GetAuthHeaderValue returns the value of the Authorization header
func GetAuthHeaderValue(username, password string, authType string) (authHeader string, err error) {

	if password == "" {
		err = fmt.Errorf("password is required for any type of authentication")
		return authHeader, err
	}

	switch authType {

	//
	case "basic":
		if username == "" {
			err = fmt.Errorf("username is required for basic authentication")
			break
		}
		authHeader = "Basic " + base64.StdEncoding.EncodeToString([]byte(username+":"+password))

	//
	case "bearer":
		authHeader = "Bearer " + password

	//
	default:
		err = fmt.Errorf("unknown authentication method: %s", authType)
	}

	return authHeader, err
}

// GetFreeRandomPort returns a random port number between min and max
func GetFreeRandomPort(min, max int) (port int) {

	port = rand.IntN(max-min+1) + min

	// TODO: logic to check if the port is truly free
	return port
}
