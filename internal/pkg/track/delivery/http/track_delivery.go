package http

import (
	"encoding/json"
	"errors"
	"fmt"
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

// TODO ERRORS

type trackCreateInput struct {
	Name      string   `json:"name"`
	AlbumID   uint32   `json:"albumID,omitempty"`
	ArtistsID []uint32 `json:"artistsID"`
	CoverSrc  string   `json:"cover"`
	RecordSrc string   `json:"record"`
}

func (tci *trackCreateInput) ToTrack() models.Track {
	return models.Track{
		Name:      tci.Name,
		AlbumID:   tci.AlbumID,
		CoverSrc:  tci.CoverSrc,
		RecordSrc: tci.RecordSrc,
	}
}

type trackCreateResponse struct {
	ID uint32 `json:"id"`
}

// @Summary		Create Track
// @Tags		Track
// @Description	Create new track by sent object
// @Accept      json
// @Produce		json
// @Param		track	body		trackCreateInput	true	"Track info"
// @Success		200		{object}	trackCreateResponse	        "Track created"
// @Failure		400		{object}	errorResponse	"Client error"
// @Failure		500		{object}	errorResponse	"Server error"
// @Router		/api/tracks/ [post]
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var tci trackCreateInput

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&tci); err != nil {
		h.logger.Info(err.Error())
		commonHttp.ErrorResponse(w, "incorrect input body", http.StatusBadRequest, h.logger)
		return
	}

	track := tci.ToTrack()

	trackID, err := h.trackServices.Create(track, tci.ArtistsID)
	if err != nil {
		h.logger.Error(err.Error())
		commonHttp.ErrorResponse(w, "can't create track", http.StatusInternalServerError, h.logger)
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
// @Failure		400		{object}	errorResponse	"Client error"
// @Failure		500		{object}	errorResponse	"Server error"
// @Router		/api/tracks/{trackID}/ [get]
func (h *Handler) Read(w http.ResponseWriter, r *http.Request) {
	userID, err := commonHttp.GetTrackIDFromRequest(r)
	if err != nil {
		h.logger.Infof("get track by id : %v", err)
		commonHttp.ErrorResponse(w, "invalid url parameter", http.StatusBadRequest, h.logger)
		return
	}

	track, err := h.trackServices.GetByID(uint32(userID))
	var errNoSuchTrack *models.NoSuchTrackError
	if errors.As(err, &errNoSuchTrack) {
		h.logger.Info(err.Error())
		commonHttp.ErrorResponse(w, "no such track", http.StatusBadRequest, h.logger)
		return
	}
	if err != nil {
		h.logger.Error(err.Error())
		commonHttp.ErrorResponse(w, "error while getting track", http.StatusInternalServerError, h.logger)
		return
	}

	tt, err := h.trackTransferFromEntry(*track)
	if err != nil {
		h.logger.Error(err.Error())
		commonHttp.ErrorResponse(w, "error while getting track", http.StatusInternalServerError, h.logger)
		return
	}

	commonHttp.SuccessResponse(w, tt, h.logger)
}

// swaggermock
func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	// ...
}

type trackDeleteResponse struct {
	Status string `json:"status"`
}

// @Summary		Delete Track
// @Tags		Track
// @Description	Delete track with chosen ID
// @Produce		json
// @Success		200		{object}	trackDeleteResponse	        "Track deleted"
// @Failure		400		{object}	errorResponse	"Client error"
// @Failure		500		{object}	errorResponse	"Server error"
// @Router		/api/tracks/{trackID}/ [delete]
func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	trackID, err := commonHttp.GetTrackIDFromRequest(r)
	if err != nil {
		h.logger.Infof("get track by id: %v", err)
		commonHttp.ErrorResponse(w, "invalid url parameter", http.StatusBadRequest, h.logger)
		return
	}

	err = h.trackServices.DeleteByID(trackID)
	var errNoSuchTrack *models.NoSuchTrackError
	if errors.As(err, &errNoSuchTrack) {
		h.logger.Info(err.Error())
		commonHttp.ErrorResponse(w, "no such track", http.StatusBadRequest, h.logger)
		return
	}
	if err != nil {
		h.logger.Error(err.Error())
		commonHttp.ErrorResponse(w, "error while deleting track", http.StatusInternalServerError, h.logger)
		return
	}

	tdr := trackDeleteResponse{Status: "ok"}

	commonHttp.SuccessResponse(w, tdr, h.logger)
}

