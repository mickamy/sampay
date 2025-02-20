// Code generated by MockGen. DO NOT EDIT.
// Source: create_direct_upload_url_use_case.go
//
// Generated by this command:
//
//	mockgen -source=create_direct_upload_url_use_case.go -destination=./mock_usecase/mock_create_direct_upload_url_use_case.go -package=mock_usecase
//

// Package mock_usecase is a generated GoMock package.
package mock_usecase

import (
	context "context"
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
	usecase "mickamy.com/sampay/internal/domain/common/usecase"
)

// MockCreateDirectUploadURL is a mock of CreateDirectUploadURL interface.
type MockCreateDirectUploadURL struct {
	ctrl     *gomock.Controller
	recorder *MockCreateDirectUploadURLMockRecorder
	isgomock struct{}
}

// MockCreateDirectUploadURLMockRecorder is the mock recorder for MockCreateDirectUploadURL.
type MockCreateDirectUploadURLMockRecorder struct {
	mock *MockCreateDirectUploadURL
}

// NewMockCreateDirectUploadURL creates a new mock instance.
func NewMockCreateDirectUploadURL(ctrl *gomock.Controller) *MockCreateDirectUploadURL {
	mock := &MockCreateDirectUploadURL{ctrl: ctrl}
	mock.recorder = &MockCreateDirectUploadURLMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockCreateDirectUploadURL) EXPECT() *MockCreateDirectUploadURLMockRecorder {
	return m.recorder
}

// Do mocks base method.
func (m *MockCreateDirectUploadURL) Do(ctx context.Context, input usecase.CreateDirectUploadURLInput) (usecase.CreateDirectUploadURLOutput, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Do", ctx, input)
	ret0, _ := ret[0].(usecase.CreateDirectUploadURLOutput)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Do indicates an expected call of Do.
func (mr *MockCreateDirectUploadURLMockRecorder) Do(ctx, input any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Do", reflect.TypeOf((*MockCreateDirectUploadURL)(nil).Do), ctx, input)
}
