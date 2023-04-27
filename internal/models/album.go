package models

import (
	"context"
)

type Album struct {
	ID          uint32  `db:"id"`
	Name        string  `db:"name"`
	Description *string `db:"description"`
	CoverSrc    string  `db:"cover_src"`
}

type AlbumTransfer struct {
	ID          uint32           `json:"id"`
	Name        string           `json:"name"`
	Artists     []ArtistTransfer `json:"artists"`
	Description *string          `json:"description,omitempty"`
	IsLiked     bool             `json:"isLiked"`
	CoverSrc    string           `json:"cover"`
}

type artistsByAlbumGetter func(ctx context.Context, albumID uint32) ([]Artist, error)
type albumLikeChecker func(ctx context.Context, albumID, userID uint32) (bool, error)

// AlbumTransferFromEntry converts Album to AlbumTransfer
func AlbumTransferFromEntry(ctx context.Context, a Album, user *User, likeChecker albumLikeChecker,
	artistLikeChecker artistLikeChecker, artistsGetter artistsByAlbumGetter) (AlbumTransfer, error) {

	artists, err := artistsGetter(ctx, a.ID)
	if err != nil {
		return AlbumTransfer{}, err
	}

	var isLiked = false
	if user != nil {
		isLiked, err = likeChecker(ctx, a.ID, user.ID)
		if err != nil {
			return AlbumTransfer{}, err
		}
	}

	at, err := ArtistTransferFromQuery(ctx, artists, user, artistLikeChecker)
	if err != nil {
		return AlbumTransfer{}, err
	}

	return AlbumTransfer{
		ID:          a.ID,
		Name:        a.Name,
		Artists:     at,
		Description: a.Description,
		IsLiked:     isLiked,
		CoverSrc:    a.CoverSrc,
	}, nil
}

// AlbumTransferFromQuery converts []Album to []AlbumTransfer
func AlbumTransferFromQuery(ctx context.Context, albums []Album, user *User, likeChecker albumLikeChecker,
	artistLikeChecker artistLikeChecker, artistsGetter artistsByAlbumGetter) ([]AlbumTransfer, error) {

	albumTransfers := make([]AlbumTransfer, 0, len(albums))
	for _, a := range albums {
		at, err := AlbumTransferFromEntry(ctx, a, user, likeChecker, artistLikeChecker, artistsGetter)
		if err != nil {
			return nil, err
		}

		albumTransfers = append(albumTransfers, at)
	}

	return albumTransfers, nil
}
