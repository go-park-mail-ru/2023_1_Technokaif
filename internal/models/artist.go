package models

import "context"

//go:generate easyjson -no_std_marshalers artist.go

type Artist struct {
	ID        uint32  `db:"id"`
	UserID    *uint32 `db:"user_id"`
	Name      string  `db:"name"`
	AvatarSrc string  `db:"avatar_src"`
	Listens   uint32  `db:"listens"`
}

//easyjson:json
type ArtistTransfer struct {
	ID        uint32 `json:"id"`
	Name      string `json:"name"`
	IsLiked   bool   `json:"isLiked"`
	AvatarSrc string `json:"cover"`
	Listens   uint32 `json:"listens"`
}

//easyjson:json
type ArtistTransfers []ArtistTransfer

type ArtistLikeChecker func(ctx context.Context, artistID, userID uint32) (bool, error)

// ArtistTransferFromEntry converts Artist to ArtistTransfer
func ArtistTransferFromEntry(ctx context.Context, a Artist, user *User,
	likeChecker ArtistLikeChecker) (ArtistTransfer, error) {

	var isLiked bool
	var err error

	if user != nil {
		isLiked, err = likeChecker(ctx, a.ID, user.ID)
		if err != nil {
			return ArtistTransfer{}, err
		}
	}

	return ArtistTransfer{
		ID:        a.ID,
		Name:      a.Name,
		IsLiked:   isLiked,
		AvatarSrc: a.AvatarSrc,
		Listens:   a.Listens,
	}, nil
}

// ArtistTransferFromList converts []Artist to []ArtistTransfer
func ArtistTransferFromList(ctx context.Context, artists []Artist, user *User,
	likeChecker ArtistLikeChecker) (ArtistTransfers, error) {

	artistTransfers := make([]ArtistTransfer, 0, len(artists))
	for _, a := range artists {
		at, err := ArtistTransferFromEntry(ctx, a, user, likeChecker)
		if err != nil {
			return nil, err
		}

		artistTransfers = append(artistTransfers, at)
	}

	return artistTransfers, nil
}
