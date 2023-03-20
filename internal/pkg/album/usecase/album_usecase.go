package album_usecase

import (
	"fmt"

	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/models"
	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/album"
	"github.com/go-park-mail-ru/2023_1_Technokaif/pkg/logger"
)

type albumUsecase struct {
	repo   album.AlbumRepository
	logger logger.Logger
}

func NewAlbumUsecase(ar album.AlbumRepository, l logger.Logger) album.AlbumUsecase {
	return &albumUsecase{repo: ar, logger: l}
}

func (au *albumUsecase) GetByID(albumID uint32) (models.Album, error) {
	album, err := au.repo.GetByID(albumID)
	if err != nil {
		return models.Album{}, fmt.Errorf("(usecase) can't get album from repository: %w", err)
	}

	return album, nil
}

func (au *albumUsecase) GetFeed() ([]models.Album, error) {
	albums, err := au.repo.GetFeed()
	if err != nil {
		return nil, fmt.Errorf("(usecase) can't get albums from repository: %w", err)
	}

	return albums, nil
}

func (au *albumUsecase) GetByArtist(artistID uint32) ([]models.Album, error) {
	albums, err := au.repo.GetByArtist(artistID)
	if err != nil {
		return nil, fmt.Errorf("(usecase) can't get albums from repository: %w", err)
	}

	return albums, nil
}

func (au *albumUsecase) GetByTrack(trackID uint32) (models.Album, error) {
	albums, err := au.repo.GetByTrack(trackID)
	if err != nil {
		return models.Album{}, fmt.Errorf("(usecase) can't get albums from repository: %w", err)
	}

	return albums, nil
}

func (au *albumUsecase) GetLikedByUser(userID uint32) ([]models.Album, error) {
	albums, err := au.repo.GetLikedByUser(userID)
	if err != nil {
		return nil, fmt.Errorf("(usecase) can't get albums from repository: %w", err)
	}

	return albums, nil
}
