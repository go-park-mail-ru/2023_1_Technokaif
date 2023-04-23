package http

import (
	"encoding/json"
	"errors"
	"net/http"

	commonHttp "github.com/go-park-mail-ru/2023_1_Technokaif/internal/common/http"
	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/models"
	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/artist"
	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/track"
	"github.com/go-park-mail-ru/2023_1_Technokaif/pkg/logger"
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
	user, err := commonHttp.GetUserFromRequest(r)
	if err != nil {
		commonHttp.ErrorResponseWithErrLogging(w, commonHttp.UnathorizedUser, http.StatusUnauthorized, h.logger, err)
		return
	}

	var tci trackCreateInput
	if err := json.NewDecoder(r.Body).Decode(&tci); err != nil {
		commonHttp.ErrorResponseWithErrLogging(w, commonHttp.IncorrectRequestBody, http.StatusBadRequest, h.logger, err)
		return
	}

	if err := tci.validateAndEscape(); err != nil {
		h.logger.Infof("track create input validation failed: %s", err.Error())
		commonHttp.ErrorResponse(w, commonHttp.IncorrectRequestBody, http.StatusBadRequest, h.logger)
		return
	}

	track := tci.ToTrack()

	trackID, err := h.trackServices.Create(track, tci.ArtistsID, user.ID)
	if err != nil {
		var errForbiddenUser *models.ForbiddenUserError
		if errors.As(err, &errForbiddenUser) {
			commonHttp.ErrorResponseWithErrLogging(w, trackCreateNorights, http.StatusForbidden, h.logger, err)
			return
		}

		commonHttp.ErrorResponseWithErrLogging(w, trackCreateServerError, http.StatusInternalServerError, h.logger, err)
		return
	}

	tcr := trackCreateResponse{ID: trackID}

	commonHttp.SuccessResponse(w, tcr, h.logger)
}

// @Summary		Get Track
// @Tags		Track
// @Description	Get track with chosen ID
// @Produce		json
// @Success		200		{object}	models.TrackTransfer "Track got"
// @Failure		400		{object}	http.Error			 "Client error"
// @Failure		401		{object}	http.Error  		 "User unathorized"
// @Failure		500		{object}	http.Error			 "Server error"
// @Router		/api/tracks/{trackID}/ [get]
func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	trackID, err := commonHttp.GetTrackIDFromRequest(r)
	if err != nil {
		h.logger.Infof("get track by id: %v", err)
		commonHttp.ErrorResponse(w, commonHttp.InvalidURLParameter, http.StatusBadRequest, h.logger)
		return
	}

	track, err := h.trackServices.GetByID(trackID)
	if err != nil {
		var errNoSuchTrack *models.NoSuchTrackError
		if errors.As(err, &errNoSuchTrack) {
			commonHttp.ErrorResponseWithErrLogging(w, trackNotFound, http.StatusBadRequest, h.logger, err)
			return
		}

		commonHttp.ErrorResponseWithErrLogging(w, trackGetServerError, http.StatusInternalServerError, h.logger, err)
		return
	}

	user, err := commonHttp.GetUserFromRequest(r)
	if err != nil && !errors.Is(err, commonHttp.ErrUnauthorized) {
		commonHttp.ErrorResponseWithErrLogging(w, trackGetServerError, http.StatusInternalServerError, h.logger, err)
		return
	}

	tt, err := models.TrackTransferFromEntry(*track, user, h.trackServices.IsLiked,
		h.artistServices.IsLiked, h.artistServices.GetByTrack)
	if err != nil {
		commonHttp.ErrorResponseWithErrLogging(w, trackGetServerError, http.StatusInternalServerError, h.logger, err)
		return
	}

	commonHttp.SuccessResponse(w, tt, h.logger)
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
	trackID, err := commonHttp.GetTrackIDFromRequest(r)
	if err != nil {
		h.logger.Infof("Get track by id: %v", err)
		commonHttp.ErrorResponse(w, commonHttp.InvalidURLParameter, http.StatusBadRequest, h.logger)
		return
	}

	user, err := commonHttp.GetUserFromRequest(r)
	if err != nil {
		commonHttp.ErrorResponseWithErrLogging(w, commonHttp.UnathorizedUser, http.StatusUnauthorized, h.logger, err)
		return
	}

	err = h.trackServices.Delete(trackID, user.ID)
	if err != nil {
		var errForbiddenUser *models.ForbiddenUserError
		if errors.As(err, &errForbiddenUser) {
			commonHttp.ErrorResponseWithErrLogging(w, trackDeleteNoRights, http.StatusForbidden, h.logger, err)
			return
		}

		var errNoSuchTrack *models.NoSuchTrackError
		if errors.As(err, &errNoSuchTrack) {
			commonHttp.ErrorResponseWithErrLogging(w, trackNotFound, http.StatusBadRequest, h.logger, err)
			return
		}

		commonHttp.ErrorResponseWithErrLogging(w, trackDeleteServerError, http.StatusInternalServerError, h.logger, err)
		return
	}

	tdr := trackDeleteResponse{Status: trackDeletedSuccessfully}

	commonHttp.SuccessResponse(w, tdr, h.logger)
}

