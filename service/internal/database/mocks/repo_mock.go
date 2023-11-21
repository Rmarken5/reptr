// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/rmarken/reptr/service/internal/database (interfaces: Repository)
//
// Generated by this command:
//
//	mockgen -destination ./mocks/repo_mock.go -package database . Repository
//
// Package database is a generated GoMock package.
package database

import (
	context "context"
	reflect "reflect"
	time "time"

	models "github.com/rmarken/reptr/service/internal/models"
	gomock "go.uber.org/mock/gomock"
)

// MockRepository is a mock of Repository interface.
type MockRepository struct {
	ctrl     *gomock.Controller
	recorder *MockRepositoryMockRecorder
}

// MockRepositoryMockRecorder is the mock recorder for MockRepository.
type MockRepositoryMockRecorder struct {
	mock *MockRepository
}

// NewMockRepository creates a new mock instance.
func NewMockRepository(ctrl *gomock.Controller) *MockRepository {
	mock := &MockRepository{ctrl: ctrl}
	mock.recorder = &MockRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockRepository) EXPECT() *MockRepositoryMockRecorder {
	return m.recorder
}

// AddUserToDownvote mocks base method.
func (m *MockRepository) AddUserToDownvote(arg0 context.Context, arg1, arg2 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddUserToDownvote", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddUserToDownvote indicates an expected call of AddUserToDownvote.
func (mr *MockRepositoryMockRecorder) AddUserToDownvote(arg0, arg1, arg2 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddUserToDownvote", reflect.TypeOf((*MockRepository)(nil).AddUserToDownvote), arg0, arg1, arg2)
}

// AddUserToUpvote mocks base method.
func (m *MockRepository) AddUserToUpvote(arg0 context.Context, arg1, arg2 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddUserToUpvote", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddUserToUpvote indicates an expected call of AddUserToUpvote.
func (mr *MockRepositoryMockRecorder) AddUserToUpvote(arg0, arg1, arg2 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddUserToUpvote", reflect.TypeOf((*MockRepository)(nil).AddUserToUpvote), arg0, arg1, arg2)
}

// GetWithCards mocks base method.
func (m *MockRepository) GetWithCards(arg0 context.Context, arg1 time.Time, arg2 *time.Time, arg3, arg4 int) ([]models.DeckWithCards, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetWithCards", arg0, arg1, arg2, arg3, arg4)
	ret0, _ := ret[0].([]models.DeckWithCards)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetWithCards indicates an expected call of GetWithCards.
func (mr *MockRepositoryMockRecorder) GetWithCards(arg0, arg1, arg2, arg3, arg4 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetWithCards", reflect.TypeOf((*MockRepository)(nil).GetWithCards), arg0, arg1, arg2, arg3, arg4)
}

// InsertCards mocks base method.
func (m *MockRepository) InsertCards(arg0 context.Context, arg1 []models.Card) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "InsertCards", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// InsertCards indicates an expected call of InsertCards.
func (mr *MockRepositoryMockRecorder) InsertCards(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "InsertCards", reflect.TypeOf((*MockRepository)(nil).InsertCards), arg0, arg1)
}

// InsertDeck mocks base method.
func (m *MockRepository) InsertDeck(arg0 context.Context, arg1 models.Deck) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "InsertDeck", arg0, arg1)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// InsertDeck indicates an expected call of InsertDeck.
func (mr *MockRepositoryMockRecorder) InsertDeck(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "InsertDeck", reflect.TypeOf((*MockRepository)(nil).InsertDeck), arg0, arg1)
}

// RemoveUserFromDownvote mocks base method.
func (m *MockRepository) RemoveUserFromDownvote(arg0 context.Context, arg1, arg2 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RemoveUserFromDownvote", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// RemoveUserFromDownvote indicates an expected call of RemoveUserFromDownvote.
func (mr *MockRepositoryMockRecorder) RemoveUserFromDownvote(arg0, arg1, arg2 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RemoveUserFromDownvote", reflect.TypeOf((*MockRepository)(nil).RemoveUserFromDownvote), arg0, arg1, arg2)
}

// RemoveUserFromUpvote mocks base method.
func (m *MockRepository) RemoveUserFromUpvote(arg0 context.Context, arg1, arg2 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RemoveUserFromUpvote", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// RemoveUserFromUpvote indicates an expected call of RemoveUserFromUpvote.
func (mr *MockRepositoryMockRecorder) RemoveUserFromUpvote(arg0, arg1, arg2 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RemoveUserFromUpvote", reflect.TypeOf((*MockRepository)(nil).RemoveUserFromUpvote), arg0, arg1, arg2)
}

// UpdateCard mocks base method.
func (m *MockRepository) UpdateCard(arg0 context.Context, arg1 models.Card) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateCard", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateCard indicates an expected call of UpdateCard.
func (mr *MockRepositoryMockRecorder) UpdateCard(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateCard", reflect.TypeOf((*MockRepository)(nil).UpdateCard), arg0, arg1)
}
