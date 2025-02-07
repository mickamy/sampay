// Code generated by MockGen. DO NOT EDIT.
// Source: usage_category_repository.go
//
// Generated by this command:
//
//	mockgen -source=usage_category_repository.go -destination=./mock_repository/mock_usage_category_repository.go -package=mock_repository
//

// Package mock_repository is a generated GoMock package.
package mock_repository

import (
	context "context"
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
	database "mickamy.com/sampay/internal/cli/infra/storage/database"
	model "mickamy.com/sampay/internal/domain/registration/model"
	repository "mickamy.com/sampay/internal/domain/registration/repository"
)

// MockUsageCategory is a mock of UsageCategory interface.
type MockUsageCategory struct {
	ctrl     *gomock.Controller
	recorder *MockUsageCategoryMockRecorder
	isgomock struct{}
}

// MockUsageCategoryMockRecorder is the mock recorder for MockUsageCategory.
type MockUsageCategoryMockRecorder struct {
	mock *MockUsageCategory
}

// NewMockUsageCategory creates a new mock instance.
func NewMockUsageCategory(ctrl *gomock.Controller) *MockUsageCategory {
	mock := &MockUsageCategory{ctrl: ctrl}
	mock.recorder = &MockUsageCategoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUsageCategory) EXPECT() *MockUsageCategoryMockRecorder {
	return m.recorder
}

// List mocks base method.
func (m *MockUsageCategory) List(ctx context.Context, scopes ...database.Scope) ([]model.UsageCategory, error) {
	m.ctrl.T.Helper()
	varargs := []any{ctx}
	for _, a := range scopes {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "List", varargs...)
	ret0, _ := ret[0].([]model.UsageCategory)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// List indicates an expected call of List.
func (mr *MockUsageCategoryMockRecorder) List(ctx any, scopes ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{ctx}, scopes...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "List", reflect.TypeOf((*MockUsageCategory)(nil).List), varargs...)
}

// Upsert mocks base method.
func (m_2 *MockUsageCategory) Upsert(ctx context.Context, m *model.UsageCategory) error {
	m_2.ctrl.T.Helper()
	ret := m_2.ctrl.Call(m_2, "Upsert", ctx, m)
	ret0, _ := ret[0].(error)
	return ret0
}

// Upsert indicates an expected call of Upsert.
func (mr *MockUsageCategoryMockRecorder) Upsert(ctx, m any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Upsert", reflect.TypeOf((*MockUsageCategory)(nil).Upsert), ctx, m)
}

// WithTx mocks base method.
func (m *MockUsageCategory) WithTx(tx *database.DB) repository.UsageCategory {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "WithTx", tx)
	ret0, _ := ret[0].(repository.UsageCategory)
	return ret0
}

// WithTx indicates an expected call of WithTx.
func (mr *MockUsageCategoryMockRecorder) WithTx(tx any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WithTx", reflect.TypeOf((*MockUsageCategory)(nil).WithTx), tx)
}
