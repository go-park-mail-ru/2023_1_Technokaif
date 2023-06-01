package usecase

import (
	"context"
	"fmt"
	"io"
	"path/filepath"

	commonFile "github.com/go-park-mail-ru/2023_1_Technokaif/internal/common/file"
	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/models"
	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/playlist"
	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/track"
	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/user"
)

const feedPlaylistsAmountLimit uint32 = 100

// Usecase implements album.Usecase
type Usecase struct {
	playlistRepo playlist.Repository
	trackRepo    track.Repository
	userRepo     user.Repository
	coverSaver   CoverSaver
}

//go:generate mockgen -source=playlist_usecase.go -destination=../mocks/saver.go -package mock_playlist
type CoverSaver interface {
	Save(ctx context.Context, cover io.Reader, objectName string, size int64) error
}

func NewUsecase(pr playlist.Repository, tr track.Repository, ur user.Repository, saver CoverSaver) *Usecase {
	return &Usecase{
		playlistRepo: pr,
		trackRepo:    tr,
		userRepo:     ur,
		coverSaver:   saver,
	}
}

func (u *Usecase) Create(ctx context.Context,
	playlist models.Playlist, usersID []uint32, userID uint32) (uint32, error) {

	userInAuthors := false
	for _, uid := range usersID {
		user, err := u.userRepo.GetByID(ctx, uid)
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

	playlistID, err := u.playlistRepo.Insert(ctx, playlist, usersID)
	if err != nil {
		return 0, fmt.Errorf("(usecase) can't insert playlist into repository: %w", err)
	}

	return playlistID, nil
}

func (u *Usecase) GetByID(ctx context.Context, playlistID uint32) (*models.Playlist, error) {
	playlist, err := u.playlistRepo.GetByID(ctx, playlistID)
	if err != nil {
		return nil, fmt.Errorf("(usecase) can't get playlist from repository: %w", err)
	}

	return playlist, nil
}

func (u *Usecase) UpdateInfoAndMembers(ctx context.Context,
	playlist models.Playlist, usersID []uint32, userID uint32) error {

	pl, err := u.playlistRepo.GetByID(ctx, playlist.ID)
	if err != nil {
		return fmt.Errorf("(usecase) can't find playlist in repository: %w", err)
	}
	playlist.CoverSrc = pl.CoverSrc

	userInAuthors, err := u.checkUserInAuthors(ctx, playlist.ID, userID)
	if err != nil {
		return err
	}
	if !userInAuthors {
		return fmt.Errorf("(usecase) playlist can't be deleted by user: %w", &models.ForbiddenUserError{})
	}

	authors, err := u.userRepo.GetByPlaylist(ctx, playlist.ID)
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

	if err := u.playlistRepo.UpdateWithMembers(ctx, playlist, newAuthorsID); err != nil {
		return fmt.Errorf("(usecase) can't update playlist in repository: %w", err)
	}

	return nil
}

func (u *Usecase) UploadCover(ctx context.Context,
	playlistID uint32, userID uint32, file io.ReadSeeker, fileSize int64, fileExtension string) error {

	playlist, err := u.playlistRepo.GetByID(ctx, playlistID)
	if err != nil {
		return fmt.Errorf("(usecase) can't find playlist: %w", err)
	}

	userInAuthors, err := u.checkUserInAuthors(ctx, playlistID, userID)
	if err != nil {
		return err
	}
	if !userInAuthors {
		return fmt.Errorf("(usecase) playlist can't be deleted by user: %w", &models.ForbiddenUserError{})
	}

	// Check format
	if fileType, err := commonFile.CheckMimeType(file, "image/png", "image/jpeg"); err != nil {
		return fmt.Errorf("(usecase) file format %s: %w", fileType, &models.CoverWrongFormatError{FileType: fileType})
	}

	filenameWithExtension, err := commonFile.FileHash(file, fileExtension)
	if err != nil {
		return fmt.Errorf("(usecase) can't get file hash: %w", err)
	}

	if err := u.coverSaver.Save(ctx, file, filenameWithExtension, fileSize); err != nil {
		return fmt.Errorf("(usecase) can't save cover: %w", err)
	}

	playlist.CoverSrc = filepath.Join(commonFile.PlaylistCoverFolder(), filenameWithExtension)
	if err := u.playlistRepo.Update(ctx, *playlist); err != nil {
		return fmt.Errorf("(usecase) can't update playlist: %w", err)
	}
	return nil
}

func (u *Usecase) Delete(ctx context.Context, playlistID uint32, userID uint32) error {
	if err := u.playlistRepo.Check(ctx, playlistID); err != nil {
		return fmt.Errorf("(usecase) can't find playlist with id #%d: %w", playlistID, err)
	}

	userInAuthors, err := u.checkUserInAuthors(ctx, playlistID, userID)
	if err != nil {
		return err
	}
	if !userInAuthors {
		return fmt.Errorf("(usecase) playlist can't be deleted by user: %w", &models.ForbiddenUserError{})
	}

	if err := u.playlistRepo.DeleteByID(ctx, playlistID); err != nil {
		return fmt.Errorf("(usecase) can't delete playlist from repository: %w", err)
	}

	return nil
}

func (u *Usecase) AddTrack(ctx context.Context, trackID, playlistID, userID uint32) error {
	if err := u.playlistRepo.Check(ctx, playlistID); err != nil {
		return fmt.Errorf("(usecase) can't find playlist with id #%d: %w", playlistID, err)
	}

	if err := u.trackRepo.Check(ctx, trackID); err != nil {
		return fmt.Errorf("(usecase) can't find track in repository: %w", err)
	}

	userInAuthors, err := u.checkUserInAuthors(ctx, playlistID, userID)
	if err != nil {
		return err
	}
	if !userInAuthors {
		return fmt.Errorf("(usecase) playlist can't be updated by user: %w", &models.ForbiddenUserError{})
	}

	if err := u.playlistRepo.AddTrack(ctx, trackID, playlistID); err != nil {
		return fmt.Errorf("(usecase) can't add track into playlist in repository: %w", err)
	}

	return nil
}

func (u *Usecase) DeleteTrack(ctx context.Context, trackID, playlistID, userID uint32) error {
	if err := u.playlistRepo.Check(ctx, playlistID); err != nil {
		return fmt.Errorf("(usecase) can't find playlist with id #%d: %w", playlistID, err)
	}

	if err := u.trackRepo.Check(ctx, trackID); err != nil {
		return fmt.Errorf("(usecase) can't find track with id #%d: %w", trackID, err)
	}

	userInAuthors, err := u.checkUserInAuthors(ctx, playlistID, userID)
	if err != nil {
		return err
	}
	if !userInAuthors {
		return fmt.Errorf("(usecase) playlist can't be updated by user: %w", &models.ForbiddenUserError{})
	}

	if err := u.playlistRepo.DeleteTrack(ctx, trackID, playlistID); err != nil {
		return fmt.Errorf("(usecase) can't delete track of playlist in repository: %w", err)
	}

	return nil
}

func (u *Usecase) GetFeed(ctx context.Context) ([]models.Playlist, error) {
	playlists, err := u.playlistRepo.GetFeed(ctx, feedPlaylistsAmountLimit)
	if err != nil {
		return nil, fmt.Errorf("(usecase) can't get feed playlists from repository: %w", err)
	}

	return playlists, nil
}

func (u *Usecase) GetByUser(ctx context.Context, userID uint32) ([]models.Playlist, error) {
	if err := u.userRepo.Check(ctx, userID); err != nil {
		return nil, fmt.Errorf("(usecase) can't find user with id #%d: %w", userID, err)
	}

	playlists, err := u.playlistRepo.GetByUser(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("(usecase) can't get playlists from repository: %w", err)
	}

	return playlists, nil
}

func (u *Usecase) GetLikedByUser(ctx context.Context, userID uint32) ([]models.Playlist, error) {
	playlists, err := u.playlistRepo.GetLikedByUser(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("(usecase) can't get playlists from repository: %w", err)
	}

	return playlists, nil
}

func (u *Usecase) SetLike(ctx context.Context, playlistID, userID uint32) (bool, error) {
	if err := u.playlistRepo.Check(ctx, playlistID); err != nil {
		return false, fmt.Errorf("(usecase) can't find playlist with id #%d: %w", playlistID, err)
	}

	isInserted, err := u.playlistRepo.InsertLike(ctx, playlistID, userID)
	if err != nil {
		return false, fmt.Errorf("(usecase) failed to set like: %w", err)
	}

	return isInserted, nil
}

func (u *Usecase) UnLike(ctx context.Context, playlistID, userID uint32) (bool, error) {
	if err := u.playlistRepo.Check(ctx, playlistID); err != nil {
		return false, fmt.Errorf("(usecase) can't find playlist with id #%d: %w", playlistID, err)
	}

	isDeleted, err := u.playlistRepo.DeleteLike(ctx, playlistID, userID)
	if err != nil {
		return false, fmt.Errorf("(usecase) failed to unset like: %w", err)
	}

	return isDeleted, nil
}

func (u *Usecase) checkUserInAuthors(ctx context.Context, playlistID, userID uint32) (bool, error) {
	userInAuthors := false
	users, err := u.userRepo.GetByPlaylist(ctx, playlistID)
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

func (u *Usecase) IsLiked(ctx context.Context, albumID, userID uint32) (bool, error) {
	isLiked, err := u.playlistRepo.IsLiked(ctx, albumID, userID)
	if err != nil {
		return false, fmt.Errorf("(usecase) can't check in repository if playlist is liked: %w", err)
	}

	return isLiked, nil
}
