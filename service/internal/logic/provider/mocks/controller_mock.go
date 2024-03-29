// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/rmarken/reptr/service/internal/logic/provider (interfaces: Controller)
//
// Generated by this command:
//
//	mockgen -destination ./mocks/controller_mock.go -package provider . Controller
//

// Package provider is a generated GoMock package.
package provider

import (
	context "context"
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
)

// MockController is a mock of Controller interface.
type MockController struct {
	ctrl     *gomock.Controller
	recorder *MockControllerMockRecorder
}

// MockControllerMockRecorder is the mock recorder for MockController.
type MockControllerMockRecorder struct {
	mock *MockController
}

// NewMockController creates a new mock instance.
func NewMockController(ctrl *gomock.Controller) *MockController {
	mock := &MockController{ctrl: ctrl}
	mock.recorder = &MockControllerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockController) EXPECT() *MockControllerMockRecorder {
	return m.recorder
}

// GetUserIDFromSubject mocks base method.
func (m *MockController) GetUserIDFromSubject(arg0 context.Context, arg1 string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserIDFromSubject", arg0, arg1)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserIDFromSubject indicates an expected call of GetUserIDFromSubject.
func (mr *MockControllerMockRecorder) GetUserIDFromSubject(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserIDFromSubject", reflect.TypeOf((*MockController)(nil).GetUserIDFromSubject), arg0, arg1)
}

// InsertPair mocks base method.
func (m *MockController) InsertPair(arg0 context.Context, arg1 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "InsertPair", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// InsertPair indicates an expected call of InsertPair.
func (mr *MockControllerMockRecorder) InsertPair(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "InsertPair", reflect.TypeOf((*MockController)(nil).InsertPair), arg0, arg1)
}

// UserExists mocks base method.
func (m *MockController) UserExists(arg0 context.Context, arg1 string) (string, bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UserExists", arg0, arg1)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(bool)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// UserExists indicates an expected call of UserExists.
func (mr *MockControllerMockRecorder) UserExists(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UserExists", reflect.TypeOf((*MockController)(nil).UserExists), arg0, arg1)
}
