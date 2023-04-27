package http

import (
	"encoding/json"
	"errors"
	"net/http"
	"path/filepath"

	commonHttp "github.com/go-park-mail-ru/2023_1_Technokaif/internal/common/http"
	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/models"
	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/playlist"
	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/track"
	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/user"
	"github.com/go-park-mail-ru/2023_1_Technokaif/pkg/logger"
)

type Handler struct {
	playlistServices playlist.Usecase
	trackServices    track.Usecase
	userServices     user.Usecase
	logger           logger.Logger
}

func NewHandler(pu playlist.Usecase, tu track.Usecase, uu user.Usecase, l logger.Logger) *Handler {
	return &Handler{
		playlistServices: pu,
		trackServices:    tu,
		userServices:     uu,

		logger: l,
	}
}

// @Summary		Create Playlist
// @Tags		Playlist
// @Description	Create new playlist by sent object
// @Accept      json
// @Produce		json
// @Param		playlist body		playlistCreateInput	true	"Playlist info"
// @Success		200		 {object}	playlistCreateResponse	    "Playlist created"
// @Failure		400		 {object}	http.Error					"Incorrect input"
// @Failure		401		 {object}	http.Error  				"User unathorized"
// @Failure		403		 {object}	http.Error					"User hasn't rights"
// @Failure		500		 {object}	http.Error					"Server error"
// @Router		/api/playlists/ [post]
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	user, err := commonHttp.GetUserFromRequest(r)
	if err != nil {
		commonHttp.ErrorResponseWithErrLogging(w, commonHttp.UnathorizedUser, http.StatusUnauthorized, h.logger, err)
		return
	}

	var pci playlistCreateInput
	if err := json.NewDecoder(r.Body).Decode(&pci); err != nil {
		commonHttp.ErrorResponseWithErrLogging(w, commonHttp.IncorrectRequestBody, http.StatusBadRequest, h.logger, err)
		return
	}

	if err := pci.validateAndEscape(); err != nil {
		h.logger.Infof("Creating playlist input validation failed: %s", err.Error())
		commonHttp.ErrorResponse(w, commonHttp.IncorrectRequestBody, http.StatusBadRequest, h.logger)
		return
	}

	playlist := pci.ToPlaylist()

	playlistID, err := h.playlistServices.Create(r.Context(), playlist, pci.UsersID, user.ID)
	if err != nil {
		var errForbiddenUser *models.ForbiddenUserError
		if errors.As(err, &errForbiddenUser) {
			commonHttp.ErrorResponseWithErrLogging(w, playlistCreateNorights, http.StatusForbidden, h.logger, err)
			return
		}

		commonHttp.ErrorResponseWithErrLogging(w, playlistCreateServerError, http.StatusInternalServerError, h.logger, err)
		return
	}

	pcr := playlistCreateResponse{ID: playlistID}

	commonHttp.SuccessResponse(w, pcr, h.logger)
}

// @Summary		Get Playlist
// @Tags		Playlist
// @Description	Get playlist with chosen ID
// @Produce		json
// @Success		200		{object}	models.PlaylistTransfer	"Playlist got"
// @Failure		400		{object}	http.Error				"Incorrect input"
// @Failure		401		{object}	http.Error  			"User unathorized"
// @Failure		500		{object}	http.Error				"Server error"
// @Router		/api/playlists/{playlistID}/ [get]
func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	playlistID, err := commonHttp.GetPlaylistIDFromRequest(r)
	if err != nil {
		h.logger.Infof("Get playlist by id: %v", err)
		commonHttp.ErrorResponse(w, commonHttp.InvalidURLParameter, http.StatusBadRequest, h.logger)
		return
	}

	playlist, err := h.playlistServices.GetByID(r.Context(), playlistID)
	if err != nil {
		var errNoSuchPlaylist *models.NoSuchPlaylistError
		if errors.As(err, &errNoSuchPlaylist) {
			commonHttp.ErrorResponseWithErrLogging(w, playlistNotFound, http.StatusBadRequest, h.logger, err)
			return
		}

		commonHttp.ErrorResponseWithErrLogging(w, playlistGetServerError, http.StatusInternalServerError, h.logger, err)
		return
	}

	user, err := commonHttp.GetUserFromRequest(r)
	if err != nil && !errors.Is(err, commonHttp.ErrUnauthorized) {
		commonHttp.ErrorResponseWithErrLogging(w, playlistGetServerError, http.StatusInternalServerError, h.logger, err)
		return
	}

	resp, err := models.PlaylistTransferFromEntry(r.Context(), *playlist, user, h.playlistServices.IsLiked, h.userServices.GetByPlaylist)
	if err != nil {
		commonHttp.ErrorResponseWithErrLogging(w, playlistGetServerError, http.StatusInternalServerError, h.logger, err)
		return
	}

	commonHttp.SuccessResponse(w, resp, h.logger)
}

