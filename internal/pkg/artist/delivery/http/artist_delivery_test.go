package http

import (
	"errors"
	"fmt"
	"net/http"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/golang/mock/gomock"

	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/models"

	commonHttp "github.com/go-park-mail-ru/2023_1_Technokaif/internal/common/http"
	commonTests "github.com/go-park-mail-ru/2023_1_Technokaif/internal/common/tests"
	artistMocks "github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/artist/mocks"
)

var correctUser = models.User{
	ID: 1,
}

func TestArtistDeliveryCreate(t *testing.T) {
	// Init
	type mockBehavior func(au *artistMocks.MockUsecase)

	c := gomock.NewController(t)

	au := artistMocks.NewMockUsecase(c)

	l := commonTests.MockLogger(c)

	h := NewHandler(au, l)

	// Routing
	r := chi.NewRouter()
	r.Post("/api/artists/", h.Create)

	// Test filling
	correctRequestBody := `{
		"name": "YARIK",
		"cover": "/artists/covers/yarik.png"
	}`

	expectedCallArtist := models.Artist{
		Name:      "YARIK",
		UserID:    &correctUser.ID,
		AvatarSrc: "/artists/covers/yarik.png",
	}

	testTable := []struct {
		name             string
		user             *models.User
		requestBody      string
		mockBehavior     mockBehavior
		expectedStatus   int
		expectedResponse string
	}{
		{
			name:        "Common",
			user:        &correctUser,
			requestBody: correctRequestBody,
			mockBehavior: func(tu *artistMocks.MockUsecase) {
				tu.EXPECT().Create(expectedCallArtist).Return(uint32(1), nil)
			},
			expectedStatus:   http.StatusOK,
			expectedResponse: `{"id": 1}`,
		},
		{
			name:             "No User",
			user:             nil,
			mockBehavior:     func(au *artistMocks.MockUsecase) {},
			expectedStatus:   http.StatusUnauthorized,
			expectedResponse: commonTests.ErrorResponse(commonHttp.UnathorizedUser),
		},
		{
			name: "Incorrect JSON",
			user: &correctUser,
			requestBody: `{
				"name":
				"cover": "/artists/covers/yarik.png"
			}`,
			mockBehavior:     func(tu *artistMocks.MockUsecase) {},
			expectedStatus:   http.StatusBadRequest,
			expectedResponse: commonTests.ErrorResponse(commonHttp.IncorrectRequestBody),
		},
		{
			name:             "Incorrect body (no cover)",
			user:             &correctUser,
			requestBody:      `{"name": "YARIK"}`,
			mockBehavior:     func(tu *artistMocks.MockUsecase) {},
			expectedStatus:   http.StatusBadRequest,
			expectedResponse: commonTests.ErrorResponse(commonHttp.IncorrectRequestBody),
		},
		{
			name:        "Server Error",
			user:        &correctUser,
			requestBody: correctRequestBody,
			mockBehavior: func(tu *artistMocks.MockUsecase) {
				tu.EXPECT().Create(expectedCallArtist).Return(uint32(0), errors.New(""))
			},
			expectedStatus:   http.StatusInternalServerError,
			expectedResponse: commonTests.ErrorResponse(artistCreateServerError),
		},
	}

	for _, tc := range testTable {
		t.Run(tc.name, func(t *testing.T) {
			// Call mock
			tc.mockBehavior(au)

			commonTests.DeliveryTestPost(t, r, "/api/artists/", tc.requestBody, tc.expectedStatus, tc.expectedResponse,
				commonTests.WrapRequestWithUserNotNilFunc(tc.user))
		})
	}
}

