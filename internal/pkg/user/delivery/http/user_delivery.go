package http

import (
	"encoding/json"
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
	} else if err != nil {
		h.logger.Info(err.Error())
		commonHttp.ErrorResponse(w, "error while getting user", http.StatusInternalServerError, h.logger)
		return
	}

	resp := models.UserTransfer{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Sex:       user.Sex,
		BirhDate:  user.BirhDate,
		AvatarSrc: user.AvatarSrc,
	}

	w.Header().Set("Content-Type", "json/application; charset=utf-8")
	encoder := json.NewEncoder(w)
	if err := encoder.Encode(&resp); err != nil {
		h.logger.Error(err.Error())
		commonHttp.ErrorResponse(w, "can't encode response into json", http.StatusInternalServerError, h.logger)
		return
	}
}

