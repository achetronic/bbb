package globals

import (
	"errors"
	"os"
	"strings"
)

var BtTemporaryDir = os.TempDir() + "/bt" // TODO

// CopyMap return a map that is a real copy of the original
// Ref: https://go.dev/blog/maps
func CopyMap(src map[string]interface{}) map[string]interface{} {
	m := make(map[string]interface{}, len(src))
	for k, v := range src {
		m[k] = v
	}
	return m
}

// SplitCommaSeparatedValues get a list of strings and return a new list
// where each element containing commas is divided in separated elements
func SplitCommaSeparatedValues(input []string) []string {
	var result []string
	for _, item := range input {
		parts := strings.Split(item, ",")
		result = append(result, parts...)
	}
	return result
}

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
