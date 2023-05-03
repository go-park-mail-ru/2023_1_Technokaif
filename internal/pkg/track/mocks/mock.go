// Code generated by MockGen. DO NOT EDIT.
// Source: track.go

// Package mock_track is a generated GoMock package.
package mock_track

import (
	context "context"
	reflect "reflect"

	models "github.com/go-park-mail-ru/2023_1_Technokaif/internal/models"
	gomock "github.com/golang/mock/gomock"
)

// MockUsecase is a mock of Usecase interface.
type MockUsecase struct {
	ctrl     *gomock.Controller
	recorder *MockUsecaseMockRecorder
}

// MockUsecaseMockRecorder is the mock recorder for MockUsecase.
type MockUsecaseMockRecorder struct {
	mock *MockUsecase
}

// NewMockUsecase creates a new mock instance.
func NewMockUsecase(ctrl *gomock.Controller) *MockUsecase {
	mock := &MockUsecase{ctrl: ctrl}
	mock.recorder = &MockUsecaseMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUsecase) EXPECT() *MockUsecaseMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockUsecase) Create(ctx context.Context, track models.Track, artistsID []uint32, userID uint32) (uint32, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", ctx, track, artistsID, userID)
	ret0, _ := ret[0].(uint32)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create.
func (mr *MockUsecaseMockRecorder) Create(ctx, track, artistsID, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockUsecase)(nil).Create), ctx, track, artistsID, userID)
}

// Delete mocks base method.
func (m *MockUsecase) Delete(ctx context.Context, trackID, userID uint32) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", ctx, trackID, userID)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockUsecaseMockRecorder) Delete(ctx, trackID, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockUsecase)(nil).Delete), ctx, trackID, userID)
}

