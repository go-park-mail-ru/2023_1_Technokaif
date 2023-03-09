package delivery

import (
	"encoding/json"
	"net/http"

	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/models"
)

type FeedResponse struct {
	Artists []models.ArtistFeed `json:"artists"`
	Tracks  []models.TrackFeed  `json:"tracks"`
	Albums  []models.AlbumFeed  `json:"albums"`
}

//	@Summary		Main Page
//	@Tags			feed
//	@Description	User's feed (Tracks, artists, albums)
//	@Accept			json
//	@Produce		json
//	@Success		200			{object}	FeedResponse	"Show feed"
//	@Failure		500			{object}	errorResponse	"Server DB error"
//	@Router			/api/feed [get]
func (h *Handler) feed(w http.ResponseWriter, r *http.Request) {
	artists, err := h.services.Artist.GetFeed()
	if err != nil {
		h.logger.Error(err.Error())
		h.errorResponse(w, "error while getting artists: "+err.Error(), http.StatusInternalServerError)
		return
	}

	tracks, err := h.services.Track.GetFeed()
	if err != nil {
		h.logger.Error(err.Error())
		h.errorResponse(w, "error while getting tracks: "+err.Error(), http.StatusInternalServerError)
		return
	}

	albums, err := h.services.Album.GetFeed()
	if err != nil {
		h.logger.Error(err.Error())
		h.errorResponse(w, "error while getting albums: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "json/application; charset=utf-8")

	fr := FeedResponse{
		Artists: artists,
		Tracks:  tracks,
		Albums:  albums}

	encoder := json.NewEncoder(w)
	if err := encoder.Encode(&fr); err != nil {
		h.logger.Error(err.Error())
		h.errorResponse(w, "can't encode response into json: "+err.Error(), http.StatusInternalServerError)
		return
	}
}
