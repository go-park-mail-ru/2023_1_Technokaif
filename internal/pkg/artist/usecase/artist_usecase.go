package usecase

import (
	"fmt"

	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/models"
	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/artist"
	"github.com/go-park-mail-ru/2023_1_Technokaif/pkg/logger"
)

// Usecase implements artist.Usecase
type Usecase struct {
	repo   artist.Repository
	logger logger.Logger
}

func NewUsecase(ar artist.Repository, l logger.Logger) *Usecase {
	return &Usecase{repo: ar, logger: l}
}

func (u *Usecase) Create(artist models.Artist) error {
	if err := u.repo.Insert(artist); err != nil {
		return fmt.Errorf("(usecase) can't insert artist into repository: %w", err)
	}

	return nil
}

func (u *Usecase) GetByID(artistID uint32) (*models.Artist, error) {
	artist, err := u.repo.GetByID(artistID)
	if err != nil {
		return &models.Artist{}, fmt.Errorf("(usecase) can't get artist from repository: %w", err)
	}

	return artist, nil
}

func (u *Usecase) Change(artist models.Artist) error {
	if err := u.repo.Update(artist); err != nil {
		return fmt.Errorf("(usecase) can't update artist in repository: %w", err)
	}

	return nil
}

func (u *Usecase) DeleteByID(artistID uint32) error {
	if err := u.repo.DeleteByID(artistID); err != nil {
		return fmt.Errorf("(usecase) can't delete artist from repository: %w", err)
	}

	return nil
}

func (u *Usecase) GetFeed() ([]models.Artist, error) {
	artists, err := u.repo.GetFeed()
	if err != nil {
		return nil, fmt.Errorf("(usecase) can't get feed artists from repository: %w", err)
	}

	return artists, nil
}

func (u *Usecase) GetByAlbum(albumID uint32) ([]models.Artist, error) {
	artists, err := u.repo.GetByAlbum(albumID)
	if err != nil {
		return nil, fmt.Errorf("(usecase) can't get feed artists from repository: %w", err)
	}

	return artists, nil
}

func (u *Usecase) GetByTrack(trackID uint32) ([]models.Artist, error) {
	artists, err := u.repo.GetByTrack(trackID)
	if err != nil {
		return nil, fmt.Errorf("(usecase) can't get artists from repository: %w", err)
	}

	return artists, nil
}

func (u *Usecase) GetLikedByUser(userID uint32) ([]models.Artist, error) {
	artists, err := u.repo.GetLikedByUser(userID)
	if err != nil {
		return nil, fmt.Errorf("(usecase) can't get artists from repository: %w", err)
	}

	return artists, nil
}
