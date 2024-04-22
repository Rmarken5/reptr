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

// AddUserToDownvoteForCard mocks base method.
func (m *MockRepository) AddUserToDownvoteForCard(arg0 context.Context, arg1, arg2 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddUserToDownvoteForCard", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddUserToDownvoteForCard indicates an expected call of AddUserToDownvoteForCard.
func (mr *MockRepositoryMockRecorder) AddUserToDownvoteForCard(arg0, arg1, arg2 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddUserToDownvoteForCard", reflect.TypeOf((*MockRepository)(nil).AddUserToDownvoteForCard), arg0, arg1, arg2)
}

// AddUserToDownvoteForDeck mocks base method.
func (m *MockRepository) AddUserToDownvoteForDeck(arg0 context.Context, arg1, arg2 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddUserToDownvoteForDeck", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddUserToDownvoteForDeck indicates an expected call of AddUserToDownvoteForDeck.
func (mr *MockRepositoryMockRecorder) AddUserToDownvoteForDeck(arg0, arg1, arg2 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddUserToDownvoteForDeck", reflect.TypeOf((*MockRepository)(nil).AddUserToDownvoteForDeck), arg0, arg1, arg2)
}

// AddUserToUpvoteForCard mocks base method.
func (m *MockRepository) AddUserToUpvoteForCard(arg0 context.Context, arg1, arg2 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddUserToUpvoteForCard", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddUserToUpvoteForCard indicates an expected call of AddUserToUpvoteForCard.
func (mr *MockRepositoryMockRecorder) AddUserToUpvoteForCard(arg0, arg1, arg2 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddUserToUpvoteForCard", reflect.TypeOf((*MockRepository)(nil).AddUserToUpvoteForCard), arg0, arg1, arg2)
}

// AddUserToUpvoteForDeck mocks base method.
func (m *MockRepository) AddUserToUpvoteForDeck(arg0 context.Context, arg1, arg2 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddUserToUpvoteForDeck", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddUserToUpvoteForDeck indicates an expected call of AddUserToUpvoteForDeck.
func (mr *MockRepositoryMockRecorder) AddUserToUpvoteForDeck(arg0, arg1, arg2 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddUserToUpvoteForDeck", reflect.TypeOf((*MockRepository)(nil).AddUserToUpvoteForDeck), arg0, arg1, arg2)
}

// CreateSessionForUserDeck mocks base method.
func (m *MockRepository) CreateSessionForUserDeck(arg0 context.Context, arg1 models.DeckSession) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateSessionForUserDeck", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateSessionForUserDeck indicates an expected call of CreateSessionForUserDeck.
func (mr *MockRepositoryMockRecorder) CreateSessionForUserDeck(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateSessionForUserDeck", reflect.TypeOf((*MockRepository)(nil).CreateSessionForUserDeck), arg0, arg1)
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
func (m *MockRepository) GetBackOfCardByID(arg0 context.Context, arg1, arg2, arg3 string) (models.BackOfCard, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetBackOfCardByID", arg0, arg1, arg2, arg3)
	ret0, _ := ret[0].(models.BackOfCard)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetBackOfCardByID indicates an expected call of GetBackOfCardByID.
func (mr *MockRepositoryMockRecorder) GetBackOfCardByID(arg0, arg1, arg2, arg3 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetBackOfCardByID", reflect.TypeOf((*MockRepository)(nil).GetBackOfCardByID), arg0, arg1, arg2, arg3)
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

// GetDecksForUser mocks base method.
func (m *MockRepository) GetDecksForUser(arg0 context.Context, arg1 string, arg2 time.Time, arg3 *time.Time, arg4, arg5 int) ([]models.GetDeckResults, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetDecksForUser", arg0, arg1, arg2, arg3, arg4, arg5)
	ret0, _ := ret[0].([]models.GetDeckResults)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetDecksForUser indicates an expected call of GetDecksForUser.
func (mr *MockRepositoryMockRecorder) GetDecksForUser(arg0, arg1, arg2, arg3, arg4, arg5 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetDecksForUser", reflect.TypeOf((*MockRepository)(nil).GetDecksForUser), arg0, arg1, arg2, arg3, arg4, arg5)
}

// GetFrontOfCardByID mocks base method.
func (m *MockRepository) GetFrontOfCardByID(arg0 context.Context, arg1, arg2, arg3 string) (models.FrontOfCard, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetFrontOfCardByID", arg0, arg1, arg2, arg3)
	ret0, _ := ret[0].(models.FrontOfCard)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetFrontOfCardByID indicates an expected call of GetFrontOfCardByID.
func (mr *MockRepositoryMockRecorder) GetFrontOfCardByID(arg0, arg1, arg2, arg3 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetFrontOfCardByID", reflect.TypeOf((*MockRepository)(nil).GetFrontOfCardByID), arg0, arg1, arg2, arg3)
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

// GetSessionForUserDeck mocks base method.
func (m *MockRepository) GetSessionForUserDeck(arg0 context.Context, arg1, arg2 string) (models.DeckSession, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSessionForUserDeck", arg0, arg1, arg2)
	ret0, _ := ret[0].(models.DeckSession)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetSessionForUserDeck indicates an expected call of GetSessionForUserDeck.
func (mr *MockRepositoryMockRecorder) GetSessionForUserDeck(arg0, arg1, arg2 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSessionForUserDeck", reflect.TypeOf((*MockRepository)(nil).GetSessionForUserDeck), arg0, arg1, arg2)
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

// RemoveUserFromDownvoteForCard mocks base method.
func (m *MockRepository) RemoveUserFromDownvoteForCard(arg0 context.Context, arg1, arg2 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RemoveUserFromDownvoteForCard", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// RemoveUserFromDownvoteForCard indicates an expected call of RemoveUserFromDownvoteForCard.
func (mr *MockRepositoryMockRecorder) RemoveUserFromDownvoteForCard(arg0, arg1, arg2 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RemoveUserFromDownvoteForCard", reflect.TypeOf((*MockRepository)(nil).RemoveUserFromDownvoteForCard), arg0, arg1, arg2)
}

// RemoveUserFromDownvoteForDeck mocks base method.
func (m *MockRepository) RemoveUserFromDownvoteForDeck(arg0 context.Context, arg1, arg2 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RemoveUserFromDownvoteForDeck", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// RemoveUserFromDownvoteForDeck indicates an expected call of RemoveUserFromDownvoteForDeck.
func (mr *MockRepositoryMockRecorder) RemoveUserFromDownvoteForDeck(arg0, arg1, arg2 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RemoveUserFromDownvoteForDeck", reflect.TypeOf((*MockRepository)(nil).RemoveUserFromDownvoteForDeck), arg0, arg1, arg2)
}

// RemoveUserFromUpvoteForCard mocks base method.
func (m *MockRepository) RemoveUserFromUpvoteForCard(arg0 context.Context, arg1, arg2 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RemoveUserFromUpvoteForCard", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// RemoveUserFromUpvoteForCard indicates an expected call of RemoveUserFromUpvoteForCard.
func (mr *MockRepositoryMockRecorder) RemoveUserFromUpvoteForCard(arg0, arg1, arg2 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RemoveUserFromUpvoteForCard", reflect.TypeOf((*MockRepository)(nil).RemoveUserFromUpvoteForCard), arg0, arg1, arg2)
}

// RemoveUserFromUpvoteForDeck mocks base method.
func (m *MockRepository) RemoveUserFromUpvoteForDeck(arg0 context.Context, arg1, arg2 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RemoveUserFromUpvoteForDeck", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// RemoveUserFromUpvoteForDeck indicates an expected call of RemoveUserFromUpvoteForDeck.
func (mr *MockRepositoryMockRecorder) RemoveUserFromUpvoteForDeck(arg0, arg1, arg2 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RemoveUserFromUpvoteForDeck", reflect.TypeOf((*MockRepository)(nil).RemoveUserFromUpvoteForDeck), arg0, arg1, arg2)
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

// UpdateSession mocks base method.
func (m *MockRepository) UpdateSession(arg0 context.Context, arg1, arg2, arg3 string, arg4 bool) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateSession", arg0, arg1, arg2, arg3, arg4)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateSession indicates an expected call of UpdateSession.
func (mr *MockRepositoryMockRecorder) UpdateSession(arg0, arg1, arg2, arg3, arg4 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateSession", reflect.TypeOf((*MockRepository)(nil).UpdateSession), arg0, arg1, arg2, arg3, arg4)
}