// @Summary		Tracks of Artist
// @Tags		Artist
// @Description	All tracks of artist with chosen ID
// @Produce		json
// @Success		200		{object}	[]models.TrackTransfer "Show tracks"
// @Failure		400		{object}	http.Error			   "Incorrect body"
// @Failure		500		{object}	http.Error			   "Server error"
// @Router		/api/artists/{artistID}/tracks [get]
func (h *Handler) GetByArtist(w http.ResponseWriter, r *http.Request) {
	artistID, err := commonHttp.GetArtistIDFromRequest(r)
	if err != nil {
		h.logger.Infof("Get by artist: %v", err)
		commonHttp.ErrorResponse(w, commonHttp.InvalidURLParameter, http.StatusBadRequest, h.logger)
		return
	}

	tracks, err := h.trackServices.GetByArtist(artistID)
	if err != nil {
		var errNoSuchArtist *models.NoSuchArtistError
		if errors.As(err, &errNoSuchArtist) {
			commonHttp.ErrorResponseWithErrLogging(w, artistNotFound, http.StatusBadRequest, h.logger, err)
			return
		}

		commonHttp.ErrorResponseWithErrLogging(w, tracksGetServerError, http.StatusInternalServerError, h.logger, err)
		return
	}

	user, err := commonHttp.GetUserFromRequest(r)
	if err != nil && !errors.Is(err, commonHttp.ErrUnauthorized) {
		commonHttp.ErrorResponseWithErrLogging(w, tracksGetServerError, http.StatusInternalServerError, h.logger, err)
		return
	}

	tt, err := models.TrackTransferFromQuery(tracks, user, h.trackServices.IsLiked,
		h.artistServices.IsLiked, h.artistServices.GetByTrack)
	if err != nil {
		commonHttp.ErrorResponseWithErrLogging(w, tracksGetServerError, http.StatusInternalServerError, h.logger, err)
		return
	}

	commonHttp.SuccessResponse(w, tt, h.logger)
}

// @Summary		Tracks of Playlist
// @Tags		Playlist
// @Description	All tracks of playlist with chosen ID
// @Produce		json
// @Success		200		{object}	[]models.TrackTransfer "Show tracks"
// @Failure		400		{object}	http.Error			   "Incorrect body"
// @Failure		500		{object}	http.Error			   "Server error"
// @Router		/api/playlists/{playlistID}/tracks [get]
func (h *Handler) GetByPlaylist(w http.ResponseWriter, r *http.Request) {
	playlistID, err := commonHttp.GetPlaylistIDFromRequest(r)
	if err != nil {
		h.logger.Infof("Get by playlist: %v", err)
		commonHttp.ErrorResponse(w, commonHttp.InvalidURLParameter, http.StatusBadRequest, h.logger)
		return
	}

	tracks, err := h.trackServices.GetByPlaylist(playlistID)
	if err != nil {
		var errNoSuchPlaylist *models.NoSuchPlaylistError
		if errors.As(err, &errNoSuchPlaylist) {
			commonHttp.ErrorResponseWithErrLogging(w, playlistNotFound, http.StatusBadRequest, h.logger, err)
			return
		}

		commonHttp.ErrorResponseWithErrLogging(w, tracksGetServerError, http.StatusInternalServerError, h.logger, err)
		return
	}

	user, err := commonHttp.GetUserFromRequest(r)
	if err != nil && !errors.Is(err, commonHttp.ErrUnauthorized) {
		commonHttp.ErrorResponseWithErrLogging(w, tracksGetServerError, http.StatusInternalServerError, h.logger, err)
		return
	}

	tt, err := models.TrackTransferFromQuery(tracks, user, h.trackServices.IsLiked,
		h.artistServices.IsLiked, h.artistServices.GetByTrack)
	if err != nil {
		commonHttp.ErrorResponseWithErrLogging(w, tracksGetServerError, http.StatusInternalServerError, h.logger, err)
		return
	}

	commonHttp.SuccessResponse(w, tt, h.logger)
}

