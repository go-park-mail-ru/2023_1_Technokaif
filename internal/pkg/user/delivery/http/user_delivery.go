package http

import (
	"encoding/json"
	"errors"
	"net/http"
	"path/filepath"

	commonHttp "github.com/go-park-mail-ru/2023_1_Technokaif/internal/common/http"
	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/models"
	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/album"
	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/artist"
	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/playlist"
	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/track"
	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/user"
	"github.com/go-park-mail-ru/2023_1_Technokaif/pkg/logger"
)

type Handler struct {
	userServices     user.Usecase
	trackServices    track.Usecase
	albumServices    album.Usecase
	playlistServices playlist.Usecase
	artistServices   artist.Usecase

	logger logger.Logger
}

func NewHandler(uu user.Usecase, tu track.Usecase, alu album.Usecase, aru artist.Usecase, l logger.Logger) *Handler {
	return &Handler{
		userServices:   uu,
		trackServices:  tu,
		albumServices:  alu,
		artistServices: aru,
		logger:         l,
	}
}

// @Summary		Get User
// @Tags		User
// @Description	Get user with chosen ID
// @Produce		json
// @Success		200		{object}	models.UserTransfer "User got"
// @Failure		400		{object}	http.Error			"Client error"
// @Failure     401    	{object}  	http.Error  		"Unauthorized user"
// @Failure     403    	{object}  	http.Error  		"Forbidden user"
// @Failure     500    	{object}  	http.Error  		"Can't get user"
// @Router	    /api/users/{userID}/ [get]
func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	user, err := commonHttp.GetUserFromRequest(r)
	if err != nil {
		commonHttp.ErrorResponseWithErrLogging(w, "can't get user", http.StatusInternalServerError, h.logger, err)
		return
	}

	ut := models.UserTransferFromEntry(*user)

	commonHttp.SuccessResponse(w, ut, h.logger)
}

// @Summary      Update Info
// @Tags         User
// @Description  Update info about user
// @Accept       json
// @Produce      json
// @Param		 user	body	  userInfoInput	true		"User info"
// @Success      200    {object}  userUploadAvatarResponse 	"User info updated"
// @Failure      400    {object}  http.Error  			   	"Invalid input"
// @Failure      401    {object}  http.Error  			   	"User Unathorized"
// @Failure      403    {object}  http.Error  			   	"User hasn't rights"
// @Failure      500    {object}  http.Error  			   	"Can't change user info"
// @Router       /api/users/{userID}/update [post]
func (h *Handler) UpdateInfo(w http.ResponseWriter, r *http.Request) {
	user, err := commonHttp.GetUserFromRequest(r)
	if err != nil {
		commonHttp.ErrorResponseWithErrLogging(w, "can't change user info", http.StatusInternalServerError, h.logger, err)
		return
	}

	var userInfo userInfoInput
	if err := json.NewDecoder(r.Body).Decode(&userInfo); err != nil {
		commonHttp.ErrorResponseWithErrLogging(w, "incorrect input body", http.StatusBadRequest, h.logger, err)
		return
	}

	if err := userInfo.validate(); err != nil {
		commonHttp.ErrorResponseWithErrLogging(w, "incorrect input body", http.StatusBadRequest, h.logger, err)
		return
	}

	if err := h.userServices.UpdateInfo(userInfo.ToUser(user)); err != nil {
		var errNoSuchUser *models.NoSuchUserError
		if errors.As(err, &errNoSuchUser) {
			commonHttp.ErrorResponseWithErrLogging(w, "no user to update", http.StatusBadRequest, h.logger, err)
			return
		}

		commonHttp.ErrorResponseWithErrLogging(w, "can't change user info", http.StatusInternalServerError, h.logger, err)
		return
	}

	uuir := userChangeInfoResponse{Status: "ok"}

	commonHttp.SuccessResponse(w, uuir, h.logger)
}

