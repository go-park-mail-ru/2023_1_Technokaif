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

// swaggermock
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	// ...
}

// swaggermock
func (h *Handler) Read(w http.ResponseWriter, r *http.Request) {
	userID, err := commonHttp.GetTrackIDFromRequest(r)
	if err != nil {
		h.logger.Infof("get track by id : %v", err.Error())
		commonHttp.ErrorResponse(w, "invalid url parameter", http.StatusBadRequest, h.logger)
		return
	}

	track, err := h.trackServices.GetByID(uint32(userID))
	var errNoSuchTrack *models.NoSuchTrackError
	if errors.As(err, &errNoSuchTrack) {
		h.logger.Error(err.Error())
		commonHttp.ErrorResponse(w, "no such track", http.StatusBadRequest, h.logger)
		return
	} else if err != nil {
		h.logger.Error(err.Error())
		commonHttp.ErrorResponse(w, "error while getting track", http.StatusInternalServerError, h.logger)
		return
	}

	resp, err := h.trackTransferFromEntry(*track)
	if err != nil {
		h.logger.Error(err.Error())
		commonHttp.ErrorResponse(w, "error while getting track", http.StatusInternalServerError, h.logger)
		return
	}

	w.Header().Set("Content-Type", "json/application; charset=utf-8")
	encoder := json.NewEncoder(w)
	if err := encoder.Encode(&resp); err != nil {
		h.logger.Error(err.Error())
		commonHttp.ErrorResponse(w, "can't encode response into json", http.StatusInternalServerError, h.logger)
		return
	}
}

// swaggermock
func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	// ...
}

// swaggermock
func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	// ...
}

//	@Summary		Track Feed
//	@Tags			track feed
//	@Description	Feed tracks for user
//	@Accept			json
//	@Produce		json
//	@Success		200		{object}	signUpResponse	"Show feed"
//	@Failure		500		{object}	errorResponse	"Server error"
//	@Router			/api/track/feed [get]
func (h *Handler) Feed(w http.ResponseWriter, r *http.Request) {
	tracks, err := h.trackServices.GetFeed()
	if err != nil {
		h.logger.Error(err.Error())
		commonHttp.ErrorResponse(w, "error while getting tracks", http.StatusInternalServerError, h.logger)
		return
	}

	tracksTransfers, err := h.trackTransferFromQuery(tracks)
	if err != nil {
		h.logger.Error(err.Error())
		commonHttp.ErrorResponse(w, "error while getting tracks", http.StatusInternalServerError, h.logger)
		return
	}

	w.Header().Set("Content-Type", "json/application; charset=utf-8")
	encoder := json.NewEncoder(w)
	if err := encoder.Encode(&tracksTransfers); err != nil {
		h.logger.Error(err.Error())
		commonHttp.ErrorResponse(w, "can't encode response into json", http.StatusInternalServerError, h.logger)
		return
	}
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

		trackTransfers = append(trackTransfers,trackTransfer)
	}

	return trackTransfers, nil
}
