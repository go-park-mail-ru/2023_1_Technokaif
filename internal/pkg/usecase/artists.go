package usecase

import (
	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/models"
	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/repository"
)

type ArtistUsecase struct {
	repo repository.Artist
}

func NewArtistUsecase(ra repository.Artist) *ArtistUsecase {
	return &ArtistUsecase{repo: ra}
}

func (a *ArtistUsecase) GetFeed() ([]models.ArtistFeed, error) {
	artists, err := a.repo.GetFeed()
	if err != nil {
		return nil, err
	}

	return artists, nil
}