// @Summary		Tracks of Album
// @Tags		Album
// @Description	All tracks of album with chosen ID
// @Produce		json
// @Success		200		{object}	[]models.TrackTransfer "Show tracks"
// @Failure		400		{object}	http.Error			   "Bad request"
// @Failure		500		{object}	http.Error			   "Server error"
// @Router		/api/albums/{albumID}/tracks [get]
func (h *Handler) GetByAlbum(w http.ResponseWriter, r *http.Request) {
	albumID, err := commonHttp.GetAlbumIDFromRequest(r)
	if err != nil {
		h.logger.Infof("Get by album: %v", err)
		commonHttp.ErrorResponse(w, commonHttp.InvalidURLParameter, http.StatusBadRequest, h.logger)
		return
	}

	tracks, err := h.trackServices.GetByAlbum(albumID)
	if err != nil {
		var errNoSuchAlbum *models.NoSuchAlbumError
		if errors.As(err, &errNoSuchAlbum) {
			commonHttp.ErrorResponseWithErrLogging(w, albumNotFound, http.StatusBadRequest, h.logger, err)
			return
		}

		commonHttp.ErrorResponseWithErrLogging(w, tracksGetServerError, http.StatusInternalServerError, h.logger, err)
		return
	}

	user, err := commonHttp.GetUserFromRequest(r)
	if err != nil && !errors.Is(err, commonHttp.ErrUnauthorized) {
		commonHttp.ErrorResponseWithErrLogging(w, tracksGetServerError, http.StatusInternalServerError, h.logger, err)
		return
	}

	tt, err := models.TrackTransferFromQuery(tracks, user, h.trackServices.IsLiked,
		h.artistServices.IsLiked, h.artistServices.GetByTrack)
	if err != nil {
		commonHttp.ErrorResponseWithErrLogging(w, tracksGetServerError, http.StatusInternalServerError, h.logger, err)
		return
	}

	commonHttp.SuccessResponse(w, tt, h.logger)
}

// @Summary		Track Feed
// @Tags		Feed
// @Description	Feed tracks
// @Produce		json
// @Success		200		{object}	[]models.TrackTransfer "Tracks feed"
// @Failure		500		{object}	http.Error			   "Server error"
// @Router		/api/tracks/feed [get]
func (h *Handler) Feed(w http.ResponseWriter, r *http.Request) {
	tracks, err := h.trackServices.GetFeed()
	if err != nil {
		commonHttp.ErrorResponseWithErrLogging(w, tracksGetServerError, http.StatusInternalServerError, h.logger, err)
		return
	}

	user, err := commonHttp.GetUserFromRequest(r)
	if err != nil && !errors.Is(err, commonHttp.ErrUnauthorized) {
		commonHttp.ErrorResponseWithErrLogging(w, tracksGetServerError, http.StatusInternalServerError, h.logger, err)
		return
	}

	tt, err := models.TrackTransferFromQuery(tracks, user, h.trackServices.IsLiked,
		h.artistServices.IsLiked, h.artistServices.GetByTrack)
	if err != nil {
		commonHttp.ErrorResponseWithErrLogging(w, tracksGetServerError, http.StatusInternalServerError, h.logger, err)
		return
	}

	commonHttp.SuccessResponse(w, tt, h.logger)
}

