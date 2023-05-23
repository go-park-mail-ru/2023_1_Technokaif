package http

import (
	"errors"
	"net/http"

	commonHTTP "github.com/go-park-mail-ru/2023_1_Technokaif/internal/common/http"
	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/models"
	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/album"
	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/artist"
	"github.com/go-park-mail-ru/2023_1_Technokaif/pkg/logger"
	"github.com/mailru/easyjson"
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

		logger: l,
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
	user, err := commonHTTP.GetUserFromRequest(r)
	if err != nil {
		commonHTTP.ErrorResponseWithErrLogging(w, r,
			commonHTTP.UnathorizedUser, http.StatusUnauthorized, h.logger, err)
		return
	}

	var aci albumCreateInput
	if err := easyjson.UnmarshalFromReader(r.Body, &aci); err != nil {
		commonHTTP.ErrorResponseWithErrLogging(w, r,
			commonHTTP.IncorrectRequestBody, http.StatusBadRequest, h.logger, err)
		return
	}

	if err := aci.validateAndEscape(); err != nil {
		commonHTTP.ErrorResponseWithErrLogging(w, r,
			commonHTTP.IncorrectRequestBody, http.StatusBadRequest, h.logger, err)
		return
	}

	album := aci.ToAlbum()

	albumID, err := h.albumServices.Create(r.Context(), album, aci.ArtistsID, user.ID)
	if err != nil {
		var errForbiddenUser *models.ForbiddenUserError
		if errors.As(err, &errForbiddenUser) {
			commonHTTP.ErrorResponseWithErrLogging(w, r,
				albumCreateNorights, http.StatusForbidden, h.logger, err)
			return
		}

		commonHTTP.ErrorResponseWithErrLogging(w, r,
			albumCreateServerError, http.StatusInternalServerError, h.logger, err)
		return
	}

	acr := albumCreateResponse{ID: albumID}

	commonHTTP.SuccessResponse(w, r, acr, h.logger)
}

// @Summary		Get Album
// @Tags		Album
// @Description	Get album with chosen ID
// @Produce		json
// @Success		200		{object}	models.AlbumTransfers	"Album got"
// @Failure		400		{object}	http.Error				"Incorrect input"
// @Failure		401		{object}	http.Error  			"User unathorized"
// @Failure		500		{object}	http.Error				"Server error"
// @Router		/api/albums/{albumID}/ [get]
func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	albumID, err := commonHTTP.GetAlbumIDFromRequest(r)
	if err != nil {
		commonHTTP.ErrorResponseWithErrLogging(w, r,
			commonHTTP.InvalidURLParameter, http.StatusBadRequest, h.logger, err)
		return
	}

	user, err := commonHTTP.GetUserFromRequest(r)
	if err != nil && !errors.Is(err, commonHTTP.ErrUnauthorized) {
		commonHTTP.ErrorResponseWithErrLogging(w, r,
			albumGetServerError, http.StatusInternalServerError, h.logger, err)
		return
	}

	album, err := h.albumServices.GetByID(r.Context(), albumID)
	if err != nil {
		var errNoSuchAlbum *models.NoSuchAlbumError
		if errors.As(err, &errNoSuchAlbum) {
			commonHTTP.ErrorResponseWithErrLogging(w, r,
				albumNotFound, http.StatusBadRequest, h.logger, err)
			return
		}

		commonHTTP.ErrorResponseWithErrLogging(w, r,
			albumGetServerError, http.StatusInternalServerError, h.logger, err)
		return
	}

	resp, err := models.AlbumTransferFromEntry(r.Context(), *album, user, h.albumServices.IsLiked,
		h.artistServices.IsLiked, h.artistServices.GetByAlbum)
	if err != nil {
		commonHTTP.ErrorResponseWithErrLogging(w, r,
			albumGetServerError, http.StatusInternalServerError, h.logger, err)
		return
	}

	commonHTTP.SuccessResponse(w, r, resp, h.logger)
}

