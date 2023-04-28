package http

import (
	"errors"
	"fmt"
	"net/http"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/golang/mock/gomock"

	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/models"

	commonHTTP "github.com/go-park-mail-ru/2023_1_Technokaif/internal/common/http"
	commonTests "github.com/go-park-mail-ru/2023_1_Technokaif/internal/common/tests"
	artistMocks "github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/artist/mocks"
	trackMocks "github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/track/mocks"
)

var correctUser = models.User{
	ID: 1,
}

const correctTrackID uint32 = 1

var correctTrackIDPath = fmt.Sprint(correctTrackID)

func TestTrackDeliveryHTTP_Create(t *testing.T) {
	// Init
	type mockBehavior func(tu *trackMocks.MockUsecase)

	c := gomock.NewController(t)

	tu := trackMocks.NewMockUsecase(c)
	au := artistMocks.NewMockUsecase(c)

	l := commonTests.MockLogger(c)

	h := NewHandler(tu, au, l)

	// Routing
	r := chi.NewRouter()
	r.Post("/api/tracks/", h.Create)

	// Test filling
	correctRequestBody := `{
		"name": "Хит",
		"artistsID": [1],
		"record": "/tracks/records/hit.wav"
	}`

	correctArtistsID := []uint32{1}

	expectedCallTrack := models.Track{
		Name:      "Хит",
		RecordSrc: "/tracks/records/hit.wav",
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
			mockBehavior: func(tu *trackMocks.MockUsecase) {
				tu.EXPECT().Create(
					gomock.Any(), expectedCallTrack, correctArtistsID, correctUser.ID,
				).Return(uint32(1), nil)
			},
			expectedStatus:   http.StatusOK,
			expectedResponse: `{"id": 1}`,
		},
		{
			name:             "No User",
			user:             nil,
			mockBehavior:     func(tu *trackMocks.MockUsecase) {},
			expectedStatus:   http.StatusUnauthorized,
			expectedResponse: commonTests.ErrorResponse(commonHTTP.UnathorizedUser),
		},
		{
			name: "Incorrect JSON",
			user: &correctUser,
			requestBody: `{
				"name":
				"artistsID": [1],
				"cover": "/tracks/covers/hit.png"
			}`,
			mockBehavior:     func(tu *trackMocks.MockUsecase) {},
			expectedStatus:   http.StatusBadRequest,
			expectedResponse: `{"message": "incorrect input body"}`,
		},
		{
			name: "Incorrect Body (no source)",
			user: &correctUser,
			requestBody: `{
				"name": "Хит",
				"artistsID": [1],
				"cover": "/tracks/covers/hit.png"
			}`,
			mockBehavior:     func(tu *trackMocks.MockUsecase) {},
			expectedStatus:   http.StatusBadRequest,
			expectedResponse: `{"message": "incorrect input body"}`,
		},
		{
			name: "Incorrect Body (albumID w/o albumPosition)",
			user: &correctUser,
			requestBody: `{
				"name": "Хит",
				"artistsID": [1],
				"albumID": 1,
				"cover": "/tracks/covers/gorgorod.png"
			}`,
			mockBehavior:     func(au *trackMocks.MockUsecase) {},
			expectedStatus:   http.StatusBadRequest,
			expectedResponse: `{"message": "incorrect input body"}`,
		},
		{
			name:        "User Has No Rights",
			user:        &correctUser,
			requestBody: correctRequestBody,
			mockBehavior: func(tu *trackMocks.MockUsecase) {
				tu.EXPECT().Create(
					gomock.Any(), expectedCallTrack, correctArtistsID, correctUser.ID,
				).Return(uint32(0), &models.ForbiddenUserError{})
			},
			expectedStatus:   http.StatusForbidden,
			expectedResponse: `{"message": "no rights to create track"}`,
		},
		{
			name:        "Server Error",
			user:        &correctUser,
			requestBody: correctRequestBody,
			mockBehavior: func(tu *trackMocks.MockUsecase) {
				tu.EXPECT().Create(
					gomock.Any(), expectedCallTrack, correctArtistsID, correctUser.ID,
				).Return(uint32(0), errors.New(""))
			},
			expectedStatus:   http.StatusInternalServerError,
			expectedResponse: `{"message": "can't create track"}`,
		},
	}

	for _, tc := range testTable {
		t.Run(tc.name, func(t *testing.T) {
			// Call mock
			tc.mockBehavior(tu)

			commonTests.DeliveryTestPost(t, r, "/api/tracks/",
				tc.requestBody, tc.expectedStatus, tc.expectedResponse,
				commonTests.WrapRequestWithUserNotNilFunc(tc.user))
		})
	}
}

