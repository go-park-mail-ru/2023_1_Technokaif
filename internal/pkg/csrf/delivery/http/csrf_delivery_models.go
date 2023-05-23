package delivery

//go:generate easyjson -no_std_marshalers csrf_delivery_models.go

// Response messages
const (
	invalidAccessToken = "invalid access token"
	csrfGetError       = "failed to get CSRF-token"
)

//easyjson:json
type getCSRFResponce struct {
	CSRF string `json:"csrf"`
}
