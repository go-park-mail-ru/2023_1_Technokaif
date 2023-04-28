package usecase

import (
	"context"
	"fmt"

	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/models"
	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/album"
	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/artist"
	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/playlist"
	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/track"
)

const feedTracksAmountLimit uint32 = 100

// Usecase implements track.Usecase
type Usecase struct {
	trackRepo    track.Repository
	artistRepo   artist.Repository
	albumRepo    album.Repository
	playlistRepo playlist.Repository
}

func NewUsecase(tr track.Repository, arr artist.Repository,
	alr album.Repository, pr playlist.Repository) *Usecase {

	return &Usecase{
		trackRepo:    tr,
		artistRepo:   arr,
		albumRepo:    alr,
		playlistRepo: pr,
	}
}

func (u *Usecase) Create(ctx context.Context,
	track models.Track, artistsID []uint32, userID uint32) (uint32, error) {

	userInArtists := false
	for _, artistID := range artistsID {
		a, err := u.artistRepo.GetByID(ctx, artistID)
		if err != nil {
			return 0, fmt.Errorf("(usecase) can't get artist with id #%d: %w", artistID, err)
		}
		if a.UserID != nil && *a.UserID == userID {
			userInArtists = true
			break
		}
	}
	if !userInArtists {
		return 0, fmt.Errorf("(usecase) track can't be created by user: %w", &models.ForbiddenUserError{})
	}

	trackID, err := u.trackRepo.Insert(ctx, track, artistsID)
	if err != nil {
		return 0, fmt.Errorf("(usecase) can't insert track into repository: %w", err)
	}

	return trackID, nil
}

func (u *Usecase) GetByID(ctx context.Context, trackID uint32) (*models.Track, error) {
	track, err := u.trackRepo.GetByID(ctx, trackID)
	if err != nil {
		return &models.Track{}, fmt.Errorf("(usecase) can't get track with id #%d: %w", trackID, err)
	}

	return track, nil
}

func (u *Usecase) Delete(ctx context.Context, trackID uint32, userID uint32) error {
	if err := u.trackRepo.Check(ctx, trackID); err != nil {
		return fmt.Errorf("(usecase) can't find track with id #%d: %w", trackID, err)
	}

	userInArtists := false
	artists, err := u.artistRepo.GetByAlbum(ctx, trackID)
	if err != nil {
		return fmt.Errorf("(usecase) can't get artists of track: %w", err)
	}
	for _, artist := range artists {
		if artist.UserID != nil && *artist.UserID == userID {
			userInArtists = true
		}
	}
	if !userInArtists {
		return fmt.Errorf("(usecase) track can't be deleted by user: %w", &models.ForbiddenUserError{})
	}

	if err := u.trackRepo.DeleteByID(ctx, trackID); err != nil {
		return fmt.Errorf("(usecase) can't delete track from repository: %w", err)
	}

	return nil
}

func (u *Usecase) GetFeed(ctx context.Context) ([]models.Track, error) {
	tracks, err := u.trackRepo.GetFeed(ctx, feedTracksAmountLimit)
	if err != nil {
		return nil, fmt.Errorf("(usecase) can't get feed tracks from repository: %w", err)
	}

	return tracks, nil
}

func (u *Usecase) GetByAlbum(ctx context.Context, albumID uint32) ([]models.Track, error) {
	if err := u.albumRepo.Check(ctx, albumID); err != nil {
		return nil, fmt.Errorf("(usecase) can't find album with id #%d: %w", albumID, err)
	}

	tracks, err := u.trackRepo.GetByAlbum(ctx, albumID)
	if err != nil {
		return nil, fmt.Errorf("(usecase) can't get tracks from repository: %w", err)
	}

	return tracks, nil
}

func (u *Usecase) GetByPlaylist(ctx context.Context, playlistID uint32) ([]models.Track, error) {
	if err := u.playlistRepo.Check(ctx, playlistID); err != nil {
		return nil, fmt.Errorf("(usecase) can't find playlist with id #%d: %w", playlistID, err)
	}

	tracks, err := u.trackRepo.GetByPlaylist(ctx, playlistID)
	if err != nil {
		return nil, fmt.Errorf("(usecase) can't get tracks from repository: %w", err)
	}

	return tracks, nil
}

func (u *Usecase) GetByArtist(ctx context.Context, artistID uint32) ([]models.Track, error) {
	if err := u.artistRepo.Check(ctx, artistID); err != nil {
		return nil, fmt.Errorf("(usecase) can't find artist with id #%d: %w", artistID, err)
	}

	tracks, err := u.trackRepo.GetByArtist(ctx, artistID)
	if err != nil {
		return nil, fmt.Errorf("(usecase) can't get tracks from repository: %w", err)
	}

	return tracks, nil
}

func (u *Usecase) GetLikedByUser(ctx context.Context, userID uint32) ([]models.Track, error) {
	tracks, err := u.trackRepo.GetLikedByUser(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("(usecase) can't get tracks from repository: %w", err)
	}

	return tracks, nil
}

func (u *Usecase) SetLike(ctx context.Context, trackID, userID uint32) (bool, error) {
	if err := u.trackRepo.Check(ctx, trackID); err != nil {
		return false, fmt.Errorf("(usecase) can't find track with id #%d: %w", trackID, err)
	}

	isInserted, err := u.trackRepo.InsertLike(ctx, trackID, userID)
	if err != nil {
		return false, fmt.Errorf("(usecase) failed to set like: %w", err)
	}

	return isInserted, nil
}

func (u *Usecase) UnLike(ctx context.Context, trackID, userID uint32) (bool, error) {
	if err := u.trackRepo.Check(ctx, trackID); err != nil {
		return false, fmt.Errorf("(usecase) can't find track with id #%d: %w", trackID, err)
	}

	isDeleted, err := u.trackRepo.DeleteLike(ctx, trackID, userID)
	if err != nil {
		return false, fmt.Errorf("(usecase) failed to unset like: %w", err)
	}

	return isDeleted, nil
}

func (u *Usecase) IsLiked(ctx context.Context, trackID, userID uint32) (bool, error) {
	isLiked, err := u.trackRepo.IsLiked(ctx, trackID, userID)
	if err != nil {
		return false, fmt.Errorf("(usecase) can't check in repository if track is liked: %w", err)
	}

	return isLiked, nil
}
