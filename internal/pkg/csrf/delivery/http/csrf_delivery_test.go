package delivery

import (
	"errors"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/models"
	tokenMocks "github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/token/mocks"
	commonTests "github.com/go-park-mail-ru/2023_1_Technokaif/internal/common/tests"
)

func TestDeliveryGetCSRF(t *testing.T) {
	// Init
	type mockBehavior func(t *tokenMocks.MockUsecase, u *models.User)

	testWrapRequestWithUser := commonTests.WrapRequestWithUser

	c := gomock.NewController(t)

	tokenMockUsecase := tokenMocks.NewMockUsecase(c)

	l := commonTests.MockLogger(c)

	h := NewHandler(tokenMockUsecase, l)

	// Routing
	r := chi.NewRouter()
	r.Get("/csrf", h.GetCSRF)

	// Test filling
	correctTestUser := &models.User{
		ID: 1,
	}

	const expectedDefaultCSRF = "csrfagjowajg"

	testTable := []struct {
		name             string
		user             *models.User
		mockBehavior     mockBehavior
		expectedStatus   int
		expectedResponse string
		doWrap           bool
	}{
		{
			name: "Common",
			user: correctTestUser,
			mockBehavior: func(t *tokenMocks.MockUsecase, u *models.User) {
				t.EXPECT().GenerateCSRFToken(u.ID).Return(expectedDefaultCSRF, nil)
			},
			expectedStatus:   200,
			expectedResponse: `{"csrf": "` + expectedDefaultCSRF + `"}`,
			doWrap:           true,
		},
		{
			name:             "No user in request",
			user:             nil,
			mockBehavior:     func(t *tokenMocks.MockUsecase, u *models.User) {},
			expectedStatus:   401,
			expectedResponse: `{"message": "invalid access token"}`,
			doWrap:           false,
		},
		{
			name:             "Nil user in request",
			user:             nil,
			mockBehavior:     func(t *tokenMocks.MockUsecase, u *models.User) {},
			expectedStatus:   401,
			expectedResponse: `{"message": "invalid access token"}`,
			doWrap:           true,
		},
		{
			name: "Failed to get CSRF",
			user: correctTestUser,
			mockBehavior: func(t *tokenMocks.MockUsecase, u *models.User) {
				t.EXPECT().GenerateCSRFToken(u.ID).Return(expectedDefaultCSRF, errors.New("server token error"))
			},
			expectedStatus:   500,
			expectedResponse: `{"message": "failed to get CSRF-token"}`,
			doWrap:           true,
		},
	}

	for _, tc := range testTable {
		t.Run(tc.name, func(t *testing.T) {
			// Call mock
			tc.mockBehavior(tokenMockUsecase, tc.user)

			// Request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/csrf", nil)

			r.ServeHTTP(w, testWrapRequestWithUser(req, tc.user, tc.doWrap))

			// Test
			assert.Equal(t, tc.expectedStatus, w.Code)
			assert.JSONEq(t, tc.expectedResponse, w.Body.String())
		})
	}
}
