// Automatically generated by MockGen. DO NOT EDIT!
// Source: hydra_client_factory.go

package mock_client

import (
	. "github.com/innotech/hydra-go-client/client"
	gomock "github.com/innotech/hydra-go-client/vendors/code.google.com/p/gomock/gomock"

	"time"
)

// Mock of clientMaker interface
type MockclientMaker struct {
	ctrl     *gomock.Controller
	recorder *_MockclientMakerRecorder
}

// Recorder for MockclientMaker (not exported)
type _MockclientMakerRecorder struct {
	mock *MockclientMaker
}

func NewMockclientMaker(ctrl *gomock.Controller) *MockclientMaker {
	mock := &MockclientMaker{ctrl: ctrl}
	mock.recorder = &_MockclientMakerRecorder{mock}
	return mock
}

func (_m *MockclientMaker) EXPECT() *_MockclientMakerRecorder {
	return _m.recorder
}

func (_m *MockclientMaker) MakeClient(seedHydraServers []string) Client {
	ret := _m.ctrl.Call(_m, "MakeClient", seedHydraServers)
	ret0, _ := ret[0].(Client)
	return ret0
}

func (_mr *_MockclientMakerRecorder) MakeClient(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "MakeClient", arg0)
}

// Mock of hydraMonitorMaker interface
type MockhydraMonitorMaker struct {
	ctrl     *gomock.Controller
	recorder *_MockhydraMonitorMakerRecorder
}

// Recorder for MockhydraMonitorMaker (not exported)
type _MockhydraMonitorMakerRecorder struct {
	mock *MockhydraMonitorMaker
}

func NewMockhydraMonitorMaker(ctrl *gomock.Controller) *MockhydraMonitorMaker {
	mock := &MockhydraMonitorMaker{ctrl: ctrl}
	mock.recorder = &_MockhydraMonitorMakerRecorder{mock}
	return mock
}

func (_m *MockhydraMonitorMaker) EXPECT() *_MockhydraMonitorMakerRecorder {
	return _m.recorder
}

func (_m *MockhydraMonitorMaker) MakeHydraMonitor(hydraClient Client, refreshTime time.Duration) CacheMonitor {
	ret := _m.ctrl.Call(_m, "MakeHydraMonitor", hydraClient, refreshTime)
	ret0, _ := ret[0].(CacheMonitor)
	return ret0
}

func (_mr *_MockhydraMonitorMakerRecorder) MakeHydraMonitor(arg0, arg1 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "MakeHydraMonitor", arg0, arg1)
}

// Mock of appsMonitorMaker interface
type MockappsMonitorMaker struct {
	ctrl     *gomock.Controller
	recorder *_MockappsMonitorMakerRecorder
}

// Recorder for MockappsMonitorMaker (not exported)
type _MockappsMonitorMakerRecorder struct {
	mock *MockappsMonitorMaker
}

func NewMockappsMonitorMaker(ctrl *gomock.Controller) *MockappsMonitorMaker {
	mock := &MockappsMonitorMaker{ctrl: ctrl}
	mock.recorder = &_MockappsMonitorMakerRecorder{mock}
	return mock
}

func (_m *MockappsMonitorMaker) EXPECT() *_MockappsMonitorMakerRecorder {
	return _m.recorder
}

func (_m *MockappsMonitorMaker) MakeAppsMonitor(hydraClient Client, refreshTime time.Duration) CacheMonitor {
	ret := _m.ctrl.Call(_m, "MakeAppsMonitor", hydraClient, refreshTime)
	ret0, _ := ret[0].(CacheMonitor)
	return ret0
}

func (_mr *_MockappsMonitorMakerRecorder) MakeAppsMonitor(arg0, arg1 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "MakeAppsMonitor", arg0, arg1)
}