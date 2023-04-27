package search

import (
	"context"

	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/models"
)

//go:generate mockgen -source=search.go -destination=mocks/mock.go

// Usecase includes bussiness logics methods to work with tracks
type Usecase interface {
	FindAlbums(ctx context.Context, query string) ([]models.Album, error)
	FindArtists(ctx context.Context, query string) ([]models.Artist, error)
	FindTracks(ctx context.Context, query string) ([]models.Track, error)
	FindPlaylists(ctx context.Context, query string) ([]models.Playlist, error)
}

// Repository includes DBMS-relatable methods to work with tracks
type Repository interface {
	FullTextSearchAlbums(ctx context.Context, query string) ([]models.Album, error)
	FullTextSearchArtists(ctx context.Context, query string) ([]models.Artist, error)
	FullTextSearchTracks(ctx context.Context, query string) ([]models.Track, error)
	FullTextSearchPlaylists(ctx context.Context, query string) ([]models.Playlist, error)
}

// Tables includes methods which return needed tables
// to work with tracks on repository-layer
type Tables interface {
	Albums() string
	Artists() string
	Tracks() string
	Playlists() string
}
