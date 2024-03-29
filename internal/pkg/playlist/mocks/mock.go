// Code generated by MockGen. DO NOT EDIT.
// Source: playlist.go

// Package mock_playlist is a generated GoMock package.
package mock_playlist

import (
	context "context"
	io "io"
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

// AddTrack mocks base method.
func (m *MockUsecase) AddTrack(ctx context.Context, trackID, playlistID, userID uint32) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddTrack", ctx, trackID, playlistID, userID)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddTrack indicates an expected call of AddTrack.
func (mr *MockUsecaseMockRecorder) AddTrack(ctx, trackID, playlistID, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddTrack", reflect.TypeOf((*MockUsecase)(nil).AddTrack), ctx, trackID, playlistID, userID)
}

// Create mocks base method.
func (m *MockUsecase) Create(ctx context.Context, playlist models.Playlist, usersID []uint32, userID uint32) (uint32, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", ctx, playlist, usersID, userID)
	ret0, _ := ret[0].(uint32)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create.
func (mr *MockUsecaseMockRecorder) Create(ctx, playlist, usersID, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockUsecase)(nil).Create), ctx, playlist, usersID, userID)
}

// Delete mocks base method.
func (m *MockUsecase) Delete(ctx context.Context, playlistID, userID uint32) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", ctx, playlistID, userID)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockUsecaseMockRecorder) Delete(ctx, playlistID, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockUsecase)(nil).Delete), ctx, playlistID, userID)
}

// DeleteTrack mocks base method.
func (m *MockUsecase) DeleteTrack(ctx context.Context, trackID, playlistID, userID uint32) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteTrack", ctx, trackID, playlistID, userID)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteTrack indicates an expected call of DeleteTrack.
func (mr *MockUsecaseMockRecorder) DeleteTrack(ctx, trackID, playlistID, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteTrack", reflect.TypeOf((*MockUsecase)(nil).DeleteTrack), ctx, trackID, playlistID, userID)
}

