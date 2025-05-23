// Code generated by MockGen. DO NOT EDIT.
// Source: pkg/repository/interfaces/admin.go

// Package mockRepository is a generated GoMock package.
package mockRepository

import (
	domain "ShowTimes/pkg/domain"
	models "ShowTimes/pkg/utils/models"
	reflect "reflect"
	time "time"

	gomock "github.com/golang/mock/gomock"
)

// MockAdminRepository is a mock of AdminRepository interface.
type MockAdminRepository struct {
	ctrl     *gomock.Controller
	recorder *MockAdminRepositoryMockRecorder
}

// MockAdminRepositoryMockRecorder is the mock recorder for MockAdminRepository.
type MockAdminRepositoryMockRecorder struct {
	mock *MockAdminRepository
}

// NewMockAdminRepository creates a new mock instance.
func NewMockAdminRepository(ctrl *gomock.Controller) *MockAdminRepository {
	mock := &MockAdminRepository{ctrl: ctrl}
	mock.recorder = &MockAdminRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockAdminRepository) EXPECT() *MockAdminRepositoryMockRecorder {
	return m.recorder
}

// DashboardAmountDetails mocks base method.
func (m *MockAdminRepository) DashboardAmountDetails() (models.DashBoardAmount, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DashboardAmountDetails")
	ret0, _ := ret[0].(models.DashBoardAmount)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DashboardAmountDetails indicates an expected call of DashboardAmountDetails.
func (mr *MockAdminRepositoryMockRecorder) DashboardAmountDetails() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DashboardAmountDetails", reflect.TypeOf((*MockAdminRepository)(nil).DashboardAmountDetails))
}

// DashboardOrderDetails mocks base method.
func (m *MockAdminRepository) DashboardOrderDetails() (models.DashBoardOrder, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DashboardOrderDetails")
	ret0, _ := ret[0].(models.DashBoardOrder)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DashboardOrderDetails indicates an expected call of DashboardOrderDetails.
func (mr *MockAdminRepositoryMockRecorder) DashboardOrderDetails() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DashboardOrderDetails", reflect.TypeOf((*MockAdminRepository)(nil).DashboardOrderDetails))
}

// DashboardProductDetails mocks base method.
func (m *MockAdminRepository) DashboardProductDetails() (models.DashBoardProduct, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DashboardProductDetails")
	ret0, _ := ret[0].(models.DashBoardProduct)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DashboardProductDetails indicates an expected call of DashboardProductDetails.
func (mr *MockAdminRepositoryMockRecorder) DashboardProductDetails() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DashboardProductDetails", reflect.TypeOf((*MockAdminRepository)(nil).DashboardProductDetails))
}

// DashboardTotalRevenueDetails mocks base method.
func (m *MockAdminRepository) DashboardTotalRevenueDetails() (models.DashBoardRevenue, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DashboardTotalRevenueDetails")
	ret0, _ := ret[0].(models.DashBoardRevenue)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DashboardTotalRevenueDetails indicates an expected call of DashboardTotalRevenueDetails.
func (mr *MockAdminRepositoryMockRecorder) DashboardTotalRevenueDetails() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DashboardTotalRevenueDetails", reflect.TypeOf((*MockAdminRepository)(nil).DashboardTotalRevenueDetails))
}

// DashboardUserDetails mocks base method.
func (m *MockAdminRepository) DashboardUserDetails() (models.DashBoardUser, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DashboardUserDetails")
	ret0, _ := ret[0].(models.DashBoardUser)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DashboardUserDetails indicates an expected call of DashboardUserDetails.
func (mr *MockAdminRepositoryMockRecorder) DashboardUserDetails() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DashboardUserDetails", reflect.TypeOf((*MockAdminRepository)(nil).DashboardUserDetails))
}

