package middleware

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/golang/mock/gomock"

	commonHttp "github.com/go-park-mail-ru/2023_1_Technokaif/internal/common/http"
	commonTests "github.com/go-park-mail-ru/2023_1_Technokaif/internal/common/tests"
	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/models"
)

func TestUserDeliveryCheckUserAuthAndResponse(t *testing.T) {
	c := gomock.NewController(t)

	l := commonTests.MockLogger(c)

	h := NewMiddleware(l)

	r := chi.NewRouter()

	correctUser := models.User{
		ID: 1,
	}
	correctUserIDPath := fmt.Sprint(correctUser.ID)

	testTable := []struct {
		name             string
		userIDPath       string
		user             *models.User
		expectedStatus   int
		expectedResponse string
	}{
		{
			name:             "Incorrect ID In Path",
			userIDPath:       "0",
			user:             &correctUser,
			expectedStatus:   http.StatusBadRequest,
			expectedResponse: commonTests.ErrorResponse(commonHttp.InvalidURLParameter),
		},
		{
			name:             "No User",
			userIDPath:       correctUserIDPath,
			user:             nil,
			expectedStatus:   http.StatusUnauthorized,
			expectedResponse: commonTests.ErrorResponse(commonHttp.UnathorizedUser),
		},
		{
			name:             "Mismatched IDs",
			userIDPath:       "2",
			user:             &correctUser,
			expectedStatus:   http.StatusForbidden,
			expectedResponse: commonTests.ErrorResponse("user has no rights"),
		},
	}

	for _, tc := range testTable {
		t.Run(tc.name, func(t *testing.T) {
			r.With(h.CheckUserAuthAndResponce).Get("/user/{userID}", func(w http.ResponseWriter, r *http.Request) {})

			commonTests.DeliveryTestGet(t, r, "/user/"+tc.userIDPath, tc.expectedStatus, tc.expectedResponse,
				commonTests.WrapRequestWithUserNotNilFunc(tc.user))
		})
	}
}
