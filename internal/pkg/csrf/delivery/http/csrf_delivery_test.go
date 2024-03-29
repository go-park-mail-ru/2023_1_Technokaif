package delivery

import (
	"errors"
	"net/http"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/golang/mock/gomock"

	commonTests "github.com/go-park-mail-ru/2023_1_Technokaif/internal/common/tests"
	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/models"
	tokenMocks "github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/token/mocks"
)

func TestCSRFDeliveryHTTP_GetCSRF(t *testing.T) {
	// Init
	type mockBehavior func(t *tokenMocks.MockUsecase, u *models.User)

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
			expectedStatus:   http.StatusOK,
			expectedResponse: `{"csrf": "` + expectedDefaultCSRF + `"}`,
			doWrap:           true,
		},
		{
			name:             "No user in request",
			user:             nil,
			mockBehavior:     func(t *tokenMocks.MockUsecase, u *models.User) {},
			expectedStatus:   http.StatusUnauthorized,
			expectedResponse: commonTests.ErrorResponse(invalidAccessToken),
			doWrap:           false,
		},
		{
			name:             "Nil user in request",
			user:             nil,
			mockBehavior:     func(t *tokenMocks.MockUsecase, u *models.User) {},
			expectedStatus:   http.StatusUnauthorized,
			expectedResponse: commonTests.ErrorResponse(invalidAccessToken),
			doWrap:           true,
		},
		{
			name: "Failed to get CSRF",
			user: correctTestUser,
			mockBehavior: func(t *tokenMocks.MockUsecase, u *models.User) {
				t.EXPECT().GenerateCSRFToken(u.ID).Return(expectedDefaultCSRF, errors.New("server token error"))
			},
			expectedStatus:   http.StatusInternalServerError,
			expectedResponse: commonTests.ErrorResponse(csrfGetError),
			doWrap:           true,
		},
	}

	for _, tc := range testTable {
		t.Run(tc.name, func(t *testing.T) {
			// Call mock
			tc.mockBehavior(tokenMockUsecase, tc.user)

			commonTests.DeliveryTestGet(t, r, "/csrf", tc.expectedStatus, tc.expectedResponse,
				commonTests.WrapRequestWithUserFunc(tc.user, tc.doWrap))
		})
	}
}
