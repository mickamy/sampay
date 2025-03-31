// Code generated by MockGen. DO NOT EDIT.
// Source: user_repository.go
//
// Generated by this command:
//
//	mockgen -source=user_repository.go -destination=./mock_repository/mock_user_repository.go -package=mock_repository
//

// Package mock_repository is a generated GoMock package.
package mock_repository

import (
	context "context"
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
	database "mickamy.com/sampay/internal/infra/storage/database"
	model "mickamy.com/sampay/internal/domain/user/model"
	repository "mickamy.com/sampay/internal/domain/user/repository"
)

// MockUser is a mock of User interface.
type MockUser struct {
	ctrl     *gomock.Controller
	recorder *MockUserMockRecorder
	isgomock struct{}
}

// MockUserMockRecorder is the mock recorder for MockUser.
type MockUserMockRecorder struct {
	mock *MockUser
}

// NewMockUser creates a new mock instance.
func NewMockUser(ctrl *gomock.Controller) *MockUser {
	mock := &MockUser{ctrl: ctrl}
	mock.recorder = &MockUserMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUser) EXPECT() *MockUserMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m_2 *MockUser) Create(ctx context.Context, m *model.User) error {
	m_2.ctrl.T.Helper()
	ret := m_2.ctrl.Call(m_2, "Create", ctx, m)
	ret0, _ := ret[0].(error)
	return ret0
}

// Create indicates an expected call of Create.
func (mr *MockUserMockRecorder) Create(ctx, m any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockUser)(nil).Create), ctx, m)
}

// Find mocks base method.
func (m *MockUser) Find(ctx context.Context, id string, scopes ...database.Scope) (*model.User, error) {
	m.ctrl.T.Helper()
	varargs := []any{ctx, id}
	for _, a := range scopes {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Find", varargs...)
	ret0, _ := ret[0].(*model.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Find indicates an expected call of Find.
func (mr *MockUserMockRecorder) Find(ctx, id any, scopes ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{ctx, id}, scopes...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Find", reflect.TypeOf((*MockUser)(nil).Find), varargs...)
}

// FindByEmail mocks base method.
func (m *MockUser) FindByEmail(ctx context.Context, email string, scopes ...database.Scope) (*model.User, error) {
	m.ctrl.T.Helper()
	varargs := []any{ctx, email}
	for _, a := range scopes {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "FindByEmail", varargs...)
	ret0, _ := ret[0].(*model.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindByEmail indicates an expected call of FindByEmail.
func (mr *MockUserMockRecorder) FindByEmail(ctx, email any, scopes ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{ctx, email}, scopes...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindByEmail", reflect.TypeOf((*MockUser)(nil).FindByEmail), varargs...)
}

// FindByEmailOrSlug mocks base method.
func (m *MockUser) FindByEmailOrSlug(ctx context.Context, emailOrSlug string, scopes ...database.Scope) (*model.User, error) {
	m.ctrl.T.Helper()
	varargs := []any{ctx, emailOrSlug}
	for _, a := range scopes {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "FindByEmailOrSlug", varargs...)
	ret0, _ := ret[0].(*model.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindByEmailOrSlug indicates an expected call of FindByEmailOrSlug.
func (mr *MockUserMockRecorder) FindByEmailOrSlug(ctx, emailOrSlug any, scopes ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{ctx, emailOrSlug}, scopes...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindByEmailOrSlug", reflect.TypeOf((*MockUser)(nil).FindByEmailOrSlug), varargs...)
}

// FindBySlug mocks base method.
func (m *MockUser) FindBySlug(ctx context.Context, slug string, scopes ...database.Scope) (*model.User, error) {
	m.ctrl.T.Helper()
	varargs := []any{ctx, slug}
	for _, a := range scopes {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "FindBySlug", varargs...)
	ret0, _ := ret[0].(*model.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindBySlug indicates an expected call of FindBySlug.
func (mr *MockUserMockRecorder) FindBySlug(ctx, slug any, scopes ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{ctx, slug}, scopes...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindBySlug", reflect.TypeOf((*MockUser)(nil).FindBySlug), varargs...)
}

// Get mocks base method.
func (m *MockUser) Get(ctx context.Context, id string, scopes ...database.Scope) (model.User, error) {
	m.ctrl.T.Helper()
	varargs := []any{ctx, id}
	for _, a := range scopes {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Get", varargs...)
	ret0, _ := ret[0].(model.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get.
func (mr *MockUserMockRecorder) Get(ctx, id any, scopes ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{ctx, id}, scopes...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockUser)(nil).Get), varargs...)
}

// Update mocks base method.
func (m_2 *MockUser) Update(ctx context.Context, m *model.User) error {
	m_2.ctrl.T.Helper()
	ret := m_2.ctrl.Call(m_2, "Update", ctx, m)
	ret0, _ := ret[0].(error)
	return ret0
}

// Update indicates an expected call of Update.
func (mr *MockUserMockRecorder) Update(ctx, m any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockUser)(nil).Update), ctx, m)
}

// WithTx mocks base method.
func (m *MockUser) WithTx(tx *database.DB) repository.User {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "WithTx", tx)
	ret0, _ := ret[0].(repository.User)
	return ret0
}

// WithTx indicates an expected call of WithTx.
func (mr *MockUserMockRecorder) WithTx(tx any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WithTx", reflect.TypeOf((*MockUser)(nil).WithTx), tx)
}
