package http

import (
	"errors"
	"net/http"

	commonHTTP "github.com/go-park-mail-ru/2023_1_Technokaif/internal/common/http"
	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/models"
	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/artist"
	"github.com/go-park-mail-ru/2023_1_Technokaif/pkg/logger"
	easyjson "github.com/mailru/easyjson"
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
	user, err := commonHTTP.GetUserFromRequest(r)
	if err != nil {
		commonHTTP.ErrorResponseWithErrLogging(w, r,
			commonHTTP.UnathorizedUser, http.StatusUnauthorized, h.logger, err)
		return
	}

	var aci artistCreateInput
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

	artist := aci.ToArtist(&user.ID)

	artistID, err := h.artistServices.Create(r.Context(), artist)
	if err != nil {
		commonHTTP.ErrorResponseWithErrLogging(w, r,
			artistCreateServerError, http.StatusInternalServerError, h.logger, err)
		return
	}

	acr := artistCreateResponse{ID: artistID}

	commonHTTP.SuccessResponse(w, r, acr, h.logger)
}

// @Summary		Get Artist
// @Tags		Artist
// @Description	Get artist with chosen ID
// @Produce		json
// @Success		200		{object}	models.ArtistTransfer 		"Artist got"
// @Failure		400		{object}	http.Error					"Incorrect body"
// @Failure		500		{object}	http.Error					"Server error"
// @Router		/api/artists/{artistID}/ [get]
func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	artistID, err := commonHTTP.GetArtistIDFromRequest(r)
	if err != nil {
		commonHTTP.ErrorResponseWithErrLogging(w, r,
			commonHTTP.InvalidURLParameter, http.StatusBadRequest, h.logger, err)
		return
	}

	user, err := commonHTTP.GetUserFromRequest(r)
	if err != nil && !errors.Is(err, commonHTTP.ErrUnauthorized) {
		commonHTTP.ErrorResponseWithErrLogging(w, r,
			artistGetServerError, http.StatusInternalServerError, h.logger, err)
		return
	}

	artist, err := h.artistServices.GetByID(r.Context(), artistID)
	if err != nil {
		var errNoSuchArtist *models.NoSuchArtistError
		if errors.As(err, &errNoSuchArtist) {
			commonHTTP.ErrorResponseWithErrLogging(w, r,
				artistNotFound, http.StatusBadRequest, h.logger, err)
			return
		}

		commonHTTP.ErrorResponseWithErrLogging(w, r,
			artistGetServerError, http.StatusInternalServerError, h.logger, err)
		return
	}

	ar, err := models.ArtistTransferFromEntry(r.Context(), *artist, user, h.artistServices.IsLiked)
	if err != nil {
		commonHTTP.ErrorResponseWithErrLogging(w, r,
			artistGetServerError, http.StatusInternalServerError, h.logger, err)
		return
	}

	commonHTTP.SuccessResponse(w, r, ar, h.logger)
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
	artistID, err := commonHTTP.GetArtistIDFromRequest(r)
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

	err = h.artistServices.Delete(r.Context(), artistID, user.ID)
	if err != nil {
		var errNoSuchArtist *models.NoSuchArtistError
		if errors.As(err, &errNoSuchArtist) {
			commonHTTP.ErrorResponseWithErrLogging(w, r,
				artistNotFound, http.StatusBadRequest, h.logger, err)
			return
		}

		var errForbiddenUser *models.ForbiddenUserError
		if errors.As(err, &errForbiddenUser) {
			commonHTTP.ErrorResponseWithErrLogging(w, r,
				artistDeleteNoRights, http.StatusForbidden, h.logger, err)
			return
		}

		commonHTTP.ErrorResponseWithErrLogging(w, r,
			artistDeleteServerError, http.StatusInternalServerError, h.logger, err)
		return
	}

	adr := artistDeleteResponse{Status: artistDeletedSuccessfully}

	commonHTTP.SuccessResponse(w, r, adr, h.logger)
}

// @Summary		Artist Feed
// @Tags		Feed
// @Description	Feed artists
// @Produce		json
// @Param		afi		body		artistFeedInput    true	"Feed info"
// @Success		200		{object}	models.ArtistTransfers	"Artists feed"
// @Failure		500		{object}	http.Error				"Server error"
// @Router		/api/artists/feed [post]
func (h *Handler) FeedTop(w http.ResponseWriter, r *http.Request) {
	var afi artistFeedInput
	if err := easyjson.UnmarshalFromReader(r.Body, &afi); err != nil {
		commonHTTP.ErrorResponseWithErrLogging(w, r,
			commonHTTP.IncorrectRequestBody, http.StatusBadRequest, h.logger, err)
		return
	}

	artists, err := h.artistServices.GetFeedTop(r.Context(), afi.Days)
	if err != nil {
		commonHTTP.ErrorResponseWithErrLogging(w, r,
			artistsGetServerError, http.StatusInternalServerError, h.logger, err)
		return
	}

	user, err := commonHTTP.GetUserFromRequest(r)
	if err != nil && !errors.Is(err, commonHTTP.ErrUnauthorized) {
		commonHTTP.ErrorResponseWithErrLogging(w, r,
			artistsGetServerError, http.StatusInternalServerError, h.logger, err)
		return
	}

	at, err := models.ArtistTransferFromList(r.Context(), artists, user, h.artistServices.IsLiked)
	if err != nil {
		commonHTTP.ErrorResponseWithErrLogging(w, r,
			artistsGetServerError, http.StatusInternalServerError, h.logger, err)
		return
	}

	commonHTTP.SuccessResponse(w, r, at, h.logger)
}

