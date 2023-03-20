package artist_usecase

import (
	"fmt"

	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/models"
	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/artist"
	"github.com/go-park-mail-ru/2023_1_Technokaif/pkg/logger"
)

type artistUsecase struct {
	repo   artist.ArtistRepository
	logger logger.Logger
}

func NewArtistUsecase(ar artist.ArtistRepository, l logger.Logger) artist.ArtistUsecase {
	return &artistUsecase{repo: ar, logger: l}
}

func (au *artistUsecase) Create(artist models.Artist) error {
	if err := au.repo.Insert(artist); err != nil {
		return fmt.Errorf("(usecase) can't create new entry of artist: %w", err)
	}

	return nil
}

func (au *artistUsecase) GetByID(artistID uint32) (models.Artist, error) {
	artist, err := au.repo.GetByID(artistID)
	if err != nil {
		return models.Artist{}, fmt.Errorf("(usecase) can't get artist from repository: %w", err)
	}

	return artist, nil
}

func (au *artistUsecase) Change(artist models.Artist) error {
	if err := au.repo.Update(artist); err != nil {
		return fmt.Errorf("(usecase) can't change entry of artist: %w", err)
	}

	return nil
}

func (au *artistUsecase) DeleteByID(artistID uint32) error {
	if err := au.repo.DeleteByID(artistID); err != nil {
		return fmt.Errorf("(usecase) can't delete entry of artist (id #%d): %w", artistID, err)
	}

	return nil
}

func (au *artistUsecase) GetFeed() ([]models.Artist, error) {
	artists, err := au.repo.GetFeed()
	if err != nil {
		return nil, fmt.Errorf("(usecase) can't get artists from repository: %w", err)
	}

	return artists, nil
}

func (au *artistUsecase) GetByAlbum(albumID uint32) ([]models.Artist, error) {
	artists, err := au.repo.GetByAlbum(albumID)
	if err != nil {
		return nil, fmt.Errorf("(usecase) can't get artists from repository: %w", err)
	}

	return artists, nil
}

func (au *artistUsecase) GetByTrack(trackID uint32) ([]models.Artist, error) {
	artists, err := au.repo.GetByTrack(trackID)
	if err != nil {
		return nil, fmt.Errorf("(usecase) can't get artists from repository: %w", err)
	}

	return artists, nil
}

func (au *artistUsecase) GetLikedByUser(userID uint32) ([]models.Artist, error) {
	artists, err := au.repo.GetLikedByUser(userID)
	if err != nil {
		return nil, fmt.Errorf("(usecase) can't get artists from repository: %w", err)
	}

	return artists, nil
}
