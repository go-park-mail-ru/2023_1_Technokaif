package http

// UploadAvatar
const maxAvatarMemory = 1 * (1 << 20)

type userUploadAvatarResponse struct {
	Status string `json:"status"`
}
