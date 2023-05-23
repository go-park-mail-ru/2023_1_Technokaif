package delivery

import (
	valid "github.com/asaskevich/govalidator"
	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/models"
)

//go:generate easyjson -no_std_marshalers search_delivery_models.go

// Response messages
const (
	albumsFindServerError    = "can't find albums"
	artistsFindServerError   = "can't find artists"
	tracksFindServerError    = "can't find tracks"
	playlistsFindServerError = "can't find playlists"
)

//easyjson:json
type searchRequest struct {
	Query  string `json:"query" valid:"required"`
	Amount uint32 `json:"amount" valid:"required,range(1|100)"`
}

func (sr *searchRequest) validate() error {
	_, err := valid.ValidateStruct(sr)
	return err
}

//easyjson:json
type searchAlbumsResponse struct {
	Albums models.AlbumTransfers `json:"albums"`
}

//easyjson:json
type searchArtistsResponse struct {
	Artists models.ArtistTransfers `json:"artists"`
}

//easyjson:json
type searchTracksResponse struct {
	Tracks models.TrackTransfers `json:"tracks"`
}

//easyjson:json
type searchPlaylistsResponse struct {
	Playlists models.PlaylistTransfers `json:"playlists"`
}
