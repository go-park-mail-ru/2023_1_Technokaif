package usecase

import (
	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/repository"
)

type AlbumUsecase struct {
	repo repository.Album
}

func (a *AlbumUsecase) GetAlbums() []Album {
	return nil
}

func NewAlbumUsecase(ra repository.Album) *AlbumUsecase {
	return &AlbumUsecase{repo: ra}
}
