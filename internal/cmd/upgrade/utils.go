package upgrade

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// UnTarGz extracts a .tar.gz file from the provided io.Reader (e.g., HTTP response body)
// and saves the contents to the specified target directory.
func UnTarGz(body io.Reader, target string) error {
	// Create a gzip reader from the body
	gzipReader, err := gzip.NewReader(body)
	if err != nil {
		return err
	}
	defer gzipReader.Close()

	// Create a tar reader
	tarReader := tar.NewReader(gzipReader)

	// Iterate over the files in the tar archive
	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		// Define the target path where the file will be extracted
		targetPath := filepath.Join(target, header.Name)

		switch header.Typeflag {
		case tar.TypeDir:
			// Create the directory if it doesn't exist
			if err := os.MkdirAll(targetPath, os.FileMode(header.Mode)); err != nil {
				return err
			}
		case tar.TypeReg:
			// Extract the regular file
			outFile, err := os.Create(targetPath)
			if err != nil {
				return err
			}
			defer outFile.Close()

			// Copy the file content from the tar archive
			if _, err := io.Copy(outFile, tarReader); err != nil {
				return err
			}

			// Set the file permissions
			if err := os.Chmod(targetPath, os.FileMode(header.Mode)); err != nil {
				return err
			}
		default:
			fmt.Printf("Unknown file type: %x in %s\n", header.Typeflag, header.Name)
		}
	}
	return nil
}
