package usecase

import (
	"context"
	"fmt"

	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/models"
	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/search"
)

// Usecase implements search.Usecase
type Usecase struct {
	searchRepo search.Repository
}

func NewUsecase(sr search.Repository) *Usecase {
	return &Usecase{
		searchRepo: sr,
	}
}

func (u *Usecase) FindAlbums(ctx context.Context, query string, amount uint32) ([]models.Album, error) {
	albums, err := u.searchRepo.FullTextSearchAlbums(ctx, query, amount)
	if err != nil {
		return nil, fmt.Errorf("(usecase) can't find albums by query: %w", err)
	}
	return albums, nil
}

func (u *Usecase) FindArtists(ctx context.Context, query string, amount uint32) ([]models.Artist, error) {
	artists, err := u.searchRepo.FullTextSearchArtists(ctx, query, amount)
	if err != nil {
		return nil, fmt.Errorf("(usecase) can't find artists by query: %w", err)
	}
	return artists, nil
}

func (u *Usecase) FindTracks(ctx context.Context, query string, amount uint32) ([]models.Track, error) {
	tracks, err := u.searchRepo.FullTextSearchTracks(ctx, query, amount)
	if err != nil {
		return nil, fmt.Errorf("(usecase) can't find tracks by query: %w", err)
	}
	return tracks, nil
}

func (u *Usecase) FindPlaylists(ctx context.Context, query string, amount uint32) ([]models.Playlist, error) {
	playlists, err := u.searchRepo.FullTextSearchPlaylists(ctx, query, amount)
	if err != nil {
		return nil, fmt.Errorf("(usecase) can't find playlists by query: %w", err)
	}
	return playlists, nil
}
