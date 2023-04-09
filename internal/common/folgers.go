package common

import (
	"errors"
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

func MediaPath() string {
	return os.Getenv("MEDIA_PATH")
}

func AvatarFolder() string {
	return os.Getenv("AVATARS_FOLDER")
}

func RecordsFolder() string {
	return os.Getenv("RECORDS_FOLDER")
}

func InitPaths() error {
	mediaPath := MediaPath()
	avatarsFolder := AvatarFolder()
	recordsFolder := RecordsFolder()

	if mediaPath == "" {
		return errors.New("MEDIA_PATH isn't set")
	}

	if avatarsFolder == "" {
		return errors.New("AVATARS_FOLDER isn't set")
	}

	if recordsFolder == "" {
		return errors.New("RECORDS_FOLDER isn't set")
	}

	var dirForUserAvatars = mediaPath + avatarsFolder
	if err := os.MkdirAll(dirForUserAvatars, os.ModePerm); err != nil {
		return fmt.Errorf("can't create dir to save avatars: %w", err)
	}

	var dirForTracks = mediaPath + recordsFolder
	if err := os.MkdirAll(dirForTracks, os.ModePerm); err != nil {
		return fmt.Errorf("can't create dir for tracks: %w", err)
	}

	return nil
}

func init() {
	godotenv.Load()
}
