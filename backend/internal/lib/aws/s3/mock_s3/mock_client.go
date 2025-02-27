// Code generated by MockGen. DO NOT EDIT.
// Source: client.go
//
// Generated by this command:
//
//	mockgen -source=client.go -destination=./mock_s3/mock_client.go -package=mock_s3
//

// Package mock_s3 is a generated GoMock package.
package mock_s3

import (
	context "context"
	io "io"
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
)

// MockClient is a mock of Client interface.
type MockClient struct {
	ctrl     *gomock.Controller
	recorder *MockClientMockRecorder
	isgomock struct{}
}

// MockClientMockRecorder is the mock recorder for MockClient.
type MockClientMockRecorder struct {
	mock *MockClient
}

// NewMockClient creates a new mock instance.
func NewMockClient(ctrl *gomock.Controller) *MockClient {
	mock := &MockClient{ctrl: ctrl}
	mock.recorder = &MockClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockClient) EXPECT() *MockClientMockRecorder {
	return m.recorder
}

// GeneratePresignedURL mocks base method.
func (m *MockClient) GeneratePresignedURL(ctx context.Context, bucket, key string, secs int64) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GeneratePresignedURL", ctx, bucket, key, secs)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GeneratePresignedURL indicates an expected call of GeneratePresignedURL.
func (mr *MockClientMockRecorder) GeneratePresignedURL(ctx, bucket, key, secs any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GeneratePresignedURL", reflect.TypeOf((*MockClient)(nil).GeneratePresignedURL), ctx, bucket, key, secs)
}

// GetObject mocks base method.
func (m *MockClient) GetObject(ctx context.Context, bucket, key string) (io.ReadCloser, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetObject", ctx, bucket, key)
	ret0, _ := ret[0].(io.ReadCloser)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetObject indicates an expected call of GetObject.
func (mr *MockClientMockRecorder) GetObject(ctx, bucket, key any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetObject", reflect.TypeOf((*MockClient)(nil).GetObject), ctx, bucket, key)
}

// PutObject mocks base method.
func (m *MockClient) PutObject(ctx context.Context, bucket, key string, body io.Reader) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PutObject", ctx, bucket, key, body)
	ret0, _ := ret[0].(error)
	return ret0
}

// PutObject indicates an expected call of PutObject.
func (mr *MockClientMockRecorder) PutObject(ctx, bucket, key, body any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PutObject", reflect.TypeOf((*MockClient)(nil).PutObject), ctx, bucket, key, body)
}
