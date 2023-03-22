package http

import (
	"encoding/json"
	"errors"
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

// @Summary		Create Album
// @Tags		Album
// @Description	Create new album by sent object
// @Accept      json
// @Produce		json
// @Param		album	body		albumCreateInput	true	"album info"
// @Success		200		{object}	albumCreateResponse	        "Album created"
// @Failure		400		{object}	commonHttp.Error	"Client error"
// @Failure		500		{object}	commonHttp.Error	"Server error"
// @Router		/api/albums/ [post]
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var aci albumCreateInput

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&aci); err != nil {
		h.logger.Info(err.Error())
		commonHttp.ErrorResponse(w, "incorrect input body", http.StatusBadRequest, h.logger)
		return
	}

	album := aci.ToAlbum()

	trackID, err := h.albumServices.Create(album, aci.ArtistsID)
	if err != nil {
		h.logger.Error(err.Error())
		commonHttp.ErrorResponse(w, "can't create album", http.StatusInternalServerError, h.logger)
		return
	}

	acr := albumCreateResponse{ID: trackID}

	commonHttp.SuccessResponse(w, acr, h.logger)
}

// @Summary		Get Album
// @Tags		Album
// @Description	Get album with chosen ID
// @Produce		json
// @Success		200		{object}	models.AlbumTransfer	    "Album got"
// @Failure		400		{object}	commonHttp.Error	"Client error"
// @Failure		500		{object}	commonHttp.Error	"Server error"
// @Router		/api/albums/{albumID}/ [get]
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
	}
	if err != nil {
		h.logger.Error(err.Error())
		commonHttp.ErrorResponse(w, "error while getting album", http.StatusInternalServerError, h.logger)
		return
	}

	resp, err := models.AlbumTransferFromEntry(*album, h.artistServices.GetByAlbum)
	if err != nil {
		h.logger.Error(err.Error())
		commonHttp.ErrorResponse(w, "error while getting album", http.StatusInternalServerError, h.logger)
		return
	}

	commonHttp.SuccessResponse(w, resp, h.logger)
}

// swaggermock
func (h *Handler) Change(w http.ResponseWriter, r *http.Request) {
	var aci albumChangeInput

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&aci); err != nil {
		h.logger.Info(err.Error())
		commonHttp.ErrorResponse(w, "incorrect input body", http.StatusBadRequest, h.logger)
		return
	}

	// album := aci.ToAlbum()
	// ...
}

// @Summary		Delete Album
// @Tags		Album
// @Description	Delete album with chosen ID
// @Produce		json
// @Success		200		{object}	albumDeleteResponse	        "Album deleted"
// @Failure		400		{object}	commonHttp.Error	"Client error"
// @Failure		500		{object}	commonHttp.Error	"Server error"
// @Router		/api/albums/{albumID}/ [delete]
func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	albumID, err := commonHttp.GetAlbumIDFromRequest(r)
	if err != nil {
		h.logger.Infof("get album by id : %v", err)
		commonHttp.ErrorResponse(w, "invalid url parameter", http.StatusBadRequest, h.logger)
		return
	}

	err = h.albumServices.DeleteByID(albumID)
	var errNoSuchAlbum *models.NoSuchAlbumError
	if errors.As(err, &errNoSuchAlbum) {
		h.logger.Info(err.Error())
		commonHttp.ErrorResponse(w, "no such album", http.StatusBadRequest, h.logger)
		return
	}
	if err != nil {
		h.logger.Error(err.Error())
		commonHttp.ErrorResponse(w, "error while deleting album", http.StatusInternalServerError, h.logger)
		return
	}

	adr := albumDeleteResponse{Status: "ok"}

	commonHttp.SuccessResponse(w, adr, h.logger)
}

// @Summary		Albums of Artist
// @Tags		Artist
// @Description	All albums of artist with chosen ID
// @Produce		json
// @Success		200		{object}	[]models.AlbumTransfer	    "Show albums"
// @Failure		400		{object}	commonHttp.Error	"Client error"
// @Failure		500		{object}	commonHttp.Error	"Server error"
// @Router		/api/artists/{artistID}/albums [get]
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
	}
	if err != nil {
		h.logger.Error(err.Error())
		commonHttp.ErrorResponse(w, "error while getting albums", http.StatusInternalServerError, h.logger)
		return
	}

	resp, err := models.AlbumTransferFromQuery(albums, h.artistServices.GetByAlbum)
	if err != nil {
		h.logger.Error(err.Error())
		commonHttp.ErrorResponse(w, "error while getting albums", http.StatusInternalServerError, h.logger)
		return
	}

	commonHttp.SuccessResponse(w, resp, h.logger)
}

// @Summary		Album Feed
// @Tags		Feed
// @Description	Feed albums
// @Produce		json
// @Success		200		{object}	[]models.AlbumTransfer	 "Albums feed"
// @Failure		500		{object}	commonHttp.Error "Server error"
// @Router		/api/albums/feed [get]
func (h *Handler) Feed(w http.ResponseWriter, r *http.Request) {
	albums, err := h.albumServices.GetFeed()
	if err != nil {
		h.logger.Error(err.Error())
		commonHttp.ErrorResponse(w, "error while getting albums", http.StatusInternalServerError, h.logger)
		return
	}

	resp, err := models.AlbumTransferFromQuery(albums, h.artistServices.GetByAlbum)
	if err != nil {
		h.logger.Error(err.Error())
		commonHttp.ErrorResponse(w, "error while getting albums", http.StatusInternalServerError, h.logger)
		return
	}

	commonHttp.SuccessResponse(w, resp, h.logger)
}
