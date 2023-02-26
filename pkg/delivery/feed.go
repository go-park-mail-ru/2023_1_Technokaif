package delivery

import (
	"net/http"
)

func (h Handler) feed(w http.ResponseWriter, rhttp.Request) {
    userId, ok := r.Context().Value("ID").(uint)
    if !ok {
        httpErrorResponce(w, "api error", http.StatusInternalServerError)
    }

}