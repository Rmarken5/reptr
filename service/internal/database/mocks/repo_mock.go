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

// AddDeckToGroup mocks base method.
func (m *MockRepository) AddDeckToGroup(arg0 context.Context, arg1, arg2 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddDeckToGroup", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddDeckToGroup indicates an expected call of AddDeckToGroup.
func (mr *MockRepositoryMockRecorder) AddDeckToGroup(arg0, arg1, arg2 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddDeckToGroup", reflect.TypeOf((*MockRepository)(nil).AddDeckToGroup), arg0, arg1, arg2)
}

// AddUserAsMemberOfGroup mocks base method.
func (m *MockRepository) AddUserAsMemberOfGroup(arg0 context.Context, arg1, arg2 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddUserAsMemberOfGroup", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddUserAsMemberOfGroup indicates an expected call of AddUserAsMemberOfGroup.
func (mr *MockRepositoryMockRecorder) AddUserAsMemberOfGroup(arg0, arg1, arg2 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddUserAsMemberOfGroup", reflect.TypeOf((*MockRepository)(nil).AddUserAsMemberOfGroup), arg0, arg1, arg2)
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

// DeleteGroup mocks base method.
func (m *MockRepository) DeleteGroup(arg0 context.Context, arg1 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteGroup", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteGroup indicates an expected call of DeleteGroup.
func (mr *MockRepositoryMockRecorder) DeleteGroup(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteGroup", reflect.TypeOf((*MockRepository)(nil).DeleteGroup), arg0, arg1)
}

// GetBackOfCardByID mocks base method.
func (m *MockRepository) GetBackOfCardByID(arg0 context.Context, arg1, arg2 string) (models.BackOfCard, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetBackOfCardByID", arg0, arg1, arg2)
	ret0, _ := ret[0].(models.BackOfCard)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetBackOfCardByID indicates an expected call of GetBackOfCardByID.
func (mr *MockRepositoryMockRecorder) GetBackOfCardByID(arg0, arg1, arg2 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetBackOfCardByID", reflect.TypeOf((*MockRepository)(nil).GetBackOfCardByID), arg0, arg1, arg2)
}

// GetDeckWithCardsByID mocks base method.
func (m *MockRepository) GetDeckWithCardsByID(arg0 context.Context, arg1 string) (models.DeckWithCards, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetDeckWithCardsByID", arg0, arg1)
	ret0, _ := ret[0].(models.DeckWithCards)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetDeckWithCardsByID indicates an expected call of GetDeckWithCardsByID.
func (mr *MockRepositoryMockRecorder) GetDeckWithCardsByID(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetDeckWithCardsByID", reflect.TypeOf((*MockRepository)(nil).GetDeckWithCardsByID), arg0, arg1)
}

// GetFrontOfCardByID mocks base method.
func (m *MockRepository) GetFrontOfCardByID(arg0 context.Context, arg1, arg2 string) (models.FrontOfCard, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetFrontOfCardByID", arg0, arg1, arg2)
	ret0, _ := ret[0].(models.FrontOfCard)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetFrontOfCardByID indicates an expected call of GetFrontOfCardByID.
func (mr *MockRepositoryMockRecorder) GetFrontOfCardByID(arg0, arg1, arg2 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetFrontOfCardByID", reflect.TypeOf((*MockRepository)(nil).GetFrontOfCardByID), arg0, arg1, arg2)
}

// GetGroupByID mocks base method.
func (m *MockRepository) GetGroupByID(arg0 context.Context, arg1 string) (models.GroupWithDecks, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetGroupByID", arg0, arg1)
	ret0, _ := ret[0].(models.GroupWithDecks)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetGroupByID indicates an expected call of GetGroupByID.
func (mr *MockRepositoryMockRecorder) GetGroupByID(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetGroupByID", reflect.TypeOf((*MockRepository)(nil).GetGroupByID), arg0, arg1)
}

// GetGroupsForUser mocks base method.
func (m *MockRepository) GetGroupsForUser(arg0 context.Context, arg1 string, arg2 time.Time, arg3 *time.Time, arg4, arg5 int) ([]models.Group, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetGroupsForUser", arg0, arg1, arg2, arg3, arg4, arg5)
	ret0, _ := ret[0].([]models.Group)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetGroupsForUser indicates an expected call of GetGroupsForUser.
func (mr *MockRepositoryMockRecorder) GetGroupsForUser(arg0, arg1, arg2, arg3, arg4, arg5 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetGroupsForUser", reflect.TypeOf((*MockRepository)(nil).GetGroupsForUser), arg0, arg1, arg2, arg3, arg4, arg5)
}

// GetGroupsWithDecks mocks base method.
func (m *MockRepository) GetGroupsWithDecks(arg0 context.Context, arg1 time.Time, arg2 *time.Time, arg3, arg4 int) ([]models.GroupWithDecks, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetGroupsWithDecks", arg0, arg1, arg2, arg3, arg4)
	ret0, _ := ret[0].([]models.GroupWithDecks)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetGroupsWithDecks indicates an expected call of GetGroupsWithDecks.
func (mr *MockRepositoryMockRecorder) GetGroupsWithDecks(arg0, arg1, arg2, arg3, arg4 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetGroupsWithDecks", reflect.TypeOf((*MockRepository)(nil).GetGroupsWithDecks), arg0, arg1, arg2, arg3, arg4)
}

// GetUserByUsername mocks base method.
func (m *MockRepository) GetUserByUsername(arg0 context.Context, arg1 string) (models.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserByUsername", arg0, arg1)
	ret0, _ := ret[0].(models.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserByUsername indicates an expected call of GetUserByUsername.
func (mr *MockRepositoryMockRecorder) GetUserByUsername(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserByUsername", reflect.TypeOf((*MockRepository)(nil).GetUserByUsername), arg0, arg1)
}

// GetUserNameForSubject mocks base method.
func (m *MockRepository) GetUserNameForSubject(arg0 context.Context, arg1 string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserNameForSubject", arg0, arg1)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserNameForSubject indicates an expected call of GetUserNameForSubject.
func (mr *MockRepositoryMockRecorder) GetUserNameForSubject(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserNameForSubject", reflect.TypeOf((*MockRepository)(nil).GetUserNameForSubject), arg0, arg1)
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

// InsertGroup mocks base method.
func (m *MockRepository) InsertGroup(arg0 context.Context, arg1 models.Group) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "InsertGroup", arg0, arg1)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// InsertGroup indicates an expected call of InsertGroup.
func (mr *MockRepositoryMockRecorder) InsertGroup(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "InsertGroup", reflect.TypeOf((*MockRepository)(nil).InsertGroup), arg0, arg1)
}

// InsertUser mocks base method.
func (m *MockRepository) InsertUser(arg0 context.Context, arg1 models.User) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "InsertUser", arg0, arg1)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// InsertUser indicates an expected call of InsertUser.
func (mr *MockRepositoryMockRecorder) InsertUser(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "InsertUser", reflect.TypeOf((*MockRepository)(nil).InsertUser), arg0, arg1)
}

// InsertUserSubjectPair mocks base method.
func (m *MockRepository) InsertUserSubjectPair(arg0 context.Context, arg1, arg2 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "InsertUserSubjectPair", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// InsertUserSubjectPair indicates an expected call of InsertUserSubjectPair.
func (mr *MockRepositoryMockRecorder) InsertUserSubjectPair(arg0, arg1, arg2 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "InsertUserSubjectPair", reflect.TypeOf((*MockRepository)(nil).InsertUserSubjectPair), arg0, arg1, arg2)
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

// UpdateGroup mocks base method.
func (m *MockRepository) UpdateGroup(arg0 context.Context, arg1 models.Group) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateGroup", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateGroup indicates an expected call of UpdateGroup.
func (mr *MockRepositoryMockRecorder) UpdateGroup(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateGroup", reflect.TypeOf((*MockRepository)(nil).UpdateGroup), arg0, arg1)
}