// @Summary		Delete Album
// @Tags		Album
// @Description	Delete album with chosen ID
// @Produce		json
// @Success		200		{object}	albumDeleteResponse	  	"Album deleted"
// @Failure		400		{object}	http.Error				"Client error"
// @Failure		401		{object}	http.Error  			"User unathorized"
// @Failure		403		{object}	http.Error				"User hasn't rights"
// @Failure		500		{object}	http.Error				"Server error"
// @Router		/api/albums/{albumID}/ [delete]
func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	albumID, err := commonHTTP.GetAlbumIDFromRequest(r)
	if err != nil {
		commonHTTP.ErrorResponseWithErrLogging(w, r,
			commonHTTP.InvalidURLParameter, http.StatusBadRequest, h.logger, err)
		return
	}

	user, err := commonHTTP.GetUserFromRequest(r)
	if err != nil {
		commonHTTP.ErrorResponseWithErrLogging(w, r,
			commonHTTP.UnathorizedUser, http.StatusUnauthorized, h.logger, err)
		return
	}

	err = h.albumServices.Delete(r.Context(), albumID, user.ID)
	if err != nil {
		var errForbiddenUser *models.ForbiddenUserError
		if errors.As(err, &errForbiddenUser) {
			commonHTTP.ErrorResponseWithErrLogging(w, r,
				albumDeleteNoRights, http.StatusForbidden, h.logger, err)
			return
		}

		var errNoSuchAlbum *models.NoSuchAlbumError
		if errors.As(err, &errNoSuchAlbum) {
			commonHTTP.ErrorResponseWithErrLogging(w, r,
				albumNotFound, http.StatusBadRequest, h.logger, err)
			return
		}

		commonHTTP.ErrorResponseWithErrLogging(w, r,
			albumDeleteServerError, http.StatusInternalServerError, h.logger, err)
		return
	}

	adr := albumDeleteResponse{Status: albumDeletedSuccessfully}

	commonHTTP.SuccessResponse(w, r, adr, h.logger)
}

// @Summary		Albums of Artist
// @Tags		Artist
// @Description	All albums of artist with chosen ID
// @Produce		json
// @Success		200		{object}	models.AlbumTransfers 	"Show albums"
// @Failure		400		{object}	http.Error				"Client error"
// @Failure		500		{object}	http.Error				"Server error"
// @Router		/api/artists/{artistID}/albums [get]
func (h *Handler) GetByArtist(w http.ResponseWriter, r *http.Request) {
	artistID, err := commonHTTP.GetArtistIDFromRequest(r)
	if err != nil {
		commonHTTP.ErrorResponseWithErrLogging(w, r,
			commonHTTP.InvalidURLParameter, http.StatusBadRequest, h.logger, err)
		return
	}

	albums, err := h.albumServices.GetByArtist(r.Context(), artistID)
	if err != nil {
		var errNoSuchArtist *models.NoSuchArtistError
		if errors.As(err, &errNoSuchArtist) {
			commonHTTP.ErrorResponseWithErrLogging(w, r,
				artistNotFound, http.StatusBadRequest, h.logger, err)
			return
		}

		commonHTTP.ErrorResponseWithErrLogging(w, r,
			albumsGetServerError, http.StatusInternalServerError, h.logger, err)
		return
	}

	user, err := commonHTTP.GetUserFromRequest(r)
	if err != nil && !errors.Is(err, commonHTTP.ErrUnauthorized) {
		commonHTTP.ErrorResponseWithErrLogging(w, r,
			albumsGetServerError, http.StatusInternalServerError, h.logger, err)
		return
	}

	resp, err := models.AlbumTransferFromList(r.Context(), albums, user,
		h.albumServices.IsLiked, h.artistServices.IsLiked, h.artistServices.GetByAlbum)
	if err != nil {
		commonHTTP.ErrorResponseWithErrLogging(w, r,
			albumsGetServerError, http.StatusInternalServerError, h.logger, err)
		return
	}

	commonHTTP.SuccessResponse(w, r, resp, h.logger)
}

// @Summary		Album Feed
// @Tags		Feed
// @Description	Feed albums
// @Produce		json
// @Success		200		{object}	models.AlbumTransfer	"Albums feed"
// @Failure		500		{object}	http.Error 				"Server error"
// @Router		/api/albums/feed [get]
func (h *Handler) Feed(w http.ResponseWriter, r *http.Request) {
	albums, err := h.albumServices.GetFeed(r.Context())
	if err != nil {
		commonHTTP.ErrorResponseWithErrLogging(w, r,
			albumsGetServerError, http.StatusInternalServerError, h.logger, err)
		return
	}

	user, err := commonHTTP.GetUserFromRequest(r)
	if err != nil && !errors.Is(err, commonHTTP.ErrUnauthorized) {
		commonHTTP.ErrorResponseWithErrLogging(w, r,
			albumsGetServerError, http.StatusInternalServerError, h.logger, err)
		return
	}

	resp, err := models.AlbumTransferFromList(r.Context(), albums, user, h.albumServices.IsLiked,
		h.artistServices.IsLiked, h.artistServices.GetByAlbum)
	if err != nil {
		commonHTTP.ErrorResponseWithErrLogging(w, r,
			albumsGetServerError, http.StatusInternalServerError, h.logger, err)
		return
	}

	commonHTTP.SuccessResponse(w, r, resp, h.logger)
}

