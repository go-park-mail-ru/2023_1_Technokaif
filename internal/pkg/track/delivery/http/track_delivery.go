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
		logger:         l,
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
		commonHttp.ErrorResponseWithErrLogging(w, "unathorized", http.StatusUnauthorized, h.logger, err)
		return
	}

	var tci trackCreateInput
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&tci); err != nil {
		commonHttp.ErrorResponseWithErrLogging(w, "incorrect input body", http.StatusBadRequest, h.logger, err)
		return
	}

	if err := tci.validate(); err != nil {
		h.logger.Infof("track create input validation failed: %s", err.Error())
		commonHttp.ErrorResponse(w, "incorrect input body", http.StatusBadRequest, h.logger)
		return
	}

	track := tci.ToTrack()

	trackID, err := h.trackServices.Create(track, tci.ArtistsID, user.ID)
	var errForbiddenUser *models.ForbiddenUserError
	if err != nil {
		if errors.As(err, &errForbiddenUser) {
			commonHttp.ErrorResponseWithErrLogging(w, "no rights to create track", http.StatusForbidden, h.logger, err)
			return
		}

		commonHttp.ErrorResponseWithErrLogging(w, "can't create track", http.StatusInternalServerError, h.logger, err)
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
// @Failure		500		{object}	http.Error			 "Server error"
// @Router		/api/tracks/{trackID}/ [get]
func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	trackID, err := commonHttp.GetTrackIDFromRequest(r)
	if err != nil {
		h.logger.Infof("get track by id: %v", err)
		commonHttp.ErrorResponse(w, "invalid url parameter", http.StatusBadRequest, h.logger)
		return
	}

	track, err := h.trackServices.GetByID(uint32(trackID))
	var errNoSuchTrack *models.NoSuchTrackError
	if errors.As(err, &errNoSuchTrack) {
		commonHttp.ErrorResponseWithErrLogging(w, "no such track", http.StatusBadRequest, h.logger, err)
		return
	}
	if err != nil {
		commonHttp.ErrorResponseWithErrLogging(w, "can't get track", http.StatusInternalServerError, h.logger, err)
		return
	}

	tt, err := models.TrackTransferFromEntry(*track, h.artistServices.GetByTrack)
	if err != nil {
		commonHttp.ErrorResponseWithErrLogging(w, "can't get track", http.StatusInternalServerError, h.logger, err)
		return
	}

	commonHttp.SuccessResponse(w, tt, h.logger)
}

// swaggermock
func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	// ...
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
		commonHttp.ErrorResponse(w, "invalid url parameter", http.StatusBadRequest, h.logger)
		return
	}

	user, err := commonHttp.GetUserFromRequest(r)
	if err != nil {
		commonHttp.ErrorResponseWithErrLogging(w, "unathorized", http.StatusUnauthorized, h.logger, err)
		return
	}

	err = h.trackServices.Delete(trackID, user.ID)
	var errNoSuchTrack *models.NoSuchTrackError
	var errForbiddenUser *models.ForbiddenUserError
	if err != nil {
		if errors.As(err, &errForbiddenUser) {
			commonHttp.ErrorResponseWithErrLogging(w, "no rights to delete track", http.StatusForbidden, h.logger, err)
			return
		}
		if errors.As(err, &errNoSuchTrack) {
			commonHttp.ErrorResponseWithErrLogging(w, "no such track", http.StatusBadRequest, h.logger, err)
			return
		}

		commonHttp.ErrorResponseWithErrLogging(w, "can't delete track", http.StatusInternalServerError, h.logger, err)
		return
	}

	tdr := trackDeleteResponse{Status: "ok"}

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
		commonHttp.ErrorResponse(w, "invalid url parameter", http.StatusBadRequest, h.logger)
		return
	}

	tracks, err := h.trackServices.GetByArtist(artistID)
	var errNoSuchArtist *models.NoSuchArtistError
	if errors.As(err, &errNoSuchArtist) {
		commonHttp.ErrorResponseWithErrLogging(w, "no such artist", http.StatusBadRequest, h.logger, err)
		return
	}
	if err != nil {
		commonHttp.ErrorResponseWithErrLogging(w, "error while getting artist tracks", http.StatusInternalServerError, h.logger, err)
		return
	}

	tt, err := models.TrackTransferFromQuery(tracks, h.artistServices.GetByTrack)
	if err != nil {
		commonHttp.ErrorResponseWithErrLogging(w, "error while getting artist tracks", http.StatusInternalServerError, h.logger, err)
		return
	}

	commonHttp.SuccessResponse(w, tt, h.logger)
}

