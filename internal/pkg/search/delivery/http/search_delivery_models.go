package delivery

import (
	valid "github.com/asaskevich/govalidator"
	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/models"
)

// Response messages
const (
	albumsFindServerError    = "can't find albums"
	artistsFindServerError   = "can't find artists"
	tracksFindServerError    = "can't find tracks"
	playlistsFindServerError = "can't find playlists"
)

type searchRequest struct {
	Query  string `json:"query" valid:"required"`
	Amount uint32 `json:"amount" valid:"required,range(1|100)"`
}

func (sr *searchRequest) validate() error {
	_, err := valid.ValidateStruct(sr)
	return err
}

type searchAlbumsResponse struct {
	Albums []models.AlbumTransfer `json:"albums"`
}

type searchArtistsResponse struct {
	Artists []models.ArtistTransfer `json:"artists"`
}

type searchTracksResponse struct {
	Tracks []models.TrackTransfer `json:"tracks"`
}

type searchPlaylistsResponse struct {
	Playlists []models.PlaylistTransfer `json:"playlists"`
}
