package album_delivery

import (
	"encoding/json"
	"fmt"
	"net/http"

	commonHttp "github.com/go-park-mail-ru/2023_1_Technokaif/internal/common/http"
	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/models"
	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/album"
	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/artist"
	"github.com/go-park-mail-ru/2023_1_Technokaif/pkg/logger"
)

type AlbumHandler struct {
	albumServices  album.AlbumUsecase
	artistServices artist.ArtistUsecase
	logger         logger.Logger
}

func NewAlbumHandler(alu album.AlbumUsecase, aru artist.ArtistUsecase, l logger.Logger) *AlbumHandler {
	return &AlbumHandler{
		albumServices:  alu,
		artistServices: aru,
		logger:         l,
	}
}

// swaggermock
func (ah *AlbumHandler) Create(w http.ResponseWriter, r *http.Request) {
	// ...
}

// swaggermock
func (ah *AlbumHandler) Read(w http.ResponseWriter, r *http.Request) {
	// ...
}

// swaggermock
func (ah *AlbumHandler) Update(w http.ResponseWriter, r *http.Request) {
	// ...
}

// swaggermock
func (ah *AlbumHandler) Delete(w http.ResponseWriter, r *http.Request) {
	// ...
}

// swaggermock
func (ah *AlbumHandler) Tracks(w http.ResponseWriter, r *http.Request) {
	// ...
}

//	@Summary		Album Feed
//	@Tags			album feed
//	@Description	Feed albums for user
//	@Accept			json
//	@Produce		json
//	@Success		200		{object}	signUpResponse	"Show feed"
//	@Failure		500		{object}	errorResponse	"Server error"
//	@Router			/api/album/feed [get]
func (ah *AlbumHandler) Feed(w http.ResponseWriter, r *http.Request) {
	albums, err := ah.albumServices.GetFeed()
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

func (ah *AlbumHandler) artistTransferFromQuery(artists []models.Artist) []models.ArtistTransfer {
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

func (ah *AlbumHandler) albumTransferFromQuery(albums []models.Album) ([]models.AlbumTransfer, error) {
	at := make([]models.AlbumTransfer, 0, len(albums))
	for _, a := range albums {
		artists, err := ah.artistServices.GetByAlbum(a.ID)
		if err != nil {
			return nil, fmt.Errorf("(delivery) can't get albums's (id #%d) artists: %w", a.ID, err)
		}

		at = append(at, models.AlbumTransfer{
			ID:          a.ID,
			Name:        a.Name,
			Artists:     ah.artistTransferFromQuery(artists),
			Description: a.Description,
			CoverSrc:    a.CoverSrc,
		})
	}

	return at, nil
}
