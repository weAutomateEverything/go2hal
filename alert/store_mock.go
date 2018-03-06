// Code generated by MockGen. DO NOT EDIT.
// Source: store.go

// Package mock_alert is a generated GoMock package.
package alert

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

// alertGroup mocks base method
func (m *MockStore) alertGroup() (int64, error) {
	ret := m.ctrl.Call(m, "alertGroup")
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// alertGroup indicates an expected call of alertGroup
func (mr *MockStoreMockRecorder) alertGroup() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "alertGroup", reflect.TypeOf((*MockStore)(nil).alertGroup))
}

// heartbeatGroup mocks base method
func (m *MockStore) heartbeatGroup() (int64, error) {
	ret := m.ctrl.Call(m, "heartbeatGroup")
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// heartbeatGroup indicates an expected call of heartbeatGroup
func (mr *MockStoreMockRecorder) heartbeatGroup() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "heartbeatGroup", reflect.TypeOf((*MockStore)(nil).heartbeatGroup))
}

// nonTechnicalGroup mocks base method
func (m *MockStore) nonTechnicalGroup() (int64, error) {
	ret := m.ctrl.Call(m, "nonTechnicalGroup")
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// nonTechnicalGroup indicates an expected call of nonTechnicalGroup
func (mr *MockStoreMockRecorder) nonTechnicalGroup() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "nonTechnicalGroup", reflect.TypeOf((*MockStore)(nil).nonTechnicalGroup))
}

// setAlertGroup mocks base method
func (m *MockStore) setAlertGroup(AlertGroupID int64) {
	m.ctrl.Call(m, "setAlertGroup", AlertGroupID)
}

// setAlertGroup indicates an expected call of setAlertGroup
func (mr *MockStoreMockRecorder) setAlertGroup(AlertGroupID interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "setAlertGroup", reflect.TypeOf((*MockStore)(nil).setAlertGroup), AlertGroupID)
}

// setHeartbeatGroup mocks base method
func (m *MockStore) setHeartbeatGroup(groupID int64) {
	m.ctrl.Call(m, "setHeartbeatGroup", groupID)
}

// setHeartbeatGroup indicates an expected call of setHeartbeatGroup
func (mr *MockStoreMockRecorder) setHeartbeatGroup(groupID interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "setHeartbeatGroup", reflect.TypeOf((*MockStore)(nil).setHeartbeatGroup), groupID)
}

// setNonTechnicalGroup mocks base method
func (m *MockStore) setNonTechnicalGroup(groupID int64) {
	m.ctrl.Call(m, "setNonTechnicalGroup", groupID)
}

// setNonTechnicalGroup indicates an expected call of setNonTechnicalGroup
func (mr *MockStoreMockRecorder) setNonTechnicalGroup(groupID interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "setNonTechnicalGroup", reflect.TypeOf((*MockStore)(nil).setNonTechnicalGroup), groupID)
}