func TestArtistDeliveryGet(t *testing.T) {
	// Init
	type mockBehavior func(au *artistMocks.MockUsecase)

	c := gomock.NewController(t)

	au := artistMocks.NewMockUsecase(c)

	l := commonTests.MockLogger(c)

	h := NewHandler(au, l)

	// Routing
	r := chi.NewRouter()
	r.Get("/api/artists/{artistID}/", h.Get)

	// Test filling
	const correctArtistID uint32 = 1
	correctArtistIDPath := fmt.Sprint(correctArtistID)

	expectedReturnArtist := models.Artist{
		ID:        1,
		Name:      "Oxxxymiron",
		AvatarSrc: "/artists/avatars/oxxxymiron.png",
	}

	correctResponse := `{
		"id": 1,
		"name": "Oxxxymiron",
		"isLiked": false,
		"cover": "/artists/avatars/oxxxymiron.png"
	}`

	testTable := []struct {
		name             string
		artistIDPath     string
		user             *models.User
		mockBehavior     mockBehavior
		expectedStatus   int
		expectedResponse string
	}{
		{
			name:         "Common",
			artistIDPath: correctArtistIDPath,
			user:         &correctUser,
			mockBehavior: func(au *artistMocks.MockUsecase) {
				au.EXPECT().GetByID(correctArtistID).Return(&expectedReturnArtist, nil)
				au.EXPECT().IsLiked(correctArtistID, correctUser.ID).Return(false, nil)
			},
			expectedStatus:   http.StatusOK,
			expectedResponse: correctResponse,
		},
		{
			name:             "Incorrect ID In Path",
			artistIDPath:     "0",
			mockBehavior:     func(au *artistMocks.MockUsecase) {},
			expectedStatus:   http.StatusBadRequest,
			expectedResponse: commonTests.ErrorResponse(commonHttp.InvalidURLParameter),
		},
		{
			name:         "No Artist To Get",
			artistIDPath: correctArtistIDPath,
			user:         &correctUser,
			mockBehavior: func(au *artistMocks.MockUsecase) {
				au.EXPECT().GetByID(correctArtistID).Return(nil, &models.NoSuchArtistError{})
			},
			expectedStatus:   http.StatusBadRequest,
			expectedResponse: commonTests.ErrorResponse(artistNotFound),
		},
		{
			name:         "Server Error",
			artistIDPath: correctArtistIDPath,
			user:         &correctUser,
			mockBehavior: func(au *artistMocks.MockUsecase) {
				au.EXPECT().GetByID(correctArtistID).Return(nil, errors.New(""))
			},
			expectedStatus:   http.StatusInternalServerError,
			expectedResponse: commonTests.ErrorResponse(artistGetServerError),
		},
	}

	for _, tc := range testTable {
		t.Run(tc.name, func(t *testing.T) {
			// Call mock
			tc.mockBehavior(au)

			commonTests.DeliveryTestGet(t, r, "/api/artists/"+tc.artistIDPath+"/", tc.expectedStatus, tc.expectedResponse,
				commonTests.WrapRequestWithUserNotNilFunc(tc.user))
		})
	}
}

func TestArtistDeliveryDelete(t *testing.T) {
	// Init
	type mockBehavior func(au *artistMocks.MockUsecase)

	c := gomock.NewController(t)

	au := artistMocks.NewMockUsecase(c)

	l := commonTests.MockLogger(c)

	h := NewHandler(au, l)

	// Routing
	r := chi.NewRouter()
	r.Delete("/api/artists/{artistID}/", h.Delete)

	// Test filing
	const correctArtistID uint32 = 1
	correctArtistIDPath := fmt.Sprint(correctArtistID)

	testTable := []struct {
		name             string
		artistIDPath     string
		user             *models.User
		mockBehavior     mockBehavior
		expectedStatus   int
		expectedResponse string
	}{
		{
			name:         "Common",
			artistIDPath: correctArtistIDPath,
			user:         &correctUser,
			mockBehavior: func(au *artistMocks.MockUsecase) {
				au.EXPECT().Delete(correctArtistID, correctUser.ID).Return(nil)
			},
			expectedStatus:   http.StatusOK,
			expectedResponse: commonTests.OKResponse(artistDeletedSuccessfully),
		},
		{
			name:             "Incorrect ID In Path",
			artistIDPath:     "incorrect",
			mockBehavior:     func(au *artistMocks.MockUsecase) {},
			expectedStatus:   http.StatusBadRequest,
			expectedResponse: commonTests.ErrorResponse(commonHttp.InvalidURLParameter),
		},
		{
			name:             "No User",
			artistIDPath:     correctArtistIDPath,
			user:             nil,
			mockBehavior:     func(au *artistMocks.MockUsecase) {},
			expectedStatus:   http.StatusUnauthorized,
			expectedResponse: commonTests.ErrorResponse(commonHttp.UnathorizedUser),
		},
		{
			name:         "User Has No Rights",
			artistIDPath: correctArtistIDPath,
			user:         &correctUser,
			mockBehavior: func(au *artistMocks.MockUsecase) {
				au.EXPECT().Delete(
					correctArtistID,
					correctUser.ID,
				).Return(&models.ForbiddenUserError{})
			},
			expectedStatus:   http.StatusForbidden,
			expectedResponse: commonTests.ErrorResponse(artistDeleteNoRights),
		},
		{
			name:         "No Artist To Delete",
			artistIDPath: correctArtistIDPath,
			user:         &correctUser,
			mockBehavior: func(au *artistMocks.MockUsecase) {
				au.EXPECT().Delete(
					correctArtistID,
					correctUser.ID,
				).Return(&models.NoSuchArtistError{})
			},
			expectedStatus:   http.StatusBadRequest,
			expectedResponse: commonTests.ErrorResponse(artistNotFound),
		},
		{
			name:         "Server Error",
			artistIDPath: correctArtistIDPath,
			user:         &correctUser,
			mockBehavior: func(au *artistMocks.MockUsecase) {
				au.EXPECT().Delete(
					correctArtistID,
					correctUser.ID,
				).Return(errors.New(""))
			},
			expectedStatus:   http.StatusInternalServerError,
			expectedResponse: commonTests.ErrorResponse(artistDeleteServerError),
		},
	}

	for _, tc := range testTable {
		t.Run(tc.name, func(t *testing.T) {
			// Call mock
			tc.mockBehavior(au)

			commonTests.DeliveryTestDelete(t, r, "/api/artists/"+tc.artistIDPath+"/", tc.expectedStatus, tc.expectedResponse,
				commonTests.WrapRequestWithUserNotNilFunc(tc.user))
		})
	}
}

