package usecase

import (
	"errors"
	"fmt"
	"io"
	"path/filepath"

	commonFile "github.com/go-park-mail-ru/2023_1_Technokaif/internal/common/file"
	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/models"
	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/playlist"
	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/track"
	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/user"
	"github.com/go-park-mail-ru/2023_1_Technokaif/pkg/logger"
)

// Usecase implements album.Usecase
type Usecase struct {
	playlistRepo playlist.Repository
	trackRepo    track.Repository
	userRepo     user.Repository
	logger       logger.Logger
}

func NewUsecase(pr playlist.Repository, tr track.Repository, ur user.Repository, l logger.Logger) *Usecase {
	return &Usecase{
		playlistRepo: pr,
		trackRepo:    tr,
		userRepo:     ur,

		logger: l}
}

func (u *Usecase) Create(playlist models.Playlist, usersID []uint32, userID uint32) (uint32, error) {
	userInAuthors := false
	for _, uid := range usersID {
		user, err := u.userRepo.GetByID(uid)
		if err != nil {
			return 0, fmt.Errorf("(usecase) can't get user with id #%d: %w", uid, err)
		}
		if user.ID == userID {
			userInAuthors = true
			break
		}
	}
	if !userInAuthors {
		return 0, fmt.Errorf("(usecase) playlist can't be created by user: %w", &models.ForbiddenUserError{})
	}

	playlistID, err := u.playlistRepo.Insert(playlist, usersID)
	if err != nil {
		return 0, fmt.Errorf("(usecase) can't insert playlist into repository: %w", err)
	}

	return playlistID, nil
}

func (u *Usecase) GetByID(playlistID uint32) (*models.Playlist, error) {
	playlist, err := u.playlistRepo.GetByID(playlistID)
	if err != nil {
		return nil, fmt.Errorf("(usecase) can't get playlist from repository: %w", err)
	}

	return playlist, nil
}

func (u *Usecase) UpdateInfoAndMembers(playlist models.Playlist, usersID []uint32, userID uint32) error {
	pl, err := u.playlistRepo.GetByID(playlist.ID)
	if err != nil {
		return fmt.Errorf("(usecase) can't find playlist in repository: %w", err)
	}
	playlist.CoverSrc = pl.CoverSrc

	userInAuthors, err := u.checkUserInAuthors(playlist.ID, userID)
	if err != nil {
		return err
	}
	if !userInAuthors {
		return fmt.Errorf("(usecase) playlist can't be deleted by user: %w", &models.ForbiddenUserError{})
	}

	authors, err := u.userRepo.GetByPlaylist(playlist.ID)
	if err != nil {
		return fmt.Errorf("(usecase) can't get authors of playlist: %w", err)
	}
	authorsMap := make(map[uint32]struct{}, len(authors))
	for _, a := range authors {
		authorsMap[a.ID] = struct{}{}
	}

	newAuthorsID := make([]uint32, 0)
	for _, uid := range usersID {
		if _, ok := authorsMap[uid]; !ok {
			newAuthorsID = append(newAuthorsID, uid)
		}
	}

	if err := u.playlistRepo.UpdateWithMembers(playlist, newAuthorsID); err != nil {
		return fmt.Errorf("(usecase) can't update playlist in repository: %w", err)
	}

	return nil
}

var dirForPlaylistCovers = filepath.Join(commonFile.MediaPath(), commonFile.PlaylistCoverFolder())

var ErrCoverWrongFormat = errors.New("wrong cover file fromat")

func (u *Usecase) UploadCoverWrongFormatError() error {
	return ErrCoverWrongFormat
}

func (u *Usecase) UploadCover(playlistID uint32, userID uint32, file io.ReadSeeker, fileExtension string) error {
	playlist, err := u.playlistRepo.GetByID(playlistID)
	if err != nil {
		return fmt.Errorf("(usecase) can't find playlist: %w", err)
	}

	userInAuthors, err := u.checkUserInAuthors(playlistID, userID)
	if err != nil {
		return err
	}
	if !userInAuthors {
		return fmt.Errorf("(usecase) playlist can't be deleted by user: %w", &models.ForbiddenUserError{})
	}

	// Check format
	if fileType, err := commonFile.CheckMimeType(file, "image/png", "image/jpeg"); err != nil {
		return fmt.Errorf("(usecase) file format %s: %w", fileType, ErrCoverWrongFormat)
	}
	if _, err := file.Seek(0, 0); err != nil {
		return fmt.Errorf("(usecase) can't do file seek: %w", err)
	}

	filenameWithExtension, _, err := commonFile.CreateFile(file, fileExtension, dirForPlaylistCovers)
	if err != nil {
		return fmt.Errorf("(usecase) can't create file: %w", err)
	}

	playlist.CoverSrc = filepath.Join(commonFile.PlaylistCoverFolder(), filenameWithExtension)
	if err := u.playlistRepo.Update(*playlist); err != nil {
		return fmt.Errorf("(usecase) can't update playlist: %w", err)
	}
	return nil
}

