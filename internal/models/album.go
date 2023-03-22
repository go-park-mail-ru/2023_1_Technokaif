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

type artistsByAlbumGetter func(albumID uint32) ([]Artist, error)

// AlbumTransferFromEntry converts Album to AlbumTransfer
func AlbumTransferFromEntry(a Album, artistsGetter artistsByAlbumGetter) (AlbumTransfer, error) {
	artists, err := artistsGetter(a.ID)
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
func AlbumTransferFromQuery(albums []Album, artistsGetter artistsByAlbumGetter) ([]AlbumTransfer, error) {
	albumTransfers := make([]AlbumTransfer, 0, len(albums))
	for _, a := range albums {
		at, err := AlbumTransferFromEntry(a, artistsGetter)
		if err != nil {
			return nil, err
		}

		albumTransfers = append(albumTransfers, at)
	}

	return albumTransfers, nil
}
