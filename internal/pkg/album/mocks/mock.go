// Code generated by MockGen. DO NOT EDIT.
// Source: album.go

// Package mock_album is a generated GoMock package.
package mock_album

import (
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

// Change mocks base method.
func (m *MockUsecase) Change(album models.Album) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Change", album)
	ret0, _ := ret[0].(error)
	return ret0
}

// Change indicates an expected call of Change.
func (mr *MockUsecaseMockRecorder) Change(album interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Change", reflect.TypeOf((*MockUsecase)(nil).Change), album)
}

// Create mocks base method.
func (m *MockUsecase) Create(album models.Album, artistsID []uint32, userID uint32) (uint32, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", album, artistsID, userID)
	ret0, _ := ret[0].(uint32)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create.
func (mr *MockUsecaseMockRecorder) Create(album, artistsID, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockUsecase)(nil).Create), album, artistsID, userID)
}

// Delete mocks base method.
func (m *MockUsecase) Delete(albumID, userID uint32) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", albumID, userID)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockUsecaseMockRecorder) Delete(albumID, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockUsecase)(nil).Delete), albumID, userID)
}

// GetByArtist mocks base method.
func (m *MockUsecase) GetByArtist(artistID uint32) ([]models.Album, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByArtist", artistID)
	ret0, _ := ret[0].([]models.Album)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByArtist indicates an expected call of GetByArtist.
func (mr *MockUsecaseMockRecorder) GetByArtist(artistID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByArtist", reflect.TypeOf((*MockUsecase)(nil).GetByArtist), artistID)
}

// GetByID mocks base method.
func (m *MockUsecase) GetByID(albumID uint32) (*models.Album, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByID", albumID)
	ret0, _ := ret[0].(*models.Album)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByID indicates an expected call of GetByID.
func (mr *MockUsecaseMockRecorder) GetByID(albumID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByID", reflect.TypeOf((*MockUsecase)(nil).GetByID), albumID)
}

// GetByTrack mocks base method.
func (m *MockUsecase) GetByTrack(trackID uint32) (*models.Album, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByTrack", trackID)
	ret0, _ := ret[0].(*models.Album)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByTrack indicates an expected call of GetByTrack.
func (mr *MockUsecaseMockRecorder) GetByTrack(trackID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByTrack", reflect.TypeOf((*MockUsecase)(nil).GetByTrack), trackID)
}

// GetFeed mocks base method.
func (m *MockUsecase) GetFeed() ([]models.Album, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetFeed")
	ret0, _ := ret[0].([]models.Album)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetFeed indicates an expected call of GetFeed.
func (mr *MockUsecaseMockRecorder) GetFeed() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetFeed", reflect.TypeOf((*MockUsecase)(nil).GetFeed))
}

// GetLikedByUser mocks base method.
func (m *MockUsecase) GetLikedByUser(userID uint32) ([]models.Album, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetLikedByUser", userID)
	ret0, _ := ret[0].([]models.Album)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetLikedByUser indicates an expected call of GetLikedByUser.
func (mr *MockUsecaseMockRecorder) GetLikedByUser(userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetLikedByUser", reflect.TypeOf((*MockUsecase)(nil).GetLikedByUser), userID)
}

// SetLike mocks base method.
func (m *MockUsecase) SetLike(albumID, userID uint32) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetLike", albumID, userID)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SetLike indicates an expected call of SetLike.
func (mr *MockUsecaseMockRecorder) SetLike(albumID, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetLike", reflect.TypeOf((*MockUsecase)(nil).SetLike), albumID, userID)
}

// UnLike mocks base method.
func (m *MockUsecase) UnLike(albumID, userID uint32) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UnLike", albumID, userID)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UnLike indicates an expected call of UnLike.
func (mr *MockUsecaseMockRecorder) UnLike(albumID, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UnLike", reflect.TypeOf((*MockUsecase)(nil).UnLike), albumID, userID)
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

// DeleteByID mocks base method.
func (m *MockRepository) DeleteByID(albumID uint32) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteByID", albumID)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteByID indicates an expected call of DeleteByID.
func (mr *MockRepositoryMockRecorder) DeleteByID(albumID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteByID", reflect.TypeOf((*MockRepository)(nil).DeleteByID), albumID)
}

// DeleteLike mocks base method.
func (m *MockRepository) DeleteLike(albumID, userID uint32) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteLike", albumID, userID)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DeleteLike indicates an expected call of DeleteLike.
func (mr *MockRepositoryMockRecorder) DeleteLike(albumID, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteLike", reflect.TypeOf((*MockRepository)(nil).DeleteLike), albumID, userID)
}

// GetByArtist mocks base method.
func (m *MockRepository) GetByArtist(artistID uint32) ([]models.Album, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByArtist", artistID)
	ret0, _ := ret[0].([]models.Album)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByArtist indicates an expected call of GetByArtist.
func (mr *MockRepositoryMockRecorder) GetByArtist(artistID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByArtist", reflect.TypeOf((*MockRepository)(nil).GetByArtist), artistID)
}

// GetByID mocks base method.
func (m *MockRepository) GetByID(albumID uint32) (*models.Album, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByID", albumID)
	ret0, _ := ret[0].(*models.Album)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByID indicates an expected call of GetByID.
func (mr *MockRepositoryMockRecorder) GetByID(albumID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByID", reflect.TypeOf((*MockRepository)(nil).GetByID), albumID)
}

// GetByTrack mocks base method.
func (m *MockRepository) GetByTrack(trackID uint32) (*models.Album, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByTrack", trackID)
	ret0, _ := ret[0].(*models.Album)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByTrack indicates an expected call of GetByTrack.
func (mr *MockRepositoryMockRecorder) GetByTrack(trackID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByTrack", reflect.TypeOf((*MockRepository)(nil).GetByTrack), trackID)
}

// GetFeed mocks base method.
func (m *MockRepository) GetFeed() ([]models.Album, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetFeed")
	ret0, _ := ret[0].([]models.Album)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetFeed indicates an expected call of GetFeed.
func (mr *MockRepositoryMockRecorder) GetFeed() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetFeed", reflect.TypeOf((*MockRepository)(nil).GetFeed))
}

// GetLikedByUser mocks base method.
func (m *MockRepository) GetLikedByUser(userID uint32) ([]models.Album, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetLikedByUser", userID)
	ret0, _ := ret[0].([]models.Album)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetLikedByUser indicates an expected call of GetLikedByUser.
func (mr *MockRepositoryMockRecorder) GetLikedByUser(userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetLikedByUser", reflect.TypeOf((*MockRepository)(nil).GetLikedByUser), userID)
}

// Insert mocks base method.
func (m *MockRepository) Insert(album models.Album, artistsID []uint32) (uint32, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Insert", album, artistsID)
	ret0, _ := ret[0].(uint32)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Insert indicates an expected call of Insert.
func (mr *MockRepositoryMockRecorder) Insert(album, artistsID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Insert", reflect.TypeOf((*MockRepository)(nil).Insert), album, artistsID)
}

// InsertLike mocks base method.
func (m *MockRepository) InsertLike(albumID, userID uint32) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "InsertLike", albumID, userID)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// InsertLike indicates an expected call of InsertLike.
func (mr *MockRepositoryMockRecorder) InsertLike(albumID, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "InsertLike", reflect.TypeOf((*MockRepository)(nil).InsertLike), albumID, userID)
}

// Update mocks base method.
func (m *MockRepository) Update(album models.Album) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", album)
	ret0, _ := ret[0].(error)
	return ret0
}

// Update indicates an expected call of Update.
func (mr *MockRepositoryMockRecorder) Update(album interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockRepository)(nil).Update), album)
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

// Albums mocks base method.
func (m *MockTables) Albums() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Albums")
	ret0, _ := ret[0].(string)
	return ret0
}

// Albums indicates an expected call of Albums.
func (mr *MockTablesMockRecorder) Albums() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Albums", reflect.TypeOf((*MockTables)(nil).Albums))
}

// ArtistsAlbums mocks base method.
func (m *MockTables) ArtistsAlbums() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ArtistsAlbums")
	ret0, _ := ret[0].(string)
	return ret0
}

// ArtistsAlbums indicates an expected call of ArtistsAlbums.
func (mr *MockTablesMockRecorder) ArtistsAlbums() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ArtistsAlbums", reflect.TypeOf((*MockTables)(nil).ArtistsAlbums))
}

// LikedAlbums mocks base method.
func (m *MockTables) LikedAlbums() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "LikedAlbums")
	ret0, _ := ret[0].(string)
	return ret0
}

// LikedAlbums indicates an expected call of LikedAlbums.
func (mr *MockTablesMockRecorder) LikedAlbums() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LikedAlbums", reflect.TypeOf((*MockTables)(nil).LikedAlbums))
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
