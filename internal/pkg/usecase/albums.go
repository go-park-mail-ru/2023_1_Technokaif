package usecase

import (
	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/logger"
	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/models"
	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/repository"
)

type AlbumUsecase struct {
	repo   repository.Album
	logger logger.Logger
}

func NewAlbumUsecase(ra repository.Album, l logger.Logger) *AlbumUsecase {
	return &AlbumUsecase{repo: ra, logger: l}
}

// GetFeed returns albums for main page
func (a *AlbumUsecase) GetFeed() ([]models.AlbumFeed, error) {
	albums, err := a.repo.GetFeed()
	if err != nil {
		return nil, err
	}

	return albums, nil
}
