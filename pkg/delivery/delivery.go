package delivery

import (
	"net/http"

	"github.com/go-chi/chi"
)

func InitRouter() *chi.Mux {
	r := chi.NewRouter()

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("I didn't hit her, I did not... Oh, hi, Mark!"))
	})

	return r
}
