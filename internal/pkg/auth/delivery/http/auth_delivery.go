package http

import (
	"encoding/json"
	"net/http"

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

// @Summary		Sign Up
// @Tags		Auth
// @Description	Create account
// @Accept		json
// @Produce		json
// @Param		user	body		models.User	true	"User info"
// @Success		200		{object}	signUpResponse		"User created"
// @Failure		400		{object}	http.Error			"Incorrect input"
// @Failure		500		{object}	http.Error			"Server error"
// @Router		/api/auth/signup [post]
func (h *Handler) SignUp(w http.ResponseWriter, r *http.Request) {
	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		commonHttp.ErrorResponseWithErrLogging(w, "incorrect input body", http.StatusBadRequest, h.logger, err)
		return
	}

	if err := user.DeliveryValidate(); err != nil {
		commonHttp.ErrorResponseWithErrLogging(w, "incorrect input body", http.StatusBadRequest, h.logger, err)
		return
	}

	id, err := h.services.SignUpUser(user)
	if err != nil {
		var errUserAlreadyExists *models.UserAlreadyExistsError
		if errors.As(err, &errUserAlreadyExists) {
			commonHttp.ErrorResponseWithErrLogging(w, "user already exists", http.StatusBadRequest, h.logger, err)
			return
		}

		commonHttp.ErrorResponseWithErrLogging(w, "server failed to sign up user", http.StatusInternalServerError, h.logger, err)
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
	if err := json.NewDecoder(r.Body).Decode(&userInput); err != nil {
		h.logger.Infof("incorrect json format: %s", err.Error())
		commonHttp.ErrorResponse(w, "incorrect input body", http.StatusBadRequest, h.logger)
		return
	}

	if err := userInput.validate(); err != nil {
		h.logger.Infof("user validation failed: %s", err.Error())
		commonHttp.ErrorResponse(w, "incorrect input body", http.StatusBadRequest, h.logger)
		return
	}

	token, err := h.services.LoginUser(userInput.Username, userInput.Password)
	if err != nil {
		var errNoSuchUser *models.NoSuchUserError
		if errors.As(err, &errNoSuchUser) {
			commonHttp.ErrorResponseWithErrLogging(w, "can't login user", http.StatusBadRequest, h.logger, err)
			return
		}

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
		commonHttp.ErrorResponse(w, "invalid token", http.StatusBadRequest, h.logger)
		return
	}

	if err = h.services.IncreaseUserVersion(user.ID); err != nil { // userVersion UP
		h.logger.Errorf("failed to logout: %s", err.Error())
		commonHttp.ErrorResponse(w, "failed to log out", http.StatusInternalServerError, h.logger)
		return
	}

	lr := logoutResponse{Status: "ok"}

	commonHttp.SuccessResponse(w, lr, h.logger)
}
