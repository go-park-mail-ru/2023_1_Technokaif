package common

import (
	"errors"
	"io"
	"net/http"
)

func CheckMimeType(file io.ReadSeeker, correctTypes ...string) (string, error) {
	curPosition, err := file.Seek(0, 1) // save current position
	if err != nil {
		return "", err
	}

	if _, err := file.Seek(0, 0); err != nil {
		return "", err
	}
	
	var fileHeader [512]byte
	if _, err := file.Read(fileHeader[:]); err != nil {
		return "", err
	}

	fileType := http.DetectContentType(fileHeader[:])
	for _, correctType := range correctTypes {
		if correctType == fileType {
			return fileType, nil
		}
	}

	if _, err = file.Seek(curPosition, 0); err != nil { // go back to curPosition
		return "", err
	}

	return fileType, errors.New("incorrect type")
}
