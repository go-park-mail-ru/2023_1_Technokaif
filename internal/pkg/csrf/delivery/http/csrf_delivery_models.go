package delivery

const (
	invalidAccessToken = "invalid access token"
	csrfGetError       = "failed to get CSRF-token"
)

type getCSRFResponce struct {
	CSRF string `json:"csrf"`
}
