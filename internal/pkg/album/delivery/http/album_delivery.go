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
// @Param		album	body		albumCreateInput	true	"Album info"
// @Success		200		{object}	albumCreateResponse	        "Album created"
// @Failure		400		{object}	http.Error					"Incorrect input"
// @Failure		401		{object}	http.Error  				"User unathorized"
// @Failure		403		{object}	http.Error					"User hasn't rights"
// @Failure		500		{object}	http.Error					"Server error"
// @Router		/api/albums/ [post]
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	user, err := commonHttp.GetUserFromRequest(r)
	if err != nil {
		commonHttp.ErrorResponseWithErrLogging(w, "unathorized", http.StatusUnauthorized, h.logger, err)
		return
	}

	var aci albumCreateInput
	if err := json.NewDecoder(r.Body).Decode(&aci); err != nil {
		commonHttp.ErrorResponseWithErrLogging(w, "incorrect input body", http.StatusBadRequest, h.logger, err)
		return
	}

	if err := aci.validate(); err != nil {
		h.logger.Infof("Creating album input validation failed: %s", err.Error())
		commonHttp.ErrorResponse(w, "incorrect input body", http.StatusBadRequest, h.logger)
		return
	}

	album := aci.ToAlbum()

	albumID, err := h.albumServices.Create(album, aci.ArtistsID, user.ID)
	if err != nil {
		var errForbiddenUser *models.ForbiddenUserError
		if errors.As(err, &errForbiddenUser) {
			commonHttp.ErrorResponseWithErrLogging(w, "no rights to crearte album", http.StatusForbidden, h.logger, err)
			return
		}

		commonHttp.ErrorResponseWithErrLogging(w, "can't create album", http.StatusInternalServerError, h.logger, err)
		return
	}

	acr := albumCreateResponse{ID: albumID}

	commonHttp.SuccessResponse(w, acr, h.logger)
}

// @Summary		Get Album
// @Tags		Album
// @Description	Get album with chosen ID
// @Produce		json
// @Success		200		{object}	models.AlbumTransfer	"Album got"
// @Failure		400		{object}	http.Error				"Incorrect input"
// @Failure		401		{object}	http.Error  			"User unathorized"
// @Failure		500		{object}	http.Error				"Server error"
// @Router		/api/albums/{albumID}/ [get]
func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	albumID, err := commonHttp.GetAlbumIDFromRequest(r)
	if err != nil {
		h.logger.Infof("Get album by id: %v", err)
		commonHttp.ErrorResponse(w, "invalid url parameter", http.StatusBadRequest, h.logger)
		return
	}

	if _, err := commonHttp.GetUserFromRequest(r); err != nil {
		commonHttp.ErrorResponseWithErrLogging(w, "unathorized", http.StatusUnauthorized, h.logger, err)
		return
	}

	album, err := h.albumServices.GetByID(albumID)
	if err != nil {
		var errNoSuchAlbum *models.NoSuchAlbumError
		if errors.As(err, &errNoSuchAlbum) {
			commonHttp.ErrorResponseWithErrLogging(w, "no such album", http.StatusBadRequest, h.logger, err)
			return
		}

		commonHttp.ErrorResponseWithErrLogging(w, "can't get album", http.StatusInternalServerError, h.logger, err)
		return
	}

	resp, err := models.AlbumTransferFromEntry(*album, h.artistServices.GetByAlbum)
	if err != nil {
		commonHttp.ErrorResponseWithErrLogging(w, "can't get album", http.StatusInternalServerError, h.logger, err)
		return
	}

	commonHttp.SuccessResponse(w, resp, h.logger)
}

