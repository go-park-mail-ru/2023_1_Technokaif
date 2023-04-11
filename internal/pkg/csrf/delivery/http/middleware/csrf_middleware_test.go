package middleware

import (
	"math/rand"
	"net/http"
	"net/http/httptest"
	"testing"
	"errors"

	"github.com/go-chi/chi/v5"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/models"
	tokenMocks "github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/token/mocks"
	commonTests "github.com/go-park-mail-ru/2023_1_Technokaif/internal/common/tests"
)

func TestAuthDeliveryCheckCSRFToken(t *testing.T) {
	type mockBehavior func(t *tokenMocks.MockUsecase, token string)

	correctTestUser := &models.User{
		ID:      uint32(rand.Intn(100) + 1),
		Version: uint32(rand.Intn(100)),
	}

	testWrapRequestWithUser := commonTests.WrapRequestWithUser

	testTable := []struct {
		name              string
		csrfHeader        string
		csrfToken         string
		userInRequest     *models.User
		wrapUser          bool
		mockBehavior      mockBehavior
		expectingResponse bool
		expectedStatus    int
		expectedResponse  string
	}{
		{
			name:          "Ok",
			csrfHeader:    csrfTokenHttpHeader,
			csrfToken:     "token",
			userInRequest: correctTestUser,
			wrapUser:      true,
			mockBehavior: func(t *tokenMocks.MockUsecase, token string) {
				t.EXPECT().CheckCSRFToken(token).Return(correctTestUser.ID, nil)
			},
			expectingResponse: false,
		},
		{
			name:          "No user in context",
			csrfHeader:    csrfTokenHttpHeader,
			csrfToken:     "token",
			wrapUser:      false,
			mockBehavior: func(t *tokenMocks.MockUsecase, token string) {},
			expectingResponse: true,
			expectedStatus: 400,
			expectedResponse: `{"message": "invalid access token"}`,
		},
		{
			name:          "Nil user in context",
			csrfHeader:    csrfTokenHttpHeader,
			csrfToken:     "token",
			userInRequest: nil,
			wrapUser:      true,
			mockBehavior: func(t *tokenMocks.MockUsecase, token string) {},
			expectingResponse: true,
			expectedStatus: 400,
			expectedResponse: `{"message": "invalid access token"}`,
		},
		{
			name:          "Invalid token",
			csrfHeader:    csrfTokenHttpHeader,
			csrfToken:     "token",
			userInRequest: correctTestUser,
			wrapUser:      true,
			mockBehavior: func(t *tokenMocks.MockUsecase, token string) {
				t.EXPECT().CheckCSRFToken(token).Return(uint32(0), errors.New("invalid signing token"))
			},
			expectingResponse: true,
			expectedStatus: 400,
			expectedResponse: `{"message": "invalid CSRF token"}`,
		},
		{
			name:          "Missing token",
			userInRequest: correctTestUser,
			wrapUser:      true,
			mockBehavior: func(t *tokenMocks.MockUsecase, token string) {},
			expectingResponse: true,
			expectedStatus: 400,
			expectedResponse: `{"message": "missing CSRF token"}`,
		},
		{
			name:          "Incorrect token payload userID",
			csrfHeader:    csrfTokenHttpHeader,
			csrfToken:     "token",
			userInRequest: correctTestUser,
			wrapUser:      true,
			mockBehavior: func(t *tokenMocks.MockUsecase, token string) {
				t.EXPECT().CheckCSRFToken(token).Return(uint32(0), nil)
			},
			expectingResponse: true,
			expectedStatus: 400,
			expectedResponse: `{"message": "invalid CSRF token"}`,
		},
	}

	for _, tc := range testTable {
		t.Run(tc.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			tokenMockUsecase := tokenMocks.NewMockUsecase(c)

			tc.mockBehavior(tokenMockUsecase, tc.csrfToken)

			l := commonTests.MockLogger(c)

			h := NewMiddleware(tokenMockUsecase, l)

			r := chi.NewRouter()
			r.With(h.CheckCSRFToken).Get("/csrf", func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(200)
			})

			// Init Test Request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/csrf", nil)
			req.Header.Set(tc.csrfHeader, tc.csrfToken)

			reqWrapped := testWrapRequestWithUser(req, tc.userInRequest, tc.wrapUser)

			r.ServeHTTP(w, reqWrapped)

			if tc.expectingResponse {
				assert.Equal(t, tc.expectedStatus, w.Code)
				assert.JSONEq(t, tc.expectedResponse, w.Body.String())
			} else {
				assert.Equal(t, 200, w.Code)
			}
		})
	}
}
