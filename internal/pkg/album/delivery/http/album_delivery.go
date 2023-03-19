package album_delivery

import (
	"encoding/json"
	"net/http"

	commonHttp "github.com/go-park-mail-ru/2023_1_Technokaif/internal/common/http"
	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/logger"
	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/album"
)

type AlbumHandler struct {
	services album.AlbumUsecase
	logger   logger.Logger
}

func NewAuthHandler(u album.AlbumUsecase, l logger.Logger) *AlbumHandler {
	return &AlbumHandler{
		services: u,
		logger:   l,
	}
}

//	@Summary		Album Feed
//	@Tags			album
//	@Description	Feed albums for user
//	@Accept			json
//	@Produce		json
//	@Param			user	body		models.User		true	"user info"
//	@Success		200		{object}	signUpResponse	"User created"
//	@Failure		400		{object}	errorResponse	"Incorrect input"
//	@Failure		500		{object}	errorResponse	"Server DB error"
//	@Router			/api/auth/signup [post]
func (ah *AlbumHandler) Feed(w http.ResponseWriter, r *http.Request) {
	albums, err := ah.services.GetFeed()
	if err != nil {
		ah.logger.Error(err.Error())
		commonHttp.ErrorResponse(w, "error while getting albums", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "json/application; charset=utf-8")

	encoder := json.NewEncoder(w)
	if err := encoder.Encode(&albums); err != nil {
		ah.logger.Error(err.Error())
		commonHttp.ErrorResponse(w, "can't encode response into json", http.StatusInternalServerError)
		return
	}
}