// FilteredSalesReport mocks base method.
func (m *MockAdminRepository) FilteredSalesReport(startTime, endTime time.Time) (models.SalesReport, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FilteredSalesReport", startTime, endTime)
	ret0, _ := ret[0].(models.SalesReport)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FilteredSalesReport indicates an expected call of FilteredSalesReport.
func (mr *MockAdminRepositoryMockRecorder) FilteredSalesReport(startTime, endTime interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FilteredSalesReport", reflect.TypeOf((*MockAdminRepository)(nil).FilteredSalesReport), startTime, endTime)
}

// GetUserByID mocks base method.
func (m *MockAdminRepository) GetUserByID(id int) (domain.Users, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserByID", id)
	ret0, _ := ret[0].(domain.Users)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserByID indicates an expected call of GetUserByID.
func (mr *MockAdminRepositoryMockRecorder) GetUserByID(id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserByID", reflect.TypeOf((*MockAdminRepository)(nil).GetUserByID), id)
}

// GetUsers mocks base method.
func (m *MockAdminRepository) GetUsers(page int) ([]models.UserDetailsAtAdmin, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUsers", page)
	ret0, _ := ret[0].([]models.UserDetailsAtAdmin)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUsers indicates an expected call of GetUsers.
func (mr *MockAdminRepositoryMockRecorder) GetUsers(page interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUsers", reflect.TypeOf((*MockAdminRepository)(nil).GetUsers), page)
}

// IsAdmin mocks base method.
func (m *MockAdminRepository) IsAdmin(mailId string) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsAdmin", mailId)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// IsAdmin indicates an expected call of IsAdmin.
func (mr *MockAdminRepositoryMockRecorder) IsAdmin(mailId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsAdmin", reflect.TypeOf((*MockAdminRepository)(nil).IsAdmin), mailId)
}

// IsUserExist mocks base method.
func (m *MockAdminRepository) IsUserExist(id int) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsUserExist", id)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// IsUserExist indicates an expected call of IsUserExist.
func (mr *MockAdminRepositoryMockRecorder) IsUserExist(id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsUserExist", reflect.TypeOf((*MockAdminRepository)(nil).IsUserExist), id)
}

// LoginHandler mocks base method.
func (m *MockAdminRepository) LoginHandler(adminDetails models.AdminLogin) (domain.Users, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "LoginHandler", adminDetails)
	ret0, _ := ret[0].(domain.Users)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// LoginHandler indicates an expected call of LoginHandler.
func (mr *MockAdminRepositoryMockRecorder) LoginHandler(adminDetails interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LoginHandler", reflect.TypeOf((*MockAdminRepository)(nil).LoginHandler), adminDetails)
}

// SalesByDay mocks base method.
func (m *MockAdminRepository) SalesByDay(yearInt, monthInt, dayInt int) ([]models.OrderDetailsAdmin, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SalesByDay", yearInt, monthInt, dayInt)
	ret0, _ := ret[0].([]models.OrderDetailsAdmin)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SalesByDay indicates an expected call of SalesByDay.
func (mr *MockAdminRepositoryMockRecorder) SalesByDay(yearInt, monthInt, dayInt interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SalesByDay", reflect.TypeOf((*MockAdminRepository)(nil).SalesByDay), yearInt, monthInt, dayInt)
}

// SalesByMonth mocks base method.
func (m *MockAdminRepository) SalesByMonth(yearInt, monthInt int) ([]models.OrderDetailsAdmin, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SalesByMonth", yearInt, monthInt)
	ret0, _ := ret[0].([]models.OrderDetailsAdmin)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SalesByMonth indicates an expected call of SalesByMonth.
func (mr *MockAdminRepositoryMockRecorder) SalesByMonth(yearInt, monthInt interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SalesByMonth", reflect.TypeOf((*MockAdminRepository)(nil).SalesByMonth), yearInt, monthInt)
}

// SalesByYear mocks base method.
func (m *MockAdminRepository) SalesByYear(yearInt int) ([]models.OrderDetailsAdmin, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SalesByYear", yearInt)
	ret0, _ := ret[0].([]models.OrderDetailsAdmin)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SalesByYear indicates an expected call of SalesByYear.
func (mr *MockAdminRepositoryMockRecorder) SalesByYear(yearInt interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SalesByYear", reflect.TypeOf((*MockAdminRepository)(nil).SalesByYear), yearInt)
}

// UpdateBlockUserByID mocks base method.
func (m *MockAdminRepository) UpdateBlockUserByID(user models.UpdateBlock) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateBlockUserByID", user)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateBlockUserByID indicates an expected call of UpdateBlockUserByID.
func (mr *MockAdminRepositoryMockRecorder) UpdateBlockUserByID(user interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateBlockUserByID", reflect.TypeOf((*MockAdminRepository)(nil).UpdateBlockUserByID), user)
}
