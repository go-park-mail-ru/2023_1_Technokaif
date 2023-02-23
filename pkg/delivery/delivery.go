package delivery

import (
	"net/http"

	"github.com/go-chi/chi"
)

type Handler struct {

}

func(h *Handler) index(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("I didn't hit her! I did not... Oh, hi, Mark!"))
}

func(h *Handler) login(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(r.Method + " Login Page"))
}

func(h *Handler) signup(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(r.Method + " Signup page"))
}

func(h *Handler) logout(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("You've been successfully logout"))
}

// InitRouter describes all app's endpoints and their handlers
func(h *Handler) InitRouter() *chi.Mux {
	r := chi.NewRouter()

	r.Get("/", h.index)

	r.Route("/auth", func(r chi.Router) {
		r.Route("/login", func(r chi.Router) {
			r.Get("/", h.login)
			r.Post("/", h.login)
		})

		r.Route("/signup", func(r chi.Router) {
			r.Get("/", h.signup)
			r.Post("/", h.signup)
		})

		r.Get("/logout", h.logout)
	})

	return r
}
