// Code generated by MockGen. DO NOT EDIT.
// Source: service.go

// Package mock_alert is a generated GoMock package.
package alert

import (
	gomock "github.com/golang/mock/gomock"
	context "golang.org/x/net/context"
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

// SendAlert mocks base method
func (m *MockService) SendAlert(ctx context.Context, chatId uint32, message string) error {
	ret := m.ctrl.Call(m, "SendAlert", ctx, chatId, message)
	ret0, _ := ret[0].(error)
	return ret0
}

// SendAlert indicates an expected call of SendAlert
func (mr *MockServiceMockRecorder) SendAlert(ctx, chatId, message interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendAlert", reflect.TypeOf((*MockService)(nil).SendAlert), ctx, chatId, message)
}

// SendImageToAlertGroup mocks base method
func (m *MockService) SendImageToAlertGroup(ctx context.Context, chatid uint32, image []byte) error {
	ret := m.ctrl.Call(m, "SendImageToAlertGroup", ctx, chatid, image)
	ret0, _ := ret[0].(error)
	return ret0
}

// SendImageToAlertGroup indicates an expected call of SendImageToAlertGroup
func (mr *MockServiceMockRecorder) SendImageToAlertGroup(ctx, chatid, image interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendImageToAlertGroup", reflect.TypeOf((*MockService)(nil).SendImageToAlertGroup), ctx, chatid, image)
}

// SendDocumentToAlertGroup mocks base method
func (m *MockService) SendDocumentToAlertGroup(ctx context.Context, chatid uint32, document []byte, extension string) error {
	ret := m.ctrl.Call(m, "SendDocumentToAlertGroup", ctx, chatid, document, extension)
	ret0, _ := ret[0].(error)
	return ret0
}

// SendDocumentToAlertGroup indicates an expected call of SendDocumentToAlertGroup
func (mr *MockServiceMockRecorder) SendDocumentToAlertGroup(ctx, chatid, document, extension interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendDocumentToAlertGroup", reflect.TypeOf((*MockService)(nil).SendDocumentToAlertGroup), ctx, chatid, document, extension)
}

// SendError mocks base method
func (m *MockService) SendError(ctx context.Context, err error) error {
	ret := m.ctrl.Call(m, "SendError", ctx, err)
	ret0, _ := ret[0].(error)
	return ret0
}

// SendError indicates an expected call of SendError
func (mr *MockServiceMockRecorder) SendError(ctx, err interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendError", reflect.TypeOf((*MockService)(nil).SendError), ctx, err)
}

// SendErrorImage mocks base method
func (m *MockService) SendErrorImage(ctx context.Context, image []byte) error {
	ret := m.ctrl.Call(m, "SendErrorImage", ctx, image)
	ret0, _ := ret[0].(error)
	return ret0
}

// SendErrorImage indicates an expected call of SendErrorImage
func (mr *MockServiceMockRecorder) SendErrorImage(ctx, image interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendErrorImage", reflect.TypeOf((*MockService)(nil).SendErrorImage), ctx, image)
}
