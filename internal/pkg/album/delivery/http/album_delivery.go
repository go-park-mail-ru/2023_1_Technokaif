package http

import (
	"encoding/json"
	"fmt"
	"net/http"

	commonHttp "github.com/go-park-mail-ru/2023_1_Technokaif/internal/common/http"
	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/models"
	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/album"
	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/artist"
	"github.com/go-park-mail-ru/2023_1_Technokaif/pkg/logger"
)

type Handler struct {
	albumServices  album.Usecase
	artistServices artist.Usecase
	logger         logger.Logger
}

func NewHandler(alu album.Usecase, aru artist.Usecase, l logger.Logger) *Handler {
	return &Handler{
		albumServices:  alu,
		artistServices: aru,
		logger:         l,
	}
}

// swaggermock
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	// ...
}

// swaggermock
func (h *Handler) Read(w http.ResponseWriter, r *http.Request) {
	// ...
}

// swaggermock
func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	// ...
}

// swaggermock
func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	// ...
}

// swaggermock
func (h *Handler) Tracks(w http.ResponseWriter, r *http.Request) {
	// ...
}

//	@Summary		Album Feed
//	@Tags			album feed
//	@Description	Feed albums for user
//	@Accept			json
//	@Produce		json
//	@Success		200		{object}	signUpResponse	"Show feed"
//	@Failure		500		{object}	errorResponse	"Server error"
//	@Router			/api/album/feed [get]
func (h *Handler) Feed(w http.ResponseWriter, r *http.Request) {
	albums, err := h.albumServices.GetFeed()
	if err != nil {
		h.logger.Error(err.Error())
		commonHttp.ErrorResponse(w, "error while getting albums", http.StatusInternalServerError, h.logger)
		return
	}
	w.Header().Set("Content-Type", "json/application; charset=utf-8")

	encoder := json.NewEncoder(w)
	if err := encoder.Encode(&albums); err != nil {
		h.logger.Error(err.Error())
		commonHttp.ErrorResponse(w, "can't encode response into json", http.StatusInternalServerError, h.logger)
		return
	}
}

func (h *Handler) artistTransferFromQuery(artists []models.Artist) []models.ArtistTransfer {
	at := make([]models.ArtistTransfer, len(artists))
	for _, a := range artists {
		at = append(at, models.ArtistTransfer{
			ID:        a.ID,
			Name:      a.Name,
			AvatarSrc: a.AvatarSrc,
		})
	}

	return at
}

func (h *Handler) albumTransferFromQuery(albums []models.Album) ([]models.AlbumTransfer, error) {
	at := make([]models.AlbumTransfer, 0, len(albums))
	for _, a := range albums {
		artists, err := h.artistServices.GetByAlbum(a.ID)
		if err != nil {
			return nil, fmt.Errorf("(delivery) can't get albums's (id #%d) artists: %w", a.ID, err)
		}

		at = append(at, models.AlbumTransfer{
			ID:          a.ID,
			Name:        a.Name,
			Artists:     h.artistTransferFromQuery(artists),
			Description: a.Description,
			CoverSrc:    a.CoverSrc,
		})
	}

	return at, nil
}
