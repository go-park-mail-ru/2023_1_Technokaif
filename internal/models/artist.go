package models

type Artist struct {
	ID        uint32 `db:"id"`
	Name      string `db:"name"`
	AvatarSrc string `db:"avatar_src"`
}

type ArtistTransfer struct {
	ID        uint32 `json:"id"`
	Name      string `json:"name"`
	AvatarSrc string `json:"cover"`
}

// Converts []Artist to []ArtistTransfer
func ArtistTransferFromQuery(artists []Artist) []ArtistTransfer {
	at := make([]ArtistTransfer, 0, len(artists))
	for _, a := range artists {
		at = append(at, ArtistTransfer{
			ID:        a.ID,
			Name:      a.Name,
			AvatarSrc: a.AvatarSrc,
		})
	}

	return at
}