func TestTrackDeliveryHTTP_Get(t *testing.T) {
	// Init
	type mockBehavior func(tu *trackMocks.MockUsecase, au *artistMocks.MockUsecase)

	c := gomock.NewController(t)

	tu := trackMocks.NewMockUsecase(c)
	au := artistMocks.NewMockUsecase(c)

	l := commonTests.MockLogger(c)

	h := NewHandler(tu, au, l)

	// Routing
	r := chi.NewRouter()
	r.Get("/api/tracks/{trackID}/", h.Get)

	// Test filling
	expectedReturnTrack := models.Track{
		ID:        correctTrackID,
		Name:      "Хит",
		CoverSrc:  "/tracks/covers/hit.png",
		Listens:   99999999,
		Duration:  180,
		RecordSrc: "/tracks/records/hit.wav",
	}

	expectedReturnArtists := []models.Artist{
		{
			ID:        1,
			Name:      "Oxxxymiron",
			AvatarSrc: "/artists/avatars/oxxxymiron.png",
		},
	}

	correctResponse := `{
		"id": 1,
		"name": "Хит",
		"artists": [
			{
				"id": 1,
				"name": "Oxxxymiron",
				"isLiked": false,
				"cover": "/artists/avatars/oxxxymiron.png"
			}
		],
		"cover": "/tracks/covers/hit.png",
		"listens": 99999999,
		"isLiked": false,
		"duration": 180,
		"recordSrc": "/tracks/records/hit.wav"
	}`

	testTable := []struct {
		name             string
		trackIDPath      string
		user             *models.User
		mockBehavior     mockBehavior
		expectedStatus   int
		expectedResponse string
	}{
		{
			name:        "Common",
			trackIDPath: correctTrackIDPath,
			user:        &correctUser,
			mockBehavior: func(tu *trackMocks.MockUsecase, au *artistMocks.MockUsecase) {
				tu.EXPECT().GetByID(gomock.Any(), correctTrackID).Return(&expectedReturnTrack, nil)
				tu.EXPECT().IsLiked(gomock.Any(), correctTrackID, correctUser.ID).Return(false, nil)
				au.EXPECT().GetByTrack(gomock.Any(), correctTrackID).Return(expectedReturnArtists, nil)
				for _, a := range expectedReturnArtists {
					au.EXPECT().IsLiked(gomock.Any(), a.ID, correctUser.ID).Return(false, nil)
				}
			},
			expectedStatus:   http.StatusOK,
			expectedResponse: correctResponse,
		},
		{
			name:             "Incorrect ID In Path",
			trackIDPath:      "-5",
			mockBehavior:     func(tu *trackMocks.MockUsecase, au *artistMocks.MockUsecase) {},
			expectedStatus:   http.StatusBadRequest,
			expectedResponse: commonTests.ErrorResponse(commonHTTP.InvalidURLParameter),
		},
		{
			name:        "No Track To Get",
			trackIDPath: correctTrackIDPath,
			user:        &correctUser,
			mockBehavior: func(tu *trackMocks.MockUsecase, au *artistMocks.MockUsecase) {
				tu.EXPECT().GetByID(gomock.Any(), correctTrackID).Return(nil, &models.NoSuchTrackError{})
			},
			expectedStatus:   http.StatusBadRequest,
			expectedResponse: `{"message": "no such track"}`,
		},
		{
			name:        "Tracks Issues",
			trackIDPath: correctTrackIDPath,
			user:        &correctUser,
			mockBehavior: func(tu *trackMocks.MockUsecase, au *artistMocks.MockUsecase) {
				tu.EXPECT().GetByID(gomock.Any(), correctTrackID).Return(nil, errors.New(""))
			},
			expectedStatus:   http.StatusInternalServerError,
			expectedResponse: `{"message": "can't get track"}`,
		},
		{
			name:        "Artists Issues",
			trackIDPath: correctTrackIDPath,
			user:        &correctUser,
			mockBehavior: func(tu *trackMocks.MockUsecase, au *artistMocks.MockUsecase) {
				tu.EXPECT().GetByID(gomock.Any(), correctTrackID).Return(&expectedReturnTrack, nil)
				au.EXPECT().GetByTrack(gomock.Any(), correctTrackID).Return(nil, errors.New(""))
			},
			expectedStatus:   http.StatusInternalServerError,
			expectedResponse: `{"message": "can't get track"}`,
		},
	}

	for _, tc := range testTable {
		t.Run(tc.name, func(t *testing.T) {
			// Call mock
			tc.mockBehavior(tu, au)

			commonTests.DeliveryTestGet(t, r,
				"/api/tracks/"+tc.trackIDPath+"/", tc.expectedStatus, tc.expectedResponse,
				commonTests.WrapRequestWithUserNotNilFunc(tc.user))
		})
	}
}

