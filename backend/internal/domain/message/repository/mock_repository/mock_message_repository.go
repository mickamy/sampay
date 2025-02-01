// Code generated by MockGen. DO NOT EDIT.
// Source: message_repository.go
//
// Generated by this command:
//
//	mockgen -source=message_repository.go -destination=./mock_repository/mock_message_repository.go -package=mock_repository
//

// Package mock_repository is a generated GoMock package.
package mock_repository

import (
	context "context"
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
	database "mickamy.com/sampay/internal/cli/infra/storage/database"
	model "mickamy.com/sampay/internal/domain/message/model"
	repository "mickamy.com/sampay/internal/domain/message/repository"
)

// MockMessage is a mock of Message interface.
type MockMessage struct {
	ctrl     *gomock.Controller
	recorder *MockMessageMockRecorder
	isgomock struct{}
}

// MockMessageMockRecorder is the mock recorder for MockMessage.
type MockMessageMockRecorder struct {
	mock *MockMessage
}

// NewMockMessage creates a new mock instance.
func NewMockMessage(ctrl *gomock.Controller) *MockMessage {
	mock := &MockMessage{ctrl: ctrl}
	mock.recorder = &MockMessageMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockMessage) EXPECT() *MockMessageMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m_2 *MockMessage) Create(ctx context.Context, m *model.Message) error {
	m_2.ctrl.T.Helper()
	ret := m_2.ctrl.Call(m_2, "Create", ctx, m)
	ret0, _ := ret[0].(error)
	return ret0
}

// Create indicates an expected call of Create.
func (mr *MockMessageMockRecorder) Create(ctx, m any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockMessage)(nil).Create), ctx, m)
}

// ListByReceiverID mocks base method.
func (m *MockMessage) ListByReceiverID(ctx context.Context, receiverID string, scopes ...database.Scope) ([]model.Message, error) {
	m.ctrl.T.Helper()
	varargs := []any{ctx, receiverID}
	for _, a := range scopes {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "ListByReceiverID", varargs...)
	ret0, _ := ret[0].([]model.Message)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListByReceiverID indicates an expected call of ListByReceiverID.
func (mr *MockMessageMockRecorder) ListByReceiverID(ctx, receiverID any, scopes ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{ctx, receiverID}, scopes...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListByReceiverID", reflect.TypeOf((*MockMessage)(nil).ListByReceiverID), varargs...)
}

// WithTx mocks base method.
func (m *MockMessage) WithTx(tx *database.DB) repository.Message {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "WithTx", tx)
	ret0, _ := ret[0].(repository.Message)
	return ret0
}

// WithTx indicates an expected call of WithTx.
func (mr *MockMessageMockRecorder) WithTx(tx any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WithTx", reflect.TypeOf((*MockMessage)(nil).WithTx), tx)
}
