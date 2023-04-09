package common

import (
	"fmt"
	"errors"
	"os"
)

var (
	mediaPath    = ""
	avatarFolder = ""
	recordsFolder = ""
)

func MediaPath() string {
	return mediaPath
}

func AvatarFolder() string {
	return avatarFolder
}

func RecordsFolder() string {
	return recordsFolder
}

func InitPaths() error {
	mediaPath = os.Getenv("MEDIA_PATH")
	if mediaPath == "" {
		return errors.New("MEDIA_PATH isn't set")
	}

	avatarFolder = os.Getenv("AVATAR_FOLDER")
	if avatarFolder == "" {
		return errors.New("AVATAR_FOLDER isn't set")
	}

	recordsFolder = os.Getenv("TRACK_FOLDER")
	if recordsFolder == "" {
		return errors.New("TRACK_FOLDER isn't set")
	}

	var dirForUserAvatars = mediaPath + avatarFolder
	if err := os.MkdirAll(dirForUserAvatars, os.ModePerm); err != nil {
		return fmt.Errorf("can't create dir to save avatars: %w", err)
	}

	var dirForTracks = mediaPath + recordsFolder
	if err := os.MkdirAll(dirForTracks, os.ModePerm); err != nil {
		return fmt.Errorf("can't create dir for tracks: %w", err)
	}

	return nil
}