func TestTrackDeliveryHTTP_Delete(t *testing.T) {
	// Init
	type mockBehavior func(au *trackMocks.MockUsecase)

	c := gomock.NewController(t)

	tu := trackMocks.NewMockUsecase(c)
	au := artistMocks.NewMockUsecase(c)

	l := commonTests.MockLogger(c)

	h := NewHandler(tu, au, l)

	// Routing
	r := chi.NewRouter()
	r.Delete("/api/tracks/{trackID}/", h.Delete)

	// Test filling
	testTable := []struct {
		name             string
		trackIDPath      string
		user             *models.User
		mockBehavior     mockBehavior
		expectedStatus   int
		expectedResponse string
	}{
		{
			name:        "Common",
			trackIDPath: correctTrackIDPath,
			user:        &correctUser,
			mockBehavior: func(au *trackMocks.MockUsecase) {
				au.EXPECT().Delete(gomock.Any(), correctTrackID, correctUser.ID).Return(nil)
			},
			expectedStatus:   http.StatusOK,
			expectedResponse: commonTests.OKResponse(trackDeletedSuccessfully),
		},
		{
			name:             "Incorrect ID In Path",
			trackIDPath:      "incorrect",
			mockBehavior:     func(au *trackMocks.MockUsecase) {},
			expectedStatus:   http.StatusBadRequest,
			expectedResponse: commonTests.ErrorResponse(commonHTTP.InvalidURLParameter),
		},
		{
			name:             "No User",
			trackIDPath:      correctTrackIDPath,
			user:             nil,
			mockBehavior:     func(au *trackMocks.MockUsecase) {},
			expectedStatus:   http.StatusUnauthorized,
			expectedResponse: commonTests.ErrorResponse(commonHTTP.UnathorizedUser),
		},
		{
			name:        "User Has No Rights",
			trackIDPath: correctTrackIDPath,
			user:        &correctUser,
			mockBehavior: func(au *trackMocks.MockUsecase) {
				au.EXPECT().Delete(
					gomock.Any(), correctTrackID, correctUser.ID,
				).Return(&models.ForbiddenUserError{})
			},
			expectedStatus:   http.StatusForbidden,
			expectedResponse: commonTests.ErrorResponse(trackDeleteNoRights),
		},
		{
			name:        "No Track To Delete",
			trackIDPath: correctTrackIDPath,
			user:        &correctUser,
			mockBehavior: func(au *trackMocks.MockUsecase) {
				au.EXPECT().Delete(
					gomock.Any(), correctTrackID, correctUser.ID,
				).Return(&models.NoSuchTrackError{})
			},
			expectedStatus:   http.StatusBadRequest,
			expectedResponse: commonTests.ErrorResponse(trackNotFound),
		},
		{
			name:        "Server Error",
			trackIDPath: correctTrackIDPath,
			user:        &correctUser,
			mockBehavior: func(au *trackMocks.MockUsecase) {
				au.EXPECT().Delete(
					gomock.Any(), correctTrackID, correctUser.ID,
				).Return(errors.New(""))
			},
			expectedStatus:   http.StatusInternalServerError,
			expectedResponse: commonTests.ErrorResponse(trackDeleteServerError),
		},
	}

	for _, tc := range testTable {
		t.Run(tc.name, func(t *testing.T) {
			// Call mock
			tc.mockBehavior(tu)

			commonTests.DeliveryTestDelete(t, r, "/api/tracks/"+tc.trackIDPath+"/",
				tc.expectedStatus, tc.expectedResponse,
				commonTests.WrapRequestWithUserNotNilFunc(tc.user))
		})
	}
}

