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

// @Summary		Create Artist
// @Tags		Artist
// @Description	Create new artist by sent object
// @Accept      json
// @Produce		json
// @Param		artist	body		artistCreateInput	true	"Track info"
// @Success		200		{object}	artistCreateResponse 		"Artist created"
// @Failure		400		{object}	http.Error	"Incorrect body"
// @Failure		401		{object}	http.Error  "User unathorized"
// @Failure		500		{object}	http.Error	"Server error"
// @Router		/api/artists/ [post]
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	user, err := commonHttp.GetUserFromRequest(r)
	if err != nil {
		commonHttp.ErrorResponseWithErrLogging(w, commonHttp.UnathorizedUser, http.StatusUnauthorized, h.logger, err)
		return
	}

	var aci artistCreateInput
	if err := json.NewDecoder(r.Body).Decode(&aci); err != nil {
		commonHttp.ErrorResponseWithErrLogging(w, commonHttp.IncorrectRequestBody, http.StatusBadRequest, h.logger, err)
		return
	}

	if err := aci.validateAndEscape(); err != nil {
		h.logger.Infof("Creating artist input validation failed: %s", err.Error())
		commonHttp.ErrorResponse(w, commonHttp.IncorrectRequestBody, http.StatusBadRequest, h.logger)
		return
	}

	artist := aci.ToArtist(&user.ID)

	artistID, err := h.artistServices.Create(artist)
	if err != nil {
		commonHttp.ErrorResponseWithErrLogging(w, artistCreateServerError, http.StatusInternalServerError, h.logger, err)
		return
	}

	acr := artistCreateResponse{ID: artistID}

	commonHttp.SuccessResponse(w, acr, h.logger)
}

// @Summary		Get Artist
// @Tags		Artist
// @Description	Get artist with chosen ID
// @Produce		json
// @Success		200		{object}	models.ArtistTransfer "Artist got"
// @Failure		400		{object}	http.Error			  "Incorrect body"
// @Failure		500		{object}	http.Error			  "Server error"
// @Router		/api/artists/{artistID}/ [get]
func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	artistID, err := commonHttp.GetArtistIDFromRequest(r)
	if err != nil {
		h.logger.Infof("Get artist by id: %v", err.Error())
		commonHttp.ErrorResponse(w, commonHttp.InvalidURLParameter, http.StatusBadRequest, h.logger)
		return
	}

	artist, err := h.artistServices.GetByID(artistID)
	if err != nil {
		var errNoSuchArtist *models.NoSuchArtistError
		if errors.As(err, &errNoSuchArtist) {
			commonHttp.ErrorResponseWithErrLogging(w, artistNotFound, http.StatusBadRequest, h.logger, err)
			return
		}

		commonHttp.ErrorResponseWithErrLogging(w, artistGetServerError, http.StatusInternalServerError, h.logger, err)
		return
	}

	user, err := commonHttp.GetUserFromRequest(r)
	if err != nil && !errors.Is(err, commonHttp.ErrUnauthorized) {
		commonHttp.ErrorResponseWithErrLogging(w, artistGetServerError, http.StatusInternalServerError, h.logger, err)
		return
	}

	ar, err := models.ArtistTransferFromEntry(*artist, user, h.artistServices.IsLiked)
	if err != nil {
		commonHttp.ErrorResponseWithErrLogging(w, artistGetServerError, http.StatusInternalServerError, h.logger, err)
		return
	}

	commonHttp.SuccessResponse(w, ar, h.logger)
}

