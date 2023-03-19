package auth_delivery

import (
	"encoding/json"
	"net/http"

	valid "github.com/asaskevich/govalidator"
	"github.com/pkg/errors"

	commonHttp "github.com/go-park-mail-ru/2023_1_Technokaif/internal/common/http"
	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/logger"
	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/models"
	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/auth"
)

type AuthHandler struct {
	services auth.AuthUsecase
	logger   logger.Logger
}

func NewAuthHandler(u auth.AuthUsecase, l logger.Logger) *AuthHandler {
	return &AuthHandler{
		services: u,
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
func (ah *AuthHandler) SignUp(w http.ResponseWriter, r *http.Request) {
	var user models.User

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&user); err != nil {
		ah.logger.Error(err.Error())
		commonHttp.ErrorResponse(w, "incorrect input body", http.StatusBadRequest)
		return
	}

	if err := user.DeliveryValidate(); err != nil {
		ah.logger.Errorf("user validation failed: %s", err.Error())
		commonHttp.ErrorResponse(w, "incorrect input body", http.StatusBadRequest)
		return
	}

	id, err := ah.services.CreateUser(user)
	var errUserAlreadyExists *models.UserAlreadyExistsError
	if errors.As(err, &errUserAlreadyExists) {
		ah.logger.Error(err.Error())
		commonHttp.ErrorResponse(w, "user already exists", http.StatusBadRequest)
		return
	} else if err != nil {
		ah.logger.Error(err.Error())
		commonHttp.ErrorResponse(w, "server error", http.StatusInternalServerError)
		return
	}

	ah.logger.Infof("user created with id: %d", id)

	// TODO maybe make wrapper for responses
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	response := signUpResponse{ID: id}
	encoder := json.NewEncoder(w)
	if err := encoder.Encode(&response); err != nil {
		ah.logger.Error(err.Error())
		commonHttp.ErrorResponse(w, "can't encode response into json", http.StatusInternalServerError)
		return
	}
}

type loginInput struct {
	Username string `json:"username" valid:"required"`
	Password string `json:"password" valid:"required"`
}

func (i *loginInput) validate() error {
	_, err := valid.ValidateStruct(i)
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
func (ah *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var userInput loginInput

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&userInput); err != nil {
		ah.logger.Errorf("incorrect json format: %s", err.Error())
		commonHttp.ErrorResponse(w, "incorrect input body", http.StatusBadRequest)
		return
	}

	if err := userInput.validate(); err != nil {
		ah.logger.Errorf("user validation failed: %s", err.Error())
		commonHttp.ErrorResponse(w, "incorrect input body", http.StatusBadRequest)
		return
	}

	token, err := ah.services.LoginUser(userInput.Username, userInput.Password)
	if err != nil {
		ah.logger.Error(err.Error())
		commonHttp.ErrorResponse(w, "can't login user", http.StatusBadRequest)
		return
	}

	ah.logger.Infof("login with token: %s", token)

	// TODO maybe make wrapper for responses
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	response := loginResponse{JWT: token}
	encoder := json.NewEncoder(w)
	if err := encoder.Encode(&response); err != nil {
		ah.logger.Error(err.Error())
		commonHttp.ErrorResponse(w, "can't encode response into json", http.StatusInternalServerError)
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
func (ah *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	user, err := commonHttp.GetUserFromAuthorization(r)
	if err != nil {
		ah.logger.Errorf("failed to logout: %s", err.Error())
		commonHttp.ErrorResponse(w, "invalid token", http.StatusBadRequest)
		return
	}
	ah.logger.Infof("userID for logout: %d", user.ID)

	if err = ah.services.IncreaseUserVersion(user.ID); err != nil { // userVersion UP
		ah.logger.Errorf("failed to logout: %s", err.Error())
		commonHttp.ErrorResponse(w, "failed to log out", http.StatusBadRequest)
		return
	}

	// TODO maybe make wrapper for responses
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	response := logoutResponse{Status: "ok"}
	encoder := json.NewEncoder(w)
	if err := encoder.Encode(&response); err != nil {
		ah.logger.Error(err.Error())
		commonHttp.ErrorResponse(w, "can't encode response into json", http.StatusInternalServerError)
		return
	}
}
