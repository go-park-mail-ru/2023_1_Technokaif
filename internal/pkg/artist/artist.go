package artist

import "github.com/go-park-mail-ru/2023_1_Technokaif/internal/models"

//go:generate mockgen -source=artist.go -destination=mocks/mock.go

// Usecase includes bussiness logics methods to work with artists
type Usecase interface {
	Create(artist models.Artist) (uint32, error)
	GetByID(artistID uint32) (*models.Artist, error)
	// Change(artist models.Artist) error
	Delete(artistID uint32, userID uint32) error
	GetFeed() ([]models.Artist, error)
	GetByAlbum(albumID uint32) ([]models.Artist, error)
	GetByTrack(trackID uint32) ([]models.Artist, error)
	GetLikedByUser(userID uint32) ([]models.Artist, error)
	SetLike(artistID, userID uint32) (bool, error)
	UnLike(artistID, userID uint32) (bool, error)
	// GetListens(artistID uint32) (uint32, error)
}

// Repository includes DBMS-relatable methods to work with artists
type Repository interface {
	// Insert creates new entry of artist in DB with given model
	Insert(artist models.Artist) (uint32, error)

	// GetByID returns one entry of artist in DB with given ID
	GetByID(artistID uint32) (*models.Artist, error)

	// Update replaces one entry of artist with given model's ID by given model
	Update(artist models.Artist) error

	// DeleteByID deletes one entry of artist with given ID
	DeleteByID(artistID uint32) error

	// GetFeed returns artist entries with biggest amount of likes per some duration
	GetFeed() ([]models.Artist, error)

	// GetByAlbum returns all artist entries related with album entry with given ID
	GetByAlbum(albumID uint32) ([]models.Artist, error)

	// GetByTrack returns all artist entries related with Track with given ID
	GetByTrack(trackID uint32) ([]models.Artist, error)

	// GetByAlbum returns all Artist entries with like entry of user with given ID
	GetLikedByUser(userID uint32) ([]models.Artist, error)

	InsertLike(artistID, userID uint32) (bool, error)

	DeleteLike(artistID, userID uint32) (bool, error)
	
	// GetLikes returns total likes related with artist with given ID
	// GetLikes(artistID uint 32) (uint32, error)

	// GetListens returns total listens of all track entries related with artist with given ID
	// GetListens(artistID uint32) (uint64, error)
}

// Tables includes methods which return needed tables
// to work with artists on repository-layer
type Tables interface {
	Artists() string
	ArtistsAlbums() string
	ArtistsTracks() string
	LikedArtists() string
}
