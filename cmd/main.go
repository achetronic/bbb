package main

import (
	"bbb/internal/fancy"
	"bbb/internal/globals"
	"context"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/google/go-github/v63/github"

	"bbb/internal/cmd"
)

func checkEnv() (err error) {
	boundaryAddress := os.Getenv("BOUNDARY_ADDR")

	if boundaryAddress == "" {
		return errors.New("BOUNDARY_ADDR environment variable not set")
	}

	return err
}

func main() {
	ctx := context.Background()
	baseName := filepath.Base(os.Args[0])

	err := checkEnv()
	if err != nil {
		fancy.Fatalf(globals.UnexpectedErrorMessage, err.Error())
	}

	// Check Boundary CLI existence and give some suggestions when not present
	_, err = exec.LookPath("boundary")
	if err != nil {

		fancy.Printf(globals.BoundaryCliNotFoundErrorMessage)

		// Ask GitHub for the latest version of Boundary CLI
		client := github.NewClient(nil)

		release, _, err := client.Repositories.GetLatestRelease(ctx, "hashicorp", "boundary")
		if err != nil {
			os.Exit(0)
		}

		releaseVersion := strings.ReplaceAll(*(release.TagName), "v", "")
		packageName := "boundary_" + releaseVersion + "_" + runtime.GOOS + "_" + runtime.GOARCH + ".zip"
		packageUrl := "https://releases.hashicorp.com/boundary/" + releaseVersion + "/" + packageName

		fancy.Fatalf(globals.BoundaryCliNotFoundSuggestionErrorMessage, runtime.GOOS, runtime.GOARCH, packageUrl)
	}

	err = cmd.NewRootCommand(baseName).Execute()
	cmd.CheckError(err)
}
