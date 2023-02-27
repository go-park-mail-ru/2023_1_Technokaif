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

func (h *Handler) feed(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	/* userId, ok := r.Context().Value("ID").(uint)
	if !ok {
		httpErrorResponce(w, "api error", http.StatusInternalServerError)
		} */

	artists, err := h.services.Artist.GetFeed()
	if err != nil {
		w.Write([]byte("Error while getting artists"))
		return
	}

	tracks, err := h.services.Track.GetFeed()
	if err != nil {
		w.Write([]byte("Error while getting tracks"))
		return
	}

	albums, err := h.services.Album.GetFeed()
	if err != nil {
		w.Write([]byte("Error while getting albums"))
		return
	}
	w.Header().Set("Content-Type", "json/application; charset=utf-8")

	fr := FeedResponse{Artists: artists,
		Tracks: tracks,
		Albums: albums}

	encoder := json.NewEncoder(w)
	if err := encoder.Encode(&fr); err != nil {
		httpErrorResponce(w, err.Error()+" json upal :(", http.StatusInternalServerError)
		return
	}

}
