package usecase

import (
	"fmt"

	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/models"
	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/album"
	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/artist"
	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/playlist"
	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/track"
	"github.com/go-park-mail-ru/2023_1_Technokaif/pkg/logger"
)

const feedTracksAmountLimit = 100

// Usecase implements track.Usecase
type Usecase struct {
	trackRepo    track.Repository
	artistRepo   artist.Repository
	albumRepo    album.Repository
	playlistRepo playlist.Repository

	logger logger.Logger
}

func NewUsecase(tr track.Repository, arr artist.Repository,
	alr album.Repository, pr playlist.Repository, l logger.Logger) *Usecase {

	return &Usecase{
		trackRepo:    tr,
		artistRepo:   arr,
		albumRepo:    alr,
		playlistRepo: pr,

		logger: l,
	}
}

func (u *Usecase) Create(track models.Track, artistsID []uint32, userID uint32) (uint32, error) {
	userInArtists := false
	for _, artistID := range artistsID {
		a, err := u.artistRepo.GetByID(artistID)
		if err != nil {
			return 0, fmt.Errorf("(usecase) can't get artist with id #%d: %w", artistID, err)
		}
		if a.UserID != nil && *a.UserID == userID {
			userInArtists = true
			break
		}
	}
	if !userInArtists {
		return 0, fmt.Errorf("(usecase) track can't be created by user: %w", &models.ForbiddenUserError{})
	}

	trackID, err := u.trackRepo.Insert(track, artistsID)
	if err != nil {
		return 0, fmt.Errorf("(usecase) can't insert track into repository: %w", err)
	}

	return trackID, nil
}

func (u *Usecase) GetByID(trackID uint32) (*models.Track, error) {
	track, err := u.trackRepo.GetByID(trackID)
	if err != nil {
		return &models.Track{}, fmt.Errorf("(usecase) can't get track with id #%d: %w", trackID, err)
	}

	return track, nil
}

func (u *Usecase) Delete(trackID uint32, userID uint32) error {
	if err := u.trackRepo.Check(trackID); err != nil {
		return fmt.Errorf("(usecase) can't find track with id #%d: %w", trackID, err)
	}

	userInArtists := false
	artists, err := u.artistRepo.GetByAlbum(trackID)
	if err != nil {
		return fmt.Errorf("(usecase) can't get artists of track: %w", err)
	}
	for _, artist := range artists {
		if artist.UserID != nil && *artist.UserID == userID {
			userInArtists = true
		}
	}
	if !userInArtists {
		return fmt.Errorf("(usecase) track can't be deleted by user: %w", &models.ForbiddenUserError{})
	}

	if err := u.trackRepo.DeleteByID(trackID); err != nil {
		return fmt.Errorf("(usecase) can't delete track from repository: %w", err)
	}

	return nil
}

func (u *Usecase) GetFeed() ([]models.Track, error) {
	tracks, err := u.trackRepo.GetFeed(feedTracksAmountLimit)
	if err != nil {
		return nil, fmt.Errorf("(usecase) can't get feed tracks from repository: %w", err)
	}

	return tracks, nil
}

func (u *Usecase) GetByAlbum(albumID uint32) ([]models.Track, error) {
	if err := u.albumRepo.Check(albumID); err != nil {
		return nil, fmt.Errorf("(usecase) can't find album with id #%d: %w", albumID, err)
	}

	tracks, err := u.trackRepo.GetByAlbum(albumID)
	if err != nil {
		return nil, fmt.Errorf("(usecase) can't get tracks from repository: %w", err)
	}

	return tracks, nil
}

func (u *Usecase) GetByPlaylist(playlistID uint32) ([]models.Track, error) {
	if err := u.playlistRepo.Check(playlistID); err != nil {
		return nil, fmt.Errorf("(usecase) can't find playlist with id #%d: %w", playlistID, err)
	}

	tracks, err := u.trackRepo.GetByPlaylist(playlistID)
	if err != nil {
		return nil, fmt.Errorf("(usecase) can't get tracks from repository: %w", err)
	}

	return tracks, nil
}

func (u *Usecase) GetByArtist(artistID uint32) ([]models.Track, error) {
	if err := u.artistRepo.Check(artistID); err != nil {
		return nil, fmt.Errorf("(usecase) can't find artist with id #%d: %w", artistID, err)
	}

	tracks, err := u.trackRepo.GetByArtist(artistID)
	if err != nil {
		return nil, fmt.Errorf("(usecase) can't get tracks from repository: %w", err)
	}

	return tracks, nil
}

func (u *Usecase) GetLikedByUser(userID uint32) ([]models.Track, error) {
	tracks, err := u.trackRepo.GetLikedByUser(userID)
	if err != nil {
		return nil, fmt.Errorf("(usecase) can't get tracks from repository: %w", err)
	}

	return tracks, nil
}

func (u *Usecase) SetLike(trackID, userID uint32) (bool, error) {
	if err := u.trackRepo.Check(trackID); err != nil {
		return false, fmt.Errorf("(usecase) can't find track with id #%d: %w", trackID, err)
	}

	isInserted, err := u.trackRepo.InsertLike(trackID, userID)
	if err != nil {
		return false, fmt.Errorf("(usecase) failed to set like: %w", err)
	}

	return isInserted, nil
}

func (u *Usecase) UnLike(trackID, userID uint32) (bool, error) {
	if err := u.trackRepo.Check(trackID); err != nil {
		return false, fmt.Errorf("(usecase) can't find track with id #%d: %w", trackID, err)
	}

	isDeleted, err := u.trackRepo.DeleteLike(trackID, userID)
	if err != nil {
		return false, fmt.Errorf("(usecase) failed to unset like: %w", err)
	}

	return isDeleted, nil
}

func (u *Usecase) IsLiked(trackID, userID uint32) (bool, error) {
	isLiked, err := u.trackRepo.IsLiked(trackID, userID)
	if err != nil {
		return false, fmt.Errorf("(usecase) can't check in repository if track is liked: %w", err)
	}

	return isLiked, nil
}