// @Summary		Delete Artist
// @Tags		Artist
// @Description	Delete artist with chosen ID
// @Produce		json
// @Success		200		{object}	artistDeleteResponse "Artist deleted"
// @Failure		400		{object}	http.Error			 "Incorrect body"
// @Failure		401		{object}	http.Error  		 "User unathorized"
// @Failure		403		{object}	http.Error			 "User hasn't rights"
// @Failure		500		{object}	http.Error			 "Server error"
// @Router		/api/artists/{artistID}/ [delete]
func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	artistID, err := commonHttp.GetArtistIDFromRequest(r)
	if err != nil {
		h.logger.Infof("Get artist by id: %v", err)
		commonHttp.ErrorResponse(w, commonHttp.InvalidURLParameter, http.StatusBadRequest, h.logger)
		return
	}

	user, err := commonHttp.GetUserFromRequest(r)
	if err != nil {
		commonHttp.ErrorResponseWithErrLogging(w, commonHttp.UnathorizedUser, http.StatusUnauthorized, h.logger, err)
		return
	}

	err = h.artistServices.Delete(artistID, user.ID)
	if err != nil {
		var errNoSuchArtist *models.NoSuchArtistError
		if errors.As(err, &errNoSuchArtist) {
			commonHttp.ErrorResponseWithErrLogging(w, artistNotFound, http.StatusBadRequest, h.logger, err)
			return
		}

		var errForbiddenUser *models.ForbiddenUserError
		if errors.As(err, &errForbiddenUser) {
			commonHttp.ErrorResponseWithErrLogging(w, artistDeleteNoRights, http.StatusForbidden, h.logger, err)
			return
		}

		commonHttp.ErrorResponseWithErrLogging(w, artistDeleteServerError, http.StatusInternalServerError, h.logger, err)
		return
	}

	adr := artistDeleteResponse{Status: artistDeletedSuccessfully}

	commonHttp.SuccessResponse(w, adr, h.logger)
}

// @Summary		Artist Feed
// @Tags		Feed
// @Description	Feed artists
// @Produce		json
// @Success		200		{object}	[]models.ArtistTransfer	"Artists feed"
// @Failure		500		{object}	http.Error				"Server error"
// @Router		/api/artists/feed [get]
func (h *Handler) Feed(w http.ResponseWriter, r *http.Request) {
	artists, err := h.artistServices.GetFeed()
	if err != nil {
		commonHttp.ErrorResponseWithErrLogging(w, artistsGetServerError, http.StatusInternalServerError, h.logger, err)
		return
	}

	user, err := commonHttp.GetUserFromRequest(r)
	if err != nil && !errors.Is(err, commonHttp.ErrUnauthorized) {
		commonHttp.ErrorResponseWithErrLogging(w, artistsGetServerError, http.StatusInternalServerError, h.logger, err)
		return
	}

	artistsTransfer, err := models.ArtistTransferFromQuery(artists, user, h.artistServices.IsLiked)
	if err != nil {
		commonHttp.ErrorResponseWithErrLogging(w, artistsGetServerError, http.StatusInternalServerError, h.logger, err)
		return
	}

	commonHttp.SuccessResponse(w, artistsTransfer, h.logger)
}

// @Summary      Favorite Artists
// @Tags         Favorite
// @Description  Get user's favorite artists
// @Produce      application/json
// @Success      200    {object}  	[]models.ArtistTransfer "Artists got"
// @Failure		 400	{object}	http.Error				"Incorrect input"
// @Failure      401    {object}  	http.Error  			"Unauthorized user"
// @Failure      403    {object}  	http.Error  			"Forbidden user"
// @Failure      500    {object}  	http.Error  			"Server error"
// @Router       /api/users/{userID}/favorite/artists [get]
func (h *Handler) GetFavorite(w http.ResponseWriter, r *http.Request) {
	user, err := commonHttp.GetUserFromRequest(r)
	if err != nil {
		commonHttp.ErrorResponseWithErrLogging(w, artistsGetServerError, http.StatusInternalServerError, h.logger, err)
		return
	}

	artists, err := h.artistServices.GetLikedByUser(user.ID)
	if err != nil {
		commonHttp.ErrorResponseWithErrLogging(w, artistsGetServerError, http.StatusInternalServerError, h.logger, err)
		return
	}

	at, err := models.ArtistTransferFromQuery(artists, user, h.artistServices.IsLiked)
	if err != nil {
		commonHttp.ErrorResponseWithErrLogging(w, artistsGetServerError, http.StatusInternalServerError, h.logger, err)
		return
	}

	commonHttp.SuccessResponse(w, at, h.logger)
}

