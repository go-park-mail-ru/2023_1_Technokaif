package usecase

import (
	"fmt"

	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/models"
	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/album"
	"github.com/go-park-mail-ru/2023_1_Technokaif/pkg/logger"
)

// Usecase implements album.Usecase
type Usecase struct {
	repo   album.Repository
	logger logger.Logger
}

func NewUsecase(ar album.Repository, l logger.Logger) *Usecase {
	return &Usecase{repo: ar, logger: l}
}

func (u *Usecase) GetByID(albumID uint32) (models.Album, error) {
	album, err := u.repo.GetByID(albumID)
	if err != nil {
		return models.Album{}, fmt.Errorf("(usecase) can't get album from repository: %w", err)
	}

	return album, nil
}

func (u *Usecase) GetFeed() ([]models.Album, error) {
	albums, err := u.repo.GetFeed()
	if err != nil {
		return nil, fmt.Errorf("(usecase) can't get albums from repository: %w", err)
	}

	return albums, nil
}

func (u *Usecase) GetByArtist(artistID uint32) ([]models.Album, error) {
	albums, err := u.repo.GetByArtist(artistID)
	if err != nil {
		return nil, fmt.Errorf("(usecase) can't get albums from repository: %w", err)
	}

	return albums, nil
}

func (u *Usecase) GetByTrack(trackID uint32) (models.Album, error) {
	albums, err := u.repo.GetByTrack(trackID)
	if err != nil {
		return models.Album{}, fmt.Errorf("(usecase) can't get albums from repository: %w", err)
	}

	return albums, nil
}

func (u *Usecase) GetLikedByUser(userID uint32) ([]models.Album, error) {
	albums, err := u.repo.GetLikedByUser(userID)
	if err != nil {
		return nil, fmt.Errorf("(usecase) can't get albums from repository: %w", err)
	}

	return albums, nil
}
