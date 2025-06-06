// Code generated by MockGen. DO NOT EDIT.
// Source: internal/pvz/delivery/http/handler.go
//
// Generated by this command:
//
//	mockgen -source=internal/pvz/delivery/http/handler.go -destination=internal/pvz/mocks/pvz_service_mock.go -package=mocks
//

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	application "github.com/0x0FACED/pvz-avito/internal/pvz/application"
	domain "github.com/0x0FACED/pvz-avito/internal/pvz/domain"
	domain0 "github.com/0x0FACED/pvz-avito/internal/reception/domain"
	gomock "go.uber.org/mock/gomock"
)

// MockPVZService is a mock of PVZService interface.
type MockPVZService struct {
	ctrl     *gomock.Controller
	recorder *MockPVZServiceMockRecorder
	isgomock struct{}
}

// MockPVZServiceMockRecorder is the mock recorder for MockPVZService.
type MockPVZServiceMockRecorder struct {
	mock *MockPVZService
}

// NewMockPVZService creates a new mock instance.
func NewMockPVZService(ctrl *gomock.Controller) *MockPVZService {
	mock := &MockPVZService{ctrl: ctrl}
	mock.recorder = &MockPVZServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockPVZService) EXPECT() *MockPVZServiceMockRecorder {
	return m.recorder
}

// CloseLastReception mocks base method.
func (m *MockPVZService) CloseLastReception(ctx context.Context, params application.CloseLastReceptionParams) (*domain0.Reception, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CloseLastReception", ctx, params)
	ret0, _ := ret[0].(*domain0.Reception)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CloseLastReception indicates an expected call of CloseLastReception.
func (mr *MockPVZServiceMockRecorder) CloseLastReception(ctx, params any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CloseLastReception", reflect.TypeOf((*MockPVZService)(nil).CloseLastReception), ctx, params)
}

// Create mocks base method.
func (m *MockPVZService) Create(ctx context.Context, params application.CreateParams) (*domain.PVZ, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", ctx, params)
	ret0, _ := ret[0].(*domain.PVZ)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create.
func (mr *MockPVZServiceMockRecorder) Create(ctx, params any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockPVZService)(nil).Create), ctx, params)
}

// DeleteLastProduct mocks base method.
func (m *MockPVZService) DeleteLastProduct(ctx context.Context, params application.DeleteLastProductParams) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteLastProduct", ctx, params)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteLastProduct indicates an expected call of DeleteLastProduct.
func (mr *MockPVZServiceMockRecorder) DeleteLastProduct(ctx, params any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteLastProduct", reflect.TypeOf((*MockPVZService)(nil).DeleteLastProduct), ctx, params)
}

// ListWithReceptions mocks base method.
func (m *MockPVZService) ListWithReceptions(ctx context.Context, params application.ListWithReceptionsParams) ([]*domain.PVZWithReceptions, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListWithReceptions", ctx, params)
	ret0, _ := ret[0].([]*domain.PVZWithReceptions)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListWithReceptions indicates an expected call of ListWithReceptions.
func (mr *MockPVZServiceMockRecorder) ListWithReceptions(ctx, params any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListWithReceptions", reflect.TypeOf((*MockPVZService)(nil).ListWithReceptions), ctx, params)
}
