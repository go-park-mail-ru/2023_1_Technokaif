package user_delivery

import (
	"net/http"

	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/user"
	"github.com/go-park-mail-ru/2023_1_Technokaif/pkg/logger"
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

func (h *UserHandler) GetById(w http.ResponseWriter, r *http.Request) {
	// userParam := chi.URLParam(r, "userID") // TODO const url param
	// userID, err := strconv.Atoi(userParam)
	// ...
}