func TestTrackDeliveryHTTP_Feed(t *testing.T) {
	// Init
	type mockBehavior func(tu *trackMocks.MockUsecase, au *artistMocks.MockUsecase)

	c := gomock.NewController(t)

	tu := trackMocks.NewMockUsecase(c)
	au := artistMocks.NewMockUsecase(c)

	l := commonTests.MockLogger(c)

	h := NewHandler(tu, au, l)

	// Routing
	r := chi.NewRouter()
	r.Get("/api/tracks/feed", h.Feed)

	// Test filling
	expectedReturnTracks := []models.Track{
		{
			ID:        1,
			Name:      "Накануне",
			CoverSrc:  "/tracks/covers/1.png",
			Listens:   2700000,
			Duration:  180,
			RecordSrc: "/tracks/records/1.wav",
		},
		{
			ID:        2,
			Name:      "LAGG OUT",
			CoverSrc:  "/tracks/covers/2.png",
			Listens:   4500000,
			Duration:  180,
			RecordSrc: "/tracks/records/2.wav",
		},
	}

	expectedReturnArtists := []models.Artist{
		{
			ID:        1,
			Name:      "Oxxxymiron",
			AvatarSrc: "/artists/avatars/1.png",
		},
		{
			ID:        2,
			Name:      "SALUKI",
			AvatarSrc: "/artists/avatars/2.png",
		},
		{
			ID:        3,
			Name:      "ATL",
			AvatarSrc: "/artists/avatars/3.png",
		},
	}

	correctResponse := `[
		{
			"id": 1,
			"name": "Накануне",
			"artists": [
				{
					"id": 1,
					"name": "Oxxxymiron",
					"isLiked": false,
					"cover": "/artists/avatars/1.png"
				}
			],
			"cover": "/tracks/covers/1.png",
			"listens": 2700000,
			"isLiked": false,
			"duration": 180,
			"recordSrc": "/tracks/records/1.wav"
		},
		{
			"id": 2,
			"name": "LAGG OUT",
			"artists": [
				{
					"id": 2,
					"name": "SALUKI",
					"isLiked": false,
					"cover": "/artists/avatars/2.png"
				},
				{
					"id": 3,
					"name": "ATL",
					"isLiked": false,
					"cover": "/artists/avatars/3.png"
				}
			],
			"cover": "/tracks/covers/2.png",
			"listens": 4500000,
			"isLiked": false,
			"duration": 180,
			"recordSrc": "/tracks/records/2.wav"
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
			mockBehavior: func(tu *trackMocks.MockUsecase, au *artistMocks.MockUsecase) {
				tu.EXPECT().GetFeed(gomock.Any()).Return(expectedReturnTracks, nil)
				au.EXPECT().GetByTrack(gomock.Any(), expectedReturnTracks[0].ID).Return(expectedReturnArtists[0:1], nil)
				au.EXPECT().GetByTrack(gomock.Any(), expectedReturnTracks[1].ID).Return(expectedReturnArtists[1:3], nil)
			},
			expectedStatus:   http.StatusOK,
			expectedResponse: correctResponse,
		},
		{
			name: "No Tracks",
			mockBehavior: func(tu *trackMocks.MockUsecase, au *artistMocks.MockUsecase) {
				tu.EXPECT().GetFeed(gomock.Any()).Return([]models.Track{}, nil)
			},
			expectedStatus:   http.StatusOK,
			expectedResponse: `[]`,
		},
		{
			name: "Tracks Issues",
			mockBehavior: func(tu *trackMocks.MockUsecase, au *artistMocks.MockUsecase) {
				tu.EXPECT().GetFeed(gomock.Any()).Return(nil, errors.New(""))
			},
			expectedStatus:   http.StatusInternalServerError,
			expectedResponse: commonTests.ErrorResponse(tracksGetServerError),
		},
		{
			name: "Artists Issues",
			mockBehavior: func(tu *trackMocks.MockUsecase, au *artistMocks.MockUsecase) {
				tu.EXPECT().GetFeed(gomock.Any()).Return(expectedReturnTracks, nil)
				au.EXPECT().GetByTrack(gomock.Any(), expectedReturnTracks[0].ID).Return(nil, errors.New(""))
			},
			expectedStatus:   http.StatusInternalServerError,
			expectedResponse: commonTests.ErrorResponse(tracksGetServerError),
		},
	}

	for _, tc := range testTable {
		t.Run(tc.name, func(t *testing.T) {
			// Call mock
			tc.mockBehavior(tu, au)

			commonTests.DeliveryTestGet(t, r, "/api/tracks/feed",
				tc.expectedStatus, tc.expectedResponse,
				commonTests.NoWrapUserFunc())
		})
	}
}

func TestTrackDeliveryHTTP_GetFavorite(t *testing.T) {
	type mockBehavior func(tu *trackMocks.MockUsecase, au *artistMocks.MockUsecase, userID uint32)

	c := gomock.NewController(t)

	tu := trackMocks.NewMockUsecase(c)
	aru := artistMocks.NewMockUsecase(c)

	l := commonTests.MockLogger(c)

	h := NewHandler(tu, aru, l)

	// Routing
	r := chi.NewRouter()
	r.Get("/api/users/{userID}/tracks", h.GetFavorite)

	// Test filling
	const correctUserID uint32 = 1
	correctUserIDPath := fmt.Sprint(correctUserID)

	expectedReturnTracks := []models.Track{
		{
			ID:        1,
			Name:      "Накануне",
			CoverSrc:  "/tracks/covers/1.png",
			Listens:   2700000,
			Duration:  180,
			RecordSrc: "/tracks/records/1.wav",
		},
		{
			ID:        2,
			Name:      "LAGG OUT",
			CoverSrc:  "/tracks/covers/2.png",
			Listens:   4500000,
			Duration:  180,
			RecordSrc: "/tracks/records/2.wav",
		},
	}

	expectedReturnArtists := []models.Artist{
		{
			ID:        1,
			Name:      "Oxxxymiron",
			AvatarSrc: "/artists/avatars/1.png",
		},
		{
			ID:        2,
			Name:      "SALUKI",
			AvatarSrc: "/artists/avatars/2.png",
		},
	}

	correctResponse := `[
		{
			"id": 1,
			"name": "Накануне",
			"artists": [
				{
					"id": 1,
					"name": "Oxxxymiron",
					"isLiked": false,
					"cover": "/artists/avatars/1.png"
				}
			],
			"cover": "/tracks/covers/1.png",
			"listens": 2700000,
			"isLiked": true,
			"duration": 180,
			"recordSrc": "/tracks/records/1.wav"
		},
		{
			"id": 2,
			"name": "LAGG OUT",
			"artists": [
				{
					"id": 2,
					"name": "SALUKI",
					"isLiked": false,
					"cover": "/artists/avatars/2.png"
				}
			],
			"cover": "/tracks/covers/2.png",
			"listens": 4500000,
			"isLiked": true,
			"duration": 180,
			"recordSrc": "/tracks/records/2.wav"
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
			mockBehavior: func(tu *trackMocks.MockUsecase, au *artistMocks.MockUsecase, userID uint32) {
				tu.EXPECT().GetLikedByUser(gomock.Any(), userID).Return(expectedReturnTracks, nil)
				for ind, track := range expectedReturnTracks {
					au.EXPECT().GetByTrack(gomock.Any(), track.ID).Return(expectedReturnArtists[ind:ind+1], nil)
					tu.EXPECT().IsLiked(gomock.Any(), track.ID, userID).Return(true, nil)
					for _, a := range expectedReturnArtists[ind : ind+1] {
						au.EXPECT().IsLiked(gomock.Any(), a.ID, correctUserID)
					}
				}
			},
			expectedStatus:   http.StatusOK,
			expectedResponse: correctResponse,
		},
		{
			name: "Tracks Issue",
			user: &correctUser,
			mockBehavior: func(tu *trackMocks.MockUsecase, au *artistMocks.MockUsecase, userID uint32) {
				tu.EXPECT().GetLikedByUser(gomock.Any(), userID).Return(nil, errors.New(""))
			},
			expectedStatus:   http.StatusInternalServerError,
			expectedResponse: commonTests.ErrorResponse(tracksGetServerError),
		},
		{
			name: "Artists Issue",
			user: &correctUser,
			mockBehavior: func(tu *trackMocks.MockUsecase, au *artistMocks.MockUsecase, userID uint32) {
				tu.EXPECT().GetLikedByUser(gomock.Any(), userID).Return(expectedReturnTracks, nil)
				au.EXPECT().GetByTrack(gomock.Any(), expectedReturnTracks[0].ID).Return(nil, errors.New(""))
			},
			expectedStatus:   http.StatusInternalServerError,
			expectedResponse: commonTests.ErrorResponse(tracksGetServerError),
		},
		{
			name: "Likes Issue",
			user: &correctUser,
			mockBehavior: func(tu *trackMocks.MockUsecase, au *artistMocks.MockUsecase, userID uint32) {
				tu.EXPECT().GetLikedByUser(gomock.Any(), userID).Return(expectedReturnTracks, nil)
				au.EXPECT().GetByTrack(gomock.Any(), expectedReturnTracks[0].ID).Return(expectedReturnArtists[0:1], nil)
				tu.EXPECT().IsLiked(gomock.Any(), expectedReturnTracks[0].ID, userID).Return(false, errors.New(""))
			},
			expectedStatus:   http.StatusInternalServerError,
			expectedResponse: commonTests.ErrorResponse(tracksGetServerError),
		},
	}

	for _, tc := range testTable {
		t.Run(tc.name, func(t *testing.T) {
			// Call mock
			tc.mockBehavior(tu, aru, tc.user.ID)

			commonTests.DeliveryTestGet(t, r, "/api/users/"+correctUserIDPath+"/tracks",
				tc.expectedStatus, tc.expectedResponse,
				commonTests.WrapRequestWithUserNotNilFunc(tc.user))
		})
	}
}

