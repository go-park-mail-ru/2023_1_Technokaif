package usecase

import (
	"github.com/pkg/errors"

	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/logger"
	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/models"
	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/repository"
)

type ArtistUsecase struct {
	repo   repository.Artist
	logger logger.Logger
}

func NewArtistUsecase(ra repository.Artist, l logger.Logger) *ArtistUsecase {
	return &ArtistUsecase{repo: ra, logger: l}
}

// GetFeed returns artists for main page
func (a *ArtistUsecase) GetFeed() ([]models.ArtistFeed, error) {
	artists, err := a.repo.GetFeed()

	return artists, errors.Wrap(err, "(Usecase) cannot get feed")
}
