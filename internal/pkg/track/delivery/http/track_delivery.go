package http

import (
	"errors"
	"net/http"

	commonHTTP "github.com/go-park-mail-ru/2023_1_Technokaif/internal/common/http"
	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/models"
	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/artist"
	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/track"
	"github.com/go-park-mail-ru/2023_1_Technokaif/pkg/logger"
	easyjson "github.com/mailru/easyjson"
)

type Handler struct {
	trackServices  track.Usecase
	artistServices artist.Usecase
	logger         logger.Logger
}

func NewHandler(tu track.Usecase, au artist.Usecase, l logger.Logger) *Handler {
	return &Handler{
		trackServices:  tu,
		artistServices: au,

		logger: l,
	}
}

// @Summary		Create Track
// @Tags		Track
// @Description	Create new track by sent object
// @Accept      json
// @Produce		json
// @Param		track	body		trackCreateInput    true "Track info"
// @Success		200		{object}	trackCreateResponse	 	 "Track created"
// @Failure		400		{object}	http.Error				 "Incorrect body"
// @Failure		401		{object}	http.Error  			 "User unathorized"
// @Failure		403		{object}	http.Error				 "User hasn't rights"
// @Failure		500		{object}	http.Error				 "Server error"
// @Router		/api/tracks/ [post]
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	user, err := commonHTTP.GetUserFromRequest(r)
	if err != nil {
		commonHTTP.ErrorResponseWithErrLogging(w, r,
			commonHTTP.UnathorizedUser, http.StatusUnauthorized, h.logger, err)
		return
	}

	var tci trackCreateInput
	if err := easyjson.UnmarshalFromReader(r.Body, &tci); err != nil {
		commonHTTP.ErrorResponseWithErrLogging(w, r,
			commonHTTP.IncorrectRequestBody, http.StatusBadRequest, h.logger, err)
		return
	}

	if err := tci.validateAndEscape(); err != nil {
		commonHTTP.ErrorResponseWithErrLogging(w, r,
			commonHTTP.IncorrectRequestBody, http.StatusBadRequest, h.logger, err)
		return
	}

	track := tci.ToTrack()

	trackID, err := h.trackServices.Create(r.Context(), track, tci.ArtistsID, user.ID)
	if err != nil {
		var errForbiddenUser *models.ForbiddenUserError
		if errors.As(err, &errForbiddenUser) {
			commonHTTP.ErrorResponseWithErrLogging(w, r,
				trackCreateNorights, http.StatusForbidden, h.logger, err)
			return
		}

		commonHTTP.ErrorResponseWithErrLogging(w, r,
			trackCreateServerError, http.StatusInternalServerError, h.logger, err)
		return
	}

	tcr := trackCreateResponse{ID: trackID}

	commonHTTP.SuccessResponse(w, r, tcr, h.logger)
}

// @Summary		Get Track
// @Tags		Track
// @Description	Get track with chosen ID
// @Produce		json
// @Success		200		{object}	models.TrackTransfer 	"Track got"
// @Failure		400		{object}	http.Error				"Client error"
// @Failure		401		{object}	http.Error  			"User unathorized"
// @Failure		500		{object}	http.Error				"Server error"
// @Router		/api/tracks/{trackID}/ [get]
func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	trackID, err := commonHTTP.GetTrackIDFromRequest(r)
	if err != nil {
		commonHTTP.ErrorResponseWithErrLogging(w, r,
			commonHTTP.InvalidURLParameter, http.StatusBadRequest, h.logger, err)
		return
	}

	user, err := commonHTTP.GetUserFromRequest(r)
	if err != nil && !errors.Is(err, commonHTTP.ErrUnauthorized) {
		commonHTTP.ErrorResponseWithErrLogging(w, r,
			trackGetServerError, http.StatusInternalServerError, h.logger, err)
		return
	}

	track, err := h.trackServices.GetByID(r.Context(), trackID)
	if err != nil {
		var errNoSuchTrack *models.NoSuchTrackError
		if errors.As(err, &errNoSuchTrack) {
			commonHTTP.ErrorResponseWithErrLogging(w, r,
				trackNotFound, http.StatusBadRequest, h.logger, err)
			return
		}

		commonHTTP.ErrorResponseWithErrLogging(w, r,
			trackGetServerError, http.StatusInternalServerError, h.logger, err)
		return
	}

	tt, err := models.TrackTransferFromEntry(r.Context(), *track, user, h.trackServices.IsLiked,
		h.artistServices.IsLiked, h.artistServices.GetByTrack)
	if err != nil {
		commonHTTP.ErrorResponseWithErrLogging(w, r,
			trackGetServerError, http.StatusInternalServerError, h.logger, err)
		return
	}

	commonHTTP.SuccessResponse(w, r, tt, h.logger)
}

