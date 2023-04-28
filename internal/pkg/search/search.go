package search

import (
	"context"

	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/models"
)

//go:generate mockgen -source=search.go -destination=mocks/mock.go

// Usecase includes bussiness logics methods to work with search
type Usecase interface {
	FindAlbums(ctx context.Context, query string) ([]models.Album, error)
	FindArtists(ctx context.Context, query string) ([]models.Artist, error)
	FindTracks(ctx context.Context, query string) ([]models.Track, error)
	FindPlaylists(ctx context.Context, query string) ([]models.Playlist, error)
}

// Repository includes DBMS-relatable methods to work with search
type Repository interface {
	FullTextSearchAlbums(ctx context.Context, query string, limit uint32) ([]models.Album, error)
	FullTextSearchArtists(ctx context.Context, query string, limit uint32) ([]models.Artist, error)
	FullTextSearchTracks(ctx context.Context, query string, limit uint32) ([]models.Track, error)
	FullTextSearchPlaylists(ctx context.Context, query string, limit uint32) ([]models.Playlist, error)
}

// Tables includes methods which return needed tables
// to work with search on repository layer
type Tables interface {
	Albums() string
	Artists() string
	Tracks() string
	Playlists() string
}