func TestArtistDeliveryFeed(t *testing.T) {
	// Init
	type mockBehavior func(au *artistMocks.MockUsecase)

	c := gomock.NewController(t)

	au := artistMocks.NewMockUsecase(c)

	l := commonTests.MockLogger(c)

	h := NewHandler(au, l)

	// Routing
	r := chi.NewRouter()
	r.Get("/api/artists/feed", h.Feed)

	// Test filling
	expectedReturnArtists := []models.Artist{
		{
			ID:        1,
			Name:      "Oxxxymiron",
			AvatarSrc: "/artists/avatars/oxxxymiron.png",
		},
		{
			ID:        2,
			Name:      "SALUKI",
			AvatarSrc: "/artists/avatars/saluki.png",
		},
		{
			ID:        3,
			Name:      "ATL",
			AvatarSrc: "/artists/avatars/atl.png",
		},
		{
			ID:        4,
			Name:      "104",
			AvatarSrc: "/artists/avatars/104.png",
		},
	}

	correctResponse := `[
		{
			"id": 1,
			"name": "Oxxxymiron",
			"isLiked": false,
			"cover": "/artists/avatars/oxxxymiron.png"
		},
		{
			"id": 2,
			"name": "SALUKI",
			"isLiked": false,
			"cover": "/artists/avatars/saluki.png"
		},
		{
			"id": 3,
			"name": "ATL",
			"isLiked": false,
			"cover": "/artists/avatars/atl.png"
		},
		{
			"id": 4,
			"name": "104",
			"isLiked": false,
			"cover": "/artists/avatars/104.png"
		}
	]`

	testTable := []struct {
		name             string
		mockBehavior     mockBehavior
		expectedStatus   int
		expectedResponse string
	}{
		{
			name: "Common",
			mockBehavior: func(au *artistMocks.MockUsecase) {
				au.EXPECT().GetFeed().Return(expectedReturnArtists, nil)
			},
			expectedStatus:   http.StatusOK,
			expectedResponse: correctResponse,
		},
		{
			name: "No Artists",
			mockBehavior: func(au *artistMocks.MockUsecase) {
				au.EXPECT().GetFeed().Return([]models.Artist{}, nil)
			},
			expectedStatus:   http.StatusOK,
			expectedResponse: `[]`,
		},
		{
			name: "Server Error",
			mockBehavior: func(au *artistMocks.MockUsecase) {
				au.EXPECT().GetFeed().Return(nil, errors.New(""))
			},
			expectedStatus:   http.StatusInternalServerError,
			expectedResponse: commonTests.ErrorResponse(artistsGetServerError),
		},
	}

	for _, tc := range testTable {
		t.Run(tc.name, func(t *testing.T) {
			// Call mock
			tc.mockBehavior(au)

			commonTests.DeliveryTestGet(t, r, "/api/artists/feed", tc.expectedStatus, tc.expectedResponse,
				commonTests.NoWrapUserFunc())
		})
	}
}

