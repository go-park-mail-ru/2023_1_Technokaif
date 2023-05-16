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
	beginPosition, err := file.Seek(0, io.SeekCurrent) // save begin position
	if err != nil {
		return "", "", fmt.Errorf("can't do file seek: %w", err)
	}

	filename, err = FileHash(file, extension)
	if err != nil {
		return "", "", fmt.Errorf("can't get fil hash: %w", err)
	}

	path = filepath.Join(dirPath, filename)

	_, err = os.Stat(path)
	if os.IsNotExist(err) {
		newFD, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, 0644)
		if err != nil {
			return "", "", fmt.Errorf("can't create file to save avatar: %w", err)
		}
		defer newFD.Close()

		if _, err := io.Copy(newFD, file); err != nil {
			return "", "", fmt.Errorf("can't write sent avatar to file: %w", err)
		}
	} else if err != nil {
		return "", "", fmt.Errorf("can't check file stat: %w", err)
	}

	if _, err = file.Seek(beginPosition, io.SeekStart); err != nil { // go back to beginPosition
		return "", "", fmt.Errorf("can't do file seek: %w", err)
	}

	return filename, path, nil
}

func FileHash(file io.ReadSeeker, extension string) (string, error) {
	beginPosition, err := file.Seek(0, io.SeekCurrent) // save begin position
	if err != nil {
		return "", fmt.Errorf("can't do file seek: %w", err)
	}

	// Create standard filename
	hasher := sha256.New()
	if _, err := io.Copy(hasher, file); err != nil {
		return "", fmt.Errorf("can't write sent avatar to hasher: %w", err)
	}
	newFileName := hex.EncodeToString(hasher.Sum(nil))

	if _, err = file.Seek(beginPosition, io.SeekStart); err != nil { // go back to beginPosition
		return "", fmt.Errorf("can't do file seek: %w", err)
	}

	return newFileName + extension, nil
}
