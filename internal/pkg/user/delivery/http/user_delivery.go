package user_delivery

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi"

	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/logger"
	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/models"
	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/user"
)

type UserHandler struct {
	services user.UserUsecase
	logger   logger.Logger
}

func NewUserHandler(u user.UserUsecase, l logger.Logger) *UserHandler {
	return &UserHandler{
		services: u,
		logger:   l,
	}
}

func(h *UserHandler) GetById(w http.ResponseWriter, r *http.Request) {
	userParam := chi.URLParam(r, "userID")  // TODO const url param
	userID, err := strconv.Atoi(userParam)
	if err != nil {
		
	}

	h.services.GetByID(userID)


}

