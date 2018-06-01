// Code generated by MockGen. DO NOT EDIT.
// Source: service.go

// Package mock_firstCall is a generated GoMock package.
package firstCall

import (
	context "context"
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockService is a mock of Service interface
type MockService struct {
	ctrl     *gomock.Controller
	recorder *MockServiceMockRecorder
}

// MockServiceMockRecorder is the mock recorder for MockService
type MockServiceMockRecorder struct {
	mock *MockService
}

// NewMockService creates a new mock instance
func NewMockService(ctrl *gomock.Controller) *MockService {
	mock := &MockService{ctrl: ctrl}
	mock.recorder = &MockServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockService) EXPECT() *MockServiceMockRecorder {
	return m.recorder
}

// GetFirstCall mocks base method
func (m *MockService) GetFirstCall(ctx context.Context, chat uint32) (string, string, error) {
	ret := m.ctrl.Call(m, "GetFirstCall", ctx, chat)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(string)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// GetFirstCall indicates an expected call of GetFirstCall
func (mr *MockServiceMockRecorder) GetFirstCall(ctx, chat interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetFirstCall", reflect.TypeOf((*MockService)(nil).GetFirstCall), ctx, chat)
}

// AddCalloutFunc mocks base method
func (m *MockService) AddCalloutFunc(function CalloutFunction) {
	m.ctrl.Call(m, "AddCalloutFunc", function)
}

// AddCalloutFunc indicates an expected call of AddCalloutFunc
func (mr *MockServiceMockRecorder) AddCalloutFunc(function interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddCalloutFunc", reflect.TypeOf((*MockService)(nil).AddCalloutFunc), function)
}

// MockCalloutFunction is a mock of CalloutFunction interface
type MockCalloutFunction struct {
	ctrl     *gomock.Controller
	recorder *MockCalloutFunctionMockRecorder
}

// MockCalloutFunctionMockRecorder is the mock recorder for MockCalloutFunction
type MockCalloutFunctionMockRecorder struct {
	mock *MockCalloutFunction
}

// NewMockCalloutFunction creates a new mock instance
func NewMockCalloutFunction(ctrl *gomock.Controller) *MockCalloutFunction {
	mock := &MockCalloutFunction{ctrl: ctrl}
	mock.recorder = &MockCalloutFunctionMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockCalloutFunction) EXPECT() *MockCalloutFunctionMockRecorder {
	return m.recorder
}

// GetFirstCallDetails mocks base method
func (m *MockCalloutFunction) GetFirstCallDetails(ctx context.Context, chat uint32) (string, string, error) {
	ret := m.ctrl.Call(m, "GetFirstCallDetails", ctx, chat)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(string)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// GetFirstCallDetails indicates an expected call of GetFirstCallDetails
func (mr *MockCalloutFunctionMockRecorder) GetFirstCallDetails(ctx, chat interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetFirstCallDetails", reflect.TypeOf((*MockCalloutFunction)(nil).GetFirstCallDetails), ctx, chat)
}
