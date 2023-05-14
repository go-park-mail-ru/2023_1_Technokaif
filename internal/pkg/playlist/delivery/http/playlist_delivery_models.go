package http

import (
	"context"
	"html"

	valid "github.com/asaskevich/govalidator"

	userHTTP "github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/user/delivery/http"

	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/models"
)

//go:generate easyjson -no_std_marshalers playlist_delivery_models.go

const MaxCoverMemory = 5 << 20
const coverFormKey = "cover"

// Response messages
const (
	playlistNotFound = "no such playlist"
	trackNotFound    = "no such track"
	userNotFound     = "no such user"

	playlistCoverInvalidData     = "invalid cover data"
	playlistCoverInvalidDataType = "invalid cover data type"
	playlistCoverUploadNoRights  = "no rights to upload cover"
	playlistCoverServerError     = "can't upload cover"

	playlistCreateNorights      = "no rights to create playlist"
	playlistUpdateNoRights      = "no rights to update playlist"
	playlistDeleteNoRights      = "no rights to delete playlist"
	playlistAddTrackNoRights    = "no rights to add track into playlist"
	playlistDeleteTrackNoRights = "no rights to delete track from playlist"

	playlistCreateServerError      = "can't create playlist"
	playlistGetServerError         = "can't get playlist"
	playlistsGetServerError        = "can't get playlists"
	playlistUpdateServerError      = "can't update playlist"
	playlistDeleteServerError      = "can't delete playlist"
	playlistAddTrackServerError    = "can't add track into playlist"
	playlistDeleteTrackServerError = "can't delete track from playlist"

	playlistUpdatedSuccessfully       = "ok"
	playlistDeletedSuccessfully       = "ok"
	playlistTrackAddedSuccessfully    = "ok"
	playlistTrackDeletedSuccessfully  = "ok"
	playlistCoverUploadedSuccessfully = "ok"
)

//easyjson:json
type PlaylistTransfer struct {
	ID          uint32                  `json:"id"`
	Name        string                  `json:"name"`
	Users       []userHTTP.UserTransfer `json:"users"`
	Description *string                 `json:"description,omitempty"`
	IsLiked     bool                    `json:"isLiked"`
	CoverSrc    string                  `json:"cover,omitempty"`
}

//easyjson:json
type PlaylistTransfers []PlaylistTransfer

type usersByPlaylistsGetter func(ctx context.Context, playlistID uint32) ([]models.User, error)
type playlistLikeChecker func(ctx context.Context, playlistID, userID uint32) (bool, error)

// PlaylistTransferFromEntry converts Playlist to PlaylistTransfer
func PlaylistTransferFromEntry(ctx context.Context, p models.Playlist, user *models.User,
	likeChecker playlistLikeChecker, usersGetter usersByPlaylistsGetter) (PlaylistTransfer, error) {

	users, err := usersGetter(ctx, p.ID)
	if err != nil {
		return PlaylistTransfer{}, err
	}

	isLiked := false
	if user != nil {
		isLiked, err = likeChecker(ctx, p.ID, user.ID)
		if err != nil {
			return PlaylistTransfer{}, err
		}
	}

	return PlaylistTransfer{
		ID:          p.ID,
		Name:        p.Name,
		Users:       userHTTP.UserTransferFromList(users),
		Description: p.Description,
		IsLiked:     isLiked,
		CoverSrc:    p.CoverSrc,
	}, nil
}

// PlaylistTransferFromList converts []Playlist to []PlaylistTransfer
func PlaylistTransferFromList(ctx context.Context, playlists []models.Playlist, user *models.User, likeChecker playlistLikeChecker,
	usersGetter usersByPlaylistsGetter) (PlaylistTransfers, error) {

	playlistTransfers := make([]PlaylistTransfer, 0, len(playlists))

	for _, p := range playlists {
		pt, err := PlaylistTransferFromEntry(ctx, p, user, likeChecker, usersGetter)
		if err != nil {
			return nil, err
		}

		playlistTransfers = append(playlistTransfers, pt)
	}

	return playlistTransfers, nil
}

// Create
//
//easyjson:json
type playlistCreateInput struct {
	Name        string   `json:"name" valid:"required"`
	UsersID     []uint32 `json:"users" valid:"required"`
	Description *string  `json:"description"`
}

func (pci *playlistCreateInput) validateAndEscape() error {
	pci.escapeHtml()

	_, err := valid.ValidateStruct(pci)

	return err
}

func (pci *playlistCreateInput) escapeHtml() {
	pci.Name = html.EscapeString(pci.Name)
	if pci.Description != nil {
		*pci.Description = html.EscapeString(*pci.Description)
	}
}

func (pci *playlistCreateInput) ToPlaylist() models.Playlist {
	return models.Playlist{
		Name:        pci.Name,
		Description: pci.Description,
	}
}

// Update
//
//easyjson:json
type playlistUpdateInput struct {
	Name        string   `json:"name" valid:"required"`
	UsersID     []uint32 `json:"users" valid:"required"`
	Description *string  `json:"description"`
}

func (pui *playlistUpdateInput) validateAndEscape() error {
	pui.escapeHtml()

	_, err := valid.ValidateStruct(pui)

	return err
}

func (pui *playlistUpdateInput) escapeHtml() {
	pui.Name = html.EscapeString(pui.Name)
	if pui.Description != nil {
		*pui.Description = html.EscapeString(*pui.Description)
	}
}

func (pui *playlistUpdateInput) ToPlaylist(playlistID uint32) models.Playlist {
	return models.Playlist{
		ID:          playlistID,
		Name:        pui.Name,
		Description: pui.Description,
	}
}

//easyjson:json
type playlistCreateResponse struct {
	ID uint32 `json:"id"`
}

//easyjson:json
type defaultResponse struct {
	Status string `json:"status"`
}
