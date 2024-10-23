package upgrade

import (
	"bbb/internal/fancy"
	"bbb/internal/globals"
	"bytes"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/minio/selfupdate"
	"github.com/spf13/cobra"

	"github.com/google/go-github/v66/github"
)

const (
	descriptionShort = `Upgrade BBB to the latest version`

	descriptionLong = `
	Upgrade BBB to the latest version.`

	//
	latestReleaseDataUrl = "https://api.github.com/repos/achetronic/bbb/releases/latest"
)

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "upgrade",
		Short: descriptionShort,
		Long:  strings.ReplaceAll(descriptionLong, "\t", ""),

		Run: RunCommand,
	}

	return cmd
}

// DISCLAIMER: THIS COMMAND IS A WORK IN PROGRESS
func RunCommand(cmd *cobra.Command, args []string) {

	// Request release data
	latestReleaseResp, err := http.Get(latestReleaseDataUrl)
	if err != nil {
		fancy.Fatalf(globals.UnexpectedErrorMessage,
			"failed getting release data: %s", err.Error())
	}

	defer latestReleaseResp.Body.Close()

	//
	latestReleaseBodyBytes, err := io.ReadAll(latestReleaseResp.Body)
	if err != nil {
		fancy.Fatalf(globals.UnexpectedErrorMessage, "failed reading JSON from releases page: %s", err.Error())
	}

	releaseObj := &github.RepositoryRelease{}
	err = json.Unmarshal(latestReleaseBodyBytes, releaseObj)
	if err != nil {
		fancy.Fatalf(globals.UnexpectedErrorMessage, "failed to decode JSON structure from release")
	}

	// Select proper asset and its checksum
	var assetUrl string
	var assetChecksumUrl string
	for _, asset := range releaseObj.Assets {

		if !strings.Contains(*asset.Name, runtime.GOOS) || !strings.Contains(*asset.Name, runtime.GOARCH) {
			continue
		}

		if strings.Contains(*asset.Name, "tar.gz.md5") {
			assetChecksumUrl = *asset.BrowserDownloadURL
			continue
		}

		assetUrl = *asset.BrowserDownloadURL
	}

	// Download the assets
	assetResp, err := http.Get(assetUrl)
	if err != nil {
		fancy.Fatalf(globals.UnexpectedErrorMessage, "failed downloading the asset: %s", err.Error())
	}
	defer assetResp.Body.Close()

	assetChecksumResp, err := http.Get(assetChecksumUrl)
	if err != nil {
		fancy.Fatalf(globals.UnexpectedErrorMessage, "failed downloading the asset's checksum: %s", err.Error())
	}
	defer assetChecksumResp.Body.Close()

	// Calculate package checksum and compare with published checksum
	publicChecksum, err := io.ReadAll(assetChecksumResp.Body)
	if err != nil {
		fancy.Fatalf(globals.UnexpectedErrorMessage, "failed loading asset's checksum into memory %s", err.Error())
	}

	calculatedChecksum := md5.New()

	assetBytes := bytes.Buffer{}
	assetBytesWriterEntity := io.Writer(&assetBytes)

	multiWriterEntity := io.MultiWriter(assetBytesWriterEntity, calculatedChecksum)

	_, err = io.Copy(multiWriterEntity, assetResp.Body)
	if err != nil {
		fancy.Fatalf(globals.UnexpectedErrorMessage, "failed copying bytes to multiwriter: %s", err.Error())
	}

	if strings.TrimSpace(string(publicChecksum)) != fmt.Sprintf("%x", calculatedChecksum.Sum(nil)) {
		fancy.Fatalf(globals.UnexpectedErrorMessage, "checksums doesn't match. Is the release corrupted?")
	}

	// Place downloaded files in a temporary place
	tmpDirPath, err := os.MkdirTemp("", "bbb-*")
	if err != nil {
		fancy.Fatalf(globals.UnexpectedErrorMessage, "failed creating temporary directory: %s", err.Error())
	}

	err = UnTarGz(bytes.NewReader(assetBytes.Bytes()), tmpDirPath)
	if err != nil {
		fancy.Fatalf(globals.UnexpectedErrorMessage, "failed un-extracting downloaded asset from release: %s", err.Error())
	}

	// Perform binary upgrade
	binaryBytes, err := os.ReadFile(filepath.Join(tmpDirPath, "bbb"))
	if err != nil {
		fancy.Fatalf(globals.UnexpectedErrorMessage, "failed loading final binary in memory: %s", err.Error())
	}

	err = selfupdate.PrepareAndCheckBinary(bytes.NewReader(binaryBytes), selfupdate.Options{})
	if err != nil {
		fancy.Fatalf(globals.UnexpectedErrorMessage, "failed preparing and checking downloaded binary: %s", err.Error())
	}

	err = selfupdate.CommitBinary(selfupdate.Options{})
	if err != nil {
		fancy.Fatalf(globals.UnexpectedErrorMessage, "failed replacing prior binary with the new one: %s", err.Error())
	}

	log.Print("Success upgrading BBB!")

}
