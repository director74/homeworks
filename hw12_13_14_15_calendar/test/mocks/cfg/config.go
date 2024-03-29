// Code generated by MockGen. DO NOT EDIT.
// Source: ./internal/cfg/config.go

// Package mock_cfg is a generated GoMock package.
package mock_cfg

import (
	reflect "reflect"

	cfg "github.com/director74/homeworks/hw12_13_14_15_calendar/internal/cfg"
	gomock "github.com/golang/mock/gomock"
)

// MockConfigurable is a mock of Configurable interface.
type MockConfigurable struct {
	ctrl     *gomock.Controller
	recorder *MockConfigurableMockRecorder
}

// MockConfigurableMockRecorder is the mock recorder for MockConfigurable.
type MockConfigurableMockRecorder struct {
	mock *MockConfigurable
}

// NewMockConfigurable creates a new mock instance.
func NewMockConfigurable(ctrl *gomock.Controller) *MockConfigurable {
	mock := &MockConfigurable{ctrl: ctrl}
	mock.recorder = &MockConfigurableMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockConfigurable) EXPECT() *MockConfigurableMockRecorder {
	return m.recorder
}

// GetAppConf mocks base method.
func (m *MockConfigurable) GetAppConf() cfg.AppConf {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAppConf")
	ret0, _ := ret[0].(cfg.AppConf)
	return ret0
}

// GetAppConf indicates an expected call of GetAppConf.
func (mr *MockConfigurableMockRecorder) GetAppConf() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAppConf", reflect.TypeOf((*MockConfigurable)(nil).GetAppConf))
}

// GetDBConf mocks base method.
func (m *MockConfigurable) GetDBConf() cfg.DatabaseConf {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetDBConf")
	ret0, _ := ret[0].(cfg.DatabaseConf)
	return ret0
}

// GetDBConf indicates an expected call of GetDBConf.
func (mr *MockConfigurableMockRecorder) GetDBConf() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetDBConf", reflect.TypeOf((*MockConfigurable)(nil).GetDBConf))
}

// GetLoggerConf mocks base method.
func (m *MockConfigurable) GetLoggerConf() cfg.LoggerConf {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetLoggerConf")
	ret0, _ := ret[0].(cfg.LoggerConf)
	return ret0
}

// GetLoggerConf indicates an expected call of GetLoggerConf.
func (mr *MockConfigurableMockRecorder) GetLoggerConf() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetLoggerConf", reflect.TypeOf((*MockConfigurable)(nil).GetLoggerConf))
}

// GetServersConf mocks base method.
func (m *MockConfigurable) GetServersConf() cfg.ServersConf {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetServersConf")
	ret0, _ := ret[0].(cfg.ServersConf)
	return ret0
}

// GetServersConf indicates an expected call of GetServersConf.
func (mr *MockConfigurableMockRecorder) GetServersConf() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetServersConf", reflect.TypeOf((*MockConfigurable)(nil).GetServersConf))
}

// Parse mocks base method.
func (m *MockConfigurable) Parse(path string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Parse", path)
	ret0, _ := ret[0].(error)
	return ret0
}

// Parse indicates an expected call of Parse.
func (mr *MockConfigurableMockRecorder) Parse(path interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Parse", reflect.TypeOf((*MockConfigurable)(nil).Parse), path)
}
