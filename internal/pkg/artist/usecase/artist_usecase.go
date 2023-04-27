package usecase

import (
	"context"
	"fmt"

	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/models"
	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/artist"
	"github.com/go-park-mail-ru/2023_1_Technokaif/pkg/logger"
)

const feedArtistsAmountLimit uint32 = 100

// Usecase implements artist.Usecase
type Usecase struct {
	repo   artist.Repository
	logger logger.Logger
}

func NewUsecase(ar artist.Repository, l logger.Logger) *Usecase {
	return &Usecase{repo: ar, logger: l}
}

func (u *Usecase) Create(ctx context.Context, artist models.Artist) (uint32, error) {
	artistID, err := u.repo.Insert(ctx, artist)
	if err != nil {
		return 0, fmt.Errorf("(usecase) can't insert artist into repository: %w", err)
	}

	return artistID, nil
}

func (u *Usecase) GetByID(ctx context.Context, artistID uint32) (*models.Artist, error) {
	artist, err := u.repo.GetByID(ctx, artistID)
	if err != nil {
		return &models.Artist{}, fmt.Errorf("(usecase) can't get artist from repository: %w", err)
	}

	return artist, nil
}

func (u *Usecase) Delete(ctx context.Context, artistID uint32, userID uint32) error {
	artist, err := u.repo.GetByID(ctx, artistID)
	if err != nil {
		return fmt.Errorf("(usecase) can't find artist in repository: %w", err)
	}

	if artist.UserID == nil {
		return fmt.Errorf("(usecase) artist can't be deleted by user: %w", &models.ForbiddenUserError{})
	}

	if *artist.UserID != userID {
		return fmt.Errorf("(usecase) artist can't be deleted by this user: %w", &models.ForbiddenUserError{})
	}

	if err := u.repo.DeleteByID(ctx, artistID); err != nil {
		return fmt.Errorf("(usecase) can't delete artist from repository: %w", err)
	}

	return nil
}

func (u *Usecase) GetFeed(ctx context.Context) ([]models.Artist, error) {
	artists, err := u.repo.GetFeed(ctx, feedArtistsAmountLimit)
	if err != nil {
		return nil, fmt.Errorf("(usecase) can't get feed artists from repository: %w", err)
	}

	return artists, nil
}

func (u *Usecase) GetByAlbum(ctx context.Context, albumID uint32) ([]models.Artist, error) {
	artists, err := u.repo.GetByAlbum(ctx, albumID)
	if err != nil {
		return nil, fmt.Errorf("(usecase) can't get feed artists from repository: %w", err)
	}

	return artists, nil
}

func (u *Usecase) GetByTrack(ctx context.Context, trackID uint32) ([]models.Artist, error) {
	artists, err := u.repo.GetByTrack(ctx, trackID)
	if err != nil {
		return nil, fmt.Errorf("(usecase) can't get artists from repository: %w", err)
	}

	return artists, nil
}

func (u *Usecase) GetLikedByUser(ctx context.Context, userID uint32) ([]models.Artist, error) {
	artists, err := u.repo.GetLikedByUser(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("(usecase) can't get artists from repository: %w", err)
	}

	return artists, nil
}

func (u *Usecase) SetLike(ctx context.Context, artistID, userID uint32) (bool, error) {
	if err := u.repo.Check(ctx, artistID); err != nil {
		return false, fmt.Errorf("(usecase) can't find artist with id #%d: %w", artistID, err)
	}

	isInserted, err := u.repo.InsertLike(ctx, artistID, userID)
	if err != nil {
		return false, fmt.Errorf("(usecase) failed to set like: %w", err)
	}

	return isInserted, nil
}

func (u *Usecase) UnLike(ctx context.Context, artistID, userID uint32) (bool, error) {
	if err := u.repo.Check(ctx, artistID); err != nil {
		return false, fmt.Errorf("(usecase) can't find artist with id #%d: %w", artistID, err)
	}

	isDeleted, err := u.repo.DeleteLike(ctx, artistID, userID)
	if err != nil {
		return false, fmt.Errorf("(usecase) failed to unset like: %w", err)
	}

	return isDeleted, nil
}

func (u *Usecase) IsLiked(ctx context.Context, artistID, userID uint32) (bool, error) {
	isLiked, err := u.repo.IsLiked(ctx, artistID, userID)
	if err != nil {
		return false, fmt.Errorf("(usecase) can't check in repository if artist is liked: %w", err)
	}

	return isLiked, nil
}
