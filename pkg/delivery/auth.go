package delivery

import (
	"encoding/json"
	"net/http"

	"github.com/go-park-mail-ru/2023_1_Technokaif/models"
)

// SIGN UP
type responseSignup struct {
	ID int `json:"id"`
}

func (h *Handler) signUp(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	user := new(models.User)
	err := decoder.Decode(user) 
	if err != nil || !user.Validate() {
		httpErrorResponce(w, "incorrect input body", http.StatusBadRequest)
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


// LOGIN
type loginInput struct {
	Username  string `json:"username"`
	Password  string `json:"password"`
}

func (i loginInput) validate() bool {
	if i.Username == "" || i.Password == "" {
		return false
	}
	return true
}

type responseLogin struct {
	JWT string `json:"jwt"`
}

func (h *Handler) login(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	userInput := new(loginInput)
	err := decoder.Decode(userInput) 
	if err != nil || !userInput.validate() {
		httpErrorResponce(w, "incorrect input body", http.StatusBadRequest)
		return
	}

	userID, err := h.services.Auth.GetUserID(userInput.Username, userInput.Password)
	if err != nil {
		httpErrorResponce(w, err.Error(), http.StatusBadRequest)  // TODO it can be repos error too
		return
	}

	token, err := h.services.Auth.GenerateToken(userID)
	if err != nil {
		httpErrorResponce(w, err.Error(), http.StatusInternalServerError)  // TODO change error message
		return
	}

	// TODO maybe make wrapper for responses
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	response := responseLogin{JWT: token}
	encoder := json.NewEncoder(w)
	if err := encoder.Encode(&response); err != nil {
		httpErrorResponce(w, err.Error(), http.StatusInternalServerError) // TODO change error message
		return
	}
}


// LOGOUT
func (h *Handler) logout(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("You've been successfully logout"))
}
