package models

type Artist struct {
	ID        uint32  `db:"id"`
	UserID    *uint32 `db:"user_id"`
	Name      string  `db:"name"`
	AvatarSrc string  `db:"avatar_src"`
}

type ArtistTransfer struct {
	ID        uint32 `json:"id"`
	Name      string `json:"name"`
	IsLiked   bool   `json:"isLiked"`
	AvatarSrc string `json:"cover"`
}

type artistLikeChecker func(artistID, userID uint32) (bool, error)

// ArtistTransferFromEntry converts Artist to ArtistTransfer
func ArtistTransferFromEntry(a Artist, user *User,
	likeChecker artistLikeChecker) (ArtistTransfer, error) {

	var isLiked bool
	var err error

	if user != nil {
		isLiked, err = likeChecker(a.ID, user.ID)
		if err != nil {
			return ArtistTransfer{}, err
		}
	}

	return ArtistTransfer{
		ID:        a.ID,
		Name:      a.Name,
		IsLiked:   isLiked,
		AvatarSrc: a.AvatarSrc,
	}, nil
}

// ArtistTransferFromQuery converts []Artist to []ArtistTransfer
func ArtistTransferFromQuery(artists []Artist, user *User,
	likeChecker artistLikeChecker) ([]ArtistTransfer, error) {

	artistTransfers := make([]ArtistTransfer, 0, len(artists))
	for _, a := range artists {
		at, err := ArtistTransferFromEntry(a, user, likeChecker)
		if err != nil {
			return nil, err
		}

		artistTransfers = append(artistTransfers, at)
	}

	return artistTransfers, nil
}
