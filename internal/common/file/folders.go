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
	MediaPath           string `valid:"required,not_empty"`
	AvatarsFolder       string `valid:"required,not_empty"`
	RecordsFolder       string `valid:"required,not_empty"`
	PlaylistCoverFolder string `valid:"required,not_empty"`
}{}

func MediaPath() string {
	return paths.MediaPath
}

func AvatarFolder() string {
	return paths.AvatarsFolder
}

func RecordsFolder() string {
	return paths.RecordsFolder
}

func PlaylistCoverFolder() string {
	return paths.PlaylistCoverFolder
}

func InitPaths() error {
	if _, err := valid.ValidateStruct(paths); err != nil {
		return err
	}

	var dirForUserAvatars = filepath.Join(paths.MediaPath, paths.AvatarsFolder)
	if err := os.MkdirAll(dirForUserAvatars, os.ModePerm); err != nil {
		return fmt.Errorf("can't create dir to save avatars: %w", err)
	}

	var dirForTracks = filepath.Join(paths.MediaPath, paths.RecordsFolder)
	if err := os.MkdirAll(dirForTracks, os.ModePerm); err != nil {
		return fmt.Errorf("can't create dir for tracks: %w", err)
	}

	var dirForPlaylistCovers = filepath.Join(paths.MediaPath, paths.PlaylistCoverFolder)
	if err := os.MkdirAll(dirForPlaylistCovers, os.ModePerm); err != nil {
		return fmt.Errorf("can't create dir for playlists: %w", err)
	}

	return nil
}

func init() {
	_ = godotenv.Load()

	paths.MediaPath = os.Getenv(mediaPathParam)
	paths.AvatarsFolder = os.Getenv(avatarFolderParam)
	paths.RecordsFolder = os.Getenv(recordsFolderParam)
	paths.PlaylistCoverFolder = os.Getenv(playlistCoverFolderParam)

	valid.TagMap["not_empty"] = valid.Validator(func(str string) bool {
		return !(valid.Trim(str, " ") == "")
	})
}
