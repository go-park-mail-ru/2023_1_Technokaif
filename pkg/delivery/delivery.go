package delivery

import (
	"net/http"

	"github.com/go-chi/chi"
)

func index(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("I didn't hit her! I did not... Oh, hi, Mark!"))
}

func login(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(r.Method + " Login Page"))
}

func signup(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(r.Method + " Signup page"))
}

func logout(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("You've been successfully logout"))
}

// InitRouter describes all app's endpoints and their handlers
func InitRouter() *chi.Mux {
	r := chi.NewRouter()

	r.Get("/", index)

	r.Route("/auth", func(r chi.Router) {
		r.Route("/login", func(r chi.Router) {
			r.Get("/", login)
			r.Post("/", login)
		})

		r.Route("/signup", func(r chi.Router) {
			r.Get("/", signup)
			r.Post("/", signup)
		})

		r.Get("/logout", logout)
	})

	return r
}
