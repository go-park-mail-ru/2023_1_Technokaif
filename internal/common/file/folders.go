package file

import (
	"fmt"
	"os"
	"path/filepath"

	valid "github.com/asaskevich/govalidator"
	"github.com/joho/godotenv"
)

const (
	mediaPathParam           = "MEDIA_PATH"
	avatarFolderParam        = "AVATARS_FOLDER"
	recordsFolderParam       = "RECORDS_FOLDER"
	playlistCoverFolderParam = "PLAYLIST_COVERS_FOLDER"
)

var paths = struct {
	mediaPath           string `valid:"required,not_empty"`
	avatarsFolder       string `valid:"required,not_empty"`
	recordsFolder       string `valid:"required,not_empty"`
	playlistCoverFolder string `valid:"required,not_empty"`
}{}

func MediaPath() string {
	return paths.mediaPath
}

func AvatarFolder() string {
	return paths.avatarsFolder
}

func RecordsFolder() string {
	return paths.recordsFolder
}

func PlaylistCoverFolder() string {
	return paths.playlistCoverFolder
}

func InitPaths() error {
	if _, err := valid.ValidateStruct(paths); err != nil {
		return err
	}

	var dirForUserAvatars = filepath.Join(paths.mediaPath, paths.avatarsFolder)
	if err := os.MkdirAll(dirForUserAvatars, os.ModePerm); err != nil {
		return fmt.Errorf("can't create dir to save avatars: %w", err)
	}

	var dirForTracks = filepath.Join(paths.mediaPath, paths.recordsFolder)
	if err := os.MkdirAll(dirForTracks, os.ModePerm); err != nil {
		return fmt.Errorf("can't create dir for tracks: %w", err)
	}

	var dirForPlaylistCovers = filepath.Join(paths.mediaPath, paths.playlistCoverFolder)
	if err := os.MkdirAll(dirForPlaylistCovers, os.ModePerm); err != nil {
		return fmt.Errorf("can't create dir for playlists: %w", err)
	}

	return nil
}

func init() {
	godotenv.Load()
	paths.mediaPath = os.Getenv(mediaPathParam)
	paths.avatarsFolder = os.Getenv(avatarFolderParam)
	paths.recordsFolder = os.Getenv(recordsFolderParam)
	paths.playlistCoverFolder = os.Getenv(playlistCoverFolderParam)

	valid.TagMap["not_empty"] = valid.Validator(func(str string) bool {
		return !(valid.Trim(str, " ") == "")
	})
}
