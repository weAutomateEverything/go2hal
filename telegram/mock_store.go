// Code generated by MockGen. DO NOT EDIT.
// Source: store.go

// Package mock_telegram is a generated GoMock package.
package telegram

import (
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockStore is a mock of Store interface
type MockStore struct {
	ctrl     *gomock.Controller
	recorder *MockStoreMockRecorder
}

// MockStoreMockRecorder is the mock recorder for MockStore
type MockStoreMockRecorder struct {
	mock *MockStore
}

// NewMockStore creates a new mock instance
func NewMockStore(ctrl *gomock.Controller) *MockStore {
	mock := &MockStore{ctrl: ctrl}
	mock.recorder = &MockStoreMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockStore) EXPECT() *MockStoreMockRecorder {
	return m.recorder
}

// SetState mocks base method
func (m *MockStore) SetState(user int, state string, field []string) error {
	ret := m.ctrl.Call(m, "SetState", user, state, field)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetState indicates an expected call of SetState
func (mr *MockStoreMockRecorder) SetState(user, state, field interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetState", reflect.TypeOf((*MockStore)(nil).SetState), user, state, field)
}

// listBots mocks base method
func (m *MockStore) listBots() []string {
	ret := m.ctrl.Call(m, "listBots")
	ret0, _ := ret[0].([]string)
	return ret0
}

// listBots indicates an expected call of listBots
func (mr *MockStoreMockRecorder) listBots() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "listBots", reflect.TypeOf((*MockStore)(nil).listBots))
}

// getState mocks base method
func (m *MockStore) getState(user int) State {
	ret := m.ctrl.Call(m, "getState", user)
	ret0, _ := ret[0].(State)
	return ret0
}

// getState indicates an expected call of getState
func (mr *MockStoreMockRecorder) getState(user interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "getState", reflect.TypeOf((*MockStore)(nil).getState), user)
}
