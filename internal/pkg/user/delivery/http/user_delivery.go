package http

import (
	"errors"
	"net/http"

	commonHttp "github.com/go-park-mail-ru/2023_1_Technokaif/internal/common/http"
	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/models"
	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/user"
	"github.com/go-park-mail-ru/2023_1_Technokaif/pkg/logger"
)

type Handler struct {
	userServices user.Usecase
	logger       logger.Logger
}

func NewHandler(uu user.Usecase, l logger.Logger) *Handler {
	return &Handler{
		userServices: uu,
		logger:       l,
	}
}

func userTransferFromUser(user models.User) models.UserTransfer {
	return models.UserTransfer{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Sex:       user.Sex,
		BirhDate:  user.BirthDate,
		AvatarSrc: user.AvatarSrc,
	}
}

func (h *Handler) Read(w http.ResponseWriter, r *http.Request) {
	userID, err := commonHttp.GetUserIDFromRequest(r)
	if err != nil {
		h.logger.Infof("get user by id: %v", err.Error())
		commonHttp.ErrorResponse(w, "invalid url parameter", http.StatusBadRequest, h.logger)
		return
	}

	user, err := h.userServices.GetByID(userID)
	var errNoSuchUser *models.NoSuchUserError
	if errors.As(err, &errNoSuchUser) {
		h.logger.Info(err.Error())
		commonHttp.ErrorResponse(w, "no such user", http.StatusBadRequest, h.logger)
		return
	}
	if err != nil {
		h.logger.Info(err.Error())
		commonHttp.ErrorResponse(w, "error while getting user", http.StatusInternalServerError, h.logger)
		return
	}

	ut := userTransferFromUser(*user)

	commonHttp.SuccessResponse(w, ut, h.logger)
}

const maxAvatarMemory = 2 * (1 << 20)

type userUploadAvatarResponse struct {
	Status string `json:"status"`
}

// swaggermock
func (h *Handler) UploadAvatar(w http.ResponseWriter, r *http.Request) {
	user, err := commonHttp.GetUserFromRequest(r)
	if err != nil {
		h.logger.Infof("unathorized user: %v", err)
		commonHttp.ErrorResponse(w, "invalid token", http.StatusUnauthorized, h.logger)
		return
	}

	urlID, err := commonHttp.GetUserIDFromRequest(r)
	if err != nil {
		h.logger.Infof("can't get user ID from URL: %v", err)
		commonHttp.ErrorResponse(w, "can't upload avatar", http.StatusInternalServerError, h.logger)
		return
	}

	if urlID != user.ID {
		h.logger.Infof("forbidden avatar upload: %v", err)
		commonHttp.ErrorResponse(w, "invalid user", http.StatusForbidden, h.logger)
		return
	}

	if err := r.ParseMultipartForm(maxAvatarMemory); err != nil {
		h.logger.Info(err.Error())
		commonHttp.ErrorResponse(w, "invalid form data", http.StatusBadRequest, h.logger)
		return
	}

	form := r.MultipartForm

	if err := h.userServices.UploadAvatar(user, form.File["avatar"][0]); err != nil {
		h.logger.Error(err.Error())
		commonHttp.ErrorResponse(w, "can't upload avatar", http.StatusInternalServerError, h.logger)
		return
	}

	uuar := userUploadAvatarResponse{Status: "ok"}

	commonHttp.SuccessResponse(w, uuar, h.logger)
}
