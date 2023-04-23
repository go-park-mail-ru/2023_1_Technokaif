package delivery

// Response messages
const (
	invalidAccessToken = "invalid access token"
	csrfGetError       = "failed to get CSRF-token"
)

type getCSRFResponce struct {
	CSRF string `json:"csrf"`
}
