package delivery

import (
	"encoding/json"
	"net/http"

	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/models"
	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/album"
	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/artist"
	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/playlist"
	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/search"
	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/track"
	"github.com/go-park-mail-ru/2023_1_Technokaif/pkg/logger"

	commonHTTP "github.com/go-park-mail-ru/2023_1_Technokaif/internal/common/http"
)

type Handler struct {
	searchServices   search.Usecase
	albumServices    album.Usecase
	artistServices   artist.Usecase
	trackServices    track.Usecase
	playlistServices playlist.Usecase
	logger           logger.Logger
}

func NewHandler(su search.Usecase, alu album.Usecase,
	aru artist.Usecase, pu playlist.Usecase, l logger.Logger) *Handler {

	return &Handler{
		searchServices: su,
		logger:         l,
	}
}

// swaggermock
func (h *Handler) FindAlbums(w http.ResponseWriter, r *http.Request) {
	user, err := commonHTTP.GetUserFromRequest(r)
	if err != nil {
		commonHTTP.ErrorResponseWithErrLogging(w, r,
			commonHTTP.UnathorizedUser, http.StatusUnauthorized, h.logger, err)
		return
	}

	var sr SearchRequest
	if err := json.NewDecoder(r.Body).Decode(&sr); err != nil {
		commonHTTP.ErrorResponseWithErrLogging(w, r,
			commonHTTP.IncorrectRequestBody, http.StatusBadRequest, h.logger, err)
		return
	}

	if err := sr.validate(); err != nil {
		commonHTTP.ErrorResponseWithErrLogging(w, r,
			commonHTTP.IncorrectRequestBody, http.StatusBadRequest, h.logger, err)
		return
	}

	albums, err := h.searchServices.FindAlbums(r.Context(), sr.Query, sr.Amount)
	if err != nil {
		commonHTTP.ErrorResponseWithErrLogging(w, r,
			albumsFindServerError, http.StatusInternalServerError, h.logger, err)
		return
	}

	resp, err := models.AlbumTransferFromQuery(r.Context(),
		albums, user, h.albumServices.IsLiked, h.artistServices.IsLiked, h.artistServices.GetByAlbum)
	if err != nil {
		commonHTTP.ErrorResponseWithErrLogging(w, r,
			albumsFindServerError, http.StatusInternalServerError, h.logger, err)
		return
	}

	commonHTTP.SuccessResponse(w, resp, h.logger)
}

// swaggermock
func (h *Handler) FindArtists(w http.ResponseWriter, r *http.Request) {

}

// swaggermock
func (h *Handler) FindTracks(w http.ResponseWriter, r *http.Request) {

}

// swaggermock
func (h *Handler) FincPlaylists(w http.ResponseWriter, r *http.Request) {

}
