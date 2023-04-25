package models

import (
	"context"
)

type Playlist struct {
	ID          uint32  `db:"id"`
	Name        string  `db:"name"`
	Description *string `db:"description"`
	CoverSrc    string  `db:"cover_src"`
}

type PlaylistTransfer struct {
	ID          uint32         `json:"id"`
	Name        string         `json:"name"`
	Users       []UserTransfer `json:"users"`
	Description *string        `json:"description,omitempty"`
	IsLiked     bool           `json:"isLiked"`
	CoverSrc    string         `json:"cover,omitempty"`
}

type usersByPlaylistsGetter func(ctx context.Context, playlistID uint32) ([]User, error)
type playlistLikeChecker func(ctx context.Context, playlistID, userID uint32) (bool, error)

// PlaylistTransferFromEntry converts Playlist to PlaylistTransfer
func PlaylistTransferFromEntry(ctx context.Context, p Playlist, user *User, likeChecker playlistLikeChecker,
	usersGetter usersByPlaylistsGetter) (PlaylistTransfer, error) {

	users, err := usersGetter(ctx, p.ID)
	if err != nil {
		return PlaylistTransfer{}, err
	}

	isLiked := false
	if user != nil {
		isLiked, err = likeChecker(ctx, p.ID, user.ID)
		if err != nil {
			return PlaylistTransfer{}, err
		}
	}

	return PlaylistTransfer{
		ID:          p.ID,
		Name:        p.Name,
		Users:       UserTransferFromQuery(users),
		Description: p.Description,
		IsLiked:     isLiked,
		CoverSrc:    p.CoverSrc,
	}, nil
}

// PlaylistTransferFromQuery converts []Playlist to []PlaylistTransfer
func PlaylistTransferFromQuery(ctx context.Context, playlists []Playlist, user *User, likeChecker playlistLikeChecker,
	usersGetter usersByPlaylistsGetter) ([]PlaylistTransfer, error) {

	playlistTransfers := make([]PlaylistTransfer, 0, len(playlists))

	for _, p := range playlists {
		pt, err := PlaylistTransferFromEntry(ctx, p, user, likeChecker, usersGetter)
		if err != nil {
			return nil, err
		}

		playlistTransfers = append(playlistTransfers, pt)
	}

	return playlistTransfers, nil
}
