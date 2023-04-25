package usecase

import (
	"fmt"

	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/models"
	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/album"
	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/artist"
	"github.com/go-park-mail-ru/2023_1_Technokaif/pkg/logger"
)

const feedAlbumsAmountLimit = 100

// Usecase implements album.Usecase
type Usecase struct {
	albumRepo  album.Repository
	artistRepo artist.Repository
	logger     logger.Logger
}

func NewUsecase(alr album.Repository, arr artist.Repository, l logger.Logger) *Usecase {
	return &Usecase{albumRepo: alr, artistRepo: arr, logger: l}
}

func (u *Usecase) Create(album models.Album, artistsID []uint32, userID uint32) (uint32, error) {
	userInArtists := false
	for _, artistID := range artistsID {
		a, err := u.artistRepo.GetByID(artistID)
		if err != nil {
			return 0, fmt.Errorf("(usecase) can't get artist with id #%d: %w", artistID, err)
		}
		if a.UserID != nil && *a.UserID == userID {
			userInArtists = true
			break
		}
	}
	if !userInArtists {
		return 0, fmt.Errorf("(usecase) album can't be created by user: %w", &models.ForbiddenUserError{})
	}

	albumID, err := u.albumRepo.Insert(album, artistsID)
	if err != nil {
		return 0, fmt.Errorf("(usecase) can't insert album into repository: %w", err)
	}

	return albumID, nil
}

func (u *Usecase) GetByID(albumID uint32) (*models.Album, error) {
	album, err := u.albumRepo.GetByID(albumID)
	if err != nil {
		return nil, fmt.Errorf("(usecase) can't get album from repository: %w", err)
	}

	return album, nil
}

func (u *Usecase) Delete(albumID uint32, userID uint32) error {
	if err := u.albumRepo.Check(albumID); err != nil {
		return fmt.Errorf("(usecase) can't find album with id #%d: %w", albumID, err)
	}

	userInArtists := false
	artists, err := u.artistRepo.GetByAlbum(albumID)
	if err != nil {
		return fmt.Errorf("(usecase) can't get artists of album: %w", err)
	}
	for _, artist := range artists {
		if artist.UserID != nil && *artist.UserID == userID {
			userInArtists = true
			break
		}
	}
	if !userInArtists {
		return fmt.Errorf("(usecase) album can't be deleted by user: %w", &models.ForbiddenUserError{})
	}

	if err := u.albumRepo.DeleteByID(albumID); err != nil {
		return fmt.Errorf("(usecase) can't delete album from repository: %w", err)
	}

	return nil
}

func (u *Usecase) GetFeed() ([]models.Album, error) {
	albums, err := u.albumRepo.GetFeed(feedAlbumsAmountLimit)
	if err != nil {
		return nil, fmt.Errorf("(usecase) can't get feed albums from repository: %w", err)
	}

	return albums, nil
}

func (u *Usecase) GetByArtist(artistID uint32) ([]models.Album, error) {
	if err := u.artistRepo.Check(artistID); err != nil {
		return nil, fmt.Errorf("(usecase) can't find artist with id #%d: %w", artistID, err)
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

func (u *Usecase) SetLike(albumID, userID uint32) (bool, error) {
	if err := u.albumRepo.Check(albumID); err != nil {
		return false, fmt.Errorf("(usecase) can't find album with id #%d: %w", albumID, err)
	}

	isInserted, err := u.albumRepo.InsertLike(albumID, userID)
	if err != nil {
		return false, fmt.Errorf("(usecase) failed to set like: %w", err)
	}

	return isInserted, nil
}

func (u *Usecase) UnLike(albumID, userID uint32) (bool, error) {
	if err := u.albumRepo.Check(albumID); err != nil {
		return false, fmt.Errorf("(usecase) can't find album with id #%d: %w", albumID, err)
	}

	isDeleted, err := u.albumRepo.DeleteLike(albumID, userID)
	if err != nil {
		return false, fmt.Errorf("(usecase) failed to unset like: %w", err)
	}

	return isDeleted, nil
}

func (u *Usecase) IsLiked(albumID, userID uint32) (bool, error) {
	isLiked, err := u.albumRepo.IsLiked(albumID, userID)
	if err != nil {
		return false, fmt.Errorf("(usecase) can't check in repository if album is liked: %w", err)
	}

	return isLiked, nil
}