// @Summary      Upload Cover
// @Tags         Playlist
// @Description  Update playlist cover
// @Accept       multipart/form-data
// @Produce      json
// @Param		 cover  formData  file true 		"Cover file"
// @Success      200    {object}  defaultResponse	"Cover updated"
// @Failure      400    {object}  http.Error  		"Invalid form data"
// @Failure      401    {object}  http.Error  		"User Unathorized"
// @Failure      403    {object}  http.Error  		"User hasn't rights"
// @Failure      500    {object}  http.Error  		"Server error"
// @Router       /api/playlists/{playlistID}/cover [post]
func (h *Handler) UploadCover(w http.ResponseWriter, r *http.Request) {
	playlistRequestID, err := commonHttp.GetPlaylistIDFromRequest(r)
	if err != nil {
		h.logger.Infof("Get playlist's id: %v", err)
		commonHttp.ErrorResponse(w, commonHttp.InvalidURLParameter, http.StatusBadRequest, h.logger)
		return
	}

	user, err := commonHttp.GetUserFromRequest(r)
	if err != nil {
		commonHttp.ErrorResponseWithErrLogging(w, commonHttp.UnathorizedUser, http.StatusUnauthorized, h.logger, err)
		return
	}

	if err := r.ParseMultipartForm(MaxCoverMemory); err != nil {
		commonHttp.ErrorResponseWithErrLogging(w, playlistCoverInvalidData, http.StatusBadRequest, h.logger, err)
		return
	}

	coverFile, coverHeader, err := r.FormFile(coverFormKey)
	if err != nil {
		commonHttp.ErrorResponseWithErrLogging(w, playlistCoverInvalidData, http.StatusBadRequest, h.logger, err)
		return
	}
	defer coverFile.Close()

	extension := filepath.Ext(coverHeader.Filename)

	err = h.playlistServices.UploadCover(r.Context(), playlistRequestID, user.ID, coverFile, extension)
	if err != nil {
		var errCoverWrongFormat *models.CoverWrongFormatError
		if errors.As(err, &errCoverWrongFormat) {
			commonHttp.ErrorResponseWithErrLogging(w, playlistCoverInvalidDataType, http.StatusBadRequest, h.logger, err)
			return
		}

		var errForbiddenUser *models.ForbiddenUserError
		if errors.As(err, &errForbiddenUser) {
			commonHttp.ErrorResponseWithErrLogging(w, playlistCoverUploadNoRights, http.StatusForbidden, h.logger, err)
			return
		}

		var errNoSuchPlaylist *models.NoSuchPlaylistError
		if errors.As(err, &errNoSuchPlaylist) {
			commonHttp.ErrorResponseWithErrLogging(w, playlistNotFound, http.StatusBadRequest, h.logger, err)
			return
		}

		commonHttp.ErrorResponseWithErrLogging(w, playlistCoverServerError, http.StatusInternalServerError, h.logger, err)
		return
	}

	resp := defaultResponse{Status: playlistCoverUploadedSuccessfully}

	commonHttp.SuccessResponse(w, resp, h.logger)

}