// @Summary		Tracks of Album
// @Tags		Album
// @Description	All tracks of album with chosen ID
// @Produce		json
// @Success		200		{object}	[]models.TrackTransfer "Show tracks"
// @Failure		400		{object}	http.Error			   "Client error"
// @Failure		500		{object}	http.Error			   "Server error"
// @Router		/api/albums/{albumID}/tracks [get]
func (h *Handler) GetByAlbum(w http.ResponseWriter, r *http.Request) {
	albumID, err := commonHttp.GetAlbumIDFromRequest(r)
	if err != nil {
		h.logger.Infof("Get by album: %v", err)
		commonHttp.ErrorResponse(w, "invalid url parameter", http.StatusBadRequest, h.logger)
		return
	}

	tracks, err := h.trackServices.GetByAlbum(albumID)
	var errNoSuchAlbum *models.NoSuchAlbumError
	if errors.As(err, &errNoSuchAlbum) {
		commonHttp.ErrorResponseWithErrLogging(w, "no such album", http.StatusBadRequest, h.logger, err)
		return
	}
	if err != nil {
		commonHttp.ErrorResponseWithErrLogging(w, "error while getting album tracks", http.StatusInternalServerError, h.logger, err)
		return
	}

	tt, err := models.TrackTransferFromQuery(tracks, h.artistServices.GetByTrack)
	if err != nil {
		commonHttp.ErrorResponseWithErrLogging(w, "error while getting artist tracks", http.StatusInternalServerError, h.logger, err)
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
		commonHttp.ErrorResponseWithErrLogging(w, "error while getting tracks", http.StatusInternalServerError, h.logger, err)
		return
	}

	tt, err := models.TrackTransferFromQuery(tracks, h.artistServices.GetByTrack)
	if err != nil {
		commonHttp.ErrorResponseWithErrLogging(w, "error while getting tracks", http.StatusInternalServerError, h.logger, err)
		return
	}

	commonHttp.SuccessResponse(w, tt, h.logger)
}

// swaggermock
func (h *Handler) Like(w http.ResponseWriter, r *http.Request) {
	trackID, err := commonHttp.GetTrackIDFromRequest(r)
	if err != nil {
		h.logger.Infof("Get track by id: %v", err)
		commonHttp.ErrorResponse(w, "invalid url parameter", http.StatusBadRequest, h.logger)
		return
	}

	user, err := commonHttp.GetUserFromRequest(r)
	if err != nil {
		commonHttp.ErrorResponseWithErrLogging(w, "unathorized", http.StatusUnauthorized, h.logger, err)
		return
	}

	notExists, err := h.trackServices.SetLike(trackID, user.ID)
	if err != nil {
		var errNoSuchTrack *models.NoSuchTrackError
		if errors.As(err, &errNoSuchTrack) {
			commonHttp.ErrorResponseWithErrLogging(w, "no such track", http.StatusBadRequest, h.logger, err)
			return
		} else {
			commonHttp.ErrorResponseWithErrLogging(w, "error while setting like", http.StatusInternalServerError, h.logger, err)
			return
		}
	}

	if notExists {
		resp := trackLikeResponse{Status: "ok"}
		commonHttp.SuccessResponse(w, resp, h.logger)
	} else {
		resp := trackLikeResponse{Status: "already liked"}
		commonHttp.SuccessResponse(w, resp, h.logger)
	}
}

// swaggermock
func (h *Handler) UnLike(w http.ResponseWriter, r *http.Request) {
	trackID, err := commonHttp.GetTrackIDFromRequest(r)
	if err != nil {
		h.logger.Infof("Get track by id: %v", err)
		commonHttp.ErrorResponse(w, "invalid url parameter", http.StatusBadRequest, h.logger)
		return
	}

	user, err := commonHttp.GetUserFromRequest(r)
	if err != nil {
		commonHttp.ErrorResponseWithErrLogging(w, "unathorized", http.StatusUnauthorized, h.logger, err)
		return
	}

	notExisted, err := h.trackServices.UnLike(trackID, user.ID)
	if err != nil {
		var errNoSuchTrack *models.NoSuchTrackError
		if errors.As(err, &errNoSuchTrack) {
			commonHttp.ErrorResponseWithErrLogging(w, "no such track", http.StatusBadRequest, h.logger, err)
			return
		} else {
			commonHttp.ErrorResponseWithErrLogging(w, "error while removing like", http.StatusInternalServerError, h.logger, err)
			return
		}
	}

	if notExisted {
		resp := trackLikeResponse{Status: "ok"}
		commonHttp.SuccessResponse(w, resp, h.logger)
	} else {
		resp := trackLikeResponse{Status: "already disliked"}
		commonHttp.SuccessResponse(w, resp, h.logger)
	}
}

// swaggermock
func (h *Handler) GetRecord(w http.ResponseWriter, r *http.Request) {
	trackID, err := commonHttp.GetTrackIDFromRequest(r)
	if err != nil {
		h.logger.Infof("Get track by id: %v", err)
		commonHttp.ErrorResponse(w, "invalid url parameter", http.StatusBadRequest, h.logger)
		return
	}

	if _, err := commonHttp.GetUserFromRequest(r); err != nil {
		commonHttp.ErrorResponseWithErrLogging(w, "unathorized", http.StatusUnauthorized, h.logger, err)
		return
	}

	track, err := h.trackServices.GetByID(uint32(trackID))
	if err != nil {
		var errNoSuchTrack *models.NoSuchTrackError
		if errors.As(err, &errNoSuchTrack) {
			commonHttp.ErrorResponseWithErrLogging(w, "no such track", http.StatusBadRequest, h.logger, err)
			return
		}

		commonHttp.ErrorResponseWithErrLogging(w, "can't get track", http.StatusInternalServerError, h.logger, err)
		return
	}

	
	http.ServeFile(w, r, "." + track.RecordSrc)
}
