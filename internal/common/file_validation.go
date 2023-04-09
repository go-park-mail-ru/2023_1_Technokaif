package common

import (
	"errors"
	"io"
	"net/http"
)

func CheckMimeType(file io.ReadSeeker, correctTypes ...string) (string, error) {
	if _, err := file.Seek(0, 0); err != nil {
		return "", err
	}
	
	fileHeader := make([]byte, 512)

	if _, err := file.Read(fileHeader); err != nil {
		return "", err
	}

	fileType := http.DetectContentType(fileHeader)
	for _, correctType := range correctTypes {
		if correctType == fileType {
			return fileType, nil
		}
	}

	return fileType, errors.New("incorrect type")
}