// @Summary		Update Playlist
// @Tags		Playlist
// @Description	Update playlist
// @Accept		json
// @Produce		json
// @Param		playlist body		playlistUpdateInput	true	"Playlist info"
// @Success		200		{object}	defaultResponse				"Playlist updated"
// @Failure		400		{object}	http.Error					"Client error"
// @Failure		401		{object}	http.Error  				"User unathorized"
// @Failure		403		{object}	http.Error					"User hasn't rights"
// @Failure		500		{object}	http.Error					"Server error"
// @Router		/api/playlists/{playlistID}/update [post]
func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	playlistRequestID, err := commonHttp.GetPlaylistIDFromRequest(r)
	if err != nil {
		h.logger.Infof("Get playlist's id: %v", err)
		commonHttp.ErrorResponse(w, commonHttp.InvalidURLParameter, http.StatusBadRequest, h.logger)
		return
	}

	user, err := commonHttp.GetUserFromRequest(r)
	if err != nil {
		commonHttp.ErrorResponseWithErrLogging(w, commonHttp.UnathorizedUser, http.StatusUnauthorized, h.logger, err)
		return
	}

	var pui playlistUpdateInput
	if err := json.NewDecoder(r.Body).Decode(&pui); err != nil {
		commonHttp.ErrorResponseWithErrLogging(w, commonHttp.IncorrectRequestBody, http.StatusBadRequest, h.logger, err)
		return
	}

	if err := pui.validateAndEscape(); err != nil {
		h.logger.Infof("Creating playlist input validation failed: %s", err.Error())
		commonHttp.ErrorResponse(w, commonHttp.IncorrectRequestBody, http.StatusBadRequest, h.logger)
		return
	}

	playlist := pui.ToPlaylist(playlistRequestID)

	err = h.playlistServices.UpdateInfoAndMembers(r.Context(), playlist, pui.UsersID, user.ID)
	if err != nil {
		var errForbiddenUser *models.ForbiddenUserError
		if errors.As(err, &errForbiddenUser) {
			commonHttp.ErrorResponseWithErrLogging(w, playlistUpdateNoRights, http.StatusForbidden, h.logger, err)
			return
		}

		var errNoSuchPlaylist *models.NoSuchPlaylistError
		if errors.As(err, &errNoSuchPlaylist) {
			commonHttp.ErrorResponseWithErrLogging(w, playlistNotFound, http.StatusBadRequest, h.logger, err)
			return
		}

		commonHttp.ErrorResponseWithErrLogging(w, playlistUpdateServerError, http.StatusInternalServerError, h.logger, err)
		return
	}

	dr := defaultResponse{Status: playlistUpdatedSuccessfully}

	commonHttp.SuccessResponse(w, dr, h.logger)
}

// @Summary		Delete Playlist
// @Tags		Playlist
// @Description	Delete playlist with chosen ID
// @Produce		json
// @Success		200		{object}	defaultResponse	"Playlist deleted"
// @Failure		400		{object}	http.Error		"Client error"
// @Failure		401		{object}	http.Error  	"User unathorized"
// @Failure		403		{object}	http.Error		"User hasn't rights"
// @Failure		500		{object}	http.Error		"Server error"
// @Router		/api/playlists/{playlistID}/ [delete]
func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	playlistID, err := commonHttp.GetPlaylistIDFromRequest(r)
	if err != nil {
		h.logger.Infof("Get playlist's id: %v", err)
		commonHttp.ErrorResponse(w, commonHttp.InvalidURLParameter, http.StatusBadRequest, h.logger)
		return
	}

	user, err := commonHttp.GetUserFromRequest(r)
	if err != nil {
		commonHttp.ErrorResponseWithErrLogging(w, commonHttp.UnathorizedUser, http.StatusUnauthorized, h.logger, err)
		return
	}

	err = h.playlistServices.Delete(r.Context(), playlistID, user.ID)
	if err != nil {
		var errForbiddenUser *models.ForbiddenUserError
		if errors.As(err, &errForbiddenUser) {
			commonHttp.ErrorResponseWithErrLogging(w, playlistDeleteNoRights, http.StatusForbidden, h.logger, err)
			return
		}

		var errNoSuchPlaylist *models.NoSuchPlaylistError
		if errors.As(err, &errNoSuchPlaylist) {
			commonHttp.ErrorResponseWithErrLogging(w, playlistNotFound, http.StatusBadRequest, h.logger, err)
			return
		}

		commonHttp.ErrorResponseWithErrLogging(w, playlistDeleteServerError, http.StatusInternalServerError, h.logger, err)
		return
	}

	dr := defaultResponse{Status: playlistDeletedSuccessfully}

	commonHttp.SuccessResponse(w, dr, h.logger)
}

