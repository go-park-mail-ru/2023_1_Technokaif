package delivery

import (
	"encoding/json"
	"net/http"

	valid "github.com/asaskevich/govalidator"

	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/models"
)

type signUpResponse struct {
	ID int `json:"id"`
}

type loginResponse struct {
	JWT string `json:"jwt"`
}

func (h *Handler) signUp(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var user models.User

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&user)
	if err != nil || !user.DeliveryValidate() {
		if err != nil { h.logger.Error(err.Error()) } else { h.logger.Error("user validation failed") } // logger
		h.errorResponce(w, "incorrect input body", http.StatusBadRequest)
		return
	}

	id, err := h.services.Auth.CreateUser(user)
	if err != nil {
		h.logger.Error(err.Error())
		h.errorResponce(w, err.Error(), http.StatusBadRequest)
		return
	}

	// TODO maybe make wrapper for responses
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	response := signUpResponse{ID: id}
	encoder := json.NewEncoder(w)
	if err := encoder.Encode(&response); err != nil {
		h.logger.Error(err.Error())
		h.errorResponce(w, err.Error(), http.StatusInternalServerError) // TODO change error message
		return
	}
}

type loginInput struct {
	Username string `json:"username" valid:"required"`
	Password string `json:"password" valid:"required"`
}

func (i *loginInput) validate() bool {
	isValid, err := valid.ValidateStruct(i)
	if err != nil || !isValid {
		return false
	}

	return true
}

func (h *Handler) login(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var userInput loginInput

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&userInput)
	if err != nil || !userInput.validate() {
		if err != nil { h.logger.Error(err.Error()) } else { h.logger.Error("user validation failed") } // logger
		h.errorResponce(w, "incorrect input body", http.StatusBadRequest)
		return
	}

	userID, err := h.services.Auth.GetUserID(userInput.Username, userInput.Password)
	if err != nil {
		h.logger.Error(err.Error())
		h.errorResponce(w, err.Error(), http.StatusBadRequest)
		return
	}

	token, err := h.services.Auth.GenerateAccessToken(userID)
	if err != nil {
		h.logger.Error(err.Error())
		h.errorResponce(w, "can't generate access token", http.StatusInternalServerError)
		return
	}

	// TODO maybe make wrapper for responses
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	response := loginResponse{JWT: token}
	encoder := json.NewEncoder(w)
	if err := encoder.Encode(&response); err != nil {
		h.logger.Error(err.Error())
		h.errorResponce(w, "can't encode responce into json", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) logout(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("You've been successfully logout"))
}
