// Code generated by MockGen. DO NOT EDIT.
// Source: delete_user_profile_image_use_case.go
//
// Generated by this command:
//
//	mockgen -source=delete_user_profile_image_use_case.go -destination=./mock_usecase/mock_delete_user_profile_image_use_case.go -package=mock_usecase
//

// Package mock_usecase is a generated GoMock package.
package mock_usecase

import (
	context "context"
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
	usecase "mickamy.com/sampay/internal/domain/user/usecase"
)

// MockDeleteUserProfileImage is a mock of DeleteUserProfileImage interface.
type MockDeleteUserProfileImage struct {
	ctrl     *gomock.Controller
	recorder *MockDeleteUserProfileImageMockRecorder
	isgomock struct{}
}

// MockDeleteUserProfileImageMockRecorder is the mock recorder for MockDeleteUserProfileImage.
type MockDeleteUserProfileImageMockRecorder struct {
	mock *MockDeleteUserProfileImage
}

// NewMockDeleteUserProfileImage creates a new mock instance.
func NewMockDeleteUserProfileImage(ctrl *gomock.Controller) *MockDeleteUserProfileImage {
	mock := &MockDeleteUserProfileImage{ctrl: ctrl}
	mock.recorder = &MockDeleteUserProfileImageMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockDeleteUserProfileImage) EXPECT() *MockDeleteUserProfileImageMockRecorder {
	return m.recorder
}

// Do mocks base method.
func (m *MockDeleteUserProfileImage) Do(ctx context.Context, input usecase.DeleteUserProfileImageInput) (usecase.DeleteUserProfileImageOutput, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Do", ctx, input)
	ret0, _ := ret[0].(usecase.DeleteUserProfileImageOutput)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Do indicates an expected call of Do.
func (mr *MockDeleteUserProfileImageMockRecorder) Do(ctx, input any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Do", reflect.TypeOf((*MockDeleteUserProfileImage)(nil).Do), ctx, input)
}
