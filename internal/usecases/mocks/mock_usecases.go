// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/alexsibrin/runbot-auth/internal/usecases (interfaces: IPasswordHasher,IAccountRepo)
//
// Generated by this command:
//
//	mockgen -destination mocks/mock_usecases.go -package usecases_test github.com/alexsibrin/runbot-auth/internal/usecases IPasswordHasher,IAccountRepo
//

// Package usecases_test is a generated GoMock package.
package usecases_test

import (
	context "context"
	reflect "reflect"

	entities "github.com/alexsibrin/runbot-auth/internal/entities"
	gomock "go.uber.org/mock/gomock"
)

// MockIPasswordHasher is a mock of IPasswordHasher interface.
type MockIPasswordHasher struct {
	ctrl     *gomock.Controller
	recorder *MockIPasswordHasherMockRecorder
}

// MockIPasswordHasherMockRecorder is the mock recorder for MockIPasswordHasher.
type MockIPasswordHasherMockRecorder struct {
	mock *MockIPasswordHasher
}

// NewMockIPasswordHasher creates a new mock instance.
func NewMockIPasswordHasher(ctrl *gomock.Controller) *MockIPasswordHasher {
	mock := &MockIPasswordHasher{ctrl: ctrl}
	mock.recorder = &MockIPasswordHasherMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockIPasswordHasher) EXPECT() *MockIPasswordHasherMockRecorder {
	return m.recorder
}

// Compare mocks base method.
func (m *MockIPasswordHasher) Compare(arg0, arg1 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Compare", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Compare indicates an expected call of Compare.
func (mr *MockIPasswordHasherMockRecorder) Compare(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Compare", reflect.TypeOf((*MockIPasswordHasher)(nil).Compare), arg0, arg1)
}

// Hash mocks base method.
func (m *MockIPasswordHasher) Hash(arg0 string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Hash", arg0)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Hash indicates an expected call of Hash.
func (mr *MockIPasswordHasherMockRecorder) Hash(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Hash", reflect.TypeOf((*MockIPasswordHasher)(nil).Hash), arg0)
}

// MockIAccountRepo is a mock of IAccountRepo interface.
type MockIAccountRepo struct {
	ctrl     *gomock.Controller
	recorder *MockIAccountRepoMockRecorder
}

// MockIAccountRepoMockRecorder is the mock recorder for MockIAccountRepo.
type MockIAccountRepoMockRecorder struct {
	mock *MockIAccountRepo
}

// NewMockIAccountRepo creates a new mock instance.
func NewMockIAccountRepo(ctrl *gomock.Controller) *MockIAccountRepo {
	mock := &MockIAccountRepo{ctrl: ctrl}
	mock.recorder = &MockIAccountRepoMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockIAccountRepo) EXPECT() *MockIAccountRepoMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockIAccountRepo) Create(arg0 context.Context, arg1 *entities.Account) (*entities.Account, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", arg0, arg1)
	ret0, _ := ret[0].(*entities.Account)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create.
func (mr *MockIAccountRepoMockRecorder) Create(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockIAccountRepo)(nil).Create), arg0, arg1)
}

// GetOneByEmail mocks base method.
func (m *MockIAccountRepo) GetOneByEmail(arg0 context.Context, arg1 string) (*entities.Account, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetOneByEmail", arg0, arg1)
	ret0, _ := ret[0].(*entities.Account)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetOneByEmail indicates an expected call of GetOneByEmail.
func (mr *MockIAccountRepoMockRecorder) GetOneByEmail(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetOneByEmail", reflect.TypeOf((*MockIAccountRepo)(nil).GetOneByEmail), arg0, arg1)
}

// GetOneByUUID mocks base method.
func (m *MockIAccountRepo) GetOneByUUID(arg0 context.Context, arg1 string) (*entities.Account, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetOneByUUID", arg0, arg1)
	ret0, _ := ret[0].(*entities.Account)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetOneByUUID indicates an expected call of GetOneByUUID.
func (mr *MockIAccountRepoMockRecorder) GetOneByUUID(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetOneByUUID", reflect.TypeOf((*MockIAccountRepo)(nil).GetOneByUUID), arg0, arg1)
}

// IsExist mocks base method.
func (m *MockIAccountRepo) IsExist(arg0 context.Context, arg1 *entities.Account) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsExist", arg0, arg1)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// IsExist indicates an expected call of IsExist.
func (mr *MockIAccountRepoMockRecorder) IsExist(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsExist", reflect.TypeOf((*MockIAccountRepo)(nil).IsExist), arg0, arg1)
}

// IsExistByUUID mocks base method.
func (m *MockIAccountRepo) IsExistByUUID(arg0 context.Context, arg1 string) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsExistByUUID", arg0, arg1)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// IsExistByUUID indicates an expected call of IsExistByUUID.
func (mr *MockIAccountRepoMockRecorder) IsExistByUUID(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsExistByUUID", reflect.TypeOf((*MockIAccountRepo)(nil).IsExistByUUID), arg0, arg1)
}

// SetAccountStatus mocks base method.
func (m *MockIAccountRepo) SetAccountStatus(arg0 context.Context, arg1 string, arg2 byte) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetAccountStatus", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetAccountStatus indicates an expected call of SetAccountStatus.
func (mr *MockIAccountRepoMockRecorder) SetAccountStatus(arg0, arg1, arg2 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetAccountStatus", reflect.TypeOf((*MockIAccountRepo)(nil).SetAccountStatus), arg0, arg1, arg2)
}
