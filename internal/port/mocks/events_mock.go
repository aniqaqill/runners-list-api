// Code generated by MockGen. DO NOT EDIT.
// Source: internal/port/events.go

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"

	domain "github.com/aniqaqill/runners-list/internal/core/domain"
	gomock "github.com/golang/mock/gomock"
)

// MockEventRepository is a mock of EventRepository interface.
type MockEventRepository struct {
	ctrl     *gomock.Controller
	recorder *MockEventRepositoryMockRecorder
}

// MockEventRepositoryMockRecorder is the mock recorder for MockEventRepository.
type MockEventRepositoryMockRecorder struct {
	mock *MockEventRepository
}

// NewMockEventRepository creates a new mock instance.
func NewMockEventRepository(ctrl *gomock.Controller) *MockEventRepository {
	mock := &MockEventRepository{ctrl: ctrl}
	mock.recorder = &MockEventRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockEventRepository) EXPECT() *MockEventRepositoryMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockEventRepository) Create(event *domain.Events) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", event)
	ret0, _ := ret[0].(error)
	return ret0
}

// Create indicates an expected call of Create.
func (mr *MockEventRepositoryMockRecorder) Create(event interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockEventRepository)(nil).Create), event)
}

// Delete mocks base method.
func (m *MockEventRepository) Delete(event *domain.Events) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", event)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockEventRepositoryMockRecorder) Delete(event interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockEventRepository)(nil).Delete), event)
}

// EventNameExists mocks base method.
func (m *MockEventRepository) EventNameExists(name string) bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "EventNameExists", name)
	ret0, _ := ret[0].(bool)
	return ret0
}

// EventNameExists indicates an expected call of EventNameExists.
func (mr *MockEventRepositoryMockRecorder) EventNameExists(name interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "EventNameExists", reflect.TypeOf((*MockEventRepository)(nil).EventNameExists), name)
}

// FindAll mocks base method.
func (m *MockEventRepository) FindAll() ([]domain.Events, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindAll")
	ret0, _ := ret[0].([]domain.Events)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindAll indicates an expected call of FindAll.
func (mr *MockEventRepositoryMockRecorder) FindAll() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindAll", reflect.TypeOf((*MockEventRepository)(nil).FindAll))
}

// FindByID mocks base method.
func (m *MockEventRepository) FindByID(id uint) (*domain.Events, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindByID", id)
	ret0, _ := ret[0].(*domain.Events)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindByID indicates an expected call of FindByID.
func (mr *MockEventRepositoryMockRecorder) FindByID(id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindByID", reflect.TypeOf((*MockEventRepository)(nil).FindByID), id)
}
