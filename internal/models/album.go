package models

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

type artistsByAlbumGetter func(albumID uint32) ([]Artist, error)
type albumLikeChecker func(albumID, userID uint32) (bool, error)

// AlbumTransferFromEntry converts Album to AlbumTransfer
func AlbumTransferFromEntry(a Album, user *User, likeChecker albumLikeChecker,
	artistLikeChecker artistLikeChecker, artistsGetter artistsByAlbumGetter) (AlbumTransfer, error) {

	artists, err := artistsGetter(a.ID)
	if err != nil {
		return AlbumTransfer{}, err
	}

	var isLiked = false
	if user != nil {
		isLiked, err = likeChecker(a.ID, user.ID)
		if err != nil {
			return AlbumTransfer{}, err
		}
	}

	at, err := ArtistTransferFromQuery(artists, user, artistLikeChecker)
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
func AlbumTransferFromQuery(albums []Album, user *User, likeChecker albumLikeChecker,
	artistLikeChecker artistLikeChecker, artistsGetter artistsByAlbumGetter) ([]AlbumTransfer, error) {

	albumTransfers := make([]AlbumTransfer, 0, len(albums))
	for _, a := range albums {
		at, err := AlbumTransferFromEntry(a, user, likeChecker, artistLikeChecker, artistsGetter)
		if err != nil {
			return nil, err
		}

		albumTransfers = append(albumTransfers, at)
	}

	return albumTransfers, nil
}
