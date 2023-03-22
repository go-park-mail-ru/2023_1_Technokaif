package models

type Album struct {
	ID          uint32 `db:"id"`
	Name        string `db:"name"`
	Description string `db:"description"`
	CoverSrc    string `db:"cover_src"`
}

type AlbumTransfer struct {
	ID          uint32           `json:"id"`
	Name        string           `json:"name"`
	Artists     []ArtistTransfer `json:"artists"`
	Description string           `json:"description"`
	CoverSrc    string           `json:"cover"`
}

type artistByAlbumGetter func(albumID uint32) ([]Artist, error)

// Converts Album to AlbumTransfer
func AlbumTransferFromEntry(a Album, artistGetter artistByAlbumGetter) (AlbumTransfer, error) {
	artists, err := artistGetter(a.ID)
	if err != nil {
		return AlbumTransfer{}, err
	}

	return AlbumTransfer{
		ID:          a.ID,
		Name:        a.Name,
		Artists:     ArtistTransferFromQuery(artists),
		Description: a.Description,
		CoverSrc:    a.CoverSrc,
	}, nil
}

// Convert []Album to []AlbumTransfer
func AlbumTransferFromQuery(albums []Album, artistGetter artistByAlbumGetter) ([]AlbumTransfer, error) {
	albumTransfers := make([]AlbumTransfer, 0, len(albums))
	for _, a := range albums {
		at, err := AlbumTransferFromEntry(a, artistGetter)
		if err != nil {
			return nil, err
		}

		albumTransfers = append(albumTransfers, at)
	}

	return albumTransfers, nil
}
