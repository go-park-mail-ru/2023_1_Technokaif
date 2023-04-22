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
	AvatarSrc string `json:"cover"`
}

// ArtistTransferFromEntry converts Artist to ArtistTransfer
func ArtistTransferFromEntry(a Artist) ArtistTransfer {
	return ArtistTransfer{
		ID:        a.ID,
		Name:      a.Name,
		AvatarSrc: a.AvatarSrc,
	}
}

// ArtistTransferFromQuery converts []Artist to []ArtistTransfer
func ArtistTransferFromQuery(artists []Artist) []ArtistTransfer {
	artistTransfers := make([]ArtistTransfer, 0, len(artists))
	for _, a := range artists {
		artistTransfers = append(artistTransfers, ArtistTransferFromEntry(a))
	}

	return artistTransfers
}
