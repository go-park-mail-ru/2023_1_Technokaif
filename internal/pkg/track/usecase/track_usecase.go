package usecase

import (
	"fmt"

	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/models"
	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/track"
	"github.com/go-park-mail-ru/2023_1_Technokaif/pkg/logger"
)

// Usecase implements track.Usecase
type Usecase struct {
	repo   track.Repository
	logger logger.Logger
}

func NewUsecase(tr track.Repository, l logger.Logger) *Usecase {
	return &Usecase{repo: tr, logger: l}
}

func (u *Usecase) Create(track models.Track) error {
	if err := u.repo.Insert(track); err != nil {
		return fmt.Errorf("(usecase) can't insert track into repository: %w", err)
	}

	return nil
}

func (u *Usecase) GetByID(trackID uint32) (models.Track, error) {
	track, err := u.repo.GetByID(trackID)
	if err != nil {
		return models.Track{}, fmt.Errorf("(usecase) can't get track from repository: %w", err)
	}

	return track, nil
}

func (u *Usecase) Change(track models.Track) error {
	if err := u.repo.Update(track); err != nil {
		return fmt.Errorf("(usecase) can't get update track in repository: %w", err)
	}

	return nil
}

func (u *Usecase) DeleteByID(trackID uint32) error {
	if err := u.repo.DeleteByID(trackID); err != nil {
		return fmt.Errorf("(usecase) can't delete track from repository: %w", err)
	}

	return nil
}

func (u *Usecase) GetFeed() ([]models.Track, error) {
	tracks, err := u.repo.GetFeed()
	if err != nil {
		return nil, fmt.Errorf("(usecase) can't get feed tracks from repository: %w", err)
	}

	return tracks, nil
}

func (u *Usecase) GetByAlbum(albumID uint32) ([]models.Track, error) {
	tracks, err := u.repo.GetByAlbum(albumID)
	if err != nil {
		return nil, fmt.Errorf("(usecase) can't get tracks from repository: %w", err)
	}

	return tracks, nil
}

func (u *Usecase) GetByArtist(artistID uint32) ([]models.Track, error) {
	tracks, err := u.repo.GetByArtist(artistID)
	if err != nil {
		return nil, fmt.Errorf("(usecase) can't get tracks from repository: %w", err)
	}

	return tracks, nil
}

func (u *Usecase) GetLikedByUser(userID uint32) ([]models.Track, error) {
	tracks, err := u.repo.GetLikedByUser(userID)
	if err != nil {
		return nil, fmt.Errorf("(usecase) can't get tracks from repository: %w", err)
	}

	return tracks, nil
}
