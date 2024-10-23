package upgrade

import (
	"bbb/internal/cmd/version"
	"bbb/internal/fancy"
	"bbb/internal/globals"
	"bufio"
	"bytes"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
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

// RunCommand TODO
func RunCommand(cmd *cobra.Command, args []string) {

	// 0. Ask the user
	fancy.Printf(UpgradeConfirmationMessage)

	inScanner := bufio.NewScanner(os.Stdin)
	inScanner.Scan()
	userAnswer := inScanner.Text()

	if userAnswer != "yes" {
		return
	}

	// 1. Retrieve release data
	latestReleaseResp, err := http.Get(latestReleaseDataUrl)
	if err != nil {
		fancy.Fatalf(globals.UnexpectedErrorMessage,
			fmt.Sprintf("failed getting release data: %s", err.Error()))
	}

	defer latestReleaseResp.Body.Close()

	//
	latestReleaseBodyBytes, err := io.ReadAll(latestReleaseResp.Body)
	if err != nil {
		fancy.Fatalf(globals.UnexpectedErrorMessage,
			fmt.Sprintf("failed reading JSON from releases page: %s", err.Error()))
	}

	releaseObj := &github.RepositoryRelease{}
	err = json.Unmarshal(latestReleaseBodyBytes, releaseObj)
	if err != nil {
		fancy.Fatalf(globals.UnexpectedErrorMessage,
			fmt.Sprintf("failed to decode JSON structure from release: %s", err.Error()))
	}

	// 2. Select proper asset and its checksum
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

	// 3. Download the assets
	assetResp, err := http.Get(assetUrl)
	if err != nil {
		fancy.Fatalf(globals.UnexpectedErrorMessage,
			fmt.Sprintf("failed downloading the asset: %s", err.Error()))
	}
	defer assetResp.Body.Close()

	assetChecksumResp, err := http.Get(assetChecksumUrl)
	if err != nil {
		fancy.Fatalf(globals.UnexpectedErrorMessage,
			fmt.Sprintf("failed downloading the asset's checksum: %s", err.Error()))
	}
	defer assetChecksumResp.Body.Close()

	// 4. Calculate package checksum and compare with published checksum
	publicChecksum, err := io.ReadAll(assetChecksumResp.Body)
	if err != nil {
		fancy.Fatalf(globals.UnexpectedErrorMessage,
			fmt.Sprintf("failed loading asset's checksum into memory %s", err.Error()))
	}

	calculatedChecksum := md5.New()

	assetBytes := bytes.Buffer{}
	assetBytesWriterEntity := io.Writer(&assetBytes)

	multiWriterEntity := io.MultiWriter(assetBytesWriterEntity, calculatedChecksum)

	_, err = io.Copy(multiWriterEntity, assetResp.Body)
	if err != nil {
		fancy.Fatalf(globals.UnexpectedErrorMessage,
			fmt.Sprintf("failed copying bytes to multiwriter: %s", err.Error()))
	}

	if strings.TrimSpace(string(publicChecksum)) != fmt.Sprintf("%x", calculatedChecksum.Sum(nil)) {
		fancy.Fatalf(globals.UnexpectedErrorMessage, "checksums doesn't match. Is the release corrupted?")
	}

	// 5. Place downloaded files in a temporary place
	tmpDirPath, err := os.MkdirTemp("", "bbb-*")
	if err != nil {
		fancy.Fatalf(globals.UnexpectedErrorMessage,
			fmt.Sprintf("failed creating temporary directory: %s", err.Error()))
	}

	err = UnTarGz(bytes.NewReader(assetBytes.Bytes()), tmpDirPath)
	if err != nil {
		fancy.Fatalf(globals.UnexpectedErrorMessage,
			fmt.Sprintf("failed un-extracting downloaded asset from release: %s", err.Error()))
	}

	// 6. Perform binary upgrade
	binaryBytes, err := os.ReadFile(filepath.Join(tmpDirPath, "bbb"))
	if err != nil {
		fancy.Fatalf(globals.UnexpectedErrorMessage,
			fmt.Sprintf("failed loading final binary in memory: %s", err.Error()))
	}

	err = selfupdate.PrepareAndCheckBinary(bytes.NewReader(binaryBytes), selfupdate.Options{})
	if err != nil {
		fancy.Fatalf(globals.UnexpectedErrorMessage,
			fmt.Sprintf("failed preparing and checking downloaded binary: %s", err.Error()))
	}

	err = selfupdate.CommitBinary(selfupdate.Options{})
	if err != nil {
		fancy.Fatalf(globals.UnexpectedErrorMessage,
			fmt.Sprintf("failed replacing prior binary with the new one: %s", err.Error()))
	}

	fancy.Printf(UpgradeSuccessfulMessage, version.Version)

}
