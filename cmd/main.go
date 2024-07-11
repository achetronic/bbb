package main

import (
	"bt/internal/cmd"
	"errors"
	"log"
	"os"
	"path/filepath"
)

func checkEnv() (err error) {
	os.Setenv("BOUNDARY_ADDR", "https://hashicorp-boundary.fpkmon.com")

	boundaryAddress := os.Getenv("BOUNDARY_ADDR")

	if boundaryAddress == "" {
		return errors.New("BOUNDARY_ADDR environment variable not set")
	}

	return err
}

func main() {
	baseName := filepath.Base(os.Args[0])

	err := checkEnv()
	if err != nil {
		log.Fatal(err.Error())
	}

	err = cmd.NewRootCommand(baseName).Execute()
	cmd.CheckError(err)
}
