// Code generated by MockGen. DO NOT EDIT.
// Source: email_verification_repository.go
//
// Generated by this command:
//
//	mockgen -source=email_verification_repository.go -destination=./mock_repository/mock_email_verification_repository.go -package=mock_repository
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

// MockEmailVerification is a mock of EmailVerification interface.
type MockEmailVerification struct {
	ctrl     *gomock.Controller
	recorder *MockEmailVerificationMockRecorder
	isgomock struct{}
}

// MockEmailVerificationMockRecorder is the mock recorder for MockEmailVerification.
type MockEmailVerificationMockRecorder struct {
	mock *MockEmailVerification
}

// NewMockEmailVerification creates a new mock instance.
func NewMockEmailVerification(ctrl *gomock.Controller) *MockEmailVerification {
	mock := &MockEmailVerification{ctrl: ctrl}
	mock.recorder = &MockEmailVerificationMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockEmailVerification) EXPECT() *MockEmailVerificationMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m_2 *MockEmailVerification) Create(ctx context.Context, m *model.EmailVerification) error {
	m_2.ctrl.T.Helper()
	ret := m_2.ctrl.Call(m_2, "Create", ctx, m)
	ret0, _ := ret[0].(error)
	return ret0
}

// Create indicates an expected call of Create.
func (mr *MockEmailVerificationMockRecorder) Create(ctx, m any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockEmailVerification)(nil).Create), ctx, m)
}

// FindByEmail mocks base method.
func (m *MockEmailVerification) FindByEmail(ctx context.Context, email string, scope ...database.Scope) (*model.EmailVerification, error) {
	m.ctrl.T.Helper()
	varargs := []any{ctx, email}
	for _, a := range scope {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "FindByEmail", varargs...)
	ret0, _ := ret[0].(*model.EmailVerification)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindByEmail indicates an expected call of FindByEmail.
func (mr *MockEmailVerificationMockRecorder) FindByEmail(ctx, email any, scope ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{ctx, email}, scope...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindByEmail", reflect.TypeOf((*MockEmailVerification)(nil).FindByEmail), varargs...)
}

// FindByRequestedTokenAndPinCode mocks base method.
func (m *MockEmailVerification) FindByRequestedTokenAndPinCode(ctx context.Context, token, pinCode string, scope ...database.Scope) (*model.EmailVerification, error) {
	m.ctrl.T.Helper()
	varargs := []any{ctx, token, pinCode}
	for _, a := range scope {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "FindByRequestedTokenAndPinCode", varargs...)
	ret0, _ := ret[0].(*model.EmailVerification)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindByRequestedTokenAndPinCode indicates an expected call of FindByRequestedTokenAndPinCode.
func (mr *MockEmailVerificationMockRecorder) FindByRequestedTokenAndPinCode(ctx, token, pinCode any, scope ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{ctx, token, pinCode}, scope...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindByRequestedTokenAndPinCode", reflect.TypeOf((*MockEmailVerification)(nil).FindByRequestedTokenAndPinCode), varargs...)
}

// FindByVerifiedToken mocks base method.
func (m *MockEmailVerification) FindByVerifiedToken(ctx context.Context, token string, scope ...database.Scope) (*model.EmailVerification, error) {
	m.ctrl.T.Helper()
	varargs := []any{ctx, token}
	for _, a := range scope {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "FindByVerifiedToken", varargs...)
	ret0, _ := ret[0].(*model.EmailVerification)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindByVerifiedToken indicates an expected call of FindByVerifiedToken.
func (mr *MockEmailVerificationMockRecorder) FindByVerifiedToken(ctx, token any, scope ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{ctx, token}, scope...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindByVerifiedToken", reflect.TypeOf((*MockEmailVerification)(nil).FindByVerifiedToken), varargs...)
}

// Update mocks base method.
func (m_2 *MockEmailVerification) Update(ctx context.Context, m *model.EmailVerification) error {
	m_2.ctrl.T.Helper()
	ret := m_2.ctrl.Call(m_2, "Update", ctx, m)
	ret0, _ := ret[0].(error)
	return ret0
}

// Update indicates an expected call of Update.
func (mr *MockEmailVerificationMockRecorder) Update(ctx, m any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockEmailVerification)(nil).Update), ctx, m)
}

// WithTx mocks base method.
func (m *MockEmailVerification) WithTx(tx *database.DB) repository.EmailVerification {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "WithTx", tx)
	ret0, _ := ret[0].(repository.EmailVerification)
	return ret0
}

// WithTx indicates an expected call of WithTx.
func (mr *MockEmailVerificationMockRecorder) WithTx(tx any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WithTx", reflect.TypeOf((*MockEmailVerification)(nil).WithTx), tx)
}
