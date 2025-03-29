// Code generated by MockGen. DO NOT EDIT.
// Source: pkg/usecase/interface/user.go

// Package mock is a generated GoMock package.
package mock

import (
	models "ShowTimes/pkg/utils/models"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockUserUseCase is a mock of UserUseCase interface.
type MockUserUseCase struct {
	ctrl     *gomock.Controller
	recorder *MockUserUseCaseMockRecorder
}

// MockUserUseCaseMockRecorder is the mock recorder for MockUserUseCase.
type MockUserUseCaseMockRecorder struct {
	mock *MockUserUseCase
}

// NewMockUserUseCase creates a new mock instance.
func NewMockUserUseCase(ctrl *gomock.Controller) *MockUserUseCase {
	mock := &MockUserUseCase{ctrl: ctrl}
	mock.recorder = &MockUserUseCaseMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUserUseCase) EXPECT() *MockUserUseCaseMockRecorder {
	return m.recorder
}

// AddAddress mocks base method.
func (m *MockUserUseCase) AddAddress(userID int, address models.AddressInfoResponse) ([]models.AddressInfoResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddAddress", userID, address)
	ret0, _ := ret[0].([]models.AddressInfoResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AddAddress indicates an expected call of AddAddress.
func (mr *MockUserUseCaseMockRecorder) AddAddress(userID, address interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddAddress", reflect.TypeOf((*MockUserUseCase)(nil).AddAddress), userID, address)
}

// ChangePassword mocks base method.
func (m *MockUserUseCase) ChangePassword(user models.ChangePassword) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ChangePassword", user)
	ret0, _ := ret[0].(error)
	return ret0
}

// ChangePassword indicates an expected call of ChangePassword.
func (mr *MockUserUseCaseMockRecorder) ChangePassword(user interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ChangePassword", reflect.TypeOf((*MockUserUseCase)(nil).ChangePassword), user)
}

// EditProfile mocks base method.
func (m *MockUserUseCase) EditProfile(user models.UsersProfileDetails) (models.UsersProfileDetails, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "EditProfile", user)
	ret0, _ := ret[0].(models.UsersProfileDetails)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// EditProfile indicates an expected call of EditProfile.
func (mr *MockUserUseCaseMockRecorder) EditProfile(user interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "EditProfile", reflect.TypeOf((*MockUserUseCase)(nil).EditProfile), user)
}

// GetAllAddress mocks base method.
func (m *MockUserUseCase) GetAllAddress(userID int) ([]models.AddressInfoResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAllAddress", userID)
	ret0, _ := ret[0].([]models.AddressInfoResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAllAddress indicates an expected call of GetAllAddress.
func (mr *MockUserUseCaseMockRecorder) GetAllAddress(userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAllAddress", reflect.TypeOf((*MockUserUseCase)(nil).GetAllAddress), userID)
}

// LoginHandler mocks base method.
func (m *MockUserUseCase) LoginHandler(user models.UserLogin) (models.TokenUsers, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "LoginHandler", user)
	ret0, _ := ret[0].(models.TokenUsers)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// LoginHandler indicates an expected call of LoginHandler.
func (mr *MockUserUseCaseMockRecorder) LoginHandler(user interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LoginHandler", reflect.TypeOf((*MockUserUseCase)(nil).LoginHandler), user)
}

// ShowUserDetails mocks base method.
func (m *MockUserUseCase) ShowUserDetails(userID int) (models.UsersProfileDetails, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ShowUserDetails", userID)
	ret0, _ := ret[0].(models.UsersProfileDetails)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ShowUserDetails indicates an expected call of ShowUserDetails.
func (mr *MockUserUseCaseMockRecorder) ShowUserDetails(userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ShowUserDetails", reflect.TypeOf((*MockUserUseCase)(nil).ShowUserDetails), userID)
}

// UserSignUp mocks base method.
func (m *MockUserUseCase) UserSignUp(user models.UserDetails) (models.TokenUsers, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UserSignUp", user)
	ret0, _ := ret[0].(models.TokenUsers)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UserSignUp indicates an expected call of UserSignUp.
func (mr *MockUserUseCaseMockRecorder) UserSignUp(user interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UserSignUp", reflect.TypeOf((*MockUserUseCase)(nil).UserSignUp), user)
}
