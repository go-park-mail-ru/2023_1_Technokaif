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
type ForbiddenUserError struct{}

func (e *ForbiddenUserError) Error() string {
	return "user has no rights"
}

type UserAlreadyExistsError struct{}

func (e *UserAlreadyExistsError) Error() string {
	return "user already exists"
}

type NoSuchUserError struct {
	UserID uint32
}

func (e *NoSuchUserError) Error() string {
	return fmt.Sprintf("user #%d doesn't exist", e.UserID)
}

type IncorrectPasswordError struct{
	UserId uint32
}

func (e *IncorrectPasswordError) Error() string {
	return fmt.Sprintf("incorrect password for user #%d", e.UserId)
}
