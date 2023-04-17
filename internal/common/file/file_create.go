package file

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func CreateFile(file io.ReadSeeker, extension string, dirPath string) (filename string, path string, _ error) {
	// Create standard filename
	hasher := sha256.New()
	if _, err := io.Copy(hasher, file); err != nil {
		return "", "", fmt.Errorf("can't write sent avatar to hasher: %w", err)
	}
	newFileName := hex.EncodeToString(hasher.Sum(nil))

	if _, err := file.Seek(0, 0); err != nil {
		return "", "", fmt.Errorf("can't do file seek: %w", err)
	}

	filename = newFileName + extension

	path = filepath.Join(dirPath, filename)

	if _, err := os.Stat(path); os.IsNotExist(err) {
		newFD, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			return "", "", fmt.Errorf("can't create file to save avatar: %w", err)
		}
		defer newFD.Close()

		if _, err := io.Copy(newFD, file); err != nil {
			return "", "", fmt.Errorf("can't write sent avatar to file: %w", err)
		}
	}

	return filename, path, nil
}
