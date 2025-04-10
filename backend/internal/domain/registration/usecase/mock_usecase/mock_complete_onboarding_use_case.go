// Code generated by MockGen. DO NOT EDIT.
// Source: complete_onboarding_use_case.go
//
// Generated by this command:
//
//	mockgen -source=complete_onboarding_use_case.go -destination=./mock_usecase/mock_complete_onboarding_use_case.go -package=mock_usecase
//

// Package mock_usecase is a generated GoMock package.
package mock_usecase

import (
	context "context"
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
	usecase "mickamy.com/sampay/internal/domain/registration/usecase"
)

// MockCompleteOnboarding is a mock of CompleteOnboarding interface.
type MockCompleteOnboarding struct {
	ctrl     *gomock.Controller
	recorder *MockCompleteOnboardingMockRecorder
	isgomock struct{}
}

// MockCompleteOnboardingMockRecorder is the mock recorder for MockCompleteOnboarding.
type MockCompleteOnboardingMockRecorder struct {
	mock *MockCompleteOnboarding
}

// NewMockCompleteOnboarding creates a new mock instance.
func NewMockCompleteOnboarding(ctrl *gomock.Controller) *MockCompleteOnboarding {
	mock := &MockCompleteOnboarding{ctrl: ctrl}
	mock.recorder = &MockCompleteOnboardingMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockCompleteOnboarding) EXPECT() *MockCompleteOnboardingMockRecorder {
	return m.recorder
}

// Do mocks base method.
func (m *MockCompleteOnboarding) Do(ctx context.Context, input usecase.CompleteOnboardingInput) (usecase.CompleteOnboardingOutput, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Do", ctx, input)
	ret0, _ := ret[0].(usecase.CompleteOnboardingOutput)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Do indicates an expected call of Do.
func (mr *MockCompleteOnboardingMockRecorder) Do(ctx, input any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Do", reflect.TypeOf((*MockCompleteOnboarding)(nil).Do), ctx, input)
}
