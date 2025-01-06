// Code generated by MockGen. DO NOT EDIT.
// Source: authentication_repository.go
//
// Generated by this command:
//
//	mockgen -source=authentication_repository.go -destination=./mock_repository/mock_authentication_repository.go -package=mock_repository
//

// Package mock_repository is a generated GoMock package.
package mock_repository

import (
	context "context"
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
	database "mickamy.com/sampay/internal/cli/infra/storage/database"
	model "mickamy.com/sampay/internal/domain/auth/model"
	repository "mickamy.com/sampay/internal/domain/auth/repository"
)

// MockAuthentication is a mock of Authentication interface.
type MockAuthentication struct {
	ctrl     *gomock.Controller
	recorder *MockAuthenticationMockRecorder
	isgomock struct{}
}

// MockAuthenticationMockRecorder is the mock recorder for MockAuthentication.
type MockAuthenticationMockRecorder struct {
	mock *MockAuthentication
}

// NewMockAuthentication creates a new mock instance.
func NewMockAuthentication(ctrl *gomock.Controller) *MockAuthentication {
	mock := &MockAuthentication{ctrl: ctrl}
	mock.recorder = &MockAuthenticationMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockAuthentication) EXPECT() *MockAuthenticationMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m_2 *MockAuthentication) Create(ctx context.Context, m *model.Authentication) error {
	m_2.ctrl.T.Helper()
	ret := m_2.ctrl.Call(m_2, "Create", ctx, m)
	ret0, _ := ret[0].(error)
	return ret0
}

// Create indicates an expected call of Create.
func (mr *MockAuthenticationMockRecorder) Create(ctx, m any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockAuthentication)(nil).Create), ctx, m)
}

// FindByKey mocks base method.
func (m *MockAuthentication) FindByKey(ctx context.Context, key repository.AuthenticationKey, scopes ...database.Scope) (*model.Authentication, error) {
	m.ctrl.T.Helper()
	varargs := []any{ctx, key}
	for _, a := range scopes {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "FindByKey", varargs...)
	ret0, _ := ret[0].(*model.Authentication)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindByKey indicates an expected call of FindByKey.
func (mr *MockAuthenticationMockRecorder) FindByKey(ctx, key any, scopes ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{ctx, key}, scopes...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindByKey", reflect.TypeOf((*MockAuthentication)(nil).FindByKey), varargs...)
}

// FindByTypeAndIdentifier mocks base method.
func (m *MockAuthentication) FindByTypeAndIdentifier(ctx context.Context, authType model.AuthenticationType, identifier string, scopes ...database.Scope) (*model.Authentication, error) {
	m.ctrl.T.Helper()
	varargs := []any{ctx, authType, identifier}
	for _, a := range scopes {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "FindByTypeAndIdentifier", varargs...)
	ret0, _ := ret[0].(*model.Authentication)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindByTypeAndIdentifier indicates an expected call of FindByTypeAndIdentifier.
func (mr *MockAuthenticationMockRecorder) FindByTypeAndIdentifier(ctx, authType, identifier any, scopes ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{ctx, authType, identifier}, scopes...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindByTypeAndIdentifier", reflect.TypeOf((*MockAuthentication)(nil).FindByTypeAndIdentifier), varargs...)
}

// ListByUserID mocks base method.
func (m *MockAuthentication) ListByUserID(ctx context.Context, userID string, scopes ...database.Scope) ([]model.Authentication, error) {
	m.ctrl.T.Helper()
	varargs := []any{ctx, userID}
	for _, a := range scopes {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "ListByUserID", varargs...)
	ret0, _ := ret[0].([]model.Authentication)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListByUserID indicates an expected call of ListByUserID.
func (mr *MockAuthenticationMockRecorder) ListByUserID(ctx, userID any, scopes ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{ctx, userID}, scopes...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListByUserID", reflect.TypeOf((*MockAuthentication)(nil).ListByUserID), varargs...)
}

// Update mocks base method.
func (m_2 *MockAuthentication) Update(ctx context.Context, m *model.Authentication) error {
	m_2.ctrl.T.Helper()
	ret := m_2.ctrl.Call(m_2, "Update", ctx, m)
	ret0, _ := ret[0].(error)
	return ret0
}

// Update indicates an expected call of Update.
func (mr *MockAuthenticationMockRecorder) Update(ctx, m any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockAuthentication)(nil).Update), ctx, m)
}

// WithTx mocks base method.
func (m *MockAuthentication) WithTx(tx *database.DB) repository.Authentication {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "WithTx", tx)
	ret0, _ := ret[0].(repository.Authentication)
	return ret0
}

// WithTx indicates an expected call of WithTx.
func (mr *MockAuthenticationMockRecorder) WithTx(tx any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WithTx", reflect.TypeOf((*MockAuthentication)(nil).WithTx), tx)
}