// @Summary      Upload Avatar
// @Tags         User
// @Description  Update user avatar
// @Accept       multipart/form-data
// @Produce      json
// @Param		 avatar formData file true 				   "Avatar file"
// @Success      200    {object}  userUploadAvatarResponse "Avatar updated"
// @Failure      400    {object}  http.Error  			   "Invalid form data"
// @Failure      401    {object}  http.Error  			   "User Unathorized"
// @Failure      403    {object}  http.Error  			   "User hasn't rights"
// @Failure      500    {object}  http.Error  			   "Server error"
// @Router       /api/users/{userID}/avatar [post]
func (h *Handler) UploadAvatar(w http.ResponseWriter, r *http.Request) {
	user, err := commonHttp.GetUserFromRequest(r)
	if err != nil {
		commonHttp.ErrorResponseWithErrLogging(w, "server error", http.StatusInternalServerError, h.logger, err)
		return
	}

	if err := r.ParseMultipartForm(MaxAvatarMemory); err != nil {
		commonHttp.ErrorResponseWithErrLogging(w, "invalid avatar data", http.StatusBadRequest, h.logger, err)
		return
	}

	avatarFile, avatarHeader, err := r.FormFile(avatarFormKey)
	if err != nil {
		commonHttp.ErrorResponseWithErrLogging(w, "invalid avatar data", http.StatusBadRequest, h.logger, err)
		return
	}
	defer avatarFile.Close()

	extension := filepath.Ext(avatarHeader.Filename)
	err = h.userServices.UploadAvatar(user, avatarFile, extension)
	if err != nil {
		if errors.Is(err, h.userServices.UploadAvatarWrongFormatError()) {
			commonHttp.ErrorResponseWithErrLogging(w, "invalid avatar data type", http.StatusBadRequest, h.logger, err)
			return
		}

		commonHttp.ErrorResponseWithErrLogging(w, "can't upload avatar", http.StatusInternalServerError, h.logger, err)
		return
	}

	uuar := userUploadAvatarResponse{Status: "ok"}

	commonHttp.SuccessResponse(w, uuar, h.logger)
}

// @Summary      Favorite Tracks
// @Tags         User
// @Description  Get ser's avorite racks
// @Produce      application/json
// @Success      200    {object}  	[]models.TrackTransfer 	"Tracks got"
// @Failure		 400	{object}	http.Error				"Incorrect input"
// @Failure      401    {object}  	http.Error  			"Unauthorized user"
// @Failure      403    {object}  	http.Error  			"Forbidden user"
// @Failure      500    {object}  	http.Error  			"Server error"
// @Router       /api/users/{userID}/favorite/tracks [get]
func (h *Handler) GetFavouriteTracks(w http.ResponseWriter, r *http.Request) {
	user, err := commonHttp.GetUserFromRequest(r)
	if err != nil {
		commonHttp.ErrorResponseWithErrLogging(w, "server error", http.StatusInternalServerError, h.logger, err)
		return
	}

	favTracks, err := h.trackServices.GetLikedByUser(user.ID)
	if err != nil {
		commonHttp.ErrorResponseWithErrLogging(w, "can't get favorite tracks", http.StatusInternalServerError, h.logger, err)
		return
	}

	tt, err := models.TrackTransferFromQuery(favTracks, user, h.trackServices.IsLiked, h.artistServices.GetByTrack)
	if err != nil {
		commonHttp.ErrorResponseWithErrLogging(w, "can't get favorite tracks", http.StatusInternalServerError, h.logger, err)
		return
	}

	commonHttp.SuccessResponse(w, tt, h.logger)
}

