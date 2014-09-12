// Automatically generated by MockGen. DO NOT EDIT!
// Source: ../hydra_client.go

package mock_client

import (
	gomock "github.com/innotech/hydra-go-client/vendors/code.google.com/p/gomock/gomock"
	time "time"

	. "github.com/innotech/hydra-go-client/client"
)

// Mock of HydraClient interface
type MockHydraClient struct {
	ctrl     *gomock.Controller
	recorder *_MockHydraClientRecorder
}

// Recorder for MockHydraClient (not exported)
type _MockHydraClientRecorder struct {
	mock *MockHydraClient
}

func NewMockHydraClient(ctrl *gomock.Controller) *MockHydraClient {
	mock := &MockHydraClient{ctrl: ctrl}
	mock.recorder = &_MockHydraClientRecorder{mock}
	return mock
}

func (_m *MockHydraClient) EXPECT() *_MockHydraClientRecorder {
	return _m.recorder
}

func (_m *MockHydraClient) Get(appId string, forceCacheRefresh bool) ([]string, error) {
	ret := _m.ctrl.Call(_m, "Get", appId, forceCacheRefresh)
	ret0, _ := ret[0].([]string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockHydraClientRecorder) Get(arg0, arg1 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "Get", arg0, arg1)
}

func (_m *MockHydraClient) ReloadAppServers() {
	_m.ctrl.Call(_m, "ReloadAppServers")
}

func (_mr *_MockHydraClientRecorder) ReloadAppServers() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "ReloadAppServers")
}

func (_m *MockHydraClient) ReloadHydraServers() {
	_m.ctrl.Call(_m, "ReloadHydraServers")
}

func (_mr *_MockHydraClientRecorder) ReloadHydraServers() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "ReloadHydraServers")
}

func (_m *MockHydraClient) SetAppsCacheMonitor(monitor *AppsCacheMonitor) {
	_m.ctrl.Call(_m, "SetAppsCacheMonitor", monitor)
}

func (_mr *_MockHydraClientRecorder) SetAppsCacheMonitor(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "SetAppsCacheMonitor", arg0)
}

func (_m *MockHydraClient) SetHydraCacheMonitor(monitor *HydraCacheMonitor) {
	_m.ctrl.Call(_m, "SetHydraCacheMonitor", monitor)
}

func (_mr *_MockHydraClientRecorder) SetHydraCacheMonitor(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "SetHydraCacheMonitor", arg0)
}

func (_m *MockHydraClient) SetMaxNumberOfRetriesPerHydraServer(numberOfRetries uint) {
	_m.ctrl.Call(_m, "SetMaxNumberOfRetriesPerHydraServer", numberOfRetries)
}

func (_mr *_MockHydraClientRecorder) SetMaxNumberOfRetriesPerHydraServer(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "SetMaxNumberOfRetriesPerHydraServer", arg0)
}

func (_m *MockHydraClient) SetWaitBetweenAllServersRetry(duration time.Duration) {
	_m.ctrl.Call(_m, "SetWaitBetweenAllServersRetry", duration)
}

func (_mr *_MockHydraClientRecorder) SetWaitBetweenAllServersRetry(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "SetWaitBetweenAllServersRetry", arg0)
}
