package usecase

import (
	"fmt"

	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/models"
	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/album"
	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/artist"
	"github.com/go-park-mail-ru/2023_1_Technokaif/pkg/logger"
)

// Usecase implements album.Usecase
type Usecase struct {
	albumRepo   album.Repository
	artistRepo	artist.Repository
	logger logger.Logger
}

func NewUsecase(alr album.Repository, arr artist.Repository, l logger.Logger) *Usecase {
	return &Usecase{albumRepo: alr, artistRepo: arr, logger: l}
}

func (u *Usecase) Create(album models.Album) error {
	if err := u.albumRepo.Insert(album); err != nil {
		return fmt.Errorf("(usecase) can't insert album into repository: %w", err)
	}

	return nil
}

func (u *Usecase) GetByID(albumID uint32) (*models.Album, error) {
	album, err := u.albumRepo.GetByID(albumID)
	if err != nil {
		return nil, fmt.Errorf("(usecase) can't get album from repository: %w", err)
	}

	return album, nil
}

func (u *Usecase) Change(albumID models.Album) error {
	if err := u.albumRepo.Update(albumID); err != nil {
		return fmt.Errorf("(usecase) can't update album in repository: %w", err)
	}

	return nil
}

func (u *Usecase) DeleteByID(albumID uint32) error {
	if err := u.albumRepo.DeleteByID(albumID); err != nil {
		return fmt.Errorf("(usecase) can't delete album from repository: %w", err)
	}

	return nil
}

func (u *Usecase) GetFeed() ([]models.Album, error) {
	albums, err := u.albumRepo.GetFeed()
	if err != nil {
		return nil, fmt.Errorf("(usecase) can't get feed albums from repository: %w", err)
	}

	return albums, nil
}

func (u *Usecase) GetByArtist(artistID uint32) ([]models.Album, error) {
	_, err := u.artistRepo.GetByID(artistID)
	if err != nil {
		return nil, fmt.Errorf("(usecase) can't get artist from repository: %w", err)
	}

	albums, err := u.albumRepo.GetByArtist(artistID)
	if err != nil {
		return nil, fmt.Errorf("(usecase) can't get albums from repository: %w", err)
	}

	return albums, nil
}

func (u *Usecase) GetByTrack(trackID uint32) (*models.Album, error) {
	album, err := u.albumRepo.GetByTrack(trackID)
	if err != nil {
		return nil, fmt.Errorf("(usecase) can't get albums from repository: %w", err)
	}

	return album, nil
}

func (u *Usecase) GetLikedByUser(userID uint32) ([]models.Album, error) {
	albums, err := u.albumRepo.GetLikedByUser(userID)
	if err != nil {
		return nil, fmt.Errorf("(usecase) can't get albums from repository: %w", err)
	}

	return albums, nil
}
