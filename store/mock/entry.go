// Code generated by MockGen. DO NOT EDIT.
// Source: store/entry.go

// Package mockdb is a generated GoMock package.
package mockdb

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	domain "github.com/pranayhere/simple-wallet/domain"
	store "github.com/pranayhere/simple-wallet/store"
)

// MockEntryRepo is a mock of EntryRepo interface.
type MockEntryRepo struct {
	ctrl     *gomock.Controller
	recorder *MockEntryRepoMockRecorder
}

// MockEntryRepoMockRecorder is the mock recorder for MockEntryRepo.
type MockEntryRepoMockRecorder struct {
	mock *MockEntryRepo
}

// NewMockEntryRepo creates a new mock instance.
func NewMockEntryRepo(ctrl *gomock.Controller) *MockEntryRepo {
	mock := &MockEntryRepo{ctrl: ctrl}
	mock.recorder = &MockEntryRepoMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockEntryRepo) EXPECT() *MockEntryRepoMockRecorder {
	return m.recorder
}

// CreateEntry mocks base method.
func (m *MockEntryRepo) CreateEntry(ctx context.Context, arg store.CreateEntryParams) (domain.Entry, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateEntry", ctx, arg)
	ret0, _ := ret[0].(domain.Entry)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateEntry indicates an expected call of CreateEntry.
func (mr *MockEntryRepoMockRecorder) CreateEntry(ctx, arg interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateEntry", reflect.TypeOf((*MockEntryRepo)(nil).CreateEntry), ctx, arg)
}

// GetEntry mocks base method.
func (m *MockEntryRepo) GetEntry(ctx context.Context, id int64) (domain.Entry, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetEntry", ctx, id)
	ret0, _ := ret[0].(domain.Entry)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetEntry indicates an expected call of GetEntry.
func (mr *MockEntryRepoMockRecorder) GetEntry(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetEntry", reflect.TypeOf((*MockEntryRepo)(nil).GetEntry), ctx, id)
}

// ListEntries mocks base method.
func (m *MockEntryRepo) ListEntries(ctx context.Context, arg store.ListEntriesParams) ([]domain.Entry, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListEntries", ctx, arg)
	ret0, _ := ret[0].([]domain.Entry)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListEntries indicates an expected call of ListEntries.
func (mr *MockEntryRepoMockRecorder) ListEntries(ctx, arg interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListEntries", reflect.TypeOf((*MockEntryRepo)(nil).ListEntries), ctx, arg)
}