// @Summary		Delete Track
// @Tags		Track
// @Description	Delete track with chosen ID
// @Produce		json
// @Success		200		{object}	trackDeleteResponse	"Track deleted"
// @Failure		400		{object}	http.Error			"No such track"
// @Failure		401		{object}	http.Error  		"User unathorized"
// @Failure		403		{object}	http.Error			"User hasn't rights"
// @Failure		500		{object}	http.Error			"Server error"
// @Router		/api/tracks/{trackID}/ [delete]
func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	trackID, err := commonHTTP.GetTrackIDFromRequest(r)
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

	err = h.trackServices.Delete(r.Context(), trackID, user.ID)
	if err != nil {
		var errForbiddenUser *models.ForbiddenUserError
		if errors.As(err, &errForbiddenUser) {
			commonHTTP.ErrorResponseWithErrLogging(w, r,
				trackDeleteNoRights, http.StatusForbidden, h.logger, err)
			return
		}

		var errNoSuchTrack *models.NoSuchTrackError
		if errors.As(err, &errNoSuchTrack) {
			commonHTTP.ErrorResponseWithErrLogging(w, r,
				trackNotFound, http.StatusBadRequest, h.logger, err)
			return
		}

		commonHTTP.ErrorResponseWithErrLogging(w, r,
			trackDeleteServerError, http.StatusInternalServerError, h.logger, err)
		return
	}

	tdr := trackDeleteResponse{Status: trackDeletedSuccessfully}

	commonHTTP.SuccessResponse(w, r, tdr, h.logger)
}

// @Summary		Tracks of Artist
// @Tags		Artist
// @Description	All tracks of artist with chosen ID
// @Produce		json
// @Success		200		{object}	models.TrackTransfers 	"Show tracks"
// @Failure		400		{object}	http.Error			   	"Incorrect body"
// @Failure		500		{object}	http.Error			   	"Server error"
// @Router		/api/artists/{artistID}/tracks [get]
func (h *Handler) GetByArtist(w http.ResponseWriter, r *http.Request) {
	artistID, err := commonHTTP.GetArtistIDFromRequest(r)
	if err != nil {
		commonHTTP.ErrorResponseWithErrLogging(w, r,
			commonHTTP.InvalidURLParameter, http.StatusBadRequest, h.logger, err)
		return
	}

	tracks, err := h.trackServices.GetByArtist(r.Context(), artistID)
	if err != nil {
		var errNoSuchArtist *models.NoSuchArtistError
		if errors.As(err, &errNoSuchArtist) {
			commonHTTP.ErrorResponseWithErrLogging(w, r,
				artistNotFound, http.StatusBadRequest, h.logger, err)
			return
		}

		commonHTTP.ErrorResponseWithErrLogging(w, r,
			tracksGetServerError, http.StatusInternalServerError, h.logger, err)
		return
	}

	user, err := commonHTTP.GetUserFromRequest(r)
	if err != nil && !errors.Is(err, commonHTTP.ErrUnauthorized) {
		commonHTTP.ErrorResponseWithErrLogging(w, r,
			tracksGetServerError, http.StatusInternalServerError, h.logger, err)
		return
	}

	tt, err := models.TrackTransferFromList(r.Context(), tracks, user, h.trackServices.IsLiked,
		h.artistServices.IsLiked, h.artistServices.GetByTrack)
	if err != nil {
		commonHTTP.ErrorResponseWithErrLogging(w, r,
			tracksGetServerError, http.StatusInternalServerError, h.logger, err)
		return
	}

	commonHTTP.SuccessResponse(w, r, tt, h.logger)
}

