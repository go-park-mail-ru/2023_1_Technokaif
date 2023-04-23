package models

import "fmt"

// Track errors

type NoSuchTrackError struct {
	TrackID uint32
}

func (e *NoSuchTrackError) Error() string {
	return fmt.Sprintf("track #%d doesn't exist", e.TrackID)
}

// Album errors

type NoSuchAlbumError struct {
	AlbumID uint32
}

func (e *NoSuchAlbumError) Error() string {
	return fmt.Sprintf("album #%d doesn't exist", e.AlbumID)
}

// Playlist errors

type NoSuchPlaylistError struct {
	PlaylistID uint32
}

func (e *NoSuchPlaylistError) Error() string {
	return fmt.Sprintf("playlist #%d doesn't exist", e.PlaylistID)
}

// Artist errors

type NoSuchArtistError struct {
	ArtistID uint32
}

func (e *NoSuchArtistError) Error() string {
	return fmt.Sprintf("artist #%d doesn't exist", e.ArtistID)
}

// Auth errors
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

type IncorrectPasswordError struct {
	UserId uint32
}

func (e *IncorrectPasswordError) Error() string {
	return fmt.Sprintf("incorrect password for user #%d", e.UserId)
}

type UnathorizedError struct{}

func (e *UnathorizedError) Error() string {
	return "unathorized"
}

type AvatarWrongFormatError struct {
	FileType string
}

func (e *AvatarWrongFormatError) Error() string {
	return fmt.Sprintf("avatar wrong format: %s", e.FileType)
}

type CoverWrongFormatError struct {
	FileType string
}

func (e *CoverWrongFormatError) Error() string {
	return fmt.Sprintf("acover wrong format: %s", e.FileType)
}
