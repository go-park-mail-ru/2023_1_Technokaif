package http

import (
	"encoding/json"
	"errors"
	"net/http"
	"path/filepath"

	commonHttp "github.com/go-park-mail-ru/2023_1_Technokaif/internal/common/http"
	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/models"
	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/playlist"
	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/user"
	"github.com/go-park-mail-ru/2023_1_Technokaif/pkg/logger"
)

type Handler struct {
	playlistServices playlist.Usecase
	userServices     user.Usecase
	logger           logger.Logger
}

func NewHandler(pu playlist.Usecase, uu user.Usecase, l logger.Logger) *Handler {
	return &Handler{
		playlistServices: pu,
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
		commonHttp.ErrorResponseWithErrLogging(w, "unathorized", http.StatusUnauthorized, h.logger, err)
		return
	}

	var pci playlistCreateInput
	if err := json.NewDecoder(r.Body).Decode(&pci); err != nil {
		commonHttp.ErrorResponseWithErrLogging(w, "incorrect input body", http.StatusBadRequest, h.logger, err)
		return
	}

	if err := pci.validateAndEscape(); err != nil {
		h.logger.Infof("Creating playlist input validation failed: %s", err.Error())
		commonHttp.ErrorResponse(w, "incorrect input body", http.StatusBadRequest, h.logger)
		return
	}

	playlist := pci.ToPlaylist()

	playlistID, err := h.playlistServices.Create(playlist, pci.UsersID, user.ID)
	if err != nil {
		var errForbiddenUser *models.ForbiddenUserError
		if errors.As(err, &errForbiddenUser) {
			commonHttp.ErrorResponseWithErrLogging(w, "no rights to create playlist", http.StatusForbidden, h.logger, err)
			return
		}

		commonHttp.ErrorResponseWithErrLogging(w, "can't create playlist", http.StatusInternalServerError, h.logger, err)
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
		commonHttp.ErrorResponse(w, "invalid url parameter", http.StatusBadRequest, h.logger)
		return
	}

	playlist, err := h.playlistServices.GetByID(playlistID)
	if err != nil {
		var errNoSuchPlaylist *models.NoSuchPlaylistError
		if errors.As(err, &errNoSuchPlaylist) {
			commonHttp.ErrorResponseWithErrLogging(w, "no such playlist", http.StatusBadRequest, h.logger, err)
			return
		}

		commonHttp.ErrorResponseWithErrLogging(w, "can't get playlist", http.StatusInternalServerError, h.logger, err)
		return
	}

	user, err := commonHttp.GetUserFromRequest(r)
	if err != nil && !errors.Is(err, commonHttp.ErrUnauthorized) {
		commonHttp.ErrorResponseWithErrLogging(w, "can't get playlists", http.StatusInternalServerError, h.logger, err)
		return
	}

	resp, err := models.PlaylistTransferFromEntry(*playlist, user, h.playlistServices.IsLiked, h.userServices.GetByPlaylist)
	if err != nil {
		commonHttp.ErrorResponseWithErrLogging(w, "can't get playlist", http.StatusInternalServerError, h.logger, err)
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
		commonHttp.ErrorResponse(w, "invalid url parameter", http.StatusBadRequest, h.logger)
		return
	}

	user, err := commonHttp.GetUserFromRequest(r)
	if err != nil {
		commonHttp.ErrorResponseWithErrLogging(w, "unathorized", http.StatusUnauthorized, h.logger, err)
		return
	}

	if err := r.ParseMultipartForm(MaxCoverMemory); err != nil {
		commonHttp.ErrorResponseWithErrLogging(w, "invalid cover data", http.StatusBadRequest, h.logger, err)
		return
	}

	coverFile, coverHeader, err := r.FormFile(coverFormKey)
	if err != nil {
		commonHttp.ErrorResponseWithErrLogging(w, "invalid cover data", http.StatusBadRequest, h.logger, err)
		return
	}
	defer coverFile.Close()

	extension := filepath.Ext(coverHeader.Filename)

	err = h.playlistServices.UploadCover(playlistRequestID, user.ID, coverFile, extension)
	if err != nil {
		if errors.Is(err, h.playlistServices.UploadCoverWrongFormatError()) {
			commonHttp.ErrorResponseWithErrLogging(w, "invalid cover data type", http.StatusBadRequest, h.logger, err)
			return
		}

		var errForbiddenUser *models.ForbiddenUserError
		if errors.As(err, &errForbiddenUser) {
			commonHttp.ErrorResponseWithErrLogging(w, "no rights to upload cover", http.StatusForbidden, h.logger, err)
			return
		}

		var errNoSuchPlaylist *models.NoSuchPlaylistError
		if errors.As(err, &errNoSuchPlaylist) {
			commonHttp.ErrorResponseWithErrLogging(w, "no such playlist", http.StatusBadRequest, h.logger, err)
			return
		}

		commonHttp.ErrorResponseWithErrLogging(w, "can't upload cover", http.StatusInternalServerError, h.logger, err)
		return
	}

	resp := defaultResponse{Status: "ok"}

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
		commonHttp.ErrorResponse(w, "invalid url parameter", http.StatusBadRequest, h.logger)
		return
	}

	user, err := commonHttp.GetUserFromRequest(r)
	if err != nil {
		commonHttp.ErrorResponseWithErrLogging(w, "unathorized", http.StatusUnauthorized, h.logger, err)
		return
	}

	var pui playlistUpdateInput
	if err := json.NewDecoder(r.Body).Decode(&pui); err != nil {
		commonHttp.ErrorResponseWithErrLogging(w, "incorrect input body", http.StatusBadRequest, h.logger, err)
		return
	}

	if err := pui.validateAndEscape(); err != nil {
		h.logger.Infof("Creating playlist input validation failed: %s", err.Error())
		commonHttp.ErrorResponse(w, "incorrect input body", http.StatusBadRequest, h.logger)
		return
	}

	playlist := pui.ToPlaylist(playlistRequestID)

	err = h.playlistServices.UpdateInfoAndMembers(playlist, pui.UsersID, user.ID)
	if err != nil {
		var errForbiddenUser *models.ForbiddenUserError
		if errors.As(err, &errForbiddenUser) {
			commonHttp.ErrorResponseWithErrLogging(w, "no rights to update playlist", http.StatusForbidden, h.logger, err)
			return
		}

		var errNoSuchPlaylist *models.NoSuchPlaylistError
		if errors.As(err, &errNoSuchPlaylist) {
			commonHttp.ErrorResponseWithErrLogging(w, "no such playlist", http.StatusBadRequest, h.logger, err)
			return
		}

		commonHttp.ErrorResponseWithErrLogging(w, "can't update playlist", http.StatusInternalServerError, h.logger, err)
		return
	}

	dr := defaultResponse{Status: "ok"}

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
		commonHttp.ErrorResponse(w, "invalid url parameter", http.StatusBadRequest, h.logger)
		return
	}

	user, err := commonHttp.GetUserFromRequest(r)
	if err != nil {
		commonHttp.ErrorResponseWithErrLogging(w, "unathorized", http.StatusUnauthorized, h.logger, err)
		return
	}

	err = h.playlistServices.Delete(playlistID, user.ID)
	if err != nil {
		var errForbiddenUser *models.ForbiddenUserError
		if errors.As(err, &errForbiddenUser) {
			commonHttp.ErrorResponseWithErrLogging(w, "no rights to delete playlist", http.StatusForbidden, h.logger, err)
			return
		}

		var errNoSuchPlaylist *models.NoSuchPlaylistError
		if errors.As(err, &errNoSuchPlaylist) {
			commonHttp.ErrorResponseWithErrLogging(w, "no such playlist", http.StatusBadRequest, h.logger, err)
			return
		}

		commonHttp.ErrorResponseWithErrLogging(w, "can't delete playlist", http.StatusInternalServerError, h.logger, err)
		return
	}

	dr := defaultResponse{Status: "ok"}

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
		commonHttp.ErrorResponse(w, "invalid url parameter", http.StatusBadRequest, h.logger)
		return
	}

	playlists, err := h.playlistServices.GetByUser(userID)
	if err != nil {
		var errNoSuchUser *models.NoSuchUserError
		if errors.As(err, &errNoSuchUser) {
			commonHttp.ErrorResponseWithErrLogging(w, "no such user", http.StatusBadRequest, h.logger, err)
			return
		}

		commonHttp.ErrorResponseWithErrLogging(w, "can't get playlists", http.StatusInternalServerError, h.logger, err)
		return
	}

	user, err := commonHttp.GetUserFromRequest(r)
	if err != nil && !errors.Is(err, commonHttp.ErrUnauthorized) {
		commonHttp.ErrorResponseWithErrLogging(w, "can't get playlists", http.StatusInternalServerError, h.logger, err)
		return
	}

	pt, err := models.PlaylistTransferFromQuery(playlists, user, h.playlistServices.IsLiked, h.userServices.GetByPlaylist)
	if err != nil {
		commonHttp.ErrorResponseWithErrLogging(w, "can't get playlists", http.StatusInternalServerError, h.logger, err)
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
		commonHttp.ErrorResponse(w, "invalid url parameter", http.StatusBadRequest, h.logger)
		return
	}

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

	if err := h.playlistServices.AddTrack(trackID, playlistID, user.ID); err != nil {
		var errForbiddenUser *models.ForbiddenUserError
		if errors.As(err, &errForbiddenUser) {
			commonHttp.ErrorResponseWithErrLogging(w, "no rights to add track into playlist",
				http.StatusForbidden, h.logger, err)

			return
		}

		var errNoSuchPlaylist *models.NoSuchPlaylistError
		if errors.As(err, &errNoSuchPlaylist) {
			commonHttp.ErrorResponseWithErrLogging(w, "no such playlist", http.StatusBadRequest, h.logger, err)
			return
		}

		var errNoSuchTrack *models.NoSuchTrackError
		if errors.As(err, &errNoSuchTrack) {
			commonHttp.ErrorResponseWithErrLogging(w, "no such track", http.StatusBadRequest, h.logger, err)
			return
		}

		commonHttp.ErrorResponseWithErrLogging(w, "can't add track into playlist",
			http.StatusInternalServerError, h.logger, err)

		return
	}

	dr := defaultResponse{Status: "ok"}

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
		commonHttp.ErrorResponse(w, "invalid url parameter", http.StatusBadRequest, h.logger)
		return
	}

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

	if err := h.playlistServices.DeleteTrack(trackID, playlistID, user.ID); err != nil {
		var errForbiddenUser *models.ForbiddenUserError
		if errors.As(err, &errForbiddenUser) {
			commonHttp.ErrorResponseWithErrLogging(w, "no rights to delete track from playlist",
				http.StatusForbidden, h.logger, err)

			return
		}

		var errNoSuchPlaylist *models.NoSuchPlaylistError
		if errors.As(err, &errNoSuchPlaylist) {
			commonHttp.ErrorResponseWithErrLogging(w, "no such playlist", http.StatusBadRequest, h.logger, err)
			return
		}

		var errNoSuchTrack *models.NoSuchTrackError
		if errors.As(err, &errNoSuchTrack) {
			commonHttp.ErrorResponseWithErrLogging(w, "no such track", http.StatusBadRequest, h.logger, err)
			return
		}

		commonHttp.ErrorResponseWithErrLogging(w, "can't delete track from playlist",
			http.StatusInternalServerError, h.logger, err)

		return
	}

	dr := defaultResponse{Status: "ok"}

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
	playlists, err := h.playlistServices.GetFeed()
	if err != nil {
		commonHttp.ErrorResponseWithErrLogging(w, "can't get playlists", http.StatusInternalServerError, h.logger, err)
		return
	}

	user, err := commonHttp.GetUserFromRequest(r)
	if err != nil && !errors.Is(err, commonHttp.ErrUnauthorized) {
		commonHttp.ErrorResponseWithErrLogging(w, "can't get playlists", http.StatusInternalServerError, h.logger, err)
		return
	}

	resp, err := models.PlaylistTransferFromQuery(playlists, user, h.playlistServices.IsLiked, h.userServices.GetByPlaylist)
	if err != nil {
		commonHttp.ErrorResponseWithErrLogging(w, "can't get playlists", http.StatusInternalServerError, h.logger, err)
		return
	}

	commonHttp.SuccessResponse(w, resp, h.logger)
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
		commonHttp.ErrorResponse(w, "invalid url parameter", http.StatusBadRequest, h.logger)
		return
	}

	user, err := commonHttp.GetUserFromRequest(r)
	if err != nil {
		commonHttp.ErrorResponseWithErrLogging(w, "unathorized", http.StatusUnauthorized, h.logger, err)
		return
	}

	notExisted, err := h.playlistServices.SetLike(playlistID, user.ID)
	if err != nil {
		var errNoSuchPlaylist *models.NoSuchPlaylistError
		if errors.As(err, &errNoSuchPlaylist) {
			commonHttp.ErrorResponseWithErrLogging(w, "no such playlist", http.StatusBadRequest, h.logger, err)
			return
		}

		commonHttp.ErrorResponseWithErrLogging(w, "can't set like", http.StatusInternalServerError, h.logger, err)
		return
	}

	dr := defaultResponse{Status: "ok"}
	if !notExisted {
		dr.Status = "already liked"
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
		commonHttp.ErrorResponse(w, "invalid url parameter", http.StatusBadRequest, h.logger)
		return
	}

	user, err := commonHttp.GetUserFromRequest(r)
	if err != nil {
		commonHttp.ErrorResponseWithErrLogging(w, "unathorized", http.StatusUnauthorized, h.logger, err)
		return
	}

	notExisted, err := h.playlistServices.UnLike(playlistID, user.ID)
	if err != nil {
		var errNoSuchPlaylist *models.NoSuchPlaylistError
		if errors.As(err, &errNoSuchPlaylist) {
			commonHttp.ErrorResponseWithErrLogging(w, "no such playlist", http.StatusBadRequest, h.logger, err)
			return
		}

		commonHttp.ErrorResponseWithErrLogging(w, "can't remove like", http.StatusInternalServerError, h.logger, err)
		return
	}

	dr := defaultResponse{Status: "ok"}
	if !notExisted {
		dr.Status = "wasn't liked"
	}
	commonHttp.SuccessResponse(w, dr, h.logger)
}
