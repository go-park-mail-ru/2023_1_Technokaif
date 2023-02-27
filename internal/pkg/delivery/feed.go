package delivery

import (
	"fmt"
	"net/http"
)

func (h *Handler) feed(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	fmt.Println("Into feed")
	/* userId, ok := r.Context().Value("ID").(uint)
	   if !ok {
	       httpErrorResponce(w, "api error", http.StatusInternalServerError)
	   } */

	tracks, err := h.services.Track.GetFeed()
	if err != nil {
		w.Write([]byte("FAIL"))
	}

	for _, track := range tracks {
		w.Write([]byte(track.Name + " " + track.ArtistName + "\n"))
	}

	w.Header().Set("Content-Type", "text; charset=utf-8")
}