// GetByID mocks base method.
func (m *MockUsecase) GetByID(ctx context.Context, playlistID uint32) (*models.Playlist, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByID", ctx, playlistID)
	ret0, _ := ret[0].(*models.Playlist)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByID indicates an expected call of GetByID.
func (mr *MockUsecaseMockRecorder) GetByID(ctx, playlistID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByID", reflect.TypeOf((*MockUsecase)(nil).GetByID), ctx, playlistID)
}

// GetByUser mocks base method.
func (m *MockUsecase) GetByUser(ctx context.Context, userID uint32) ([]models.Playlist, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByUser", ctx, userID)
	ret0, _ := ret[0].([]models.Playlist)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByUser indicates an expected call of GetByUser.
func (mr *MockUsecaseMockRecorder) GetByUser(ctx, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByUser", reflect.TypeOf((*MockUsecase)(nil).GetByUser), ctx, userID)
}

// GetFeed mocks base method.
func (m *MockUsecase) GetFeed(ctx context.Context) ([]models.Playlist, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetFeed", ctx)
	ret0, _ := ret[0].([]models.Playlist)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetFeed indicates an expected call of GetFeed.
func (mr *MockUsecaseMockRecorder) GetFeed(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetFeed", reflect.TypeOf((*MockUsecase)(nil).GetFeed), ctx)
}

// GetLikedByUser mocks base method.
func (m *MockUsecase) GetLikedByUser(ctx context.Context, userID uint32) ([]models.Playlist, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetLikedByUser", ctx, userID)
	ret0, _ := ret[0].([]models.Playlist)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetLikedByUser indicates an expected call of GetLikedByUser.
func (mr *MockUsecaseMockRecorder) GetLikedByUser(ctx, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetLikedByUser", reflect.TypeOf((*MockUsecase)(nil).GetLikedByUser), ctx, userID)
}

// IsLiked mocks base method.
func (m *MockUsecase) IsLiked(ctx context.Context, artistID, userID uint32) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsLiked", ctx, artistID, userID)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// IsLiked indicates an expected call of IsLiked.
func (mr *MockUsecaseMockRecorder) IsLiked(ctx, artistID, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsLiked", reflect.TypeOf((*MockUsecase)(nil).IsLiked), ctx, artistID, userID)
}

// SetLike mocks base method.
func (m *MockUsecase) SetLike(ctx context.Context, playlistID, userID uint32) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetLike", ctx, playlistID, userID)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SetLike indicates an expected call of SetLike.
func (mr *MockUsecaseMockRecorder) SetLike(ctx, playlistID, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetLike", reflect.TypeOf((*MockUsecase)(nil).SetLike), ctx, playlistID, userID)
}

// UnLike mocks base method.
func (m *MockUsecase) UnLike(ctx context.Context, playlistID, userID uint32) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UnLike", ctx, playlistID, userID)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UnLike indicates an expected call of UnLike.
func (mr *MockUsecaseMockRecorder) UnLike(ctx, playlistID, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UnLike", reflect.TypeOf((*MockUsecase)(nil).UnLike), ctx, playlistID, userID)
}

// UpdateInfoAndMembers mocks base method.
func (m *MockUsecase) UpdateInfoAndMembers(ctx context.Context, playlist models.Playlist, usersID []uint32, userID uint32) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateInfoAndMembers", ctx, playlist, usersID, userID)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateInfoAndMembers indicates an expected call of UpdateInfoAndMembers.
func (mr *MockUsecaseMockRecorder) UpdateInfoAndMembers(ctx, playlist, usersID, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateInfoAndMembers", reflect.TypeOf((*MockUsecase)(nil).UpdateInfoAndMembers), ctx, playlist, usersID, userID)
}

// UploadCover mocks base method.
func (m *MockUsecase) UploadCover(ctx context.Context, playlistID, userID uint32, file io.ReadSeeker, fileSize int64, fileExtension string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UploadCover", ctx, playlistID, userID, file, fileSize, fileExtension)
	ret0, _ := ret[0].(error)
	return ret0
}

// UploadCover indicates an expected call of UploadCover.
func (mr *MockUsecaseMockRecorder) UploadCover(ctx, playlistID, userID, file, fileSize, fileExtension interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UploadCover", reflect.TypeOf((*MockUsecase)(nil).UploadCover), ctx, playlistID, userID, file, fileSize, fileExtension)
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

// AddTrack mocks base method.
func (m *MockRepository) AddTrack(ctx context.Context, trackID, playlistID uint32) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddTrack", ctx, trackID, playlistID)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddTrack indicates an expected call of AddTrack.
func (mr *MockRepositoryMockRecorder) AddTrack(ctx, trackID, playlistID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddTrack", reflect.TypeOf((*MockRepository)(nil).AddTrack), ctx, trackID, playlistID)
}

// Check mocks base method.
func (m *MockRepository) Check(ctx context.Context, playlistID uint32) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Check", ctx, playlistID)
	ret0, _ := ret[0].(error)
	return ret0
}

// Check indicates an expected call of Check.
func (mr *MockRepositoryMockRecorder) Check(ctx, playlistID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Check", reflect.TypeOf((*MockRepository)(nil).Check), ctx, playlistID)
}

// DeleteByID mocks base method.
func (m *MockRepository) DeleteByID(ctx context.Context, playlistID uint32) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteByID", ctx, playlistID)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteByID indicates an expected call of DeleteByID.
func (mr *MockRepositoryMockRecorder) DeleteByID(ctx, playlistID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteByID", reflect.TypeOf((*MockRepository)(nil).DeleteByID), ctx, playlistID)
}

// DeleteLike mocks base method.
func (m *MockRepository) DeleteLike(ctx context.Context, playlistID, userID uint32) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteLike", ctx, playlistID, userID)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DeleteLike indicates an expected call of DeleteLike.
func (mr *MockRepositoryMockRecorder) DeleteLike(ctx, playlistID, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteLike", reflect.TypeOf((*MockRepository)(nil).DeleteLike), ctx, playlistID, userID)
}

// DeleteTrack mocks base method.
func (m *MockRepository) DeleteTrack(ctx context.Context, trackID, playlistID uint32) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteTrack", ctx, trackID, playlistID)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteTrack indicates an expected call of DeleteTrack.
func (mr *MockRepositoryMockRecorder) DeleteTrack(ctx, trackID, playlistID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteTrack", reflect.TypeOf((*MockRepository)(nil).DeleteTrack), ctx, trackID, playlistID)
}

// GetByID mocks base method.
func (m *MockRepository) GetByID(ctx context.Context, playlistID uint32) (*models.Playlist, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByID", ctx, playlistID)
	ret0, _ := ret[0].(*models.Playlist)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByID indicates an expected call of GetByID.
func (mr *MockRepositoryMockRecorder) GetByID(ctx, playlistID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByID", reflect.TypeOf((*MockRepository)(nil).GetByID), ctx, playlistID)
}

// GetByUser mocks base method.
func (m *MockRepository) GetByUser(ctx context.Context, userID uint32) ([]models.Playlist, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByUser", ctx, userID)
	ret0, _ := ret[0].([]models.Playlist)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByUser indicates an expected call of GetByUser.
func (mr *MockRepositoryMockRecorder) GetByUser(ctx, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByUser", reflect.TypeOf((*MockRepository)(nil).GetByUser), ctx, userID)
}

// GetFeed mocks base method.
func (m *MockRepository) GetFeed(ctx context.Context, limit uint32) ([]models.Playlist, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetFeed", ctx, limit)
	ret0, _ := ret[0].([]models.Playlist)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetFeed indicates an expected call of GetFeed.
func (mr *MockRepositoryMockRecorder) GetFeed(ctx, limit interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetFeed", reflect.TypeOf((*MockRepository)(nil).GetFeed), ctx, limit)
}

// GetLikedByUser mocks base method.
func (m *MockRepository) GetLikedByUser(ctx context.Context, userID uint32) ([]models.Playlist, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetLikedByUser", ctx, userID)
	ret0, _ := ret[0].([]models.Playlist)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetLikedByUser indicates an expected call of GetLikedByUser.
func (mr *MockRepositoryMockRecorder) GetLikedByUser(ctx, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetLikedByUser", reflect.TypeOf((*MockRepository)(nil).GetLikedByUser), ctx, userID)
}

// Insert mocks base method.
func (m *MockRepository) Insert(ctx context.Context, playlist models.Playlist, usersID []uint32) (uint32, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Insert", ctx, playlist, usersID)
	ret0, _ := ret[0].(uint32)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Insert indicates an expected call of Insert.
func (mr *MockRepositoryMockRecorder) Insert(ctx, playlist, usersID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Insert", reflect.TypeOf((*MockRepository)(nil).Insert), ctx, playlist, usersID)
}

// InsertLike mocks base method.
func (m *MockRepository) InsertLike(ctx context.Context, playlistID, userID uint32) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "InsertLike", ctx, playlistID, userID)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// InsertLike indicates an expected call of InsertLike.
func (mr *MockRepositoryMockRecorder) InsertLike(ctx, playlistID, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "InsertLike", reflect.TypeOf((*MockRepository)(nil).InsertLike), ctx, playlistID, userID)
}

// IsLiked mocks base method.
func (m *MockRepository) IsLiked(ctx context.Context, artistID, userID uint32) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsLiked", ctx, artistID, userID)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// IsLiked indicates an expected call of IsLiked.
func (mr *MockRepositoryMockRecorder) IsLiked(ctx, artistID, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsLiked", reflect.TypeOf((*MockRepository)(nil).IsLiked), ctx, artistID, userID)
}

// Update mocks base method.
func (m *MockRepository) Update(ctx context.Context, playlist models.Playlist) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", ctx, playlist)
	ret0, _ := ret[0].(error)
	return ret0
}

// Update indicates an expected call of Update.
func (mr *MockRepositoryMockRecorder) Update(ctx, playlist interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockRepository)(nil).Update), ctx, playlist)
}

// UpdateWithMembers mocks base method.
func (m *MockRepository) UpdateWithMembers(ctx context.Context, playlist models.Playlist, usersID []uint32) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateWithMembers", ctx, playlist, usersID)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateWithMembers indicates an expected call of UpdateWithMembers.
func (mr *MockRepositoryMockRecorder) UpdateWithMembers(ctx, playlist, usersID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateWithMembers", reflect.TypeOf((*MockRepository)(nil).UpdateWithMembers), ctx, playlist, usersID)
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

// LikedPlaylists mocks base method.
func (m *MockTables) LikedPlaylists() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "LikedPlaylists")
	ret0, _ := ret[0].(string)
	return ret0
}

// LikedPlaylists indicates an expected call of LikedPlaylists.
func (mr *MockTablesMockRecorder) LikedPlaylists() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LikedPlaylists", reflect.TypeOf((*MockTables)(nil).LikedPlaylists))
}

// Playlists mocks base method.
func (m *MockTables) Playlists() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Playlists")
	ret0, _ := ret[0].(string)
	return ret0
}

// Playlists indicates an expected call of Playlists.
func (mr *MockTablesMockRecorder) Playlists() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Playlists", reflect.TypeOf((*MockTables)(nil).Playlists))
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

// UsersPlaylists mocks base method.
func (m *MockTables) UsersPlaylists() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UsersPlaylists")
	ret0, _ := ret[0].(string)
	return ret0
}

// UsersPlaylists indicates an expected call of UsersPlaylists.
func (mr *MockTablesMockRecorder) UsersPlaylists() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UsersPlaylists", reflect.TypeOf((*MockTables)(nil).UsersPlaylists))
}
