// Code generated by MockGen. DO NOT EDIT.
// Source: user.go

// Package mock_user is a generated GoMock package.
package mock_user

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

// GetByID mocks base method.
func (m *MockUsecase) GetByID(ctx context.Context, userID uint32) (*models.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByID", ctx, userID)
	ret0, _ := ret[0].(*models.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByID indicates an expected call of GetByID.
func (mr *MockUsecaseMockRecorder) GetByID(ctx, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByID", reflect.TypeOf((*MockUsecase)(nil).GetByID), ctx, userID)
}

// GetByPlaylist mocks base method.
func (m *MockUsecase) GetByPlaylist(ctx context.Context, playlistID uint32) ([]models.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByPlaylist", ctx, playlistID)
	ret0, _ := ret[0].([]models.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByPlaylist indicates an expected call of GetByPlaylist.
func (mr *MockUsecaseMockRecorder) GetByPlaylist(ctx, playlistID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByPlaylist", reflect.TypeOf((*MockUsecase)(nil).GetByPlaylist), ctx, playlistID)
}

// UpdateInfo mocks base method.
func (m *MockUsecase) UpdateInfo(ctx context.Context, user *models.User) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateInfo", ctx, user)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateInfo indicates an expected call of UpdateInfo.
func (mr *MockUsecaseMockRecorder) UpdateInfo(ctx, user interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateInfo", reflect.TypeOf((*MockUsecase)(nil).UpdateInfo), ctx, user)
}

// UploadAvatar mocks base method.
func (m *MockUsecase) UploadAvatar(ctx context.Context, userID uint32, file io.ReadSeeker, fileExtension string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UploadAvatar", ctx, userID, file, fileExtension)
	ret0, _ := ret[0].(error)
	return ret0
}

// UploadAvatar indicates an expected call of UploadAvatar.
func (mr *MockUsecaseMockRecorder) UploadAvatar(ctx, userID, file, fileExtension interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UploadAvatar", reflect.TypeOf((*MockUsecase)(nil).UploadAvatar), ctx, userID, file, fileExtension)
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
func (m *MockRepository) Check(ctx context.Context, userID uint32) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Check", ctx, userID)
	ret0, _ := ret[0].(error)
	return ret0
}

// Check indicates an expected call of Check.
func (mr *MockRepositoryMockRecorder) Check(ctx, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Check", reflect.TypeOf((*MockRepository)(nil).Check), ctx, userID)
}

// CreateUser mocks base method.
func (m *MockRepository) CreateUser(ctx context.Context, user models.User) (uint32, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateUser", ctx, user)
	ret0, _ := ret[0].(uint32)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateUser indicates an expected call of CreateUser.
func (mr *MockRepositoryMockRecorder) CreateUser(ctx, user interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateUser", reflect.TypeOf((*MockRepository)(nil).CreateUser), ctx, user)
}

// GetByID mocks base method.
func (m *MockRepository) GetByID(ctx context.Context, userID uint32) (*models.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByID", ctx, userID)
	ret0, _ := ret[0].(*models.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByID indicates an expected call of GetByID.
func (mr *MockRepositoryMockRecorder) GetByID(ctx, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByID", reflect.TypeOf((*MockRepository)(nil).GetByID), ctx, userID)
}

// GetByPlaylist mocks base method.
func (m *MockRepository) GetByPlaylist(ctx context.Context, playlistID uint32) ([]models.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByPlaylist", ctx, playlistID)
	ret0, _ := ret[0].([]models.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByPlaylist indicates an expected call of GetByPlaylist.
func (mr *MockRepositoryMockRecorder) GetByPlaylist(ctx, playlistID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByPlaylist", reflect.TypeOf((*MockRepository)(nil).GetByPlaylist), ctx, playlistID)
}

// GetUserByUsername mocks base method.
func (m *MockRepository) GetUserByUsername(ctx context.Context, username string) (*models.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserByUsername", ctx, username)
	ret0, _ := ret[0].(*models.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserByUsername indicates an expected call of GetUserByUsername.
func (mr *MockRepositoryMockRecorder) GetUserByUsername(ctx, username interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserByUsername", reflect.TypeOf((*MockRepository)(nil).GetUserByUsername), ctx, username)
}

// UpdateAvatarSrc mocks base method.
func (m *MockRepository) UpdateAvatarSrc(ctx context.Context, userID uint32, avatarSrc string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateAvatarSrc", ctx, userID, avatarSrc)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateAvatarSrc indicates an expected call of UpdateAvatarSrc.
func (mr *MockRepositoryMockRecorder) UpdateAvatarSrc(ctx, userID, avatarSrc interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateAvatarSrc", reflect.TypeOf((*MockRepository)(nil).UpdateAvatarSrc), ctx, userID, avatarSrc)
}

// UpdateInfo mocks base method.
func (m *MockRepository) UpdateInfo(ctx context.Context, user *models.User) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateInfo", ctx, user)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateInfo indicates an expected call of UpdateInfo.
func (mr *MockRepositoryMockRecorder) UpdateInfo(ctx, user interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateInfo", reflect.TypeOf((*MockRepository)(nil).UpdateInfo), ctx, user)
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

// Users mocks base method.
func (m *MockTables) Users() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Users")
	ret0, _ := ret[0].(string)
	return ret0
}

// Users indicates an expected call of Users.
func (mr *MockTablesMockRecorder) Users() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Users", reflect.TypeOf((*MockTables)(nil).Users))
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