// @Summary		Tracks of Artist
// @Tags		Artist
// @Description	All tracks of artist with chosen ID
// @Produce		json
// @Success		200		{object}	[]models.TrackTransfer	    "Show tracks"
// @Failure		400		{object}	errorResponse	"Client error"
// @Failure		500		{object}	errorResponse	"Server error"
// @Router		/api/artists/{artistID}/tracks [get]
func (h *Handler) ReadByArtist(w http.ResponseWriter, r *http.Request) {
	artistID, err := commonHttp.GetArtistIDFromRequest(r)
	if err != nil {
		h.logger.Infof("read by artist: %v", err)
		commonHttp.ErrorResponse(w, "invalid url parameter", http.StatusBadRequest, h.logger)
		return
	}

	tracks, err := h.trackServices.GetByArtist(artistID)
	var errNoSuchArtist *models.NoSuchArtistError
	if errors.As(err, &errNoSuchArtist) {
		h.logger.Info(err.Error())
		commonHttp.ErrorResponse(w, "no such artist", http.StatusBadRequest, h.logger)
		return
	}
	if err != nil {
		h.logger.Error(err.Error())
		commonHttp.ErrorResponse(w, "error while getting artist tracks", http.StatusInternalServerError, h.logger)
		return
	}

	tt, err := h.trackTransferFromQuery(tracks)
	if err != nil {
		h.logger.Error(err.Error())
		commonHttp.ErrorResponse(w, "error while getting artist tracks", http.StatusInternalServerError, h.logger)
		return
	}

	commonHttp.SuccessResponse(w, tt, h.logger)
}

// @Summary		Tracks of Album
// @Tags		Album
// @Description	All tracks of album with chosen ID
// @Produce		json
// @Success		200		{object}	[]models.TrackTransfer	    "Show tracks"
// @Failure		400		{object}	errorResponse	"Client error"
// @Failure		500		{object}	errorResponse	"Server error"
// @Router		/api/albums/{albumID}/tracks [get]
func (h *Handler) ReadByAlbum(w http.ResponseWriter, r *http.Request) {
	albumID, err := commonHttp.GetAlbumIDFromRequest(r)
	if err != nil {
		h.logger.Infof("read by album : %v", err)
		commonHttp.ErrorResponse(w, "invalid url parameter", http.StatusBadRequest, h.logger)
		return
	}

	tracks, err := h.trackServices.GetByAlbum(albumID)
	var errNoSuchAlbum *models.NoSuchAlbumError
	if errors.As(err, &errNoSuchAlbum) {
		h.logger.Info(err.Error())
		commonHttp.ErrorResponse(w, "no such album", http.StatusBadRequest, h.logger)
		return
	}
	if err != nil {
		h.logger.Error(err.Error())
		commonHttp.ErrorResponse(w, "error while getting album tracks", http.StatusInternalServerError, h.logger)
		return
	}

	tt, err := h.trackTransferFromQuery(tracks)
	if err != nil {
		h.logger.Error(err.Error())
		commonHttp.ErrorResponse(w, "error while getting artist tracks", http.StatusInternalServerError, h.logger)
		return
	}

	commonHttp.SuccessResponse(w, tt, h.logger)
}

// @Summary		Track Feed
// @Tags		Feed
// @Description	Feed tracks
// @Produce		json
// @Success		200		{object}	[]models.TrackTransfer	"Tracks feed"
// @Failure		500		{object}	errorResponse	"Server error"
// @Router		/api/tracks/feed [get]
func (h *Handler) Feed(w http.ResponseWriter, r *http.Request) {
	tracks, err := h.trackServices.GetFeed()
	if err != nil {
		h.logger.Error(err.Error())
		commonHttp.ErrorResponse(w, "error while getting tracks", http.StatusInternalServerError, h.logger)
		return
	}

	tt, err := h.trackTransferFromQuery(tracks)
	if err != nil {
		h.logger.Error(err.Error())
		commonHttp.ErrorResponse(w, "error while getting tracks", http.StatusInternalServerError, h.logger)
		return
	}

	commonHttp.SuccessResponse(w, tt, h.logger)
}

// Converts Artist to ArtistTransfer
func (h *Handler) artistTransferFromQuery(artists []models.Artist) []models.ArtistTransfer {
	at := make([]models.ArtistTransfer, 0, len(artists))
	for _, a := range artists {
		at = append(at, models.ArtistTransfer{
			ID:        a.ID,
			Name:      a.Name,
			AvatarSrc: a.AvatarSrc,
		})
	}

	return at
}

// Converts Track to TrackTransfer
func (h *Handler) trackTransferFromEntry(track models.Track) (models.TrackTransfer, error) {

	artists, err := h.artistServices.GetByTrack(track.ID)
	if err != nil {
		return models.TrackTransfer{}, fmt.Errorf("(delivery) can't get track's (id #%d) artists: %w", track.ID, err)
	}

	return models.TrackTransfer{
		ID:        track.ID,
		Name:      track.Name,
		AlbumID:   track.AlbumID,
		Artists:   h.artistTransferFromQuery(artists),
		CoverSrc:  track.CoverSrc,
		RecordSrc: track.RecordSrc,
	}, nil
}

func (h *Handler) trackTransferFromQuery(tracks []models.Track) ([]models.TrackTransfer, error) {
	trackTransfers := make([]models.TrackTransfer, 0, len(tracks))
	for _, t := range tracks {
		trackTransfer, err := h.trackTransferFromEntry(t)
		if err != nil {
			return nil, err
		}

		trackTransfers = append(trackTransfers, trackTransfer)
	}

	return trackTransfers, nil
}

// For swagger, but how to fix?
type errorResponse struct {
	Message string `json:"message"`
}
