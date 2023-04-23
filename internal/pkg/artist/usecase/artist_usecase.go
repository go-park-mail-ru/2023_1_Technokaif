package usecase

import (
	"fmt"

	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/models"
	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/artist"
	"github.com/go-park-mail-ru/2023_1_Technokaif/pkg/logger"
)

// Usecase implements artist.Usecase
type Usecase struct {
	repo   artist.Repository
	logger logger.Logger
}

func NewUsecase(ar artist.Repository, l logger.Logger) *Usecase {
	return &Usecase{repo: ar, logger: l}
}

func (u *Usecase) Create(artist models.Artist) (uint32, error) {
	artistID, err := u.repo.Insert(artist)
	if err != nil {
		return 0, fmt.Errorf("(usecase) can't insert artist into repository: %w", err)
	}

	return artistID, nil
}

func (u *Usecase) GetByID(artistID uint32) (*models.Artist, error) {
	artist, err := u.repo.GetByID(artistID)
	if err != nil {
		return &models.Artist{}, fmt.Errorf("(usecase) can't get artist from repository: %w", err)
	}

	return artist, nil
}

func (u *Usecase) Delete(artistID uint32, userID uint32) error {
	artist, err := u.repo.GetByID(artistID)
	if err != nil {
		return fmt.Errorf("(usecase) can't find artist in repository: %w", err)
	}

	if artist.UserID == nil {
		return fmt.Errorf("(usecase) artist can't be deleted by user: %w", &models.ForbiddenUserError{})
	}

	if *artist.UserID != userID {
		return fmt.Errorf("(usecase) artist can't be deleted by this user: %w", &models.ForbiddenUserError{})
	}

	if err := u.repo.DeleteByID(artistID); err != nil {
		return fmt.Errorf("(usecase) can't delete artist from repository: %w", err)
	}

	return nil
}

func (u *Usecase) GetFeed() ([]models.Artist, error) {
	artists, err := u.repo.GetFeed()
	if err != nil {
		return nil, fmt.Errorf("(usecase) can't get feed artists from repository: %w", err)
	}

	return artists, nil
}

func (u *Usecase) GetByAlbum(albumID uint32) ([]models.Artist, error) {
	artists, err := u.repo.GetByAlbum(albumID)
	if err != nil {
		return nil, fmt.Errorf("(usecase) can't get feed artists from repository: %w", err)
	}

	return artists, nil
}

func (u *Usecase) GetByTrack(trackID uint32) ([]models.Artist, error) {
	artists, err := u.repo.GetByTrack(trackID)
	if err != nil {
		return nil, fmt.Errorf("(usecase) can't get artists from repository: %w", err)
	}

	return artists, nil
}

func (u *Usecase) GetLikedByUser(userID uint32) ([]models.Artist, error) {
	artists, err := u.repo.GetLikedByUser(userID)
	if err != nil {
		return nil, fmt.Errorf("(usecase) can't get artists from repository: %w", err)
	}

	return artists, nil
}

func (u *Usecase) SetLike(artistID, userID uint32) (bool, error) {
	if _, err := u.repo.GetByID(artistID); err != nil {
		return false, fmt.Errorf("(usecase) can't get artist: %w", err)
	}

	iSinserted, err := u.repo.InsertLike(artistID, userID)
	if err != nil {
		return false, fmt.Errorf("(usecase) failed to set like: %w", err)
	}

	return iSinserted, nil
}

func (u *Usecase) UnLike(artistID, userID uint32) (bool, error) {
	if _, err := u.repo.GetByID(artistID); err != nil {
		return false, fmt.Errorf("(usecase) can't get artist: %w", err)
	}

	isDeleted, err := u.repo.DeleteLike(artistID, userID)
	if err != nil {
		return false, fmt.Errorf("(usecase) failed to unset like: %w", err)
	}

	return isDeleted, nil
}

func (u *Usecase) IsLiked(artistID, userID uint32) (bool, error) {
	isLiked, err := u.repo.IsLiked(artistID, userID)
	if err != nil {
		return false, fmt.Errorf("(usecase) can't check in repository if artist is liked: %w", err)
	}

	return isLiked, nil
}
