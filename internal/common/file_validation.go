package common

import (
	"errors"
	"fmt"
	"io"
	"net/http"
)

func CheckMimeType(file io.ReadSeeker, correctTypes ...string) (string, error) {
	
	fileHeader := make([]byte, 512)

	if _, err := file.Read(fileHeader); err != nil {
		return "", err
	}

	if _, err := file.Seek(0, 0); err != nil {
		return "", err
	}

	fmt.Println()
	fileType := http.DetectContentType(fileHeader)
	for _, correctType := range correctTypes {
		if correctType == fileType {
			return fileType, nil
		}
	}

	return fileType, errors.New("incorrect type")
} 