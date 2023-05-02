package delivery

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/models"
	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/album"
	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/artist"
	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/playlist"
	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/search"
	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/track"
	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/user"
	"github.com/go-park-mail-ru/2023_1_Technokaif/pkg/logger"

	commonHTTP "github.com/go-park-mail-ru/2023_1_Technokaif/internal/common/http"
)

type Handler struct {
	searchServices   search.Usecase
	albumServices    album.Usecase
	artistServices   artist.Usecase
	trackServices    track.Usecase
	playlistServices playlist.Usecase
	userServices     user.Usecase
	logger           logger.Logger
}

func NewHandler(su search.Usecase, alu album.Usecase, aru artist.Usecase,
	tu track.Usecase, pu playlist.Usecase, uu user.Usecase, l logger.Logger) *Handler {

	return &Handler{
		searchServices:   su,
		albumServices:    alu,
		artistServices:   aru,
		trackServices:    tu,
		playlistServices: pu,
		userServices:     uu,

		logger: l,
	}
}

// @Summary		Find Albums
// @Tags		Search
// @Description	Find amount of albums by search-query
// @Accept      json
// @Produce		json
// @Param		query	body		searchRequest    	true "Query for search"
// @Success		200		{object}	searchAlbumsResponse	 "Albums found"
// @Failure		400		{object}	http.Error				 "Incorrect body"
// @Failure		401		{object}	http.Error  			 "User unathorized"
// @Failure		500		{object}	http.Error				 "Server error"
// @Router		/api/albums/search [post]
func (h *Handler) FindAlbums(w http.ResponseWriter, r *http.Request) {
	user, err := commonHTTP.GetUserFromRequest(r)
	if err != nil && !errors.Is(err, commonHTTP.ErrUnauthorized) {
		commonHTTP.ErrorResponseWithErrLogging(w, r,
			albumsFindServerError, http.StatusInternalServerError, h.logger, err)
		return
	}

	var sr searchRequest
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

	at, err := models.AlbumTransferFromList(r.Context(),
		albums, user, h.albumServices.IsLiked, h.artistServices.IsLiked, h.artistServices.GetByAlbum)
	if err != nil {
		commonHTTP.ErrorResponseWithErrLogging(w, r,
			albumsFindServerError, http.StatusInternalServerError, h.logger, err)
		return
	}

	resp := searchAlbumsResponse{Albums: at}

	commonHTTP.SuccessResponse(w, r, resp, h.logger)
}

// @Summary		Find Artists
// @Tags		Search
// @Description	Find amount of artists by search-query
// @Accept      json
// @Produce		json
// @Param		query	body		searchRequest    	true "Query for search"
// @Success		200		{object}	searchArtistsResponse	 "Artists found"
// @Failure		400		{object}	http.Error				 "Incorrect body"
// @Failure		401		{object}	http.Error  			 "User unathorized"
// @Failure		500		{object}	http.Error				 "Server error"
// @Router		/api/artists/search [post]
func (h *Handler) FindArtists(w http.ResponseWriter, r *http.Request) {
	user, err := commonHTTP.GetUserFromRequest(r)
	if err != nil && !errors.Is(err, commonHTTP.ErrUnauthorized) {
		commonHTTP.ErrorResponseWithErrLogging(w, r,
			artistsFindServerError, http.StatusInternalServerError, h.logger, err)
		return
	}

	var sr searchRequest
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

	artists, err := h.searchServices.FindArtists(r.Context(), sr.Query, sr.Amount)
	if err != nil {
		commonHTTP.ErrorResponseWithErrLogging(w, r,
			artistsFindServerError, http.StatusInternalServerError, h.logger, err)
		return
	}

	at, err := models.ArtistTransferFromList(r.Context(), artists, user, h.artistServices.IsLiked)
	if err != nil {
		commonHTTP.ErrorResponseWithErrLogging(w, r,
			artistsFindServerError, http.StatusInternalServerError, h.logger, err)
		return
	}

	resp := searchArtistsResponse{Artists: at}

	commonHTTP.SuccessResponse(w, r, resp, h.logger)
}

// @Summary		Find Tracks
// @Tags		Search
// @Description	Find amount of tracks by search-query
// @Accept      json
// @Produce		json
// @Param		query	body		searchRequest    	true "Query for search"
// @Success		200		{object}	searchTracksResponse	 "Tracks found"
// @Failure		400		{object}	http.Error				 "Incorrect body"
// @Failure		401		{object}	http.Error  			 "User unathorized"
// @Failure		500		{object}	http.Error				 "Server error"
// @Router		/api/tracks/search [post]
func (h *Handler) FindTracks(w http.ResponseWriter, r *http.Request) {
	user, err := commonHTTP.GetUserFromRequest(r)
	if err != nil && !errors.Is(err, commonHTTP.ErrUnauthorized) {
		commonHTTP.ErrorResponseWithErrLogging(w, r,
			tracksFindServerError, http.StatusInternalServerError, h.logger, err)
		return
	}

	var sr searchRequest
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

	tracks, err := h.searchServices.FindTracks(r.Context(), sr.Query, sr.Amount)
	if err != nil {
		commonHTTP.ErrorResponseWithErrLogging(w, r,
			tracksFindServerError, http.StatusInternalServerError, h.logger, err)
		return
	}

	tt, err := models.TrackTransferFromList(r.Context(),
		tracks, user, h.trackServices.IsLiked, h.artistServices.IsLiked, h.artistServices.GetByTrack)
	if err != nil {
		commonHTTP.ErrorResponseWithErrLogging(w, r,
			tracksFindServerError, http.StatusInternalServerError, h.logger, err)
		return
	}

	resp := searchTracksResponse{Tracks: tt}

	commonHTTP.SuccessResponse(w, r, resp, h.logger)
}

// @Summary		Find Playlists
// @Tags		Search
// @Description	Find amount of playlists by search-query
// @Accept      json
// @Produce		json
// @Param		query	body		searchRequest    	true "Query for search"
// @Success		200		{object}	searchPlaylistsResponse	 "Playlists found"
// @Failure		400		{object}	http.Error				 "Incorrect body"
// @Failure		401		{object}	http.Error  			 "User unathorized"
// @Failure		500		{object}	http.Error				 "Server error"
// @Router		/api/playlists/search [post]
func (h *Handler) FindPlaylists(w http.ResponseWriter, r *http.Request) {
	user, err := commonHTTP.GetUserFromRequest(r)
	if err != nil && !errors.Is(err, commonHTTP.ErrUnauthorized) {
		commonHTTP.ErrorResponseWithErrLogging(w, r,
			playlistsFindServerError, http.StatusInternalServerError, h.logger, err)
		return
	}

	var sr searchRequest
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

	playlists, err := h.searchServices.FindPlaylists(r.Context(), sr.Query, sr.Amount)
	if err != nil {
		commonHTTP.ErrorResponseWithErrLogging(w, r,
			albumsFindServerError, http.StatusInternalServerError, h.logger, err)
		return
	}

	pt, err := models.PlaylistTransferFromList(r.Context(),
		playlists, user, h.playlistServices.IsLiked, h.userServices.GetByPlaylist)
	if err != nil {
		commonHTTP.ErrorResponseWithErrLogging(w, r,
			albumsFindServerError, http.StatusInternalServerError, h.logger, err)
		return
	}

	resp := searchPlaylistsResponse{Playlists: pt}

	commonHTTP.SuccessResponse(w, r, resp, h.logger)
}
