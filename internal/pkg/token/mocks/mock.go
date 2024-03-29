// Code generated by MockGen. DO NOT EDIT.
// Source: token.go

// Package mock_token is a generated GoMock package.
package mock_token

import (
	reflect "reflect"

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

// CheckAccessToken mocks base method.
func (m *MockUsecase) CheckAccessToken(acessToken string) (uint32, uint32, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CheckAccessToken", acessToken)
	ret0, _ := ret[0].(uint32)
	ret1, _ := ret[1].(uint32)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// CheckAccessToken indicates an expected call of CheckAccessToken.
func (mr *MockUsecaseMockRecorder) CheckAccessToken(acessToken interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CheckAccessToken", reflect.TypeOf((*MockUsecase)(nil).CheckAccessToken), acessToken)
}

// CheckCSRFToken mocks base method.
func (m *MockUsecase) CheckCSRFToken(csrfToken string) (uint32, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CheckCSRFToken", csrfToken)
	ret0, _ := ret[0].(uint32)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CheckCSRFToken indicates an expected call of CheckCSRFToken.
func (mr *MockUsecaseMockRecorder) CheckCSRFToken(csrfToken interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CheckCSRFToken", reflect.TypeOf((*MockUsecase)(nil).CheckCSRFToken), csrfToken)
}

// GenerateAccessToken mocks base method.
func (m *MockUsecase) GenerateAccessToken(userID, userVersion uint32) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GenerateAccessToken", userID, userVersion)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GenerateAccessToken indicates an expected call of GenerateAccessToken.
func (mr *MockUsecaseMockRecorder) GenerateAccessToken(userID, userVersion interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GenerateAccessToken", reflect.TypeOf((*MockUsecase)(nil).GenerateAccessToken), userID, userVersion)
}

// GenerateCSRFToken mocks base method.
func (m *MockUsecase) GenerateCSRFToken(userID uint32) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GenerateCSRFToken", userID)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GenerateCSRFToken indicates an expected call of GenerateCSRFToken.
func (mr *MockUsecaseMockRecorder) GenerateCSRFToken(userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GenerateCSRFToken", reflect.TypeOf((*MockUsecase)(nil).GenerateCSRFToken), userID)
}
