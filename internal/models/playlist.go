package models

import "context"

//go:generate easyjson -no_std_marshalers playlist.go

type Playlist struct {
	ID          uint32  `db:"id"`
	Name        string  `db:"name"`
	Description *string `db:"description"`
	CoverSrc    string  `db:"cover_src"`
}

//easyjson:json
type PlaylistTransfer struct {
	ID          uint32        `json:"id"`
	Name        string        `json:"name"`
	Users       UserTransfers `json:"users"`
	Description *string       `json:"description,omitempty"`
	IsLiked     bool          `json:"isLiked"`
	CoverSrc    string        `json:"cover,omitempty"`
}

//easyjson:json
type PlaylistTransfers []PlaylistTransfer

type usersByPlaylistsGetter func(ctx context.Context, playlistID uint32) ([]User, error)
type playlistLikeChecker func(ctx context.Context, playlistID, userID uint32) (bool, error)

// PlaylistTransferFromEntry converts Playlist to PlaylistTransfer
func PlaylistTransferFromEntry(ctx context.Context, p Playlist, user *User,
	likeChecker playlistLikeChecker, usersGetter usersByPlaylistsGetter) (PlaylistTransfer, error) {

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
		Users:       UserTransferFromList(users),
		Description: p.Description,
		IsLiked:     isLiked,
		CoverSrc:    p.CoverSrc,
	}, nil
}

// PlaylistTransferFromList converts []Playlist to []PlaylistTransfer
func PlaylistTransferFromList(ctx context.Context, playlists []Playlist, user *User, likeChecker playlistLikeChecker,
	usersGetter usersByPlaylistsGetter) (PlaylistTransfers, error) {

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
