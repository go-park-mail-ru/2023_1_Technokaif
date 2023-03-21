package http

import (
	"fmt"
	"net/http"
	"errors"

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
	albumID, err := commonHttp.GetAlbumIDFromRequest(r)
	if err != nil {
		h.logger.Infof("get album by id : %v", err)
		commonHttp.ErrorResponse(w, "invalid url parameter", http.StatusBadRequest, h.logger)
		return
	}

	album, err := h.albumServices.GetByID(albumID)
	var errNoSuchAlbum *models.NoSuchAlbumError
	if errors.As(err, &errNoSuchAlbum) {
		h.logger.Info(err.Error())
		commonHttp.ErrorResponse(w, "no such album", http.StatusBadRequest, h.logger)
		return
	} else if err != nil {
		h.logger.Error(err.Error())
		commonHttp.ErrorResponse(w, "error while getting album", http.StatusInternalServerError, h.logger)
		return
	}

	resp, err := h.albumTransferFromEntry(*album)
	if err != nil {
		h.logger.Error(err.Error())
		commonHttp.ErrorResponse(w, "error while getting album", http.StatusInternalServerError, h.logger)
		return
	}
	
	commonHttp.SuccessResponse(w, resp, h.logger)
}

// swaggermock
func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	// ...
}

// swaggermock
func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	// ...
}

func (h *Handler) ReadByArtist(w http.ResponseWriter, r *http.Request) {
	artistID, err := commonHttp.GetArtistIDFromRequest(r)
	if err != nil {
		h.logger.Infof("get artist by id : %v", err)
		commonHttp.ErrorResponse(w, "invalid url parameter", http.StatusBadRequest, h.logger)
		return
	}

	albums, err := h.albumServices.GetByArtist(artistID)
	var errNoSuchArtist *models.NoSuchArtistError
	if errors.As(err, &errNoSuchArtist) {
		h.logger.Info(err.Error())
		commonHttp.ErrorResponse(w, "no such artist", http.StatusBadRequest, h.logger)
		return
	} else if err != nil {
		h.logger.Error(err.Error())
		commonHttp.ErrorResponse(w, "error while getting albums", http.StatusInternalServerError, h.logger)
		return
	}

	resp, err := h.albumTransferFromQuery(albums)
	if err != nil {
		h.logger.Error(err.Error())
		commonHttp.ErrorResponse(w, "error while getting albums", http.StatusInternalServerError, h.logger)
		return
	}

	commonHttp.SuccessResponse(w, resp, h.logger)	
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

	resp, err := h.albumTransferFromQuery(albums)
	if err != nil {
		h.logger.Error(err.Error())
		commonHttp.ErrorResponse(w, "error while getting albums", http.StatusInternalServerError, h.logger)
		return
	}

	commonHttp.SuccessResponse(w, resp, h.logger)
}

// Converts Artist to ArtistTransfer
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

// Converts Album to AlbumTransfer
func (h *Handler) albumTransferFromEntry(a models.Album) (models.AlbumTransfer, error) {
	artists, err := h.artistServices.GetByAlbum(a.ID)
	if err != nil {
		return models.AlbumTransfer{}, fmt.Errorf("(delivery) can't get albums's (id #%d) artists: %w", a.ID, err)
	}

	return models.AlbumTransfer{
		ID:          a.ID,
		Name:        a.Name,
		Artists:     h.artistTransferFromQuery(artists),
		Description: a.Description,
		CoverSrc:    a.CoverSrc,
	}, nil
}

func (h *Handler) albumTransferFromQuery(albums []models.Album) ([]models.AlbumTransfer, error) {
	albumTransfers := make([]models.AlbumTransfer, 0, len(albums))
	for _, a := range albums {
		albumTransfer, err := h.albumTransferFromEntry(a)
		if err != nil {
			return nil, err
		}

		albumTransfers = append(albumTransfers, albumTransfer)
	}

	return albumTransfers, nil
}