// GetByAlbum mocks base method.
func (m *MockUsecase) GetByAlbum(ctx context.Context, albumID uint32) ([]models.Track, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByAlbum", ctx, albumID)
	ret0, _ := ret[0].([]models.Track)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByAlbum indicates an expected call of GetByAlbum.
func (mr *MockUsecaseMockRecorder) GetByAlbum(ctx, albumID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByAlbum", reflect.TypeOf((*MockUsecase)(nil).GetByAlbum), ctx, albumID)
}

// GetByArtist mocks base method.
func (m *MockUsecase) GetByArtist(ctx context.Context, artistID uint32) ([]models.Track, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByArtist", ctx, artistID)
	ret0, _ := ret[0].([]models.Track)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByArtist indicates an expected call of GetByArtist.
func (mr *MockUsecaseMockRecorder) GetByArtist(ctx, artistID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByArtist", reflect.TypeOf((*MockUsecase)(nil).GetByArtist), ctx, artistID)
}

// GetByID mocks base method.
func (m *MockUsecase) GetByID(ctx context.Context, trackID uint32) (*models.Track, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByID", ctx, trackID)
	ret0, _ := ret[0].(*models.Track)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByID indicates an expected call of GetByID.
func (mr *MockUsecaseMockRecorder) GetByID(ctx, trackID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByID", reflect.TypeOf((*MockUsecase)(nil).GetByID), ctx, trackID)
}

// GetByPlaylist mocks base method.
func (m *MockUsecase) GetByPlaylist(ctx context.Context, playlistID uint32) ([]models.Track, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByPlaylist", ctx, playlistID)
	ret0, _ := ret[0].([]models.Track)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByPlaylist indicates an expected call of GetByPlaylist.
func (mr *MockUsecaseMockRecorder) GetByPlaylist(ctx, playlistID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByPlaylist", reflect.TypeOf((*MockUsecase)(nil).GetByPlaylist), ctx, playlistID)
}

// GetFeed mocks base method.
func (m *MockUsecase) GetFeed(ctx context.Context) ([]models.Track, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetFeed", ctx)
	ret0, _ := ret[0].([]models.Track)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetFeed indicates an expected call of GetFeed.
func (mr *MockUsecaseMockRecorder) GetFeed(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetFeed", reflect.TypeOf((*MockUsecase)(nil).GetFeed), ctx)
}

// GetLikedByUser mocks base method.
func (m *MockUsecase) GetLikedByUser(ctx context.Context, userID uint32) ([]models.Track, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetLikedByUser", ctx, userID)
	ret0, _ := ret[0].([]models.Track)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetLikedByUser indicates an expected call of GetLikedByUser.
func (mr *MockUsecaseMockRecorder) GetLikedByUser(ctx, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetLikedByUser", reflect.TypeOf((*MockUsecase)(nil).GetLikedByUser), ctx, userID)
}

// IsLiked mocks base method.
func (m *MockUsecase) IsLiked(ctx context.Context, trackID, userID uint32) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsLiked", ctx, trackID, userID)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// IsLiked indicates an expected call of IsLiked.
func (mr *MockUsecaseMockRecorder) IsLiked(ctx, trackID, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsLiked", reflect.TypeOf((*MockUsecase)(nil).IsLiked), ctx, trackID, userID)
}

// SetLike mocks base method.
func (m *MockUsecase) SetLike(ctx context.Context, trackID, userID uint32) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetLike", ctx, trackID, userID)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SetLike indicates an expected call of SetLike.
func (mr *MockUsecaseMockRecorder) SetLike(ctx, trackID, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetLike", reflect.TypeOf((*MockUsecase)(nil).SetLike), ctx, trackID, userID)
}

// UnLike mocks base method.
func (m *MockUsecase) UnLike(ctx context.Context, trackID, userID uint32) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UnLike", ctx, trackID, userID)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UnLike indicates an expected call of UnLike.
func (mr *MockUsecaseMockRecorder) UnLike(ctx, trackID, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UnLike", reflect.TypeOf((*MockUsecase)(nil).UnLike), ctx, trackID, userID)
}

// MockRepository is a mock of Repository interface.
type MockRepository struct {
	ctrl     *gomock.Controller
	recorder *MockRepositoryMockRecorder
}

// MockRepositoryMockRecorder is the mock recorder for MockRepository.
type MockRepositoryMockRecorder struct {
	mock *MockRepository
}

// NewMockRepository creates a new mock instance.
func NewMockRepository(ctrl *gomock.Controller) *MockRepository {
	mock := &MockRepository{ctrl: ctrl}
	mock.recorder = &MockRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockRepository) EXPECT() *MockRepositoryMockRecorder {
	return m.recorder
}

// Check mocks base method.
func (m *MockRepository) Check(ctx context.Context, trackID uint32) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Check", ctx, trackID)
	ret0, _ := ret[0].(error)
	return ret0
}

// Check indicates an expected call of Check.
func (mr *MockRepositoryMockRecorder) Check(ctx, trackID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Check", reflect.TypeOf((*MockRepository)(nil).Check), ctx, trackID)
}

// DeleteByID mocks base method.
func (m *MockRepository) DeleteByID(ctx context.Context, trackID uint32) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteByID", ctx, trackID)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteByID indicates an expected call of DeleteByID.
func (mr *MockRepositoryMockRecorder) DeleteByID(ctx, trackID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteByID", reflect.TypeOf((*MockRepository)(nil).DeleteByID), ctx, trackID)
}

// DeleteLike mocks base method.
func (m *MockRepository) DeleteLike(ctx context.Context, trackID, userID uint32) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteLike", ctx, trackID, userID)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DeleteLike indicates an expected call of DeleteLike.
func (mr *MockRepositoryMockRecorder) DeleteLike(ctx, trackID, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteLike", reflect.TypeOf((*MockRepository)(nil).DeleteLike), ctx, trackID, userID)
}

// GetByAlbum mocks base method.
func (m *MockRepository) GetByAlbum(ctx context.Context, albumID uint32) ([]models.Track, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByAlbum", ctx, albumID)
	ret0, _ := ret[0].([]models.Track)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByAlbum indicates an expected call of GetByAlbum.
func (mr *MockRepositoryMockRecorder) GetByAlbum(ctx, albumID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByAlbum", reflect.TypeOf((*MockRepository)(nil).GetByAlbum), ctx, albumID)
}

// GetByArtist mocks base method.
func (m *MockRepository) GetByArtist(ctx context.Context, artistID uint32) ([]models.Track, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByArtist", ctx, artistID)
	ret0, _ := ret[0].([]models.Track)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByArtist indicates an expected call of GetByArtist.
func (mr *MockRepositoryMockRecorder) GetByArtist(ctx, artistID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByArtist", reflect.TypeOf((*MockRepository)(nil).GetByArtist), ctx, artistID)
}

// GetByID mocks base method.
func (m *MockRepository) GetByID(ctx context.Context, trackID uint32) (*models.Track, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByID", ctx, trackID)
	ret0, _ := ret[0].(*models.Track)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByID indicates an expected call of GetByID.
func (mr *MockRepositoryMockRecorder) GetByID(ctx, trackID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByID", reflect.TypeOf((*MockRepository)(nil).GetByID), ctx, trackID)
}

// GetByPlaylist mocks base method.
func (m *MockRepository) GetByPlaylist(ctx context.Context, playlistID uint32) ([]models.Track, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByPlaylist", ctx, playlistID)
	ret0, _ := ret[0].([]models.Track)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByPlaylist indicates an expected call of GetByPlaylist.
func (mr *MockRepositoryMockRecorder) GetByPlaylist(ctx, playlistID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByPlaylist", reflect.TypeOf((*MockRepository)(nil).GetByPlaylist), ctx, playlistID)
}

// GetFeed mocks base method.
func (m *MockRepository) GetFeed(ctx context.Context, limit uint32) ([]models.Track, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetFeed", ctx, limit)
	ret0, _ := ret[0].([]models.Track)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetFeed indicates an expected call of GetFeed.
func (mr *MockRepositoryMockRecorder) GetFeed(ctx, limit interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetFeed", reflect.TypeOf((*MockRepository)(nil).GetFeed), ctx, limit)
}

// GetLikedByUser mocks base method.
func (m *MockRepository) GetLikedByUser(ctx context.Context, userID uint32) ([]models.Track, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetLikedByUser", ctx, userID)
	ret0, _ := ret[0].([]models.Track)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetLikedByUser indicates an expected call of GetLikedByUser.
func (mr *MockRepositoryMockRecorder) GetLikedByUser(ctx, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetLikedByUser", reflect.TypeOf((*MockRepository)(nil).GetLikedByUser), ctx, userID)
}

// Insert mocks base method.
func (m *MockRepository) Insert(ctx context.Context, track models.Track, artistsID []uint32) (uint32, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Insert", ctx, track, artistsID)
	ret0, _ := ret[0].(uint32)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Insert indicates an expected call of Insert.
func (mr *MockRepositoryMockRecorder) Insert(ctx, track, artistsID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Insert", reflect.TypeOf((*MockRepository)(nil).Insert), ctx, track, artistsID)
}

// InsertLike mocks base method.
func (m *MockRepository) InsertLike(ctx context.Context, trackID, userID uint32) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "InsertLike", ctx, trackID, userID)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// InsertLike indicates an expected call of InsertLike.
func (mr *MockRepositoryMockRecorder) InsertLike(ctx, trackID, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "InsertLike", reflect.TypeOf((*MockRepository)(nil).InsertLike), ctx, trackID, userID)
}

// IsLiked mocks base method.
func (m *MockRepository) IsLiked(ctx context.Context, trackID, userID uint32) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsLiked", ctx, trackID, userID)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// IsLiked indicates an expected call of IsLiked.
func (mr *MockRepositoryMockRecorder) IsLiked(ctx, trackID, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsLiked", reflect.TypeOf((*MockRepository)(nil).IsLiked), ctx, trackID, userID)
}

// MockTables is a mock of Tables interface.
type MockTables struct {
	ctrl     *gomock.Controller
	recorder *MockTablesMockRecorder
}

// MockTablesMockRecorder is the mock recorder for MockTables.
type MockTablesMockRecorder struct {
	mock *MockTables
}

// NewMockTables creates a new mock instance.
func NewMockTables(ctrl *gomock.Controller) *MockTables {
	mock := &MockTables{ctrl: ctrl}
	mock.recorder = &MockTablesMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockTables) EXPECT() *MockTablesMockRecorder {
	return m.recorder
}

// ArtistsTracks mocks base method.
func (m *MockTables) ArtistsTracks() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ArtistsTracks")
	ret0, _ := ret[0].(string)
	return ret0
}

// ArtistsTracks indicates an expected call of ArtistsTracks.
func (mr *MockTablesMockRecorder) ArtistsTracks() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ArtistsTracks", reflect.TypeOf((*MockTables)(nil).ArtistsTracks))
}

// LikedTracks mocks base method.
func (m *MockTables) LikedTracks() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "LikedTracks")
	ret0, _ := ret[0].(string)
	return ret0
}

// LikedTracks indicates an expected call of LikedTracks.
func (mr *MockTablesMockRecorder) LikedTracks() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LikedTracks", reflect.TypeOf((*MockTables)(nil).LikedTracks))
}

// PlaylistsTracks mocks base method.
func (m *MockTables) PlaylistsTracks() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PlaylistsTracks")
	ret0, _ := ret[0].(string)
	return ret0
}

// PlaylistsTracks indicates an expected call of PlaylistsTracks.
func (mr *MockTablesMockRecorder) PlaylistsTracks() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PlaylistsTracks", reflect.TypeOf((*MockTables)(nil).PlaylistsTracks))
}

// Tracks mocks base method.
func (m *MockTables) Tracks() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Tracks")
	ret0, _ := ret[0].(string)
	return ret0
}

// Tracks indicates an expected call of Tracks.
func (mr *MockTablesMockRecorder) Tracks() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Tracks", reflect.TypeOf((*MockTables)(nil).Tracks))
}
