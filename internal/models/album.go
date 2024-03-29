package models

import "context"

//go:generate easyjson -no_std_marshalers album.go

type Album struct {
	ID          uint32  `db:"id"`
	Name        string  `db:"name"`
	Description *string `db:"description"`
	CoverSrc    string  `db:"cover_src"`
}

//easyjson:json
type AlbumTransfer struct {
	ID          uint32          `json:"id"`
	Name        string          `json:"name"`
	Artists     ArtistTransfers `json:"artists"`
	Description *string         `json:"description,omitempty"`
	IsLiked     bool            `json:"isLiked"`
	CoverSrc    string          `json:"cover"`
}

//easyjson:json
type AlbumTransfers []AlbumTransfer

type artistsByAlbumGetter func(ctx context.Context, albumID uint32) ([]Artist, error)
type albumLikeChecker func(ctx context.Context, albumID, userID uint32) (bool, error)

// AlbumTransferFromEntry converts Album to AlbumTransfer
func AlbumTransferFromEntry(ctx context.Context, a Album, user *User, likeChecker albumLikeChecker,
	artistLikeChecker ArtistLikeChecker, artistsGetter artistsByAlbumGetter) (AlbumTransfer, error) {

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

	at, err := ArtistTransferFromList(ctx, artists, user, artistLikeChecker)
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

// AlbumTransferFromList converts []Album to []AlbumTransfer
func AlbumTransferFromList(ctx context.Context, albums []Album, user *User, likeChecker albumLikeChecker,
	artistLikeChecker ArtistLikeChecker, artistsGetter artistsByAlbumGetter) (AlbumTransfers, error) {

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
