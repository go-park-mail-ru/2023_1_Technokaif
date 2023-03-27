package http

import (
	"errors"
	"net/http"
	"strings"

	commonHttp "github.com/go-park-mail-ru/2023_1_Technokaif/internal/common/http"
	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/models"
	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/album"
	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/artist"
	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/track"
	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/user"
	"github.com/go-park-mail-ru/2023_1_Technokaif/pkg/logger"
)

type Handler struct {
	userServices   user.Usecase
	trackServices  track.Usecase
	albumServices  album.Usecase
	artistServices artist.Usecase
	logger         logger.Logger
}

func NewHandler(uu user.Usecase, tu track.Usecase, alu album.Usecase, aru artist.Usecase, l logger.Logger) *Handler {
	return &Handler{
		userServices: uu,
		trackServices: tu,
		albumServices: alu,
		artistServices: aru,
		logger:       l,
	}
}

// @Summary		Get User
// @Tags		User
// @Description	Get user with chosen ID
// @Produce		json
// @Success		200		{object}	models.UserTransfer "User got"
// @Failure		400		{object}	http.Error	"Client error"
// @Failure     401    	{object}  	http.Error  "Unauthorized user"
// @Failure     403    	{object}  	http.Error  "Forbidden user"
// @Failure     500    	{object}  	http.Error  "Server error"
// @Router	    /api/users/{userID}/ [get]
func (h *Handler) Read(w http.ResponseWriter, r *http.Request) {
	user, err := h.checkUserAuthAndResponce(w, r)
	if err != nil {
		return
	}

	ut := models.UserTransferFromUser(*user)

	commonHttp.SuccessResponse(w, ut, h.logger)
}

// @Summary      Upload Avatar
// @Tags         User
// @Description  Update user avatar
// @Accept       multipart/form-data
// @Produce      application/json
// @Param		 avatar formData file true "avatar file"
// @Success      200    {object}  userUploadAvatarResponse "Avatar updated"
// @Failure      400    {object}  http.Error  "Invalid form data"
// @Failure      401    {object}  http.Error  "Unauthorized user"
// @Failure      403    {object}  http.Error  "Forbidden user"
// @Failure      500    {object}  http.Error  "Server error"
// @Router       /api/users/{userID}/avatar [post]
func (h *Handler) UploadAvatar(w http.ResponseWriter, r *http.Request) {
	user, err := h.checkUserAuthAndResponce(w, r)
	if err != nil {
		return
	}

	if err := r.ParseMultipartForm(maxAvatarMemory); err != nil {
		h.logger.Info(err.Error())
		commonHttp.ErrorResponse(w, "invalid avatar data", http.StatusBadRequest, h.logger)
		return
	}
	r.Body = http.MaxBytesReader(w, r.Body, maxAvatarMemory)
	avatarFile, avatarHeader, err := r.FormFile("avatar")
	if err != nil {
		h.logger.Info(err.Error())
		commonHttp.ErrorResponse(w, "invalid avatar data", http.StatusBadRequest, h.logger)
		return
	}
	defer avatarFile.Close()

	fileNameParts := strings.Split(avatarHeader.Filename, ".")
	extension := fileNameParts[len(fileNameParts)-1]
	err = h.userServices.UploadAvatar(user, avatarFile, extension)
	if errors.Is(err, h.userServices.UploadAvatarWrongFormatError()) {
		h.logger.Info(err.Error())
		commonHttp.ErrorResponse(w, "invalid avatar data type", http.StatusBadRequest, h.logger)
		return
	} else if err != nil {
		h.logger.Error(err.Error())
		commonHttp.ErrorResponse(w, "can't upload avatar", http.StatusInternalServerError, h.logger)
		return
	}

	uuar := userUploadAvatarResponse{Status: "ok"}

	commonHttp.SuccessResponse(w, uuar, h.logger)
}

func (h *Handler) ReadFavouriteTracks(w http.ResponseWriter, r *http.Request) {
	user, err := h.checkUserAuthAndResponce(w, r)
	if err != nil {
		return
	}

	favTracks, err := h.trackServices.GetLikedByUser(user.ID)
	if err != nil {
		h.logger.Error(err.Error())
		commonHttp.ErrorResponse(w, "error while getting favourite tracks", http.StatusInternalServerError, h.logger)
		return
	}

	tt, err := models.TrackTransferFromQuery(favTracks, h.artistServices.GetByTrack)
	if err != nil {
		h.logger.Error(err.Error())
		commonHttp.ErrorResponse(w, "error while getting favourite tracks", http.StatusInternalServerError, h.logger)
		return
	}

	commonHttp.SuccessResponse(w, tt, h.logger)
}

func (h *Handler) ReadFavouriteAlbums(w http.ResponseWriter, r *http.Request) {
	user, err := h.checkUserAuthAndResponce(w, r)
	if err != nil {
		return
	}

	favAlbums, err := h.albumServices.GetLikedByUser(user.ID)
	if err != nil {
		h.logger.Error(err.Error())
		commonHttp.ErrorResponse(w, "error while getting favourite albums", http.StatusInternalServerError, h.logger)
		return
	}

	resp, err := models.AlbumTransferFromQuery(favAlbums, h.artistServices.GetByAlbum)
	if err != nil {
		h.logger.Error(err.Error())
		commonHttp.ErrorResponse(w, "error while getting favourite albums", http.StatusInternalServerError, h.logger)
		return
	}

	commonHttp.SuccessResponse(w, resp, h.logger)
}

func (h *Handler) ReadFavouriteArtists(w http.ResponseWriter, r *http.Request) {
	user, err := h.checkUserAuthAndResponce(w, r)
	if err != nil {
		return
	}

	favArtists, err := h.artistServices.GetLikedByUser(user.ID)
	if err != nil {
		h.logger.Error(err.Error())
		commonHttp.ErrorResponse(w, "error while getting favourite albums", http.StatusInternalServerError, h.logger)
		return
	}

	artistsTransfer := models.ArtistTransferFromQuery(favArtists)

	commonHttp.SuccessResponse(w, artistsTransfer, h.logger)
}

// help func
func (h *Handler) checkUserAuthAndResponce(w http.ResponseWriter, r *http.Request) (*models.User, error) {
	authFailedError := errors.New("user auth failed")

	user, err := commonHttp.GetUserFromRequest(r)
	if err != nil {
		h.logger.Infof("unathorized user: %v", err)
		commonHttp.ErrorResponse(w, "invalid token", http.StatusUnauthorized, h.logger)
		return nil, authFailedError
	}

	urlID, err := commonHttp.GetUserIDFromRequest(r)
	if err != nil {
		h.logger.Infof("invalid url parameter: %v", err.Error())
		commonHttp.ErrorResponse(w, "invalid url parameter", http.StatusBadRequest, h.logger)
		return nil, authFailedError
	}

	if urlID != user.ID {
		h.logger.Infof("forbidden user with id #%d", urlID)
		commonHttp.ErrorResponse(w, "invalid user", http.StatusForbidden, h.logger)
		return nil, authFailedError
	}

	return user, nil
}