// @Summary		Tracks of Playlist
// @Tags		Playlist
// @Description	All tracks of playlist with chosen ID
// @Produce		json
// @Success		200		{object}	models.TrackTransfers  "Show tracks"
// @Failure		400		{object}	http.Error			   "Incorrect body"
// @Failure		500		{object}	http.Error			   "Server error"
// @Router		/api/playlists/{playlistID}/tracks [get]
func (h *Handler) GetByPlaylist(w http.ResponseWriter, r *http.Request) {
	playlistID, err := commonHTTP.GetPlaylistIDFromRequest(r)
	if err != nil {
		commonHTTP.ErrorResponseWithErrLogging(w, r,
			commonHTTP.InvalidURLParameter, http.StatusBadRequest, h.logger, err)
		return
	}

	tracks, err := h.trackServices.GetByPlaylist(r.Context(), playlistID)
	if err != nil {
		var errNoSuchPlaylist *models.NoSuchPlaylistError
		if errors.As(err, &errNoSuchPlaylist) {
			commonHTTP.ErrorResponseWithErrLogging(w, r,
				playlistNotFound, http.StatusBadRequest, h.logger, err)
			return
		}

		commonHTTP.ErrorResponseWithErrLogging(w, r,
			tracksGetServerError, http.StatusInternalServerError, h.logger, err)
		return
	}

	user, err := commonHTTP.GetUserFromRequest(r)
	if err != nil && !errors.Is(err, commonHTTP.ErrUnauthorized) {
		commonHTTP.ErrorResponseWithErrLogging(w, r,
			tracksGetServerError, http.StatusInternalServerError, h.logger, err)
		return
	}

	tt, err := models.TrackTransferFromList(r.Context(), tracks, user, h.trackServices.IsLiked,
		h.artistServices.IsLiked, h.artistServices.GetByTrack)
	if err != nil {
		commonHTTP.ErrorResponseWithErrLogging(w, r,
			tracksGetServerError, http.StatusInternalServerError, h.logger, err)
		return
	}

	commonHTTP.SuccessResponse(w, r, tt, h.logger)
}

// @Summary		Tracks of Album
// @Tags		Album
// @Description	All tracks of album with chosen ID
// @Produce		json
// @Success		200		{object}	models.TrackTransfers  "Show tracks"
// @Failure		400		{object}	http.Error			   "Bad request"
// @Failure		500		{object}	http.Error			   "Server error"
// @Router		/api/albums/{albumID}/tracks [get]
func (h *Handler) GetByAlbum(w http.ResponseWriter, r *http.Request) {
	albumID, err := commonHTTP.GetAlbumIDFromRequest(r)
	if err != nil {
		commonHTTP.ErrorResponseWithErrLogging(w, r,
			commonHTTP.InvalidURLParameter, http.StatusBadRequest, h.logger, err)
		return
	}

	tracks, err := h.trackServices.GetByAlbum(r.Context(), albumID)
	if err != nil {
		var errNoSuchAlbum *models.NoSuchAlbumError
		if errors.As(err, &errNoSuchAlbum) {
			commonHTTP.ErrorResponseWithErrLogging(w, r,
				albumNotFound, http.StatusBadRequest, h.logger, err)
			return
		}

		commonHTTP.ErrorResponseWithErrLogging(w, r,
			tracksGetServerError, http.StatusInternalServerError, h.logger, err)
		return
	}

	user, err := commonHTTP.GetUserFromRequest(r)
	if err != nil && !errors.Is(err, commonHTTP.ErrUnauthorized) {
		commonHTTP.ErrorResponseWithErrLogging(w, r,
			tracksGetServerError, http.StatusInternalServerError, h.logger, err)
		return
	}

	tt, err := models.TrackTransferFromList(r.Context(), tracks, user, h.trackServices.IsLiked,
		h.artistServices.IsLiked, h.artistServices.GetByTrack)
	if err != nil {
		commonHTTP.ErrorResponseWithErrLogging(w, r,
			tracksGetServerError, http.StatusInternalServerError, h.logger, err)
		return
	}

	commonHTTP.SuccessResponse(w, r, tt, h.logger)
}

// @Summary		Track Feed
// @Tags		Feed
// @Description	Feed tracks
// @Produce		json
// @Success		200		{object}	models.TrackTransfers  "Tracks feed"
// @Failure		500		{object}	http.Error			   "Server error"
// @Router		/api/tracks/feed [get]
func (h *Handler) Feed(w http.ResponseWriter, r *http.Request) {
	tracks, err := h.trackServices.GetFeed(r.Context())
	if err != nil {
		commonHTTP.ErrorResponseWithErrLogging(w, r,
			tracksGetServerError, http.StatusInternalServerError, h.logger, err)
		return
	}

	user, err := commonHTTP.GetUserFromRequest(r)
	if err != nil && !errors.Is(err, commonHTTP.ErrUnauthorized) {
		commonHTTP.ErrorResponseWithErrLogging(w, r,
			tracksGetServerError, http.StatusInternalServerError, h.logger, err)
		return
	}

	tt, err := models.TrackTransferFromList(r.Context(), tracks, user, h.trackServices.IsLiked,
		h.artistServices.IsLiked, h.artistServices.GetByTrack)
	if err != nil {
		commonHTTP.ErrorResponseWithErrLogging(w, r,
			tracksGetServerError, http.StatusInternalServerError, h.logger, err)
		return
	}

	commonHTTP.SuccessResponse(w, r, tt, h.logger)
}

