package delivery

import (
	"net/http"

	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/search"
	"github.com/go-park-mail-ru/2023_1_Technokaif/pkg/logger"
)

type Handler struct {
	searchUsecase search.Usecase
	logger        logger.Logger
}

func NewHandler(su search.Usecase, l logger.Logger) *Handler {
	return &Handler{
		searchUsecase: su,
		logger:        l,
	}
}

// swaggermock
func (h *Handler) FindAlbums(w http.ResponseWriter, r *http.Request) {

}

// swaggermock
func (h *Handler) FindArtists(w http.ResponseWriter, r *http.Request) {

}

// swaggermock
func (h *Handler) FindTracks(w http.ResponseWriter, r *http.Request) {

}

// swaggermock
func (h *Handler) FincPlaylists(w http.ResponseWriter, r *http.Request) {

}
