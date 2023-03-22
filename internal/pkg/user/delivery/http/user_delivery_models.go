package http

// UploadAvatar
const maxAvatarMemory = 2 * (1 << 20)

type userUploadAvatarResponse struct {
	Status string `json:"status"`
}
