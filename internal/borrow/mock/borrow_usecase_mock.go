// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/sgraham785/gocleanarch-example/internal/borrow/usecase (interfaces: BorrowUseCase)

// Package mock is a generated GoMock package.
package mock

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	entity "github.com/sgraham785/gocleanarch-example/internal/book/entity"
	entity0 "github.com/sgraham785/gocleanarch-example/internal/user/entity"
)

// MockBorrowUseCase is a mock of BorrowUseCase interface.
type MockBorrowUseCase struct {
	ctrl     *gomock.Controller
	recorder *MockBorrowUseCaseMockRecorder
}

// MockBorrowUseCaseMockRecorder is the mock recorder for MockBorrowUseCase.
type MockBorrowUseCaseMockRecorder struct {
	mock *MockBorrowUseCase
}

// NewMockBorrowUseCase creates a new mock instance.
func NewMockBorrowUseCase(ctrl *gomock.Controller) *MockBorrowUseCase {
	mock := &MockBorrowUseCase{ctrl: ctrl}
	mock.recorder = &MockBorrowUseCaseMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockBorrowUseCase) EXPECT() *MockBorrowUseCaseMockRecorder {
	return m.recorder
}

// Borrow mocks base method.
func (m *MockBorrowUseCase) Borrow(arg0 *entity0.User, arg1 *entity.Book) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Borrow", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Borrow indicates an expected call of Borrow.
func (mr *MockBorrowUseCaseMockRecorder) Borrow(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Borrow", reflect.TypeOf((*MockBorrowUseCase)(nil).Borrow), arg0, arg1)
}

// Return mocks base method.
func (m *MockBorrowUseCase) Return(arg0 *entity.Book) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Return", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Return indicates an expected call of Return.
func (mr *MockBorrowUseCaseMockRecorder) Return(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Return", reflect.TypeOf((*MockBorrowUseCase)(nil).Return), arg0)
}