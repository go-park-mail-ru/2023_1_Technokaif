package http

import (
	"net/http"

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

func (h *Handler) GetByID(w http.ResponseWriter, r *http.Request) {
	// userParam := chi.URLParam(r, "userID") // TODO const url param
	// userID, err := strconv.Atoi(userParam)
	// ...
}

func (h *Handler) GetBriefByID(w http.ResponseWriter, r *http.Request) {
	// userParam := chi.URLParam(r, "userID") // TODO const url param
	// userID, err := strconv.Atoi(userParam)
	// ...
}