func (u *Usecase) Delete(playlistID uint32, userID uint32) error {
	if _, err := u.playlistRepo.GetByID(playlistID); err != nil {
		return fmt.Errorf("(usecase) can't find playlist in repository: %w", err)
	}

	userInAuthors, err := u.checkUserInAuthors(playlistID, userID)
	if err != nil {
		return err
	}
	if !userInAuthors {
		return fmt.Errorf("(usecase) playlist can't be deleted by user: %w", &models.ForbiddenUserError{})
	}

	if err := u.playlistRepo.DeleteByID(playlistID); err != nil {
		return fmt.Errorf("(usecase) can't delete playlist from repository: %w", err)
	}

	return nil
}

func (u *Usecase) AddTrack(trackID, playlistID, userID uint32) error {
	if _, err := u.playlistRepo.GetByID(playlistID); err != nil {
		return fmt.Errorf("(usecase) can't find playlist in repository: %w", err)
	}

	if _, err := u.trackRepo.GetByID(trackID); err != nil {
		return fmt.Errorf("(usecase) can't find track in repository: %w", err)
	}

	userInAuthors, err := u.checkUserInAuthors(playlistID, userID)
	if err != nil {
		return err
	}
	if !userInAuthors {
		return fmt.Errorf("(usecase) playlist can't be updated by user: %w", &models.ForbiddenUserError{})
	}

	if err := u.playlistRepo.AddTrack(trackID, playlistID); err != nil {
		return fmt.Errorf("(usecase) can't add track into playlist in repository: %w", err)
	}

	return nil
}

func (u *Usecase) DeleteTrack(trackID, playlistID, userID uint32) error {
	if _, err := u.playlistRepo.GetByID(playlistID); err != nil {
		return fmt.Errorf("(usecase) can't find playlist in repository: %w", err)
	}

	if _, err := u.trackRepo.GetByID(trackID); err != nil {
		return fmt.Errorf("(usecase) can't find track in repository: %w", err)
	}

	userInAuthors, err := u.checkUserInAuthors(playlistID, userID)
	if err != nil {
		return err
	}
	if !userInAuthors {
		return fmt.Errorf("(usecase) playlist can't be updated by user: %w", &models.ForbiddenUserError{})
	}

	if err := u.playlistRepo.DeleteTrack(trackID, playlistID); err != nil {
		return fmt.Errorf("(usecase) can't delete track of playlist in repository: %w", err)
	}

	return nil
}

func (u *Usecase) GetFeed() ([]models.Playlist, error) {
	playlists, err := u.playlistRepo.GetFeed()
	if err != nil {
		return nil, fmt.Errorf("(usecase) can't get feed playlists from repository: %w", err)
	}

	return playlists, nil
}

func (u *Usecase) GetByUser(userID uint32) ([]models.Playlist, error) {
	_, err := u.playlistRepo.GetByID(userID)
	if err != nil {
		return nil, fmt.Errorf("(usecase) can't get authors from repository: %w", err)
	}

	playlists, err := u.playlistRepo.GetByUser(userID)
	if err != nil {
		return nil, fmt.Errorf("(usecase) can't get playlists from repository: %w", err)
	}

	return playlists, nil
}

func (u *Usecase) GetLikedByUser(userID uint32) ([]models.Playlist, error) {
	playlists, err := u.playlistRepo.GetLikedByUser(userID)
	if err != nil {
		return nil, fmt.Errorf("(usecase) can't get playlists from repository: %w", err)
	}

	return playlists, nil
}

func (u *Usecase) SetLike(playlistID, userID uint32) (bool, error) {
	if _, err := u.playlistRepo.GetByID(playlistID); err != nil {
		return false, fmt.Errorf("(usecase) can't get playlist: %w", err)
	}

	isInserted, err := u.playlistRepo.InsertLike(playlistID, userID)
	if err != nil {
		return false, fmt.Errorf("(usecase) failed to set like: %w", err)
	}

	return isInserted, nil
}

func (u *Usecase) UnLike(playlistID, userID uint32) (bool, error) {
	if _, err := u.playlistRepo.GetByID(playlistID); err != nil {
		return false, fmt.Errorf("(usecase) can't get playlist: %w", err)
	}

	isDeleted, err := u.playlistRepo.DeleteLike(playlistID, userID)
	if err != nil {
		return false, fmt.Errorf("(usecase) failed to unset like: %w", err)
	}

	return isDeleted, nil
}

func (u *Usecase) checkUserInAuthors(playlistID, userID uint32) (bool, error) {
	userInAuthors := false
	users, err := u.userRepo.GetByPlaylist(playlistID)
	if err != nil {
		return false, fmt.Errorf("(usecase) can't get authors of playlist: %w", err)
	}
	for _, user := range users {
		if user.ID == userID {
			userInAuthors = true
			break
		}
	}

	return userInAuthors, nil
}

func (u *Usecase) IsLiked(albumID, userID uint32) (bool, error) {
	isLiked, err := u.playlistRepo.IsLiked(albumID, userID)
	if err != nil {
		return false, fmt.Errorf("(usecase) can't check in repository if playlist is liked: %w", err)
	}

	return isLiked, nil
}
