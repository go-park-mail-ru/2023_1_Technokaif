package usecase

import (
	"fmt"

	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/models"
	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/album"
	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/artist"
	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/track"
	"github.com/go-park-mail-ru/2023_1_Technokaif/pkg/logger"
)

// Usecase implements track.Usecase
type Usecase struct {
	trackRepo  track.Repository
	artistRepo artist.Repository
	albumRepo  album.Repository

	logger logger.Logger
}

func NewUsecase(tr track.Repository, arr artist.Repository, alr album.Repository, l logger.Logger) *Usecase {
	return &Usecase{
		trackRepo:  tr,
		artistRepo: arr,
		albumRepo:  alr,
		logger:     l,
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
		return &models.Track{}, fmt.Errorf("(usecase) can't get track from repository: %w", err)
	}

	return track, nil
}

func (u *Usecase) Delete(trackID uint32, userID uint32) error {
	if _, err := u.trackRepo.GetByID(trackID); err != nil {
		return fmt.Errorf("(usecase) can't find track in repository: %w", err)
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
	tracks, err := u.trackRepo.GetFeed()
	if err != nil {
		return nil, fmt.Errorf("(usecase) can't get feed tracks from repository: %w", err)
	}

	return tracks, nil
}

func (u *Usecase) GetByAlbum(albumID uint32) ([]models.Track, error) {
	_, err := u.albumRepo.GetByID(albumID)
	if err != nil {
		return nil, fmt.Errorf("(usecase) can't get album with id #%d: %w", albumID, err)
	}

	tracks, err := u.trackRepo.GetByAlbum(albumID)
	if err != nil {
		return nil, fmt.Errorf("(usecase) can't get tracks from repository: %w", err)
	}

	return tracks, nil
}

func (u *Usecase) GetByArtist(artistID uint32) ([]models.Track, error) {
	_, err := u.artistRepo.GetByID(artistID)
	if err != nil {
		return nil, fmt.Errorf("(usecase) can't get artist with id #%d: %w", artistID, err)
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
	if _, err := u.trackRepo.GetByID(trackID); err != nil {
		return false, fmt.Errorf("(usecase) can't get track: %w", err)
	}

	iSinserted, err := u.trackRepo.InsertLike(trackID, userID)
	if err != nil {
		return false, fmt.Errorf("(usecase) failed to set like: %w", err)
	}

	return iSinserted, nil
}

func (u *Usecase) UnLike(trackID, userID uint32) (bool, error) {
	if _, err := u.trackRepo.GetByID(trackID); err != nil {
		return false, fmt.Errorf("(usecase) can't get track: %w", err)
	}

	iSdeleted, err := u.trackRepo.DeleteLike(trackID, userID)
	if err != nil {
		return false, fmt.Errorf("(usecase) failed to unset like: %w", err)
	}

	return iSdeleted, nil
}

func (u *Usecase) IsLiked(trackID, userID uint32) (bool, error) {
	isLiked, err := u.trackRepo.IsLiked(trackID, userID)
	if err != nil {
		return false, fmt.Errorf("(usecase) can't check in repository if track is liked: %w", err)
	}

	return isLiked, nil
}
