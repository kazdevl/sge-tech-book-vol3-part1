// Code generated by MockGen. DO NOT EDIT.
// Source: cache_db.go

// Package cachedb is a generated GoMock package.
package cachedb

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	sqlx "github.com/jmoiron/sqlx"
)

// MockBulkExecutor is a mock of BulkExecutor interface.
type MockBulkExecutor struct {
	ctrl     *gomock.Controller
	recorder *MockBulkExecutorMockRecorder
}

// MockBulkExecutorMockRecorder is the mock recorder for MockBulkExecutor.
type MockBulkExecutorMockRecorder struct {
	mock *MockBulkExecutor
}

// NewMockBulkExecutor creates a new mock instance.
func NewMockBulkExecutor(ctrl *gomock.Controller) *MockBulkExecutor {
	mock := &MockBulkExecutor{ctrl: ctrl}
	mock.recorder = &MockBulkExecutorMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockBulkExecutor) EXPECT() *MockBulkExecutorMockRecorder {
	return m.recorder
}

// BulkDelete mocks base method.
func (m *MockBulkExecutor) BulkDelete(ctx context.Context, tx *sqlx.Tx, contents []CacheContent) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "BulkDelete", ctx, tx, contents)
	ret0, _ := ret[0].(error)
	return ret0
}

// BulkDelete indicates an expected call of BulkDelete.
func (mr *MockBulkExecutorMockRecorder) BulkDelete(ctx, tx, contents interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "BulkDelete", reflect.TypeOf((*MockBulkExecutor)(nil).BulkDelete), ctx, tx, contents)
}

// BulkInsert mocks base method.
func (m *MockBulkExecutor) BulkInsert(ctx context.Context, tx *sqlx.Tx, contents []CacheContent) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "BulkInsert", ctx, tx, contents)
	ret0, _ := ret[0].(error)
	return ret0
}

// BulkInsert indicates an expected call of BulkInsert.
func (mr *MockBulkExecutorMockRecorder) BulkInsert(ctx, tx, contents interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "BulkInsert", reflect.TypeOf((*MockBulkExecutor)(nil).BulkInsert), ctx, tx, contents)
}

// BulkUpdate mocks base method.
func (m *MockBulkExecutor) BulkUpdate(ctx context.Context, tx *sqlx.Tx, contents []CacheContent) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "BulkUpdate", ctx, tx, contents)
	ret0, _ := ret[0].(error)
	return ret0
}

// BulkUpdate indicates an expected call of BulkUpdate.
func (mr *MockBulkExecutorMockRecorder) BulkUpdate(ctx, tx, contents interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "BulkUpdate", reflect.TypeOf((*MockBulkExecutor)(nil).BulkUpdate), ctx, tx, contents)
}
