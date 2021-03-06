// Code generated by MockGen. DO NOT EDIT.
// Source: repository.go

// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	reflect "reflect"

	links "github.com/buni/scraper/internal/api/links"
	gomock "github.com/golang/mock/gomock"
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

// CreateLinksJob mocks base method.
func (m *MockRepository) CreateLinksJob(ctx context.Context, job links.Job) (links.Job, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateLinksJob", ctx, job)
	ret0, _ := ret[0].(links.Job)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateLinksJob indicates an expected call of CreateLinksJob.
func (mr *MockRepositoryMockRecorder) CreateLinksJob(ctx, job interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateLinksJob", reflect.TypeOf((*MockRepository)(nil).CreateLinksJob), ctx, job)
}

// CreateLinksJobResult mocks base method.
func (m *MockRepository) CreateLinksJobResult(ctx context.Context, results []links.JobResult) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateLinksJobResult", ctx, results)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateLinksJobResult indicates an expected call of CreateLinksJobResult.
func (mr *MockRepositoryMockRecorder) CreateLinksJobResult(ctx, results interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateLinksJobResult", reflect.TypeOf((*MockRepository)(nil).CreateLinksJobResult), ctx, results)
}

// FinishLinksJob mocks base method.
func (m *MockRepository) FinishLinksJob(ctx context.Context, jobID string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FinishLinksJob", ctx, jobID)
	ret0, _ := ret[0].(error)
	return ret0
}

// FinishLinksJob indicates an expected call of FinishLinksJob.
func (mr *MockRepositoryMockRecorder) FinishLinksJob(ctx, jobID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FinishLinksJob", reflect.TypeOf((*MockRepository)(nil).FinishLinksJob), ctx, jobID)
}

// GetLinksJob mocks base method.
func (m *MockRepository) GetLinksJob(ctx context.Context, jobID string) (links.Job, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetLinksJob", ctx, jobID)
	ret0, _ := ret[0].(links.Job)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetLinksJob indicates an expected call of GetLinksJob.
func (mr *MockRepositoryMockRecorder) GetLinksJob(ctx, jobID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetLinksJob", reflect.TypeOf((*MockRepository)(nil).GetLinksJob), ctx, jobID)
}

// GetLinksJobResult mocks base method.
func (m *MockRepository) GetLinksJobResult(ctx context.Context, jobID string) ([]links.JobResult, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetLinksJobResult", ctx, jobID)
	ret0, _ := ret[0].([]links.JobResult)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetLinksJobResult indicates an expected call of GetLinksJobResult.
func (mr *MockRepositoryMockRecorder) GetLinksJobResult(ctx, jobID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetLinksJobResult", reflect.TypeOf((*MockRepository)(nil).GetLinksJobResult), ctx, jobID)
}
