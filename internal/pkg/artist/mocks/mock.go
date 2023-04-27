// Code generated by MockGen. DO NOT EDIT.
// Source: artist.go

// Package mock_artist is a generated GoMock package.
package mock_artist

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
func (m *MockUsecase) Create(ctx context.Context, artist models.Artist) (uint32, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", ctx, artist)
	ret0, _ := ret[0].(uint32)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create.
func (mr *MockUsecaseMockRecorder) Create(ctx, artist interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockUsecase)(nil).Create), ctx, artist)
}

// Delete mocks base method.
func (m *MockUsecase) Delete(ctx context.Context, artistID, userID uint32) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", ctx, artistID, userID)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockUsecaseMockRecorder) Delete(ctx, artistID, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockUsecase)(nil).Delete), ctx, artistID, userID)
}

// GetByAlbum mocks base method.
func (m *MockUsecase) GetByAlbum(ctx context.Context, albumID uint32) ([]models.Artist, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByAlbum", ctx, albumID)
	ret0, _ := ret[0].([]models.Artist)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByAlbum indicates an expected call of GetByAlbum.
func (mr *MockUsecaseMockRecorder) GetByAlbum(ctx, albumID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByAlbum", reflect.TypeOf((*MockUsecase)(nil).GetByAlbum), ctx, albumID)
}

// GetByID mocks base method.
func (m *MockUsecase) GetByID(ctx context.Context, artistID uint32) (*models.Artist, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByID", ctx, artistID)
	ret0, _ := ret[0].(*models.Artist)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByID indicates an expected call of GetByID.
func (mr *MockUsecaseMockRecorder) GetByID(ctx, artistID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByID", reflect.TypeOf((*MockUsecase)(nil).GetByID), ctx, artistID)
}

// GetByTrack mocks base method.
func (m *MockUsecase) GetByTrack(ctx context.Context, trackID uint32) ([]models.Artist, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByTrack", ctx, trackID)
	ret0, _ := ret[0].([]models.Artist)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByTrack indicates an expected call of GetByTrack.
func (mr *MockUsecaseMockRecorder) GetByTrack(ctx, trackID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByTrack", reflect.TypeOf((*MockUsecase)(nil).GetByTrack), ctx, trackID)
}

// GetFeed mocks base method.
func (m *MockUsecase) GetFeed(ctx context.Context) ([]models.Artist, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetFeed", ctx)
	ret0, _ := ret[0].([]models.Artist)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetFeed indicates an expected call of GetFeed.
func (mr *MockUsecaseMockRecorder) GetFeed(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetFeed", reflect.TypeOf((*MockUsecase)(nil).GetFeed), ctx)
}

// GetLikedByUser mocks base method.
func (m *MockUsecase) GetLikedByUser(ctx context.Context, userID uint32) ([]models.Artist, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetLikedByUser", ctx, userID)
	ret0, _ := ret[0].([]models.Artist)
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
func (m *MockUsecase) SetLike(ctx context.Context, artistID, userID uint32) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetLike", ctx, artistID, userID)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SetLike indicates an expected call of SetLike.
func (mr *MockUsecaseMockRecorder) SetLike(ctx, artistID, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetLike", reflect.TypeOf((*MockUsecase)(nil).SetLike), ctx, artistID, userID)
}

// UnLike mocks base method.
func (m *MockUsecase) UnLike(ctx context.Context, artistID, userID uint32) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UnLike", ctx, artistID, userID)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UnLike indicates an expected call of UnLike.
func (mr *MockUsecaseMockRecorder) UnLike(ctx, artistID, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UnLike", reflect.TypeOf((*MockUsecase)(nil).UnLike), ctx, artistID, userID)
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
func (m *MockRepository) Check(ctx context.Context, artistID uint32) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Check", ctx, artistID)
	ret0, _ := ret[0].(error)
	return ret0
}

// Check indicates an expected call of Check.
func (mr *MockRepositoryMockRecorder) Check(ctx, artistID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Check", reflect.TypeOf((*MockRepository)(nil).Check), ctx, artistID)
}

// DeleteByID mocks base method.
func (m *MockRepository) DeleteByID(ctx context.Context, artistID uint32) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteByID", ctx, artistID)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteByID indicates an expected call of DeleteByID.
func (mr *MockRepositoryMockRecorder) DeleteByID(ctx, artistID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteByID", reflect.TypeOf((*MockRepository)(nil).DeleteByID), ctx, artistID)
}

// DeleteLike mocks base method.
func (m *MockRepository) DeleteLike(ctx context.Context, artistID, userID uint32) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteLike", ctx, artistID, userID)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DeleteLike indicates an expected call of DeleteLike.
func (mr *MockRepositoryMockRecorder) DeleteLike(ctx, artistID, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteLike", reflect.TypeOf((*MockRepository)(nil).DeleteLike), ctx, artistID, userID)
}

// GetByAlbum mocks base method.
func (m *MockRepository) GetByAlbum(ctx context.Context, albumID uint32) ([]models.Artist, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByAlbum", ctx, albumID)
	ret0, _ := ret[0].([]models.Artist)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByAlbum indicates an expected call of GetByAlbum.
func (mr *MockRepositoryMockRecorder) GetByAlbum(ctx, albumID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByAlbum", reflect.TypeOf((*MockRepository)(nil).GetByAlbum), ctx, albumID)
}

// GetByID mocks base method.
func (m *MockRepository) GetByID(ctx context.Context, artistID uint32) (*models.Artist, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByID", ctx, artistID)
	ret0, _ := ret[0].(*models.Artist)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByID indicates an expected call of GetByID.
func (mr *MockRepositoryMockRecorder) GetByID(ctx, artistID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByID", reflect.TypeOf((*MockRepository)(nil).GetByID), ctx, artistID)
}

// GetByTrack mocks base method.
func (m *MockRepository) GetByTrack(ctx context.Context, trackID uint32) ([]models.Artist, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByTrack", ctx, trackID)
	ret0, _ := ret[0].([]models.Artist)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByTrack indicates an expected call of GetByTrack.
func (mr *MockRepositoryMockRecorder) GetByTrack(ctx, trackID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByTrack", reflect.TypeOf((*MockRepository)(nil).GetByTrack), ctx, trackID)
}

// GetFeed mocks base method.
func (m *MockRepository) GetFeed(ctx context.Context, amountLimit int) ([]models.Artist, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetFeed", ctx, amountLimit)
	ret0, _ := ret[0].([]models.Artist)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetFeed indicates an expected call of GetFeed.
func (mr *MockRepositoryMockRecorder) GetFeed(ctx, amountLimit interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetFeed", reflect.TypeOf((*MockRepository)(nil).GetFeed), ctx, amountLimit)
}

// GetLikedByUser mocks base method.
func (m *MockRepository) GetLikedByUser(ctx context.Context, userID uint32) ([]models.Artist, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetLikedByUser", ctx, userID)
	ret0, _ := ret[0].([]models.Artist)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetLikedByUser indicates an expected call of GetLikedByUser.
func (mr *MockRepositoryMockRecorder) GetLikedByUser(ctx, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetLikedByUser", reflect.TypeOf((*MockRepository)(nil).GetLikedByUser), ctx, userID)
}

// Insert mocks base method.
func (m *MockRepository) Insert(ctx context.Context, artist models.Artist) (uint32, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Insert", ctx, artist)
	ret0, _ := ret[0].(uint32)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Insert indicates an expected call of Insert.
func (mr *MockRepositoryMockRecorder) Insert(ctx, artist interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Insert", reflect.TypeOf((*MockRepository)(nil).Insert), ctx, artist)
}

// InsertLike mocks base method.
func (m *MockRepository) InsertLike(ctx context.Context, artistID, userID uint32) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "InsertLike", ctx, artistID, userID)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// InsertLike indicates an expected call of InsertLike.
func (mr *MockRepositoryMockRecorder) InsertLike(ctx, artistID, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "InsertLike", reflect.TypeOf((*MockRepository)(nil).InsertLike), ctx, artistID, userID)
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

// Artists mocks base method.
func (m *MockTables) Artists() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Artists")
	ret0, _ := ret[0].(string)
	return ret0
}

// Artists indicates an expected call of Artists.
func (mr *MockTablesMockRecorder) Artists() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Artists", reflect.TypeOf((*MockTables)(nil).Artists))
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

// LikedArtists mocks base method.
func (m *MockTables) LikedArtists() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "LikedArtists")
	ret0, _ := ret[0].(string)
	return ret0
}

// LikedArtists indicates an expected call of LikedArtists.
func (mr *MockTablesMockRecorder) LikedArtists() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LikedArtists", reflect.TypeOf((*MockTables)(nil).LikedArtists))
}
