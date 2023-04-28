package delivery

import "github.com/go-park-mail-ru/2023_1_Technokaif/internal/models"

// Response messages
const (
	albumsFindNoRights    = "no rights to find albums"
	artistsFindNoRights   = "no rights to find artists"
	tracksFindNoRights    = "no rights to find tracks"
	playlistsFindNoRights = "no rights to find playlists"

	albumsFindServerError    = "can't find albums"
	artistsFindServerError   = "can't find artists"
	tracksFindServerError    = "can't find tracks"
	playlistsFindServerError = "can't find playlists"
)

type SearchRequest struct {
	Query string `json:"query" valid:"required"`
	Limit uint32 `json:"limit" valid:"required,range(1|100)"`
}

type SearchAlbumsResponse struct {
	Albums []models.AlbumTransfer `json:"albums"`
}

type SearchArtistsResponse struct {
	Artists []models.ArtistTransfer `json:"artists"`
}

type SearchTracksResponse struct {
	Tracks []models.Track `json:"tracks"`
}

type SearchPlaylistsResponse struct {
	Playlists []models.Playlist `json:"playlists"`
}
