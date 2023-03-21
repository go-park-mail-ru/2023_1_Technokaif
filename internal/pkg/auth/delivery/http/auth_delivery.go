package http

import (
	"encoding/json"
	"net/http"

	valid "github.com/asaskevich/govalidator"
	"github.com/pkg/errors"

	commonHttp "github.com/go-park-mail-ru/2023_1_Technokaif/internal/common/http"
	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/models"
	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/auth"
	"github.com/go-park-mail-ru/2023_1_Technokaif/pkg/logger"
)

type Handler struct {
	services auth.Usecase
	logger   logger.Logger
}

func NewHandler(au auth.Usecase, l logger.Logger) *Handler {
	return &Handler{
		services: au,
		logger:   l,
	}
}

type signUpResponse struct {
	ID uint32 `json:"id"`
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
func (h *Handler) SignUp(w http.ResponseWriter, r *http.Request) {
	var user models.User

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&user); err != nil {
		h.logger.Error(err.Error())
		commonHttp.ErrorResponse(w, "incorrect input body", http.StatusBadRequest, h.logger)
		return
	}

	if err := user.DeliveryValidate(); err != nil {
		h.logger.Errorf("user validation failed: %s", err.Error())
		commonHttp.ErrorResponse(w, "incorrect input body", http.StatusBadRequest, h.logger)
		return
	}

	id, err := h.services.SignUpUser(user)
	var errUserAlreadyExists *models.UserAlreadyExistsError
	if errors.As(err, &errUserAlreadyExists) {
		h.logger.Error(err.Error())
		commonHttp.ErrorResponse(w, "user already exists", http.StatusBadRequest, h.logger)
		return
	} else if err != nil {
		h.logger.Error(err.Error())
		commonHttp.ErrorResponse(w, "server error", http.StatusInternalServerError, h.logger)
		return
	}

	h.logger.Infof("user created with id: %d", id)

	// TODO maybe make wrapper for responses
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	response := signUpResponse{ID: id}
	encoder := json.NewEncoder(w)
	if err := encoder.Encode(&response); err != nil {
		h.logger.Error(err.Error())
		commonHttp.ErrorResponse(w, "can't encode response into json", http.StatusInternalServerError, h.logger)
		return
	}
}

type loginInput struct {
	Username string `json:"username" valid:"required"`
	Password string `json:"password" valid:"required"`
}

func (li *loginInput) validate() error {
	_, err := valid.ValidateStruct(li)
	return err
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
func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var userInput loginInput

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&userInput); err != nil {
		h.logger.Errorf("incorrect json format: %s", err.Error())
		commonHttp.ErrorResponse(w, "incorrect input body", http.StatusBadRequest, h.logger)
		return
	}

	if err := userInput.validate(); err != nil {
		h.logger.Errorf("user validation failed: %s", err.Error())
		commonHttp.ErrorResponse(w, "incorrect input body", http.StatusBadRequest, h.logger)
		return
	}

	token, err := h.services.LoginUser(userInput.Username, userInput.Password)
	if err != nil {
		h.logger.Error(err.Error())
		commonHttp.ErrorResponse(w, "can't login user", http.StatusBadRequest, h.logger)
		return
	}

	h.logger.Infof("login with token: %s", token)

	// TODO maybe make wrapper for responses
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	response := loginResponse{JWT: token}
	encoder := json.NewEncoder(w)
	if err := encoder.Encode(&response); err != nil {
		h.logger.Error(err.Error())
		commonHttp.ErrorResponse(w, "can't encode response into json", http.StatusInternalServerError, h.logger)
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
func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	user, err := commonHttp.GetUserFromRequest(r)
	if err != nil {
		h.logger.Errorf("failed to logout: %s", err.Error())
		commonHttp.ErrorResponse(w, "invalid token", http.StatusBadRequest, h.logger)
		return
	}
	h.logger.Infof("userID for logout: %d", user.ID)

	if err = h.services.IncreaseUserVersion(user.ID); err != nil { // userVersion UP
		h.logger.Errorf("failed to logout: %s", err.Error())
		commonHttp.ErrorResponse(w, "failed to log out", http.StatusBadRequest, h.logger)
		return
	}

	// TODO maybe make wrapper for responses
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	response := logoutResponse{Status: "ok"}
	encoder := json.NewEncoder(w)
	if err := encoder.Encode(&response); err != nil {
		h.logger.Error(err.Error())
		commonHttp.ErrorResponse(w, "can't encode response into json", http.StatusInternalServerError, h.logger)
		return
	}
}