// @Summary      Favorite Albums
// @Tags         Favorite
// @Description  Get user's favorite albums
// @Produce      json
// @Success      200    {object}  	models.AlbumTransfers 	"Albums got"
// @Failure		 400	{object}	http.Error				"Incorrect input"
// @Failure      401    {object}  	http.Error  			"Unauthorized user"
// @Failure      403    {object}  	http.Error  			"Forbidden user"
// @Failure      500    {object}  	http.Error  			"Server error"
// @Router       /api/users/{userID}/favorite/albums [get]
func (h *Handler) GetFavorite(w http.ResponseWriter, r *http.Request) {
	user, err := commonHTTP.GetUserFromRequest(r)
	if err != nil {
		commonHTTP.ErrorResponseWithErrLogging(w, r,
			albumsGetServerError, http.StatusInternalServerError, h.logger, err)
		return
	}

	favAlbums, err := h.albumServices.GetLikedByUser(r.Context(), user.ID)
	if err != nil {
		commonHTTP.ErrorResponseWithErrLogging(w, r,
			albumsGetServerError, http.StatusInternalServerError, h.logger, err)
		return
	}

	at, err := models.AlbumTransferFromList(r.Context(), favAlbums, user, h.albumServices.IsLiked,
		h.artistServices.IsLiked, h.artistServices.GetByAlbum)
	if err != nil {
		commonHTTP.ErrorResponseWithErrLogging(w, r,
			albumsGetServerError, http.StatusInternalServerError, h.logger, err)
		return
	}

	commonHTTP.SuccessResponse(w, r, at, h.logger)
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
	albumID, err := commonHTTP.GetAlbumIDFromRequest(r)
	if err != nil {
		commonHTTP.ErrorResponseWithErrLogging(w, r,
			commonHTTP.InvalidURLParameter, http.StatusBadRequest, h.logger, err)
		return
	}

	user, err := commonHTTP.GetUserFromRequest(r)
	if err != nil {
		commonHTTP.ErrorResponseWithErrLogging(w, r,
			commonHTTP.UnathorizedUser, http.StatusUnauthorized, h.logger, err)
		return
	}

	notExisted, err := h.albumServices.SetLike(r.Context(), albumID, user.ID)
	if err != nil {
		var errNoSuchAlbum *models.NoSuchAlbumError
		if errors.As(err, &errNoSuchAlbum) {
			commonHTTP.ErrorResponseWithErrLogging(w, r,
				albumNotFound, http.StatusBadRequest, h.logger, err)
			return
		}

		commonHTTP.ErrorResponseWithErrLogging(w, r,
			commonHTTP.SetLikeServerError, http.StatusInternalServerError, h.logger, err)
		return
	}

	alr := albumLikeResponse{Status: commonHTTP.LikeSuccess}
	if !notExisted {
		alr.Status = commonHTTP.LikeAlreadyExists
	}
	commonHTTP.SuccessResponse(w, r, alr, h.logger)
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
	albumID, err := commonHTTP.GetAlbumIDFromRequest(r)
	if err != nil {
		commonHTTP.ErrorResponseWithErrLogging(w, r,
			commonHTTP.InvalidURLParameter, http.StatusBadRequest, h.logger, err)
		return
	}

	user, err := commonHTTP.GetUserFromRequest(r)
	if err != nil {
		commonHTTP.ErrorResponseWithErrLogging(w, r,
			commonHTTP.UnathorizedUser, http.StatusUnauthorized, h.logger, err)
		return
	}

	notExisted, err := h.albumServices.UnLike(r.Context(), albumID, user.ID)
	if err != nil {
		var errNoSuchAlbum *models.NoSuchAlbumError
		if errors.As(err, &errNoSuchAlbum) {
			commonHTTP.ErrorResponseWithErrLogging(w, r,
				albumNotFound, http.StatusBadRequest, h.logger, err)
			return
		}

		commonHTTP.ErrorResponseWithErrLogging(w, r,
			commonHTTP.DeleteLikeServerError, http.StatusInternalServerError, h.logger, err)
		return
	}

	alr := albumLikeResponse{Status: commonHTTP.UnLikeSuccess}
	if !notExisted {
		alr.Status = commonHTTP.LikeDoesntExist
	}
	commonHTTP.SuccessResponse(w, r, alr, h.logger)
}