func TestArtistDeliveryGetFavorite(t *testing.T) {
	type mockBehavior func(aru *artistMocks.MockUsecase, userID uint32)

	c := gomock.NewController(t)

	au := artistMocks.NewMockUsecase(c)

	l := commonTests.MockLogger(c)

	h := NewHandler(au, l)

	// Routing
	r := chi.NewRouter()
	r.Get("/api/users/{userID}/artists", h.GetFavorite)

	// Test filling
	const correctUserID uint32 = 1
	correctUserIDPath := fmt.Sprint(correctUserID)

	expectedReturnArtists := []models.Artist{
		{
			ID:        1,
			Name:      "Oxxxymiron",
			AvatarSrc: "/artists/avatars/oxxxymiron.png",
		},
		{
			ID:        2,
			Name:      "SALUKI",
			AvatarSrc: "/artists/avatars/saluki.png",
		},
	}

	correctResponse := `[
		{
			"id": 1,
			"name": "Oxxxymiron",
			"isLiked": true,
			"cover": "/artists/avatars/oxxxymiron.png"
		},
		{
			"id": 2,
			"name": "SALUKI",
			"isLiked": true,
			"cover": "/artists/avatars/saluki.png"
		}
	]`

	testTable := []struct {
		name             string
		user             *models.User
		mockBehavior     mockBehavior
		expectedStatus   int
		expectedResponse string
	}{
		{
			name: "Common",
			user: &correctUser,
			mockBehavior: func(au *artistMocks.MockUsecase, userID uint32) {
				au.EXPECT().GetLikedByUser(userID).Return(expectedReturnArtists, nil)
				for _, a := range expectedReturnArtists {
					au.EXPECT().IsLiked(a.ID, userID).Return(true, nil)
				}
			},
			expectedStatus:   http.StatusOK,
			expectedResponse: correctResponse,
		},
		{
			name: "Artists Issue",
			user: &correctUser,
			mockBehavior: func(au *artistMocks.MockUsecase, userID uint32) {
				au.EXPECT().GetLikedByUser(userID).Return(nil, errors.New(""))
			},
			expectedStatus:   http.StatusInternalServerError,
			expectedResponse: commonTests.ErrorResponse(artistsGetServerError),
		},
	}

	for _, tc := range testTable {
		t.Run(tc.name, func(t *testing.T) {
			// Call mock
			tc.mockBehavior(au, tc.user.ID)

			commonTests.DeliveryTestGet(t, r, "/api/users/"+correctUserIDPath+"/artists", tc.expectedStatus, tc.expectedResponse,
				commonTests.WrapRequestWithUserNotNilFunc(tc.user))
		})
	}
}

func TestArtistDeliveryLike(t *testing.T) {
	// Init
	type mockBehavior func(au *artistMocks.MockUsecase)

	c := gomock.NewController(t)

	au := artistMocks.NewMockUsecase(c)

	l := commonTests.MockLogger(c)

	h := NewHandler(au, l)

	// Routing
	r := chi.NewRouter()
	r.Get("/api/artists/{artistID}/like", h.Like)

	// Test filling
	const correctArtistID uint32 = 1
	correctArtistIDPath := fmt.Sprint(correctArtistID)

	testTable := []struct {
		name             string
		artistIDPath     string
		user             *models.User
		mockBehavior     mockBehavior
		expectedStatus   int
		expectedResponse string
	}{
		{
			name:         "Common",
			artistIDPath: correctArtistIDPath,
			user:         &correctUser,
			mockBehavior: func(au *artistMocks.MockUsecase) {
				au.EXPECT().SetLike(correctArtistID, correctUser.ID).Return(true, nil)
			},
			expectedStatus:   http.StatusOK,
			expectedResponse: commonTests.OKResponse(commonHttp.LikeSuccess),
		},
		{
			name:         "Already liked (Anyway Success)",
			artistIDPath: correctArtistIDPath,
			user:         &correctUser,
			mockBehavior: func(au *artistMocks.MockUsecase) {
				au.EXPECT().SetLike(correctArtistID, correctUser.ID).Return(false, nil)
			},
			expectedStatus:   http.StatusOK,
			expectedResponse: commonTests.OKResponse(commonHttp.LikeAlreadyExists),
		},
		{
			name:             "Incorrect ID In Path",
			artistIDPath:     "0",
			user:             &correctUser,
			mockBehavior:     func(au *artistMocks.MockUsecase) {},
			expectedStatus:   http.StatusBadRequest,
			expectedResponse: commonTests.ErrorResponse(commonHttp.InvalidURLParameter),
		},
		{
			name:             "No User",
			artistIDPath:     correctArtistIDPath,
			user:             nil,
			mockBehavior:     func(au *artistMocks.MockUsecase) {},
			expectedStatus:   http.StatusUnauthorized,
			expectedResponse: commonTests.ErrorResponse(commonHttp.UnathorizedUser),
		},
		{
			name:         "No Artist To Like",
			artistIDPath: correctArtistIDPath,
			user:         &correctUser,
			mockBehavior: func(au *artistMocks.MockUsecase) {
				au.EXPECT().SetLike(correctArtistID, correctUser.ID).Return(false, &models.NoSuchArtistError{})
			},
			expectedStatus:   http.StatusBadRequest,
			expectedResponse: commonTests.ErrorResponse(artistNotFound),
		},
		{
			name:         "Server Error",
			artistIDPath: correctArtistIDPath,
			user:         &correctUser,
			mockBehavior: func(au *artistMocks.MockUsecase) {
				au.EXPECT().SetLike(correctArtistID, correctUser.ID).Return(false, errors.New(""))
			},
			expectedStatus:   http.StatusInternalServerError,
			expectedResponse: commonTests.ErrorResponse(commonHttp.SetLikeServerError),
		},
	}

	for _, tc := range testTable {
		t.Run(tc.name, func(t *testing.T) {
			// Call mock
			tc.mockBehavior(au)

			commonTests.DeliveryTestGet(t, r, "/api/artists/"+tc.artistIDPath+"/like", tc.expectedStatus, tc.expectedResponse,
				commonTests.WrapRequestWithUserNotNilFunc(tc.user))
		})
	}
}