// @Summary		Playlists of User
// @Tags		User
// @Description	All playlists of user with chosen ID
// @Produce		json
// @Success		200		{object}	[]models.PlaylistTransfer	"Show playlists"
// @Failure		400		{object}	http.Error					"Client error"
// @Failure		500		{object}	http.Error					"Server error"
// @Router		/api/users/{userID}/playlists [get]
func (h *Handler) GetByUser(w http.ResponseWriter, r *http.Request) {
	userID, err := commonHttp.GetUserIDFromRequest(r)
	if err != nil {
		h.logger.Infof("Get user by id: %v", err)
		commonHttp.ErrorResponse(w, commonHttp.InvalidURLParameter, http.StatusBadRequest, h.logger)
		return
	}

	playlists, err := h.playlistServices.GetByUser(r.Context(), userID)
	if err != nil {
		var errNoSuchUser *models.NoSuchUserError
		if errors.As(err, &errNoSuchUser) {
			commonHttp.ErrorResponseWithErrLogging(w, userNotFound, http.StatusBadRequest, h.logger, err)
			return
		}

		commonHttp.ErrorResponseWithErrLogging(w, playlistsGetServerError, http.StatusInternalServerError, h.logger, err)
		return
	}

	user, err := commonHttp.GetUserFromRequest(r)
	if err != nil && !errors.Is(err, commonHttp.ErrUnauthorized) {
		commonHttp.ErrorResponseWithErrLogging(w, playlistsGetServerError, http.StatusInternalServerError, h.logger, err)
		return
	}

	pt, err := models.PlaylistTransferFromQuery(r.Context(), playlists, user, h.playlistServices.IsLiked, h.userServices.GetByPlaylist)
	if err != nil {
		commonHttp.ErrorResponseWithErrLogging(w, playlistsGetServerError, http.StatusInternalServerError, h.logger, err)
		return
	}

	commonHttp.SuccessResponse(w, pt, h.logger)
}

// @Summary		Add Track
// @Tags		Playlist
// @Description	Add track into playlist
// @Produce		json
// @Success		200		 {object}	playlistCreateResponse	    "Track added"
// @Failure		400		 {object}	http.Error					"Incorrect input"
// @Failure		401		 {object}	http.Error  				"User unathorized"
// @Failure		403		 {object}	http.Error					"User hasn't rights"
// @Failure		500		 {object}	http.Error					"Server error"
// @Router		/api/playlists/{playlistID}/tracks/{trackID} [post]
func (h *Handler) AddTrack(w http.ResponseWriter, r *http.Request) {
	playlistID, err := commonHttp.GetPlaylistIDFromRequest(r)
	if err != nil {
		h.logger.Infof("Get playlist by id: %v", err)
		commonHttp.ErrorResponse(w, commonHttp.InvalidURLParameter, http.StatusBadRequest, h.logger)
		return
	}

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

	if err := h.playlistServices.AddTrack(r.Context(), trackID, playlistID, user.ID); err != nil {
		var errForbiddenUser *models.ForbiddenUserError
		if errors.As(err, &errForbiddenUser) {
			commonHttp.ErrorResponseWithErrLogging(w, playlistAddTrackNoRights,
				http.StatusForbidden, h.logger, err)

			return
		}

		var errNoSuchPlaylist *models.NoSuchPlaylistError
		if errors.As(err, &errNoSuchPlaylist) {
			commonHttp.ErrorResponseWithErrLogging(w, playlistNotFound, http.StatusBadRequest, h.logger, err)
			return
		}

		var errNoSuchTrack *models.NoSuchTrackError
		if errors.As(err, &errNoSuchTrack) {
			commonHttp.ErrorResponseWithErrLogging(w, trackNotFound, http.StatusBadRequest, h.logger, err)
			return
		}

		commonHttp.ErrorResponseWithErrLogging(w, playlistAddTrackServerError,
			http.StatusInternalServerError, h.logger, err)

		return
	}

	dr := defaultResponse{Status: playlistTrackAddedSuccessfully}

	commonHttp.SuccessResponse(w, dr, h.logger)
}