// @Summary		Set like
// @Tags		Artist
// @Description	Set like by user to chosen artist (add to favorite)
// @Produce		json
// @Success		200		{object}	artistLikeResponse	"Like set"
// @Failure		400		{object}	http.Error			"Client error"
// @Failure		401		{object}	http.Error  		"User unathorized"
// @Failure		500		{object}	http.Error			"Server error"
// @Router		/api/artists/{artistID}/like [post]
func (h *Handler) Like(w http.ResponseWriter, r *http.Request) {
	artistID, err := commonHttp.GetArtistIDFromRequest(r)
	if err != nil {
		h.logger.Infof("Get artist by id: %v", err)
		commonHttp.ErrorResponse(w, commonHttp.InvalidURLParameter, http.StatusBadRequest, h.logger)
		return
	}

	user, err := commonHttp.GetUserFromRequest(r)
	if err != nil {
		commonHttp.ErrorResponseWithErrLogging(w, commonHttp.UnathorizedUser, http.StatusUnauthorized, h.logger, err)
		return
	}

	notExisted, err := h.artistServices.SetLike(artistID, user.ID)
	if err != nil {
		var errNoSuchArtist *models.NoSuchArtistError
		if errors.As(err, &errNoSuchArtist) {
			commonHttp.ErrorResponseWithErrLogging(w, artistNotFound, http.StatusBadRequest, h.logger, err)
			return
		}

		commonHttp.ErrorResponseWithErrLogging(w, commonHttp.SetLikeServerError, http.StatusInternalServerError, h.logger, err)
		return
	}

	alr := artistLikeResponse{Status: commonHttp.LikeSuccess}
	if !notExisted {
		alr.Status = commonHttp.LikeAlreadyExists // "already liked"
	}
	commonHttp.SuccessResponse(w, alr, h.logger)
}

// @Summary		Remove like
// @Tags		Artist
// @Description	Remove like by user from chosen artist (remove from favorite)
// @Produce		json
// @Success		200		{object}	artistLikeResponse	"Like removed"
// @Failure		400		{object}	http.Error			"Client error"
// @Failure		401		{object}	http.Error  		"User unathorized"
// @Failure		500		{object}	http.Error			"Server error"
// @Router		/api/artists/{artistID}/unlike [post]
func (h *Handler) UnLike(w http.ResponseWriter, r *http.Request) {
	user, err := commonHttp.GetUserFromRequest(r)
	if err != nil {
		commonHttp.ErrorResponseWithErrLogging(w, commonHttp.UnathorizedUser, http.StatusUnauthorized, h.logger, err)
		return
	}

	artistID, err := commonHttp.GetArtistIDFromRequest(r)
	if err != nil {
		h.logger.Infof("Get artist by id: %v", err)
		commonHttp.ErrorResponse(w, commonHttp.InvalidURLParameter, http.StatusBadRequest, h.logger)
		return
	}

	notExisted, err := h.artistServices.UnLike(artistID, user.ID)
	if err != nil {
		var errNoSuchArtist *models.NoSuchArtistError
		if errors.As(err, &errNoSuchArtist) {
			commonHttp.ErrorResponseWithErrLogging(w, artistNotFound, http.StatusBadRequest, h.logger, err)
			return
		}

		commonHttp.ErrorResponseWithErrLogging(w, commonHttp.DeleteLikeServerError, http.StatusInternalServerError, h.logger, err)
		return
	}

	alr := artistLikeResponse{Status: commonHttp.UnLikeSuccess}
	if !notExisted {
		alr.Status = commonHttp.LikeDoesntExist
	}
	commonHttp.SuccessResponse(w, alr, h.logger)
}
