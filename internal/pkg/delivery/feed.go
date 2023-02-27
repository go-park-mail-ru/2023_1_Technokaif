package delivery

import (
	"net/http"
)

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

	for _, artist := range artists {
		w.Write([]byte(artist.Name + "\n"))
	}

	for _, track := range tracks {
		w.Write([]byte(track.Name + " " + track.ArtistName + "\n"))
	}

	for _, album := range albums {
		w.Write([]byte(album.Name + " " + album.ArtistName + "\n"))
	}

	w.Header().Set("Content-Type", "json/application; charset=utf-8")
}