// @Summary		Artist Feed
// @Tags		Feed
// @Description	Feed artists
// @Produce		json
// @Success		200		{object}	models.ArtistTransfers	"Artists feed"
// @Failure		500		{object}	http.Error				"Server error"
// @Router		/api/artists/feed [get]
func (h *Handler) Feed(w http.ResponseWriter, r *http.Request) {
	artists, err := h.artistServices.GetFeed(r.Context())
	if err != nil {
		commonHTTP.ErrorResponseWithErrLogging(w, r,
			artistsGetServerError, http.StatusInternalServerError, h.logger, err)
		return
	}

	user, err := commonHTTP.GetUserFromRequest(r)
	if err != nil && !errors.Is(err, commonHTTP.ErrUnauthorized) {
		commonHTTP.ErrorResponseWithErrLogging(w, r,
			artistsGetServerError, http.StatusInternalServerError, h.logger, err)
		return
	}

	at, err := models.ArtistTransferFromList(r.Context(), artists, user, h.artistServices.IsLiked)
	if err != nil {
		commonHTTP.ErrorResponseWithErrLogging(w, r,
			artistsGetServerError, http.StatusInternalServerError, h.logger, err)
		return
	}

	commonHTTP.SuccessResponse(w, r, at, h.logger)
}

// @Summary      Favorite Artists
// @Tags         Favorite
// @Description  Get user's favorite artists
// @Produce      json
// @Success      200    {object}  	models.ArtistTransfers 	"Artists got"
// @Failure		 400	{object}	http.Error				"Incorrect input"
// @Failure      401    {object}  	http.Error  			"Unauthorized user"
// @Failure      403    {object}  	http.Error  			"Forbidden user"
// @Failure      500    {object}  	http.Error  			"Server error"
// @Router       /api/users/{userID}/favorite/artists [get]
func (h *Handler) GetFavorite(w http.ResponseWriter, r *http.Request) {
	user, err := commonHTTP.GetUserFromRequest(r)
	if err != nil {
		commonHTTP.ErrorResponseWithErrLogging(w, r,
			artistsGetServerError, http.StatusInternalServerError, h.logger, err)
		return
	}

	artists, err := h.artistServices.GetLikedByUser(r.Context(), user.ID)
	if err != nil {
		commonHTTP.ErrorResponseWithErrLogging(w, r,
			artistsGetServerError, http.StatusInternalServerError, h.logger, err)
		return
	}

	at, err := models.ArtistTransferFromList(r.Context(), artists, user, h.artistServices.IsLiked)
	if err != nil {
		commonHTTP.ErrorResponseWithErrLogging(w, r,
			artistsGetServerError, http.StatusInternalServerError, h.logger, err)
		return
	}

	commonHTTP.SuccessResponse(w, r, at, h.logger)
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
	artistID, err := commonHTTP.GetArtistIDFromRequest(r)
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

	notExisted, err := h.artistServices.SetLike(r.Context(), artistID, user.ID)
	if err != nil {
		var errNoSuchArtist *models.NoSuchArtistError
		if errors.As(err, &errNoSuchArtist) {
			commonHTTP.ErrorResponseWithErrLogging(w, r,
				artistNotFound, http.StatusBadRequest, h.logger, err)
			return
		}

		commonHTTP.ErrorResponseWithErrLogging(w, r,
			commonHTTP.SetLikeServerError, http.StatusInternalServerError, h.logger, err)
		return
	}

	alr := artistLikeResponse{Status: commonHTTP.LikeSuccess}
	if !notExisted {
		alr.Status = commonHTTP.LikeAlreadyExists // "already liked"
	}
	commonHTTP.SuccessResponse(w, r, alr, h.logger)
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
	user, err := commonHTTP.GetUserFromRequest(r)
	if err != nil {
		commonHTTP.ErrorResponseWithErrLogging(w, r,
			commonHTTP.UnathorizedUser, http.StatusUnauthorized, h.logger, err)
		return
	}

	artistID, err := commonHTTP.GetArtistIDFromRequest(r)
	if err != nil {
		commonHTTP.ErrorResponseWithErrLogging(w, r,
			commonHTTP.InvalidURLParameter, http.StatusBadRequest, h.logger, err)
		return
	}

	notExisted, err := h.artistServices.UnLike(r.Context(), artistID, user.ID)
	if err != nil {
		var errNoSuchArtist *models.NoSuchArtistError
		if errors.As(err, &errNoSuchArtist) {
			commonHTTP.ErrorResponseWithErrLogging(w, r,
				artistNotFound, http.StatusBadRequest, h.logger, err)
			return
		}

		commonHTTP.ErrorResponseWithErrLogging(w, r,
			commonHTTP.DeleteLikeServerError, http.StatusInternalServerError, h.logger, err)
		return
	}

	alr := artistLikeResponse{Status: commonHTTP.UnLikeSuccess}
	if !notExisted {
		alr.Status = commonHTTP.LikeDoesntExist
	}
	commonHTTP.SuccessResponse(w, r, alr, h.logger)
}
