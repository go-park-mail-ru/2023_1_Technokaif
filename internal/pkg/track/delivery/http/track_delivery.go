package track_delivery

import (
	"encoding/json"
	"fmt"
	"net/http"

	commonHttp "github.com/go-park-mail-ru/2023_1_Technokaif/internal/common/http"
	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/logger"
	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/models"
	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/artist"
	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/track"
)

type TrackHandler struct {
	trackServices  track.TrackUsecase
	artistServices artist.ArtistUsecase
	logger         logger.Logger
}

func NewTrackHandler(tu track.TrackUsecase, au artist.ArtistUsecase, l logger.Logger) *TrackHandler {
	return &TrackHandler{
		trackServices:  tu,
		artistServices: au,
		logger:         l,
	}
}

// swaggermock
func (th *TrackHandler) Create(w http.ResponseWriter, r *http.Request) {
	// ...
}

// swaggermock
func (th *TrackHandler) Read(w http.ResponseWriter, r *http.Request) {
	// ...
}

// swaggermock
func (th *TrackHandler) Update(w http.ResponseWriter, r *http.Request) {
	// ...
}

// swaggermock
func (th *TrackHandler) Delete(w http.ResponseWriter, r *http.Request) {
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
func (th *TrackHandler) Feed(w http.ResponseWriter, r *http.Request) {
	tracks, err := th.trackServices.GetFeed()
	if err != nil {
		th.logger.Error(err.Error())
		commonHttp.ErrorResponse(w, "error while getting tracks", http.StatusInternalServerError)
		return
	}

	tracksTransfer, err := th.trackTransferFromQuery(tracks)
	if err != nil {
		th.logger.Error(err.Error())
		commonHttp.ErrorResponse(w, "error while getting tracks", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "json/application; charset=utf-8")
	encoder := json.NewEncoder(w)
	if err := encoder.Encode(&tracksTransfer); err != nil {
		th.logger.Error(err.Error())
		commonHttp.ErrorResponse(w, "can't encode response into json", http.StatusInternalServerError)
		return
	}
}

func (th *TrackHandler) artistTransferFromQuery(artists []models.Artist) []models.ArtistTransfer {
	at := make([]models.ArtistTransfer, len(artists))
	for _, a := range artists {
		at = append(at, models.ArtistTransfer{
			ID:        a.ID,
			Name:      a.Name,
			AvatarSrc: a.AvatarSrc,
		})
	}

	return at
}

func (th *TrackHandler) trackTransferFromQuery(tracks []models.Track) ([]models.TrackTransfer, error) {
	tt := make([]models.TrackTransfer, len(tracks))
	for _, t := range tracks {
		artists, err := th.artistServices.GetByTrack(t.ID)
		if err != nil {
			return nil, fmt.Errorf("(delivery) can't get track's (id #%d) artists: %w", t.ID, err)
		}

		tt = append(tt, models.TrackTransfer{
			ID:        t.ID,
			Name:      t.Name,
			Artists:   th.artistTransferFromQuery(artists),
			CoverSrc:  t.CoverSrc,
			RecordSrc: t.RecordSrc,
		})
	}

	return tt, nil
}
