package middleware

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/models"
	logMocks "github.com/go-park-mail-ru/2023_1_Technokaif/pkg/logger/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

var wrapRequestWithUser = func(r *http.Request, user *models.User) *http.Request {
	if user == nil {
		return r
	}
	ctx := context.WithValue(r.Context(), models.ContextKeyUserType{}, user)
	return r.WithContext(ctx)
}

var correctUser = models.User{
	ID: 1,
}

func TestUserDeliveryCheckUserAuthAndResponse(t *testing.T) {
	c := gomock.NewController(t)

	l := logMocks.NewMockLogger(c)
	l.EXPECT().Error(gomock.Any()).AnyTimes()
	l.EXPECT().Info(gomock.Any()).AnyTimes()
	l.EXPECT().Errorf(gomock.Any(), gomock.Any()).AnyTimes()
	l.EXPECT().Infof(gomock.Any(), gomock.Any()).AnyTimes()

	h := NewMiddleware(l)

	r := chi.NewRouter()

	correctUserID := uint32(1)
	correctUserIDPath := fmt.Sprint(correctUserID)

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
			expectedStatus:   400,
			expectedResponse: `{"message": "invalid url parameter"}`,
		},
		{
			name:             "No User",
			userIDPath:       correctUserIDPath,
			user:             nil,
			expectedStatus:   401,
			expectedResponse: `{"message": "unathorized"}`,
		},
		{
			name:             "Mismatched IDs",
			userIDPath:       "2",
			user:             &correctUser,
			expectedStatus:   403,
			expectedResponse: `{"message": "user has no rights"}`,
		},
	}

	for _, tc := range testTable {
		t.Run(tc.name, func(t *testing.T) {
			r.With(h.CheckUserAuthAndResponce).Get("/user/{userID}", func(w http.ResponseWriter, r *http.Request) {})

			// Request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/user/"+tc.userIDPath, nil)
			r.ServeHTTP(w, wrapRequestWithUser(req, tc.user))

			// Test
			assert.Equal(t, tc.expectedStatus, w.Code)
			assert.JSONEq(t, tc.expectedResponse, w.Body.String())
		})
	}
}