// @Summary		Delete Album
// @Tags		Album
// @Description	Delete album with chosen ID
// @Produce		json
// @Success		200		{object}	albumDeleteResponse	        "Album deleted"
// @Failure		400		{object}	http.Error	"Client error"
// @Failure		401		{object}	http.Error  "User unathorized"
// @Failure		403		{object}	http.Error	"User hasn't rights"
// @Failure		500		{object}	http.Error	"Server error"
// @Router		/api/albums/{albumID}/ [delete]
func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	albumID, err := commonHttp.GetAlbumIDFromRequest(r)
	if err != nil {
		h.logger.Infof("Get album's id: %v", err)
		commonHttp.ErrorResponse(w, "invalid url parameter", http.StatusBadRequest, h.logger)
		return
	}

	user, err := commonHttp.GetUserFromRequest(r)
	if err != nil {
		commonHttp.ErrorResponseWithErrLogging(w, "unathorized", http.StatusUnauthorized, h.logger, err)
		return
	}

	err = h.albumServices.Delete(albumID, user.ID)
	if err != nil {
		var errForbiddenUser *models.ForbiddenUserError
		if errors.As(err, &errForbiddenUser) {
			commonHttp.ErrorResponseWithErrLogging(w, "no rights to delete album", http.StatusForbidden, h.logger, err)
			return
		}

		var errNoSuchAlbum *models.NoSuchAlbumError
		if errors.As(err, &errNoSuchAlbum) {
			commonHttp.ErrorResponseWithErrLogging(w, "no such album", http.StatusBadRequest, h.logger, err)
			return
		}

		commonHttp.ErrorResponseWithErrLogging(w, "can't delete album", http.StatusInternalServerError, h.logger, err)
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
// @Failure		400		{object}	http.Error	"Client error"
// @Failure		500		{object}	http.Error	"Server error"
// @Router		/api/artists/{artistID}/albums [get]
func (h *Handler) GetByArtist(w http.ResponseWriter, r *http.Request) {
	artistID, err := commonHttp.GetArtistIDFromRequest(r)
	if err != nil {
		h.logger.Infof("Get artist by id: %v", err)
		commonHttp.ErrorResponse(w, "invalid url parameter", http.StatusBadRequest, h.logger)
		return
	}

	albums, err := h.albumServices.GetByArtist(artistID)
	if err != nil {
		var errNoSuchArtist *models.NoSuchArtistError
		if errors.As(err, &errNoSuchArtist) {
			commonHttp.ErrorResponseWithErrLogging(w, "no such artist", http.StatusBadRequest, h.logger, err)
			return
		}

		commonHttp.ErrorResponseWithErrLogging(w, "can't get albums", http.StatusInternalServerError, h.logger, err)
		return
	}

	resp, err := models.AlbumTransferFromQuery(albums, h.artistServices.GetByAlbum)
	if err != nil {
		commonHttp.ErrorResponseWithErrLogging(w, "can't get albums", http.StatusInternalServerError, h.logger, err)
		return
	}

	commonHttp.SuccessResponse(w, resp, h.logger)
}

// @Summary		Album Feed
// @Tags		Feed
// @Description	Feed albums
// @Produce		json
// @Success		200		{object}	[]models.AlbumTransfer	 "Albums feed"
// @Failure		500		{object}	http.Error "Server error"
// @Router		/api/albums/feed [get]
func (h *Handler) Feed(w http.ResponseWriter, r *http.Request) {
	albums, err := h.albumServices.GetFeed()
	if err != nil {
		commonHttp.ErrorResponseWithErrLogging(w, "can't get albums", http.StatusInternalServerError, h.logger, err)
		return
	}

	resp, err := models.AlbumTransferFromQuery(albums, h.artistServices.GetByAlbum)
	if err != nil {
		commonHttp.ErrorResponseWithErrLogging(w, "can't get albums", http.StatusInternalServerError, h.logger, err)
		return
	}

	commonHttp.SuccessResponse(w, resp, h.logger)
}

// @Summary		Set like
// @Tags		Album
// @Description	Set like by user to chosen album (add to favorite)
// @Produce		json
// @Success		200		{object}	albumLikeResponse	"Like set"
// @Failure		400		{object}	http.Error			"Client error"
// @Failure		401		{object}	http.Error  		"User unathorized"
// @Failure		500		{object}	http.Error			"Server error"
// @Router		/api/albums/{albumID}/like [post]
func (h *Handler) Like(w http.ResponseWriter, r *http.Request) {
	albumID, err := commonHttp.GetAlbumIDFromRequest(r)
	if err != nil {
		h.logger.Infof("Get album by id: %v", err)
		commonHttp.ErrorResponse(w, "invalid url parameter", http.StatusBadRequest, h.logger)
		return
	}

	user, err := commonHttp.GetUserFromRequest(r)
	if err != nil {
		commonHttp.ErrorResponseWithErrLogging(w, "unathorized", http.StatusUnauthorized, h.logger, err)
		return
	}

	notExisted, err := h.albumServices.SetLike(albumID, user.ID)
	if err != nil {
		var errNoSuchAlbum *models.NoSuchAlbumError
		if errors.As(err, &errNoSuchAlbum) {
			commonHttp.ErrorResponseWithErrLogging(w, "no such album", http.StatusBadRequest, h.logger, err)
			return
		}

		commonHttp.ErrorResponseWithErrLogging(w, "can't set like", http.StatusInternalServerError, h.logger, err)
		return
	}

	alr := albumLikeResponse{Status: "ok"}
	if !notExisted {
		alr.Status = "already liked"
	}
	commonHttp.SuccessResponse(w, alr, h.logger)
}

// @Summary		Remove like
// @Tags		Album
// @Description	Remove like by user from chosen album (remove from favorite)
// @Produce		json
// @Success		200		{object}	albumLikeResponse	"Like removed"
// @Failure		400		{object}	http.Error			"Client error"
// @Failure		401		{object}	http.Error  		"User unathorized"
// @Failure		500		{object}	http.Error			"Server error"
// @Router		/api/albums/{albumID}/unlike [post]
func (h *Handler) UnLike(w http.ResponseWriter, r *http.Request) {
	albumID, err := commonHttp.GetAlbumIDFromRequest(r)
	if err != nil {
		h.logger.Infof("Get album by id: %v", err)
		commonHttp.ErrorResponse(w, "invalid url parameter", http.StatusBadRequest, h.logger)
		return
	}

	user, err := commonHttp.GetUserFromRequest(r)
	if err != nil {
		commonHttp.ErrorResponseWithErrLogging(w, "unathorized", http.StatusUnauthorized, h.logger, err)
		return
	}

	notExisted, err := h.albumServices.UnLike(albumID, user.ID)
	if err != nil {
		var errNoSuchAlbum *models.NoSuchAlbumError
		if errors.As(err, &errNoSuchAlbum) {
			commonHttp.ErrorResponseWithErrLogging(w, "no such album", http.StatusBadRequest, h.logger, err)
			return
		}

		commonHttp.ErrorResponseWithErrLogging(w, "can't remove like", http.StatusInternalServerError, h.logger, err)
		return
	}

	alr := albumLikeResponse{Status: "ok"}
	if !notExisted {
		alr.Status = "wasn't liked"
	}
	commonHttp.SuccessResponse(w, alr, h.logger)
}
