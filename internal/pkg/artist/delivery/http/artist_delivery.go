package http

import (
	"encoding/json"
	"errors"
	"net/http"

	commonHttp "github.com/go-park-mail-ru/2023_1_Technokaif/internal/common/http"
	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/models"
	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/artist"
	"github.com/go-park-mail-ru/2023_1_Technokaif/pkg/logger"
)

type Handler struct {
	artistServices artist.Usecase
	logger         logger.Logger
}

func NewHandler(au artist.Usecase, logger logger.Logger) *Handler {
	return &Handler{
		artistServices: au,
		logger:         logger,
	}
}

// swaggermock
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var artist models.Artist

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&artist); err != nil {
		commonHttp.ErrorResponse(w, "incorrect input body", http.StatusBadRequest, h.logger)
		return
	}

	if err := h.artistServices.Create(artist); err != nil {
		h.logger.Error(err.Error())
		commonHttp.ErrorResponse(w, "can't create artist", http.StatusInternalServerError, h.logger)
	}

	// ...
}

// swaggermock
func (h *Handler) Read(w http.ResponseWriter, r *http.Request) {
	artistID, err := commonHttp.GetArtistIDFromRequest(r)
	if err != nil {
		h.logger.Infof("get artist by id: %v", err.Error())
		commonHttp.ErrorResponse(w, "invalid url parameter", http.StatusBadRequest, h.logger)
		return
	}

	artist, err := h.artistServices.GetByID(artistID)
	var errNoSuchArtist *models.NoSuchArtistError
	if errors.As(err, &errNoSuchArtist) {
		h.logger.Info(err.Error())
		commonHttp.ErrorResponse(w, "no such artist", http.StatusBadRequest, h.logger)
		return
	} else if err != nil {
		h.logger.Error(err.Error())
		commonHttp.ErrorResponse(w, "error while getting artist", http.StatusInternalServerError, h.logger)
		return
	}

	artistResponse := h.artistTransferFromEntry(*artist)

	commonHttp.SuccessResponse(w, artistResponse, h.logger)
}

// swaggermock
func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	// ...
}

// swaggermock
func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	// ...
}

// @Summary		Artist Feed
// @Tags			artist feed
// @Description	Feed albums for user
// @Accept			json
// @Produce		json
// @Success		200		{object}	signUpResponse	"Show feed"
// @Failure		500		{object}	errorResponse	"Server error"
// @Router			/api/artist/feed [get]
func (h *Handler) Feed(w http.ResponseWriter, r *http.Request) {
	artists, err := h.artistServices.GetFeed()
	if err != nil {
		h.logger.Error(err.Error())
		commonHttp.ErrorResponse(w, "error while getting albums", http.StatusInternalServerError, h.logger)
		return
	}

	artistsTransfer := h.artistTransferFromQuery(artists)

	commonHttp.SuccessResponse(w, artistsTransfer, h.logger)
}

func (h *Handler) artistTransferFromEntry(artist models.Artist) models.ArtistTransfer {
	return models.ArtistTransfer{
		ID:        artist.ID,
		Name:      artist.Name,
		AvatarSrc: artist.AvatarSrc,
	}
}

func (h *Handler) artistTransferFromQuery(artists []models.Artist) []models.ArtistTransfer {
	at := make([]models.ArtistTransfer, 0, len(artists))
	for _, a := range artists {
		at = append(at, h.artistTransferFromEntry(a))
	}

	return at
}
