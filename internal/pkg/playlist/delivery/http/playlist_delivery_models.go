package http

import (
	"html"

	valid "github.com/asaskevich/govalidator"

	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/models"
)

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

// Create
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

type playlistCreateResponse struct {
	ID uint32 `json:"id"`
}

type defaultResponse struct {
	Status string `json:"status"`
}