// @Summary      Favorite Tracks
// @Tags         Favorite
// @Description  Get ser's avorite tracks
// @Produce      json
// @Success      200    {object}  	models.TrackTransfers 	"Tracks got"
// @Failure		 400	{object}	http.Error				"Incorrect input"
// @Failure      401    {object}  	http.Error  			"Unauthorized user"
// @Failure      403    {object}  	http.Error  			"Forbidden user"
// @Failure      500    {object}  	http.Error  			"Server error"
// @Router       /api/users/{userID}/favorite/tracks [get]
func (h *Handler) GetFavorite(w http.ResponseWriter, r *http.Request) {
	user, err := commonHTTP.GetUserFromRequest(r)
	if err != nil {
		commonHTTP.ErrorResponseWithErrLogging(w, r,
			tracksGetServerError, http.StatusInternalServerError, h.logger, err)
		return
	}

	favTracks, err := h.trackServices.GetLikedByUser(r.Context(), user.ID)
	if err != nil {
		commonHTTP.ErrorResponseWithErrLogging(w, r,
			tracksGetServerError, http.StatusInternalServerError, h.logger, err)
		return
	}

	tt, err := models.TrackTransferFromList(r.Context(), favTracks, user, h.trackServices.IsLiked,
		h.artistServices.IsLiked, h.artistServices.GetByTrack)
	if err != nil {
		commonHTTP.ErrorResponseWithErrLogging(w, r,
			tracksGetServerError, http.StatusInternalServerError, h.logger, err)
		return
	}

	commonHTTP.SuccessResponse(w, r, tt, h.logger)
}

// @Summary		Set like
// @Tags		Track
// @Description	Set like by user to chosen track (add to favorite)
// @Produce		json
// @Success		200		{object}	trackLikeResponse	"Like set"
// @Failure		400		{object}	http.Error			"Client error"
// @Failure		401		{object}	http.Error  		"User unathorized"
// @Failure		500		{object}	http.Error			"Server error"
// @Router		/api/tracks/{trackID}/like [post]
func (h *Handler) Like(w http.ResponseWriter, r *http.Request) {
	trackID, err := commonHTTP.GetTrackIDFromRequest(r)
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

	notExisted, err := h.trackServices.SetLike(r.Context(), trackID, user.ID)
	if err != nil {
		var errNoSuchTrack *models.NoSuchTrackError
		if errors.As(err, &errNoSuchTrack) {
			commonHTTP.ErrorResponseWithErrLogging(w, r,
				trackNotFound, http.StatusBadRequest, h.logger, err)
			return
		}

		commonHTTP.ErrorResponseWithErrLogging(w, r,
			commonHTTP.SetLikeServerError, http.StatusInternalServerError, h.logger, err)
		return
	}

	tlr := trackLikeResponse{Status: commonHTTP.LikeSuccess}
	if !notExisted {
		tlr.Status = commonHTTP.LikeAlreadyExists
	}
	commonHTTP.SuccessResponse(w, r, tlr, h.logger)
}

// @Summary		Remove like
// @Tags		Track
// @Description	Remove like by user from chosen track (remove from favorite)
// @Produce		json
// @Success		200		{object}	trackLikeResponse	"Like removed"
// @Failure		400		{object}	http.Error			"Client error"
// @Failure		401		{object}	http.Error  		"User unathorized"
// @Failure		500		{object}	http.Error			"Server error"
// @Router		/api/tracks/{trackID}/unlike [post]
func (h *Handler) UnLike(w http.ResponseWriter, r *http.Request) {
	trackID, err := commonHTTP.GetTrackIDFromRequest(r)
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

	notExisted, err := h.trackServices.UnLike(r.Context(), trackID, user.ID)
	if err != nil {
		var errNoSuchTrack *models.NoSuchTrackError
		if errors.As(err, &errNoSuchTrack) {
			commonHTTP.ErrorResponseWithErrLogging(w, r,
				trackNotFound, http.StatusBadRequest, h.logger, err)
			return
		}

		commonHTTP.ErrorResponseWithErrLogging(w, r,
			commonHTTP.DeleteLikeServerError, http.StatusInternalServerError, h.logger, err)
		return
	}

	tlr := trackLikeResponse{Status: commonHTTP.UnLikeSuccess}
	if !notExisted {
		tlr.Status = commonHTTP.LikeDoesntExist
	}
	commonHTTP.SuccessResponse(w, r, tlr, h.logger)
}
