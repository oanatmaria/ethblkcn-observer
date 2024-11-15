// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/oanatmaria/ethblkcn-observer/client (interfaces: Client)

// Package client is a generated GoMock package.
package client

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockClient is a mock of Client interface.
type MockClient struct {
	ctrl     *gomock.Controller
	recorder *MockClientMockRecorder
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

// GetBlockByNumber mocks base method.
func (m *MockClient) GetBlockByNumber(arg0 int) (Block, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetBlockByNumber", arg0)
	ret0, _ := ret[0].(Block)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetBlockByNumber indicates an expected call of GetBlockByNumber.
func (mr *MockClientMockRecorder) GetBlockByNumber(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetBlockByNumber", reflect.TypeOf((*MockClient)(nil).GetBlockByNumber), arg0)
}

// GetLatestBlockNumber mocks base method.
func (m *MockClient) GetLatestBlockNumber() (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetLatestBlockNumber")
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetLatestBlockNumber indicates an expected call of GetLatestBlockNumber.
func (mr *MockClientMockRecorder) GetLatestBlockNumber() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetLatestBlockNumber", reflect.TypeOf((*MockClient)(nil).GetLatestBlockNumber))
}