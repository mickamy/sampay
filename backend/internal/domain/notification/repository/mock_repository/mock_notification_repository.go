// Code generated by MockGen. DO NOT EDIT.
// Source: notification_repository.go
//
// Generated by this command:
//
//	mockgen -source=notification_repository.go -destination=./mock_repository/mock_notification_repository.go -package=mock_repository
//

// Package mock_repository is a generated GoMock package.
package mock_repository

import (
	context "context"
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
	database "mickamy.com/sampay/internal/cli/infra/storage/database"
	model "mickamy.com/sampay/internal/domain/notification/model"
	repository "mickamy.com/sampay/internal/domain/notification/repository"
)

// MockNotification is a mock of Notification interface.
type MockNotification struct {
	ctrl     *gomock.Controller
	recorder *MockNotificationMockRecorder
	isgomock struct{}
}

// MockNotificationMockRecorder is the mock recorder for MockNotification.
type MockNotificationMockRecorder struct {
	mock *MockNotification
}

// NewMockNotification creates a new mock instance.
func NewMockNotification(ctrl *gomock.Controller) *MockNotification {
	mock := &MockNotification{ctrl: ctrl}
	mock.recorder = &MockNotificationMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockNotification) EXPECT() *MockNotificationMockRecorder {
	return m.recorder
}

// CountByUserID mocks base method.
func (m *MockNotification) CountByUserID(ctx context.Context, userID string, scopes ...database.Scope) (int, error) {
	m.ctrl.T.Helper()
	varargs := []any{ctx, userID}
	for _, a := range scopes {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "CountByUserID", varargs...)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CountByUserID indicates an expected call of CountByUserID.
func (mr *MockNotificationMockRecorder) CountByUserID(ctx, userID any, scopes ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{ctx, userID}, scopes...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CountByUserID", reflect.TypeOf((*MockNotification)(nil).CountByUserID), varargs...)
}

// Create mocks base method.
func (m_2 *MockNotification) Create(ctx context.Context, m *model.Notification) error {
	m_2.ctrl.T.Helper()
	ret := m_2.ctrl.Call(m_2, "Create", ctx, m)
	ret0, _ := ret[0].(error)
	return ret0
}

// Create indicates an expected call of Create.
func (mr *MockNotificationMockRecorder) Create(ctx, m any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockNotification)(nil).Create), ctx, m)
}

// Find mocks base method.
func (m *MockNotification) Find(ctx context.Context, id string, scopes ...database.Scope) (*model.Notification, error) {
	m.ctrl.T.Helper()
	varargs := []any{ctx, id}
	for _, a := range scopes {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Find", varargs...)
	ret0, _ := ret[0].(*model.Notification)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Find indicates an expected call of Find.
func (mr *MockNotificationMockRecorder) Find(ctx, id any, scopes ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{ctx, id}, scopes...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Find", reflect.TypeOf((*MockNotification)(nil).Find), varargs...)
}

// ListByUserID mocks base method.
func (m *MockNotification) ListByUserID(ctx context.Context, userID string, scopes ...database.Scope) ([]model.Notification, error) {
	m.ctrl.T.Helper()
	varargs := []any{ctx, userID}
	for _, a := range scopes {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "ListByUserID", varargs...)
	ret0, _ := ret[0].([]model.Notification)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListByUserID indicates an expected call of ListByUserID.
func (mr *MockNotificationMockRecorder) ListByUserID(ctx, userID any, scopes ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{ctx, userID}, scopes...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListByUserID", reflect.TypeOf((*MockNotification)(nil).ListByUserID), varargs...)
}

// Update mocks base method.
func (m_2 *MockNotification) Update(ctx context.Context, m *model.Notification) error {
	m_2.ctrl.T.Helper()
	ret := m_2.ctrl.Call(m_2, "Update", ctx, m)
	ret0, _ := ret[0].(error)
	return ret0
}

// Update indicates an expected call of Update.
func (mr *MockNotificationMockRecorder) Update(ctx, m any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockNotification)(nil).Update), ctx, m)
}

// WithTx mocks base method.
func (m *MockNotification) WithTx(tx *database.DB) repository.Notification {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "WithTx", tx)
	ret0, _ := ret[0].(repository.Notification)
	return ret0
}

// WithTx indicates an expected call of WithTx.
func (mr *MockNotificationMockRecorder) WithTx(tx any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WithTx", reflect.TypeOf((*MockNotification)(nil).WithTx), tx)
}
