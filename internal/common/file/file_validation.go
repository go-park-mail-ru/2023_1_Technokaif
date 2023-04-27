package file

import (
	"errors"
	"io"
	"net/http"
)

func CheckMimeType(file io.ReadSeeker, correctTypes ...string) (string, error) {
	curPosition, err := file.Seek(0, io.SeekCurrent) // save current position
	if err != nil {
		return "", err
	}

	var fileHeader [512]byte
	if _, err := file.Read(fileHeader[:]); err != nil {
		return "", err
	}

	if _, err = file.Seek(curPosition, io.SeekStart); err != nil { // go back to curPosition
		return "", err
	}

	fileType := http.DetectContentType(fileHeader[:])
	for _, correctType := range correctTypes {
		if correctType == fileType {
			return fileType, nil
		}
	}

	return fileType, errors.New("incorrect type")
}
