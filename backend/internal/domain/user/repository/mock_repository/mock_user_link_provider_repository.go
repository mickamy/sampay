// Code generated by MockGen. DO NOT EDIT.
// Source: user_link_provider_repository.go
//
// Generated by this command:
//
//	mockgen -source=user_link_provider_repository.go -destination=./mock_repository/mock_user_link_provider_repository.go -package=mock_repository
//

// Package mock_repository is a generated GoMock package.
package mock_repository

import (
	context "context"
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
	database "mickamy.com/sampay/internal/cli/infra/storage/database"
	model "mickamy.com/sampay/internal/domain/user/model"
	repository "mickamy.com/sampay/internal/domain/user/repository"
)

// MockUserLinkProvider is a mock of UserLinkProvider interface.
type MockUserLinkProvider struct {
	ctrl     *gomock.Controller
	recorder *MockUserLinkProviderMockRecorder
	isgomock struct{}
}

// MockUserLinkProviderMockRecorder is the mock recorder for MockUserLinkProvider.
type MockUserLinkProviderMockRecorder struct {
	mock *MockUserLinkProvider
}

// NewMockUserLinkProvider creates a new mock instance.
func NewMockUserLinkProvider(ctrl *gomock.Controller) *MockUserLinkProvider {
	mock := &MockUserLinkProvider{ctrl: ctrl}
	mock.recorder = &MockUserLinkProviderMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUserLinkProvider) EXPECT() *MockUserLinkProviderMockRecorder {
	return m.recorder
}

// Upsert mocks base method.
func (m_2 *MockUserLinkProvider) Upsert(ctx context.Context, m *model.UserLinkProvider) error {
	m_2.ctrl.T.Helper()
	ret := m_2.ctrl.Call(m_2, "Upsert", ctx, m)
	ret0, _ := ret[0].(error)
	return ret0
}

// Upsert indicates an expected call of Upsert.
func (mr *MockUserLinkProviderMockRecorder) Upsert(ctx, m any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Upsert", reflect.TypeOf((*MockUserLinkProvider)(nil).Upsert), ctx, m)
}

// WithTx mocks base method.
func (m *MockUserLinkProvider) WithTx(tx *database.DB) repository.UserLinkProvider {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "WithTx", tx)
	ret0, _ := ret[0].(repository.UserLinkProvider)
	return ret0
}

// WithTx indicates an expected call of WithTx.
func (mr *MockUserLinkProviderMockRecorder) WithTx(tx any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WithTx", reflect.TypeOf((*MockUserLinkProvider)(nil).WithTx), tx)
}
