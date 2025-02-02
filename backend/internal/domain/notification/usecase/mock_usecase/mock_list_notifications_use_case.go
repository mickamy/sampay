// Code generated by MockGen. DO NOT EDIT.
// Source: list_notifications_use_case.go
//
// Generated by this command:
//
//	mockgen -source=list_notifications_use_case.go -destination=./mock_usecase/mock_list_notifications_use_case.go -package=mock_usecase
//

// Package mock_usecase is a generated GoMock package.
package mock_usecase

import (
	context "context"
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
	usecase "mickamy.com/sampay/internal/domain/notification/usecase"
)

// MockListNotifications is a mock of ListNotifications interface.
type MockListNotifications struct {
	ctrl     *gomock.Controller
	recorder *MockListNotificationsMockRecorder
	isgomock struct{}
}

// MockListNotificationsMockRecorder is the mock recorder for MockListNotifications.
type MockListNotificationsMockRecorder struct {
	mock *MockListNotifications
}

// NewMockListNotifications creates a new mock instance.
func NewMockListNotifications(ctrl *gomock.Controller) *MockListNotifications {
	mock := &MockListNotifications{ctrl: ctrl}
	mock.recorder = &MockListNotificationsMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockListNotifications) EXPECT() *MockListNotificationsMockRecorder {
	return m.recorder
}

// Do mocks base method.
func (m *MockListNotifications) Do(ctx context.Context, input usecase.ListNotificationsInput) (usecase.ListNotificationsOutput, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Do", ctx, input)
	ret0, _ := ret[0].(usecase.ListNotificationsOutput)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Do indicates an expected call of Do.
func (mr *MockListNotificationsMockRecorder) Do(ctx, input any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Do", reflect.TypeOf((*MockListNotifications)(nil).Do), ctx, input)
}