func TestArtistDeliveryUnLike(t *testing.T) {
	// Init
	type mockBehavior func(au *artistMocks.MockUsecase)

	c := gomock.NewController(t)

	au := artistMocks.NewMockUsecase(c)

	l := commonTests.MockLogger(c)

	h := NewHandler(au, l)

	// Routing
	r := chi.NewRouter()
	r.Get("/api/artists/{artistID}/unlike", h.UnLike)

	// Test filling
	const correctArtistID uint32 = 1
	correctArtistIDPath := fmt.Sprint(correctArtistID)

	testTable := []struct {
		name             string
		artistIDPath     string
		user             *models.User
		mockBehavior     mockBehavior
		expectedStatus   int
		expectedResponse string
	}{
		{
			name:         "Common",
			artistIDPath: correctArtistIDPath,
			user:         &correctUser,
			mockBehavior: func(au *artistMocks.MockUsecase) {
				au.EXPECT().UnLike(correctArtistID, correctUser.ID).Return(true, nil)
			},
			expectedStatus:   http.StatusOK,
			expectedResponse: commonTests.OKResponse(commonHttp.UnLikeSuccess),
		},
		{
			name:         "Wasn't Liked (Anyway Success)",
			artistIDPath: correctArtistIDPath,
			user:         &correctUser,
			mockBehavior: func(au *artistMocks.MockUsecase) {
				au.EXPECT().UnLike(correctArtistID, correctUser.ID).Return(false, nil)
			},
			expectedStatus:   http.StatusOK,
			expectedResponse: commonTests.OKResponse(commonHttp.LikeDoesntExist),
		},
		{
			name:             "Incorrect ID In Path",
			artistIDPath:     "0",
			user:             &correctUser,
			mockBehavior:     func(au *artistMocks.MockUsecase) {},
			expectedStatus:   http.StatusBadRequest,
			expectedResponse: commonTests.ErrorResponse(commonHttp.InvalidURLParameter),
		},
		{
			name:             "No User",
			artistIDPath:     correctArtistIDPath,
			user:             nil,
			mockBehavior:     func(au *artistMocks.MockUsecase) {},
			expectedStatus:   http.StatusUnauthorized,
			expectedResponse: commonTests.ErrorResponse(commonHttp.UnathorizedUser),
		},
		{
			name:         "No Artist To Unlike",
			artistIDPath: correctArtistIDPath,
			user:         &correctUser,
			mockBehavior: func(au *artistMocks.MockUsecase) {
				au.EXPECT().UnLike(correctArtistID, correctUser.ID).Return(false, &models.NoSuchArtistError{})
			},
			expectedStatus:   http.StatusBadRequest,
			expectedResponse: commonTests.ErrorResponse(artistNotFound),
		},
		{
			name:         "Server Error",
			artistIDPath: correctArtistIDPath,
			user:         &correctUser,
			mockBehavior: func(au *artistMocks.MockUsecase) {
				au.EXPECT().UnLike(correctArtistID, correctUser.ID).Return(false, errors.New(""))
			},
			expectedStatus:   http.StatusInternalServerError,
			expectedResponse: commonTests.ErrorResponse(commonHttp.DeleteLikeServerError),
		},
	}

	for _, tc := range testTable {
		t.Run(tc.name, func(t *testing.T) {
			// Call mock
			tc.mockBehavior(au)

			commonTests.DeliveryTestGet(t, r, "/api/artists/"+tc.artistIDPath+"/unlike", tc.expectedStatus, tc.expectedResponse,
				commonTests.WrapRequestWithUserNotNilFunc(tc.user))
		})
	}
}
