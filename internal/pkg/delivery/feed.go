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

// Main Page
func (h *Handler) feed(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	/* userId, ok := r.Context().Value("ID").(uint)
	if !ok {
		httpErrorResponce(w, "api error", http.StatusInternalServerError)
	} */

	artists, err := h.services.Artist.GetFeed()
	if err != nil {
		h.logger.Error(err.Error())
		h.errorResponce(w, "error while getting artists", http.StatusInternalServerError)
		return
	}

	tracks, err := h.services.Track.GetFeed()
	if err != nil {
		h.logger.Error(err.Error())
		h.errorResponce(w, "error while getting tracks", http.StatusInternalServerError)
		return
	}

	albums, err := h.services.Album.GetFeed()
	if err != nil {
		h.logger.Error(err.Error())
		h.errorResponce(w, "error while getting albums", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "json/application; charset=utf-8")

	fr := FeedResponse{
		Artists: artists,
		Tracks:  tracks,
		Albums:  albums}

	encoder := json.NewEncoder(w)
	if err := encoder.Encode(&fr); err != nil {
		h.errorResponce(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
