package artist_delivery

import (
	"encoding/json"
	"net/http"

	commonHttp "github.com/go-park-mail-ru/2023_1_Technokaif/internal/common/http"
	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/models"
	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/artist"
	"github.com/go-park-mail-ru/2023_1_Technokaif/pkg/logger"
)

type ArtistHandler struct {
	artistServices artist.ArtistUsecase
	logger         logger.Logger
}

func NewArtistHandler(au artist.ArtistUsecase, logger logger.Logger) *ArtistHandler {
	return &ArtistHandler{
		artistServices: au,
		logger:         logger,
	}
}

// swaggermock
func (ah *ArtistHandler) Create(w http.ResponseWriter, r *http.Request) {
	var artist models.Artist

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&artist); err != nil {
		commonHttp.ErrorResponse(w, "incorrect input body", http.StatusBadRequest)
		return
	}

	if err := ah.artistServices.Create(artist); err != nil {
		ah.logger.Error(err.Error())
		commonHttp.ErrorResponse(w, "can't create artist", http.StatusInternalServerError)
	}

	// ...
}

// swaggermock
func (ah *ArtistHandler) Read(w http.ResponseWriter, r *http.Request) {
	// ...
}

// swaggermock
func (ah *ArtistHandler) Update(w http.ResponseWriter, r *http.Request) {
	// ...
}

// swaggermock
func (ah *ArtistHandler) Delete(w http.ResponseWriter, r *http.Request) {
	// ...
}

// swaggermock
func (ah *ArtistHandler) Tracks(w http.ResponseWriter, r *http.Request) {
	// ...
}

// swaggermock
func (ah *ArtistHandler) Albums(w http.ResponseWriter, r *http.Request) {
	// ...
}

//	@Summary		Artist Feed
//	@Tags			artist feed
//	@Description	Feed albums for user
//	@Accept			json
//	@Produce		json
//	@Success		200		{object}	signUpResponse	"Show feed"
//	@Failure		500		{object}	errorResponse	"Server error"
//	@Router			/api/artist/feed [get]
func (ah *ArtistHandler) Feed(w http.ResponseWriter, r *http.Request) {
	artists, err := ah.artistServices.GetFeed()
	if err != nil {
		ah.logger.Error(err.Error())
		commonHttp.ErrorResponse(w, "error while getting albums", http.StatusInternalServerError)
		return
	}

	artistsTransfer := ah.artistTransferFromQuery(artists)

	w.Header().Set("Content-Type", "json/application; charset=utf-8")
	encoder := json.NewEncoder(w)
	if err := encoder.Encode(&artistsTransfer); err != nil {
		ah.logger.Error(err.Error())
		commonHttp.ErrorResponse(w, "can't encode response into json", http.StatusInternalServerError)
		return
	}
}

func (ah *ArtistHandler) artistTransferFromQuery(artists []models.Artist) []models.ArtistTransfer {
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