// @Summary		Delete Track
// @Tags		Playlist
// @Description	Delete track from playlist
// @Produce		json
// @Success		200		 {object}	playlistCreateResponse	    "Track deleted"
// @Failure		400		 {object}	http.Error					"Incorrect input"
// @Failure		401		 {object}	http.Error  				"User unathorized"
// @Failure		403		 {object}	http.Error					"User hasn't rights"
// @Failure		500		 {object}	http.Error					"Server error"
// @Router		/api/playlists/{playlistID}/tracks/{trackID} [delete]
func (h *Handler) DeleteTrack(w http.ResponseWriter, r *http.Request) {
	playlistID, err := commonHttp.GetPlaylistIDFromRequest(r)
	if err != nil {
		h.logger.Infof("Get playlist by id: %v", err)
		commonHttp.ErrorResponse(w, commonHttp.InvalidURLParameter, http.StatusBadRequest, h.logger)
		return
	}

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

	if err := h.playlistServices.DeleteTrack(r.Context(), trackID, playlistID, user.ID); err != nil {
		var errForbiddenUser *models.ForbiddenUserError
		if errors.As(err, &errForbiddenUser) {
			commonHttp.ErrorResponseWithErrLogging(w, playlistDeleteTrackNoRights,
				http.StatusForbidden, h.logger, err)

			return
		}

		var errNoSuchPlaylist *models.NoSuchPlaylistError
		if errors.As(err, &errNoSuchPlaylist) {
			commonHttp.ErrorResponseWithErrLogging(w, playlistNotFound, http.StatusBadRequest, h.logger, err)
			return
		}

		var errNoSuchTrack *models.NoSuchTrackError
		if errors.As(err, &errNoSuchTrack) {
			commonHttp.ErrorResponseWithErrLogging(w, trackNotFound, http.StatusBadRequest, h.logger, err)
			return
		}

		commonHttp.ErrorResponseWithErrLogging(w, playlistDeleteTrackServerError,
			http.StatusInternalServerError, h.logger, err)

		return
	}

	dr := defaultResponse{Status: playlistTrackDeletedSuccessfully}

	commonHttp.SuccessResponse(w, dr, h.logger)
}

// @Summary		Playlist Feed
// @Tags		Feed
// @Description	Feed playlists
// @Produce		json
// @Success		200		{object}	[]models.PlaylistTransfer	 "Playlist feed"
// @Failure		500		{object}	http.Error "Server error"
// @Router		/api/playlists/feed [get]
func (h *Handler) Feed(w http.ResponseWriter, r *http.Request) {
	playlists, err := h.playlistServices.GetFeed(r.Context())
	if err != nil {
		commonHttp.ErrorResponseWithErrLogging(w, playlistsGetServerError, http.StatusInternalServerError, h.logger, err)
		return
	}

	user, err := commonHttp.GetUserFromRequest(r)
	if err != nil && !errors.Is(err, commonHttp.ErrUnauthorized) {
		commonHttp.ErrorResponseWithErrLogging(w, playlistsGetServerError, http.StatusInternalServerError, h.logger, err)
		return
	}

	resp, err := models.PlaylistTransferFromQuery(r.Context(), playlists, user, h.playlistServices.IsLiked, h.userServices.GetByPlaylist)
	if err != nil {
		commonHttp.ErrorResponseWithErrLogging(w, playlistsGetServerError, http.StatusInternalServerError, h.logger, err)
		return
	}

	commonHttp.SuccessResponse(w, resp, h.logger)
}

// @Summary      Favorite Playlists
// @Tags         Favorite
// @Description  Get user's favorite playlists
// @Produce      application/json
// @Success      200    {object}  	[]models.PlaylistTransfer 	"Playlists got"
// @Failure		 400	{object}	http.Error					"Incorrect input"
// @Failure      401    {object}  	http.Error  				"Unauthorized user"
// @Failure      403    {object}  	http.Error  				"Forbidden user"
// @Failure      500    {object}  	http.Error  				"Server error"
// @Router       /api/users/{userID}/favorite/playlists [get]
func (h *Handler) GetFavorite(w http.ResponseWriter, r *http.Request) {
	user, err := commonHttp.GetUserFromRequest(r)
	if err != nil {
		commonHttp.ErrorResponseWithErrLogging(w, playlistsGetServerError, http.StatusInternalServerError, h.logger, err)
		return
	}

	favPlaylists, err := h.playlistServices.GetLikedByUser(r.Context(), user.ID)
	if err != nil {
		commonHttp.ErrorResponseWithErrLogging(w, playlistsGetServerError, http.StatusInternalServerError, h.logger, err)
		return
	}

	at, err := models.PlaylistTransferFromQuery(r.Context(), favPlaylists, user,
		h.trackServices.IsLiked, h.userServices.GetByPlaylist)
	if err != nil {
		commonHttp.ErrorResponseWithErrLogging(w, playlistsGetServerError, http.StatusInternalServerError, h.logger, err)
		return
	}

	commonHttp.SuccessResponse(w, at, h.logger)
}