func TestTrackDeliveryHTTP_Like(t *testing.T) {
	// Init
	type mockBehavior func(tu *trackMocks.MockUsecase)

	c := gomock.NewController(t)

	tu := trackMocks.NewMockUsecase(c)
	au := artistMocks.NewMockUsecase(c)

	l := commonTests.MockLogger(c)

	h := NewHandler(tu, au, l)

	// Routing
	r := chi.NewRouter()
	r.Get("/api/tracks/{trackID}/like", h.Like)

	// Test filling
	testTable := []struct {
		name             string
		trackIDPath      string
		user             *models.User
		mockBehavior     mockBehavior
		expectedStatus   int
		expectedResponse string
	}{
		{
			name:        "Common",
			trackIDPath: correctTrackIDPath,
			user:        &correctUser,
			mockBehavior: func(tu *trackMocks.MockUsecase) {
				tu.EXPECT().SetLike(gomock.Any(), correctTrackID, correctUser.ID).Return(true, nil)
			},
			expectedStatus:   http.StatusOK,
			expectedResponse: commonTests.OKResponse(commonHTTP.LikeSuccess),
		},
		{
			name:        "Already Liked (Anyway Success)",
			trackIDPath: correctTrackIDPath,
			user:        &correctUser,
			mockBehavior: func(tu *trackMocks.MockUsecase) {
				tu.EXPECT().SetLike(gomock.Any(), correctTrackID, correctUser.ID).Return(false, nil)
			},
			expectedStatus:   http.StatusOK,
			expectedResponse: commonTests.OKResponse(commonHTTP.LikeAlreadyExists),
		},
		{
			name:             "Incorrect ID In Path",
			trackIDPath:      "0",
			user:             &correctUser,
			mockBehavior:     func(tu *trackMocks.MockUsecase) {},
			expectedStatus:   http.StatusBadRequest,
			expectedResponse: commonTests.ErrorResponse(commonHTTP.InvalidURLParameter),
		},
		{
			name:             "No User",
			trackIDPath:      correctTrackIDPath,
			user:             nil,
			mockBehavior:     func(tu *trackMocks.MockUsecase) {},
			expectedStatus:   http.StatusUnauthorized,
			expectedResponse: commonTests.ErrorResponse(commonHTTP.UnathorizedUser),
		},
		{
			name:        "No Track To Like",
			trackIDPath: correctTrackIDPath,
			user:        &correctUser,
			mockBehavior: func(tu *trackMocks.MockUsecase) {
				tu.EXPECT().SetLike(gomock.Any(), correctTrackID, correctUser.ID).
					Return(false, &models.NoSuchTrackError{})
			},
			expectedStatus:   http.StatusBadRequest,
			expectedResponse: commonTests.ErrorResponse(trackNotFound),
		},
		{
			name:        "Server Error",
			trackIDPath: correctTrackIDPath,
			user:        &correctUser,
			mockBehavior: func(tu *trackMocks.MockUsecase) {
				tu.EXPECT().SetLike(gomock.Any(), correctTrackID, correctUser.ID).
					Return(false, errors.New(""))
			},
			expectedStatus:   http.StatusInternalServerError,
			expectedResponse: commonTests.ErrorResponse(commonHTTP.SetLikeServerError),
		},
	}

	for _, tc := range testTable {
		t.Run(tc.name, func(t *testing.T) {
			// Call mock
			tc.mockBehavior(tu)

			commonTests.DeliveryTestGet(t, r, "/api/tracks/"+tc.trackIDPath+"/like",
				tc.expectedStatus, tc.expectedResponse,
				commonTests.WrapRequestWithUserNotNilFunc(tc.user))
		})
	}
}

