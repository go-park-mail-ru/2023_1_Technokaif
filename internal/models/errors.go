package models

import "fmt"

// TRACK ERRORS
type NoSuchTrackError struct {
	TrackID uint32
}

func (e *NoSuchTrackError) Error() string {
	return fmt.Sprintf("track #%d doesn't exist", e.TrackID)
}

// ALBUM ERRORS
type NoSuchAlbumError struct {
	AlbumID uint32
}

func (e *NoSuchAlbumError) Error() string {
	return fmt.Sprintf("album #%d doesn't exist", e.AlbumID)
}

// ARTIST ERRORS
type NoSuchArtistError struct {
	ArtistID uint32
}

func (e *NoSuchArtistError) Error() string {
	return fmt.Sprintf("artist #%d doesn't exist", e.ArtistID)
}

// AUTH ERRORS
type UserAlreadyExistsError struct{}
type NoSuchUserError struct{}

func (e *UserAlreadyExistsError) Error() string {
	return "user already exists"
}

func (e *NoSuchUserError) Error() string {
	return "no such user"
}