// @Summary		Set like
// @Tags		Playlist
// @Description	Set like by user to chosen playlist (add to favorite)
// @Produce		json
// @Success		200		{object}	defaultResponse	"Like set"
// @Failure		400		{object}	http.Error		"Client error"
// @Failure		401		{object}	http.Error  	"User unathorized"
// @Failure		500		{object}	http.Error		"Server error"
// @Router		/api/playlists/{playlistID}/like [post]
func (h *Handler) Like(w http.ResponseWriter, r *http.Request) {
	playlistID, err := commonHttp.GetPlaylistIDFromRequest(r)
	if err != nil {
		h.logger.Infof("Get playlist by id: %v", err)
		commonHttp.ErrorResponse(w, commonHttp.InvalidURLParameter, http.StatusBadRequest, h.logger)
		return
	}

	user, err := commonHttp.GetUserFromRequest(r)
	if err != nil {
		commonHttp.ErrorResponseWithErrLogging(w, commonHttp.UnathorizedUser, http.StatusUnauthorized, h.logger, err)
		return
	}

	notExisted, err := h.playlistServices.SetLike(r.Context(), playlistID, user.ID)
	if err != nil {
		var errNoSuchPlaylist *models.NoSuchPlaylistError
		if errors.As(err, &errNoSuchPlaylist) {
			commonHttp.ErrorResponseWithErrLogging(w, playlistNotFound, http.StatusBadRequest, h.logger, err)
			return
		}

		commonHttp.ErrorResponseWithErrLogging(w, commonHttp.SetLikeServerError, http.StatusInternalServerError, h.logger, err)
		return
	}

	dr := defaultResponse{Status: commonHttp.LikeSuccess}
	if !notExisted {
		dr.Status = commonHttp.LikeAlreadyExists
	}
	commonHttp.SuccessResponse(w, dr, h.logger)
}

// @Summary		Remove like
// @Tags		Playlist
// @Description	Remove like by user from chosen playlist (remove from favorite)
// @Produce		json
// @Success		200		{object}	defaultResponse	"Like removed"
// @Failure		400		{object}	http.Error		"Client error"
// @Failure		401		{object}	http.Error  	"User unathorized"
// @Failure		500		{object}	http.Error		"Server error"
// @Router		/api/playlists/{playlistID}/unlike [post]
func (h *Handler) UnLike(w http.ResponseWriter, r *http.Request) {
	playlistID, err := commonHttp.GetPlaylistIDFromRequest(r)
	if err != nil {
		h.logger.Infof("Get playlist by id: %v", err)
		commonHttp.ErrorResponse(w, commonHttp.InvalidURLParameter, http.StatusBadRequest, h.logger)
		return
	}

	user, err := commonHttp.GetUserFromRequest(r)
	if err != nil {
		commonHttp.ErrorResponseWithErrLogging(w, commonHttp.UnathorizedUser, http.StatusUnauthorized, h.logger, err)
		return
	}

	notExisted, err := h.playlistServices.UnLike(r.Context(), playlistID, user.ID)
	if err != nil {
		var errNoSuchPlaylist *models.NoSuchPlaylistError
		if errors.As(err, &errNoSuchPlaylist) {
			commonHttp.ErrorResponseWithErrLogging(w, playlistNotFound, http.StatusBadRequest, h.logger, err)
			return
		}

		commonHttp.ErrorResponseWithErrLogging(w, commonHttp.DeleteLikeServerError, http.StatusInternalServerError, h.logger, err)
		return
	}

	dr := defaultResponse{Status: commonHttp.UnLikeSuccess}
	if !notExisted {
		dr.Status = commonHttp.LikeDoesntExist
	}
	commonHttp.SuccessResponse(w, dr, h.logger)
}
