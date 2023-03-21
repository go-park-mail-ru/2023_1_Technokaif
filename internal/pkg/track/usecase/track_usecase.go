package usecase

import (
	"fmt"

	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/models"
	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/artist"
	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/track"
	"github.com/go-park-mail-ru/2023_1_Technokaif/pkg/logger"
)

// Usecase implements track.Usecase
type Usecase struct {
	trackRepo  track.Repository
	artistRepo artist.Repository
	logger     logger.Logger
}

func NewUsecase(tr track.Repository, ar artist.Repository, l logger.Logger) *Usecase {
	return &Usecase{
		trackRepo:  tr,
		artistRepo: ar,
		logger:     l,
	}
}

func (u *Usecase) Create(track models.Track, artistsID []uint32) (uint32, error) {
	for _, artistID := range artistsID {
		if _, err := u.artistRepo.GetByID(artistID); err != nil {
			return 0, fmt.Errorf("(usecase) no such artist with id #%d: %w", artistID, err)
		}
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

func (u *Usecase) Change(track models.Track) error {
	if err := u.trackRepo.Update(track); err != nil {
		return fmt.Errorf("(usecase) can't get update track in repository: %w", err)
	}

	return nil
}

func (u *Usecase) DeleteByID(trackID uint32) error {
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
	tracks, err := u.trackRepo.GetByAlbum(albumID)
	if err != nil {
		return nil, fmt.Errorf("(usecase) can't get tracks from repository: %w", err)
	}

	return tracks, nil
}

func (u *Usecase) GetByArtist(artistID uint32) ([]models.Track, error) {
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