// @Summary      Favorite Albums
// @Tags         User
// @Description  Get user's favorite albums
// @Produce      application/json
// @Success      200    {object}  	[]models.AlbumTransfer 	"Albums got"
// @Failure		 400	{object}	http.Error				"Incorrect input"
// @Failure      401    {object}  	http.Error  			"Unauthorized user"
// @Failure      403    {object}  	http.Error  			"Forbidden user"
// @Failure      500    {object}  	http.Error  			"Server error"
// @Router       /api/users/{userID}/favorite/albums [get]
func (h *Handler) GetFavouriteAlbums(w http.ResponseWriter, r *http.Request) {
	user, err := commonHttp.GetUserFromRequest(r)
	if err != nil {
		commonHttp.ErrorResponseWithErrLogging(w, "server error", http.StatusInternalServerError, h.logger, err)
		return
	}

	favAlbums, err := h.albumServices.GetLikedByUser(user.ID)
	if err != nil {
		commonHttp.ErrorResponseWithErrLogging(w, "can't get favorite albums", http.StatusInternalServerError, h.logger, err)
		return
	}

	at, err := models.AlbumTransferFromQuery(favAlbums, h.artistServices.GetByAlbum)
	if err != nil {
		commonHttp.ErrorResponseWithErrLogging(w, "can't get favorite albums", http.StatusInternalServerError, h.logger, err)
		return
	}

	commonHttp.SuccessResponse(w, at, h.logger)
}

// @Summary      Favorite Playlists
// @Tags         User
// @Description  Get user's favorite playlists
// @Produce      application/json
// @Success      200    {object}  	[]models.PlaylistTransfer 	"Playlists got"
// @Failure		 400	{object}	http.Error					"Incorrect input"
// @Failure      401    {object}  	http.Error  				"Unauthorized user"
// @Failure      403    {object}  	http.Error  				"Forbidden user"
// @Failure      500    {object}  	http.Error  				"Server error"
// @Router       /api/users/{userID}/favorite/playlists [get]
func (h *Handler) GetFavouritePlaylists(w http.ResponseWriter, r *http.Request) {
	user, err := commonHttp.GetUserFromRequest(r)
	if err != nil {
		commonHttp.ErrorResponseWithErrLogging(w, "server error", http.StatusInternalServerError, h.logger, err)
		return
	}

	favPlaylists, err := h.playlistServices.GetLikedByUser(user.ID)
	if err != nil {
		commonHttp.ErrorResponseWithErrLogging(w, "can't get favorite playlists", http.StatusInternalServerError, h.logger, err)
		return
	}

	at, err := models.PlaylistTransferFromQuery(favPlaylists, h.userServices.GetByPlaylist)
	if err != nil {
		commonHttp.ErrorResponseWithErrLogging(w, "can't get favorite playlists", http.StatusInternalServerError, h.logger, err)
		return
	}

	commonHttp.SuccessResponse(w, at, h.logger)
}

// @Summary      Favorite Artists
// @Tags         User
// @Description  Get user's favorite artists
// @Produce      application/json
// @Success      200    {object}  	[]models.ArtistTransfer "Artists got"
// @Failure		 400	{object}	http.Error				"Incorrect input"
// @Failure      401    {object}  	http.Error  			"Unauthorized user"
// @Failure      403    {object}  	http.Error  			"Forbidden user"
// @Failure      500    {object}  	http.Error  			"Server error"
// @Router       /api/users/{userID}/favorite/artists [get]
func (h *Handler) GetFavouriteArtists(w http.ResponseWriter, r *http.Request) {
	user, err := commonHttp.GetUserFromRequest(r)
	if err != nil {
		commonHttp.ErrorResponseWithErrLogging(w, "server error", http.StatusInternalServerError, h.logger, err)
		return
	}

	favArtists, err := h.artistServices.GetLikedByUser(user.ID)
	if err != nil {
		commonHttp.ErrorResponseWithErrLogging(w, "can't get favorite artists", http.StatusInternalServerError, h.logger, err)
		return
	}

	at := models.ArtistTransferFromQuery(favArtists)

	commonHttp.SuccessResponse(w, at, h.logger)
}
