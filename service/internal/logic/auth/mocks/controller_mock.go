// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/rmarken/reptr/service/internal/logic/auth (interfaces: Authentication)
//
// Generated by this command:
//
//	mockgen -destination ./mocks/controller_mock.go -package auth . Authentication
//
// Package auth is a generated GoMock package.
package auth

import (
	context "context"
	reflect "reflect"

	oidc "github.com/coreos/go-oidc/v3/oidc"
	models "github.com/rmarken/reptr/service/internal/models"
	gomock "go.uber.org/mock/gomock"
	oauth2 "golang.org/x/oauth2"
)

// MockAuthentication is a mock of Authentication interface.
type MockAuthentication struct {
	ctrl     *gomock.Controller
	recorder *MockAuthenticationMockRecorder
}

// MockAuthenticationMockRecorder is the mock recorder for MockAuthentication.
type MockAuthenticationMockRecorder struct {
	mock *MockAuthentication
}

// NewMockAuthentication creates a new mock instance.
func NewMockAuthentication(ctrl *gomock.Controller) *MockAuthentication {
	mock := &MockAuthentication{ctrl: ctrl}
	mock.recorder = &MockAuthenticationMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockAuthentication) EXPECT() *MockAuthenticationMockRecorder {
	return m.recorder
}

// PasswordCredentialsToken mocks base method.
func (m *MockAuthentication) PasswordCredentialsToken(arg0 context.Context, arg1, arg2 string) (*oauth2.Token, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PasswordCredentialsToken", arg0, arg1, arg2)
	ret0, _ := ret[0].(*oauth2.Token)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// PasswordCredentialsToken indicates an expected call of PasswordCredentialsToken.
func (mr *MockAuthenticationMockRecorder) PasswordCredentialsToken(arg0, arg1, arg2 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PasswordCredentialsToken", reflect.TypeOf((*MockAuthentication)(nil).PasswordCredentialsToken), arg0, arg1, arg2)
}

// RegisterUser mocks base method.
func (m *MockAuthentication) RegisterUser(arg0 context.Context, arg1, arg2 string) (models.RegistrationUser, models.RegistrationError, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RegisterUser", arg0, arg1, arg2)
	ret0, _ := ret[0].(models.RegistrationUser)
	ret1, _ := ret[1].(models.RegistrationError)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// RegisterUser indicates an expected call of RegisterUser.
func (mr *MockAuthenticationMockRecorder) RegisterUser(arg0, arg1, arg2 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RegisterUser", reflect.TypeOf((*MockAuthentication)(nil).RegisterUser), arg0, arg1, arg2)
}

// VerifyIDToken mocks base method.
func (m *MockAuthentication) VerifyIDToken(arg0 context.Context, arg1 string) (*oidc.IDToken, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "VerifyIDToken", arg0, arg1)
	ret0, _ := ret[0].(*oidc.IDToken)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// VerifyIDToken indicates an expected call of VerifyIDToken.
func (mr *MockAuthenticationMockRecorder) VerifyIDToken(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "VerifyIDToken", reflect.TypeOf((*MockAuthentication)(nil).VerifyIDToken), arg0, arg1)
}
