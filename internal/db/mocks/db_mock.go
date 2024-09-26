// Code generated by MockGen. DO NOT EDIT.
// Source: db.go
//
// Generated by this command:
//
//	mockgen -source=db.go -destination=mocks/db_mock.go
//

// Package mock_db is a generated GoMock package.
package mock_db

import (
	reflect "reflect"
	entity "service-chat/internal/db/entity"

	gomock "go.uber.org/mock/gomock"
)

// MockAuthorization is a mock of Authorization interface.
type MockAuthorization struct {
	ctrl     *gomock.Controller
	recorder *MockAuthorizationMockRecorder
}

// MockAuthorizationMockRecorder is the mock recorder for MockAuthorization.
type MockAuthorizationMockRecorder struct {
	mock *MockAuthorization
}

// NewMockAuthorization creates a new mock instance.
func NewMockAuthorization(ctrl *gomock.Controller) *MockAuthorization {
	mock := &MockAuthorization{ctrl: ctrl}
	mock.recorder = &MockAuthorizationMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockAuthorization) EXPECT() *MockAuthorizationMockRecorder {
	return m.recorder
}

// CreateUser mocks base method.
func (m *MockAuthorization) CreateUser(user entity.User) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateUser", user)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateUser indicates an expected call of CreateUser.
func (mr *MockAuthorizationMockRecorder) CreateUser(user any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateUser", reflect.TypeOf((*MockAuthorization)(nil).CreateUser), user)
}

// GetUser mocks base method.
func (m *MockAuthorization) GetUser(user entity.User) (*entity.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUser", user)
	ret0, _ := ret[0].(*entity.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUser indicates an expected call of GetUser.
func (mr *MockAuthorizationMockRecorder) GetUser(user any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUser", reflect.TypeOf((*MockAuthorization)(nil).GetUser), user)
}

// MockChat is a mock of Chat interface.
type MockChat struct {
	ctrl     *gomock.Controller
	recorder *MockChatMockRecorder
}

// MockChatMockRecorder is the mock recorder for MockChat.
type MockChatMockRecorder struct {
	mock *MockChat
}

// NewMockChat creates a new mock instance.
func NewMockChat(ctrl *gomock.Controller) *MockChat {
	mock := &MockChat{ctrl: ctrl}
	mock.recorder = &MockChatMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockChat) EXPECT() *MockChatMockRecorder {
	return m.recorder
}

// CreateChat mocks base method.
func (m *MockChat) CreateChat(in entity.ChatAdd) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateChat", in)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateChat indicates an expected call of CreateChat.
func (mr *MockChatMockRecorder) CreateChat(in any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateChat", reflect.TypeOf((*MockChat)(nil).CreateChat), in)
}

// DeleteChat mocks base method.
func (m *MockChat) DeleteChat(in entity.ChatDelete) ([]entity.DeletedChats, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteChat", in)
	ret0, _ := ret[0].([]entity.DeletedChats)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DeleteChat indicates an expected call of DeleteChat.
func (mr *MockChatMockRecorder) DeleteChat(in any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteChat", reflect.TypeOf((*MockChat)(nil).DeleteChat), in)
}

// GetChat mocks base method.
func (m *MockChat) GetChat(in entity.ChatGet) ([]entity.Chat, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetChat", in)
	ret0, _ := ret[0].([]entity.Chat)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetChat indicates an expected call of GetChat.
func (mr *MockChatMockRecorder) GetChat(in any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetChat", reflect.TypeOf((*MockChat)(nil).GetChat), in)
}

// MockMessage is a mock of Message interface.
type MockMessage struct {
	ctrl     *gomock.Controller
	recorder *MockMessageMockRecorder
}

// MockMessageMockRecorder is the mock recorder for MockMessage.
type MockMessageMockRecorder struct {
	mock *MockMessage
}

// NewMockMessage creates a new mock instance.
func NewMockMessage(ctrl *gomock.Controller) *MockMessage {
	mock := &MockMessage{ctrl: ctrl}
	mock.recorder = &MockMessageMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockMessage) EXPECT() *MockMessageMockRecorder {
	return m.recorder
}

// AddMessage mocks base method.
func (m *MockMessage) AddMessage(in entity.MessageAdd) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddMessage", in)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AddMessage indicates an expected call of AddMessage.
func (mr *MockMessageMockRecorder) AddMessage(in any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddMessage", reflect.TypeOf((*MockMessage)(nil).AddMessage), in)
}

// DeleteMessage mocks base method.
func (m *MockMessage) DeleteMessage(in entity.MessageDel) ([]entity.DelMsg, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteMessage", in)
	ret0, _ := ret[0].([]entity.DelMsg)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DeleteMessage indicates an expected call of DeleteMessage.
func (mr *MockMessageMockRecorder) DeleteMessage(in any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteMessage", reflect.TypeOf((*MockMessage)(nil).DeleteMessage), in)
}

// GetMessage mocks base method.
func (m *MockMessage) GetMessage(in entity.MessageGet) ([]entity.Message, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetMessage", in)
	ret0, _ := ret[0].([]entity.Message)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetMessage indicates an expected call of GetMessage.
func (mr *MockMessageMockRecorder) GetMessage(in any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetMessage", reflect.TypeOf((*MockMessage)(nil).GetMessage), in)
}

// UpdateMessage mocks base method.
func (m *MockMessage) UpdateMessage(in entity.MessageUpdate) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateMessage", in)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateMessage indicates an expected call of UpdateMessage.
func (mr *MockMessageMockRecorder) UpdateMessage(in any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateMessage", reflect.TypeOf((*MockMessage)(nil).UpdateMessage), in)
}
