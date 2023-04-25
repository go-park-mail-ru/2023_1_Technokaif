package http

import (
	"encoding/json"
	"errors"
	"net/http"
	"path/filepath"

	commonHttp "github.com/go-park-mail-ru/2023_1_Technokaif/internal/common/http"
	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/models"
	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/user"
	"github.com/go-park-mail-ru/2023_1_Technokaif/pkg/logger"
)

type Handler struct {
	userServices user.Usecase

	logger logger.Logger
}

func NewHandler(uu user.Usecase, l logger.Logger) *Handler {
	return &Handler{
		userServices: uu,
		logger:       l,
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
		commonHttp.ErrorResponseWithErrLogging(w, userGetServerError, http.StatusInternalServerError, h.logger, err)
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
		commonHttp.ErrorResponseWithErrLogging(w, userUpdateInfoServerError, http.StatusInternalServerError, h.logger, err)
		return
	}

	var userInfo userInfoInput
	if err := json.NewDecoder(r.Body).Decode(&userInfo); err != nil {
		commonHttp.ErrorResponseWithErrLogging(w, commonHttp.IncorrectRequestBody, http.StatusBadRequest, h.logger, err)
		return
	}

	if err := userInfo.validateAndEscape(); err != nil {
		commonHttp.ErrorResponseWithErrLogging(w, commonHttp.IncorrectRequestBody, http.StatusBadRequest, h.logger, err)
		return
	}

	if err := h.userServices.UpdateInfo(r.Context(), userInfo.ToUser(user)); err != nil {
		var errNoSuchUser *models.NoSuchUserError
		if errors.As(err, &errNoSuchUser) {
			commonHttp.ErrorResponseWithErrLogging(w, userNotFound, http.StatusBadRequest, h.logger, err)
			return
		}

		commonHttp.ErrorResponseWithErrLogging(w, userUpdateInfoServerError, http.StatusInternalServerError, h.logger, err)
		return
	}

	uuir := userChangeInfoResponse{Status: userUpdatedInfoSuccessfully}

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
		commonHttp.ErrorResponseWithErrLogging(w, userAvatarUploadServerError, http.StatusUnauthorized, h.logger, err)
		return
	}

	if err := r.ParseMultipartForm(MaxAvatarMemory); err != nil {
		commonHttp.ErrorResponseWithErrLogging(w, userAvatarUploadInvalidData, http.StatusBadRequest, h.logger, err)
		return
	}

	avatarFile, avatarHeader, err := r.FormFile(avatarFormKey)
	if err != nil {
		commonHttp.ErrorResponseWithErrLogging(w, userAvatarUploadInvalidData, http.StatusBadRequest, h.logger, err)
		return
	}
	defer avatarFile.Close()

	extension := filepath.Ext(avatarHeader.Filename)
	err = h.userServices.UploadAvatar(r.Context(), user.ID, avatarFile, extension)
	if err != nil {
		if errors.Is(err, h.userServices.UploadAvatarWrongFormatError()) {
			commonHttp.ErrorResponseWithErrLogging(w, userAvatarUploadInvalidDataType, http.StatusBadRequest, h.logger, err)
			return
		}

		commonHttp.ErrorResponseWithErrLogging(w, userAvatarUploadServerError, http.StatusInternalServerError, h.logger, err)
		return
	}

	uuar := userUploadAvatarResponse{Status: userAvatarUploadedSuccessfully}

	commonHttp.SuccessResponse(w, uuar, h.logger)
}
