package usecase

import (
	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/repository"
)

type ArtistUsecase struct {
	repo repository.Artist
}

func (a *ArtistUsecase) GetArtists() []Artist {
	return nil
}

func NewArtistUsecase(ra repository.Artist) *ArtistUsecase {
	return &ArtistUsecase{repo: ra}
}