// @Summary      Favorite Tracks
// @Tags         Favorite
// @Description  Get ser's avorite tracks
// @Produce      application/json
// @Success      200    {object}  	[]models.TrackTransfer 	"Tracks got"
// @Failure		 400	{object}	http.Error				"Incorrect input"
// @Failure      401    {object}  	http.Error  			"Unauthorized user"
// @Failure      403    {object}  	http.Error  			"Forbidden user"
// @Failure      500    {object}  	http.Error  			"Server error"
// @Router       /api/users/{userID}/favorite/tracks [get]
func (h *Handler) GetFavorite(w http.ResponseWriter, r *http.Request) {
	user, err := commonHttp.GetUserFromRequest(r)
	if err != nil {
		commonHttp.ErrorResponseWithErrLogging(w, tracksGetServerError, http.StatusInternalServerError, h.logger, err)
		return
	}

	favTracks, err := h.trackServices.GetLikedByUser(user.ID)
	if err != nil {
		commonHttp.ErrorResponseWithErrLogging(w, tracksGetServerError, http.StatusInternalServerError, h.logger, err)
		return
	}

	tt, err := models.TrackTransferFromQuery(favTracks, user, h.trackServices.IsLiked,
		h.artistServices.IsLiked, h.artistServices.GetByTrack)
	if err != nil {
		commonHttp.ErrorResponseWithErrLogging(w, tracksGetServerError, http.StatusInternalServerError, h.logger, err)
		return
	}

	commonHttp.SuccessResponse(w, tt, h.logger)
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
	trackID, err := commonHttp.GetTrackIDFromRequest(r)
	if err != nil {
		h.logger.Infof("Get track by id: %v", err)
		commonHttp.ErrorResponse(w, commonHttp.InvalidURLParameter, http.StatusBadRequest, h.logger)
		return
	}

	user, err := commonHttp.GetUserFromRequest(r)
	if err != nil {
		commonHttp.ErrorResponseWithErrLogging(w, commonHttp.UnathorizedUser, http.StatusUnauthorized, h.logger, err)
		return
	}

	notExisted, err := h.trackServices.SetLike(trackID, user.ID)
	if err != nil {
		var errNoSuchTrack *models.NoSuchTrackError
		if errors.As(err, &errNoSuchTrack) {
			commonHttp.ErrorResponseWithErrLogging(w, trackNotFound, http.StatusBadRequest, h.logger, err)
			return
		}

		commonHttp.ErrorResponseWithErrLogging(w, commonHttp.SetLikeServerError, http.StatusInternalServerError, h.logger, err)
		return
	}

	tlr := trackLikeResponse{Status: commonHttp.LikeSuccess}
	if !notExisted {
		tlr.Status = commonHttp.LikeAlreadyExists
	}
	commonHttp.SuccessResponse(w, tlr, h.logger)
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
	trackID, err := commonHttp.GetTrackIDFromRequest(r)
	if err != nil {
		h.logger.Infof("Get track by id: %v", err)
		commonHttp.ErrorResponse(w, commonHttp.InvalidURLParameter, http.StatusBadRequest, h.logger)
		return
	}

	user, err := commonHttp.GetUserFromRequest(r)
	if err != nil {
		commonHttp.ErrorResponseWithErrLogging(w, commonHttp.UnathorizedUser, http.StatusUnauthorized, h.logger, err)
		return
	}

	notExisted, err := h.trackServices.UnLike(trackID, user.ID)
	if err != nil {
		var errNoSuchTrack *models.NoSuchTrackError
		if errors.As(err, &errNoSuchTrack) {
			commonHttp.ErrorResponseWithErrLogging(w, trackNotFound, http.StatusBadRequest, h.logger, err)
			return
		}

		commonHttp.ErrorResponseWithErrLogging(w, commonHttp.DeleteLikeServerError, http.StatusInternalServerError, h.logger, err)
		return
	}

	tlr := trackLikeResponse{Status: commonHttp.UnLikeSuccess}
	if !notExisted {
		tlr.Status = commonHttp.LikeDoesntExist
	}
	commonHttp.SuccessResponse(w, tlr, h.logger)
}