func TestTrackDeliveryHTTP_UnLike(t *testing.T) {
	// Init
	type mockBehavior func(tu *trackMocks.MockUsecase)

	c := gomock.NewController(t)

	tu := trackMocks.NewMockUsecase(c)
	au := artistMocks.NewMockUsecase(c)

	l := commonTests.MockLogger(c)

	h := NewHandler(tu, au, l)

	// Routing
	r := chi.NewRouter()
	r.Get("/api/tracks/{trackID}/unlike", h.UnLike)

	// Test filling
	testTable := []struct {
		name             string
		trackIDPath      string
		user             *models.User
		mockBehavior     mockBehavior
		expectedStatus   int
		expectedResponse string
	}{
		{
			name:        "Common",
			trackIDPath: correctTrackIDPath,
			user:        &correctUser,
			mockBehavior: func(tu *trackMocks.MockUsecase) {
				tu.EXPECT().UnLike(gomock.Any(), correctTrackID, correctUser.ID).Return(true, nil)
			},
			expectedStatus:   http.StatusOK,
			expectedResponse: commonTests.OKResponse(commonHTTP.UnLikeSuccess),
		},
		{
			name:        "Wasn't Liked (Anyway Success)",
			trackIDPath: correctTrackIDPath,
			user:        &correctUser,
			mockBehavior: func(tu *trackMocks.MockUsecase) {
				tu.EXPECT().UnLike(gomock.Any(), correctTrackID, correctUser.ID).Return(false, nil)
			},
			expectedStatus:   http.StatusOK,
			expectedResponse: commonTests.OKResponse(commonHTTP.LikeDoesntExist),
		},
		{
			name:             "Incorrect ID In Path",
			trackIDPath:      "0",
			user:             &correctUser,
			mockBehavior:     func(tu *trackMocks.MockUsecase) {},
			expectedStatus:   http.StatusBadRequest,
			expectedResponse: commonTests.ErrorResponse(commonHTTP.InvalidURLParameter),
		},
		{
			name:             "No User",
			trackIDPath:      correctTrackIDPath,
			user:             nil,
			mockBehavior:     func(tu *trackMocks.MockUsecase) {},
			expectedStatus:   http.StatusUnauthorized,
			expectedResponse: commonTests.ErrorResponse(commonHTTP.UnathorizedUser),
		},
		{
			name:        "No Track To Unlike",
			trackIDPath: correctTrackIDPath,
			user:        &correctUser,
			mockBehavior: func(tu *trackMocks.MockUsecase) {
				tu.EXPECT().UnLike(gomock.Any(), correctTrackID, correctUser.ID).
					Return(false, &models.NoSuchTrackError{})
			},
			expectedStatus:   http.StatusBadRequest,
			expectedResponse: commonTests.ErrorResponse(trackNotFound),
		},
		{
			name:        "Server Error",
			trackIDPath: correctTrackIDPath,
			user:        &correctUser,
			mockBehavior: func(tu *trackMocks.MockUsecase) {
				tu.EXPECT().UnLike(gomock.Any(), correctTrackID, correctUser.ID).
					Return(false, errors.New(""))
			},
			expectedStatus:   http.StatusInternalServerError,
			expectedResponse: commonTests.ErrorResponse(commonHTTP.DeleteLikeServerError),
		},
	}

	for _, tc := range testTable {
		t.Run(tc.name, func(t *testing.T) {
			// Call mock
			tc.mockBehavior(tu)

			commonTests.DeliveryTestGet(t, r, "/api/tracks/"+tc.trackIDPath+"/unlike",
				tc.expectedStatus, tc.expectedResponse,
				commonTests.WrapRequestWithUserNotNilFunc(tc.user))
		})
	}
}
