package http

import (
	"encoding/json"
	"errors"
	"net/http"

	commonHTTP "github.com/go-park-mail-ru/2023_1_Technokaif/internal/common/http"
	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/models"
	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/auth"
	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/token"
	"github.com/go-park-mail-ru/2023_1_Technokaif/pkg/logger"
)

type Handler struct {
	authServices  auth.Usecase
	tokenServices token.Usecase
	logger        logger.Logger
}

func NewHandler(au auth.Usecase, tu token.Usecase, l logger.Logger) *Handler {
	return &Handler{
		authServices:  au,
		tokenServices: tu,

		logger: l,
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
		commonHTTP.ErrorResponseWithErrLogging(w, r,
			commonHTTP.IncorrectRequestBody, http.StatusBadRequest, h.logger, err)
		return
	}

	if err := userAuthDeliveryValidate(&user); err != nil {
		commonHTTP.ErrorResponseWithErrLogging(w, r,
			commonHTTP.IncorrectRequestBody, http.StatusBadRequest, h.logger, err)
		return
	}

	id, err := h.authServices.SignUpUser(r.Context(), user)
	if err != nil {
		var errUserAlreadyExists *models.UserAlreadyExistsError
		if errors.As(err, &errUserAlreadyExists) {
			commonHTTP.ErrorResponseWithErrLogging(w, r,
				userAlreadyExists, http.StatusBadRequest, h.logger, err)
			return
		}

		commonHTTP.ErrorResponseWithErrLogging(w, r,
			userSignUpServerError, http.StatusInternalServerError, h.logger, err)
		return
	}

	h.logger.Infof("user created with id: %d", id)

	sur := signUpResponse{ID: id}

	commonHTTP.SuccessResponse(w, sur, h.logger)
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
		commonHTTP.ErrorResponseWithErrLogging(w, r,
			commonHTTP.IncorrectRequestBody, http.StatusBadRequest, h.logger, err)
		return
	}

	if err := userInput.validateAndEscape(); err != nil {
		commonHTTP.ErrorResponseWithErrLogging(w, r,
			commonHTTP.IncorrectRequestBody, http.StatusBadRequest, h.logger, err)
		return
	}

	user, err := h.authServices.GetUserByCreds(r.Context(), userInput.Username, userInput.Password)
	if err != nil {
		var errNoSuchUser *models.NoSuchUserError
		if errors.As(err, &errNoSuchUser) {
			commonHTTP.ErrorResponseWithErrLogging(w, r,
				userNotFound, http.StatusBadRequest, h.logger, err)
			return
		}

		var errIncorrectPassword *models.IncorrectPasswordError
		if errors.As(err, &errIncorrectPassword) {
			commonHTTP.ErrorResponseWithErrLogging(w, r,
				passwordMismatch, http.StatusBadRequest, h.logger, err)
			return
		}

		commonHTTP.ErrorResponseWithErrLogging(w, r,
			userLoginServerError, http.StatusInternalServerError, h.logger, err)
		return
	}

	token, err := h.tokenServices.GenerateAccessToken(user.ID, user.Version)
	if err != nil {
		commonHTTP.ErrorResponseWithErrLogging(w, r,
			userLoginServerError, http.StatusInternalServerError, h.logger, err)
		return
	}

	h.logger.Infof("login with token: %s", token)

	lr := loginResponse{UserID: user.ID}

	commonHTTP.SetAccessTokenCookie(w, token)
	commonHTTP.SuccessResponse(w, lr, h.logger)
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
	user, err := commonHTTP.GetUserFromRequest(r)
	if err != nil {
		commonHTTP.ErrorResponseWithErrLogging(w, r,
			invalidToken, http.StatusUnauthorized, h.logger, err)
		return
	}

	if err = h.authServices.IncreaseUserVersion(r.Context(), user.ID); err != nil {
		commonHTTP.ErrorResponseWithErrLogging(w, r,
			userLogoutServerError, http.StatusInternalServerError, h.logger, err)
		return
	}

	lr := logoutResponse{Status: userLogedOutSuccessfully}

	commonHTTP.SetAccessTokenCookie(w, "")
	commonHTTP.SuccessResponse(w, lr, h.logger)
}

// swaggermock
func (h *Handler) ChangePassword(w http.ResponseWriter, r *http.Request) {
	user, err := commonHTTP.GetUserFromRequest(r)
	if err != nil {
		commonHTTP.ErrorResponseWithErrLogging(w, r,
			invalidToken, http.StatusUnauthorized, h.logger, err)
		return
	}

	var passwordsInput changePassInput
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&passwordsInput); err != nil {
		commonHTTP.ErrorResponseWithErrLogging(w, r,
			commonHTTP.IncorrectRequestBody, http.StatusBadRequest, h.logger, err)
		return
	}

	if err := passwordsInput.validate(); err != nil {
		commonHTTP.ErrorResponseWithErrLogging(w, r,
			commonHTTP.IncorrectRequestBody, http.StatusBadRequest, h.logger, err)
		return
	}

	if _, err := h.authServices.GetUserByCreds(r.Context(), user.Username, passwordsInput.OldPassword); err != nil {
		var errIncorrectPassword *models.IncorrectPasswordError
		if errors.As(err, &errIncorrectPassword) {
			commonHTTP.ErrorResponseWithErrLogging(w, r,
				passwordMismatch, http.StatusBadRequest, h.logger, err)
			return
		}

		commonHTTP.ErrorResponseWithErrLogging(w, r,
			userGetServerError, http.StatusInternalServerError, h.logger, err)
		return
	}

	if err := h.authServices.ChangePassword(r.Context(), user.ID, passwordsInput.NewPassword); err != nil {
		commonHTTP.ErrorResponseWithErrLogging(w, r,
			userChangePasswordError, http.StatusInternalServerError, h.logger, err)
		return
	}

	if err = h.authServices.IncreaseUserVersion(r.Context(), user.ID); err != nil {
		commonHTTP.ErrorResponseWithErrLogging(w, r,
			userChangePasswordError, http.StatusInternalServerError, h.logger, err)
		return
	}

	token, err := h.tokenServices.GenerateAccessToken(user.ID, user.Version+1)
	if err != nil {
		commonHTTP.ErrorResponseWithErrLogging(w, r,
			tokenGenerateServerError, http.StatusInternalServerError, h.logger, err)
		return
	}

	resp := changePassResponse{Status: userChangedPasswordSuccessfully}

	commonHTTP.SetAccessTokenCookie(w, token)
	commonHTTP.SuccessResponse(w, resp, h.logger)
}

// swaggermock
func (h *Handler) IsAuthenticated(w http.ResponseWriter, r *http.Request) {
	iar := isAuthenticatedResponse{}

	if _, err := commonHTTP.GetUserFromRequest(r); err == nil {
		iar.Authenticated = true
	}

	commonHTTP.SuccessResponse(w, iar, h.logger)
}

// swaggermock
func (h *Handler) Auth(w http.ResponseWriter, r *http.Request) {
	if _, err := commonHTTP.GetUserFromRequest(r); err != nil {
		commonHTTP.ErrorResponseWithErrLogging(w, r, userForbidden, http.StatusForbidden, h.logger, err)
		return
	}

	commonHTTP.SuccessResponse(w, isAuthenticatedResponse{Authenticated: true}, h.logger)
}
