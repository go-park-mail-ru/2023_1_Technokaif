package usecase

import (
	"context"
	"fmt"

	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/models"
	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/search"
)

// More likely API-user will set this limits later
const searchAlbumsAmountLimit uint32 = 3
const searchArtistsAmountLimit uint32 = 3
const searchTracksAmountLimit uint32 = 3
const searchPlaylistsAmountLimit uint32 = 3

// Usecase implements search.Usecase
type Usecase struct {
	searchRepo search.Repository
}

func NewUsecase(sr search.Repository) *Usecase {
	return &Usecase{
		searchRepo: sr,
	}
}

func (u *Usecase) FindAlbums(ctx context.Context, query string) ([]models.Album, error) {
	albums, err := u.searchRepo.FullTextSearchAlbums(ctx, query, searchAlbumsAmountLimit)
	if err != nil {
		return nil, fmt.Errorf("(usecase) can't find albums by query: %w", err)
	}
	return albums, nil
}

func (u *Usecase) FindArtists(ctx context.Context, query string) ([]models.Artist, error) {
	artists, err := u.searchRepo.FullTextSearchArtists(ctx, query, searchArtistsAmountLimit)
	if err != nil {
		return nil, fmt.Errorf("(usecase) can't find artists by query: %w", err)
	}
	return artists, nil
}

func (u *Usecase) FindTracks(ctx context.Context, query string) ([]models.Track, error) {
	tracks, err := u.searchRepo.FullTextSearchTracks(ctx, query, searchTracksAmountLimit)
	if err != nil {
		return nil, fmt.Errorf("(usecase) can't find tracks by query: %w", err)
	}
	return tracks, nil
}

func (u *Usecase) FindPlaylists(ctx context.Context, query string) ([]models.Playlist, error) {
	playlists, err := u.searchRepo.FullTextSearchPlaylists(ctx, query, searchPlaylistsAmountLimit)
	if err != nil {
		return nil, fmt.Errorf("(usecase) can't find playlists by query: %w", err)
	}
	return playlists, nil
}
