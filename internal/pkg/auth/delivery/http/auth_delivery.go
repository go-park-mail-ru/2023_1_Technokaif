package http

import (
	"encoding/json"
	"net/http"
	"errors"

	commonHttp "github.com/go-park-mail-ru/2023_1_Technokaif/internal/common/http"
	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/models"
	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/auth"
	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/token"
	"github.com/go-park-mail-ru/2023_1_Technokaif/pkg/logger"
)

type Handler struct {
	authServices 	auth.Usecase
	tokenServices 	token.Usecase
	logger   		logger.Logger
}

func NewHandler(au auth.Usecase, tu token.Usecase, l logger.Logger) *Handler {
	return &Handler{
		authServices: 	au,
		tokenServices: 	tu,
		logger:   		l,
	}
}

// @Summary		Sign Up
// @Tags		Auth
// @Description	Create account
// @Accept		json
// @Produce		json
// @Param		user	body		models.User		true	"user info"
// @Success		200		{object}	signUpResponse	"User created"
// @Failure		400		{object}	http.Error	"Incorrect input"
// @Failure		500		{object}	http.Error	"Server error"
// @Router		/api/auth/signup [post]
func (h *Handler) SignUp(w http.ResponseWriter, r *http.Request) {
	var user models.User

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&user); err != nil {
		commonHttp.ErrorResponseWithErrLogging(w, "incorrect input body", http.StatusBadRequest, h.logger, err)
		return
	}

	if err := user.DeliveryValidate(); err != nil {
		commonHttp.ErrorResponseWithErrLogging(w, "incorrect input body", http.StatusBadRequest, h.logger, err)
		return
	}

	id, err := h.authServices.SignUpUser(user)
	var errUserAlreadyExists *models.UserAlreadyExistsError
	if errors.As(err, &errUserAlreadyExists) {
		commonHttp.ErrorResponseWithErrLogging(w, "user already exists", http.StatusBadRequest, h.logger, err)
		return
	} else if err != nil {
		commonHttp.ErrorResponseWithErrLogging(w, "server error", http.StatusInternalServerError, h.logger, err)
		return
	}

	h.logger.Infof("user created with id: %d", id)

	sur := signUpResponse{ID: id}

	commonHttp.SuccessResponse(w, sur, h.logger)
}

// @Summary		Sign In
// @Tags		Auth
// @Description	Login account
// @Accept		json
// @Produce		json
// @Param		userInput	body		loginInput		true	"username and password"
// @Success		200			{object}	loginResponse	"User created"
// @Failure		400			{object}	http.Error	"Incorrect input"
// @Failure		500			{object}	http.Error	"Server error"
// @Router		/api/auth/login [post]
func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var userInput loginInput

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&userInput); err != nil {
		h.logger.Infof("incorrect json format: %s", err.Error())
		commonHttp.ErrorResponse(w, "incorrect input body", http.StatusBadRequest, h.logger)
		return
	}

	if err := userInput.validate(); err != nil {
		h.logger.Infof("user validation failed: %s", err.Error())
		commonHttp.ErrorResponse(w, "incorrect input body", http.StatusBadRequest, h.logger)
		return
	}

	user, err := h.authServices.GetUserByCreds(userInput.Username, userInput.Password)
	if err != nil {
		var errNoSuchUser *models.NoSuchUserError
		if errors.As(err, &errNoSuchUser) {
			commonHttp.ErrorResponseWithErrLogging(w, "no such user", http.StatusBadRequest, h.logger, err)
			return
		}

		var errIncorrectPassword *models.IncorrectPasswordError
		if errors.As(err, &errIncorrectPassword) {
			commonHttp.ErrorResponseWithErrLogging(w, "incorrect password", http.StatusBadRequest, h.logger, err)
			return
		}

		commonHttp.ErrorResponseWithErrLogging(w, "server failed to login user", http.StatusInternalServerError, h.logger, err)
		return
	}

	token, err := h.tokenServices.GenerateAccessToken(user.ID, user.Version)
	if err != nil {
		commonHttp.ErrorResponseWithErrLogging(w, "server failed to login user", http.StatusInternalServerError, h.logger, err)
		return
	}

	h.logger.Infof("login with token: %s", token)

	lr := loginResponse{JWT: token}

	commonHttp.SuccessResponse(w, lr, h.logger)
}

// @Summary		Log Out
// @Tags		Auth
// @Description	Logout account
// @Accept		json
// @Produce		json
// @Success		200	{object}	logoutResponse	"User loged out"
// @Failure		400	{object}	http.Error	"Logout fail"
// @Failure		500	{object}	http.Error	"Server error"
// @Router		/api/auth/logout [get]
func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	user, err := commonHttp.GetUserFromRequest(r)
	if err != nil {
		h.logger.Infof("failed to logout: %s", err.Error())
		commonHttp.ErrorResponse(w, "invalid token", http.StatusUnauthorized, h.logger)
		return
	}
	h.logger.Infof("userID for logout: %d", user.ID)

	if err = h.authServices.IncreaseUserVersion(user.ID); err != nil { // userVersion UP
		h.logger.Errorf("failed to logout: %s", err.Error())
		commonHttp.ErrorResponse(w, "failed to log out", http.StatusInternalServerError, h.logger)
		return
	}

	lr := logoutResponse{Status: "ok"}

	commonHttp.SuccessResponse(w, lr, h.logger)
}

func (h *Handler) ChangePassword(w http.ResponseWriter, r *http.Request) {
	user, err := commonHttp.GetUserFromRequest(r)
	if err != nil {
		h.logger.Infof("failed to change password: %s", err.Error())
		commonHttp.ErrorResponse(w, "invalid token", http.StatusUnauthorized, h.logger)
		return
	}

	var passwordsInput changePassInput
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&passwordsInput); err != nil {
		h.logger.Infof("incorrect json format: %s", err.Error())
		commonHttp.ErrorResponse(w, "incorrect input body", http.StatusBadRequest, h.logger)
		return
	}

	if _, err := h.authServices.GetUserByCreds(user.Username, passwordsInput.OldPassword); err != nil {
		var errIncorrectPassword *models.IncorrectPasswordError
		if errors.As(err, &errIncorrectPassword) {
			commonHttp.ErrorResponseWithErrLogging(w, "incorrect password", http.StatusBadRequest, h.logger, err)
			return
		}

		commonHttp.ErrorResponseWithErrLogging(w, "server failed to get user", http.StatusInternalServerError, h.logger, err)
		return
	}

	if err := h.authServices.ChangePassword(user.ID, passwordsInput.NewPassword); err != nil {
		commonHttp.ErrorResponseWithErrLogging(w, "server failed to change password", http.StatusInternalServerError, h.logger, err)
		return
	}

	resp := changePassResponse{Status: "ok"}
	commonHttp.SuccessResponse(w, resp, h.logger)
}
