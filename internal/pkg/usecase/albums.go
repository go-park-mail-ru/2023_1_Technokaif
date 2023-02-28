package usecase

import (
	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/models"
	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/repository"
)

type AlbumUsecase struct {
	repo repository.Album
}

func NewAlbumUsecase(ra repository.Album) *AlbumUsecase {
	return &AlbumUsecase{repo: ra}
}

// GetFeed returns albums for main page
func (a *AlbumUsecase) GetFeed() ([]models.AlbumFeed, error) {
	albums, err := a.repo.GetFeed()
	if err != nil {
		return nil, err
	}

	return albums, nil
}
