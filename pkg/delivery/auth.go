package delivery

import (
	"encoding/json"
	"net/http"

	"github.com/go-park-mail-ru/2023_1_Technokaif/models"
)

type responseSignup struct {
	ID int `json:"id"`
}

func (h *Handler) signUp(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	user := new(models.User)
	if err := decoder.Decode(user); err != nil {
		httpErrorResponce(w, "incorrect input body", http.StatusBadRequest)
		return
	}

	if !user.Validate() {
		httpErrorResponce(w, "incorrect fileds values", http.StatusBadRequest)
		return
	}

	id, err := h.services.Auth.CreateUser(*user)
	if err != nil {
		httpErrorResponce(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// TODO maybe make wrapper for responses
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	response := responseSignup{ID: id}
	encoder := json.NewEncoder(w)
	if err := encoder.Encode(&response); err != nil {
		httpErrorResponce(w, err.Error(), http.StatusInternalServerError) // TODO change error message
		return
	}
}

func (h *Handler) login(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(r.Method + " Login Page"))
}

func (h *Handler) logout(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("You've been successfully logout"))
}
