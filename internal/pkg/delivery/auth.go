package delivery

import (
	"encoding/json"
	"net/http"
	"strconv"

	valid "github.com/asaskevich/govalidator"

	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/models"
)

type signUpResponse struct {
	ID int `json:"id"`
}

type loginResponse struct {
	JWT string `json:"jwt"`
}

type logoutResponse struct {
	Status string `json:"status"`
}

//	@Summary		Sign Up
//	@Tags			auth
//	@Description	create account
//	@Accept			json
//	@Produce		json
//	@Param			user	body		models.User		true	"user info"
//	@Success		200		{object}	signUpResponse	"User created"
//	@Failure		400		{object}	errorResponse	"Incorrect input"
//	@Failure		500		{object}	errorResponse	"Server DB error"
//	@Router			/api/auth/signup [post]
func (h *Handler) signUp(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var user models.User

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&user)
	if err != nil || !user.DeliveryValidate() {
		if err != nil {
			h.logger.Error(err.Error())
		} else {
			h.logger.Error("user validation failed")
		}
		h.errorResponce(w, "incorrect input body", http.StatusBadRequest)
		return
	}

	id, err := h.services.Auth.CreateUser(user)
	if err != nil {
		h.logger.Error(err.Error())
		h.errorResponce(w, err.Error(), http.StatusBadRequest)
		return
	}

	h.logger.Info("user created, Id : " + strconv.Itoa(id))

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

//	@Summary		Sign In
//	@Tags			auth
//	@Description	login account
//	@Accept			json
//	@Produce		json
//	@Param			userInput	body		loginInput		true	"username and password"
//	@Success		200			{object}	loginResponse	"User created"
//	@Failure		400			{object}	errorResponse	"Incorrect input"
//	@Failure		500			{object}	errorResponse	"Server DB error"
//	@Router			/api/auth/login [post]
func (h *Handler) login(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var userInput loginInput

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&userInput)
	if err != nil || !userInput.validate() {
		if err != nil {
			h.logger.Error(err.Error())
		} else {
			h.logger.Error("user validation failed")
		}
		h.errorResponce(w, "incorrect input body", http.StatusBadRequest)
		return
	}

	user, err := h.services.Auth.GetUserByCreds(userInput.Username, userInput.Password)
	if err != nil {
		h.logger.Error(err.Error())
		h.errorResponce(w, err.Error(), http.StatusBadRequest)
		return
	}

	token, err := h.services.Auth.GenerateAccessToken(user.ID, user.Version)
	if err != nil {
		h.logger.Error(err.Error())
		h.errorResponce(w, "can't generate access token", http.StatusInternalServerError)
		return
	}

	h.logger.Info("login returned token : " + token)

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

//	@Summary		Log Out
//	@Tags			auth
//	@Description	logout account
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	logoutResponse	"User loged out"
//	@Failure		400	{object}	errorResponse	"Logout fail"
//	@Failure		500	{object}	errorResponse	"Server DB error"
//	@Router			/api/auth/logout [get]
func (h *Handler) logout(w http.ResponseWriter, r *http.Request) {
	user, err := h.GetUserFromAuthorization(r)
	if err != nil {
		h.errorResponce(w, "invalid token", http.StatusBadRequest)
		return
	}
	h.logger.Info("UserID for logout : " + strconv.Itoa(int(user.ID)))

	if err = h.services.ChangeUserVersion(user.ID); err != nil {
		h.logger.Error(err.Error())
		h.errorResponce(w, "failed to log out", http.StatusBadRequest)
		return
	}

	// TODO maybe make wrapper for responses
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	response := logoutResponse{Status: "ok"}
	encoder := json.NewEncoder(w)
	if err := encoder.Encode(&response); err != nil {
		h.logger.Error(err.Error())
		h.errorResponce(w, "can't encode responce into json", http.StatusInternalServerError)
		return
	}
}
