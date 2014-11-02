// Automatically generated by MockGen. DO NOT EDIT!
// Source: hydra_requester.go

package mock_client

import (
	gomock "github.com/innotech/hydra-go-client/vendors/code.google.com/p/gomock/gomock"
)

// Mock of Requester interface
type MockRequester struct {
	ctrl     *gomock.Controller
	recorder *_MockRequesterRecorder
}

// Recorder for MockRequester (not exported)
type _MockRequesterRecorder struct {
	mock *MockRequester
}

func NewMockRequester(ctrl *gomock.Controller) *MockRequester {
	mock := &MockRequester{ctrl: ctrl}
	mock.recorder = &_MockRequesterRecorder{mock}
	return mock
}

func (_m *MockRequester) EXPECT() *_MockRequesterRecorder {
	return _m.recorder
}

func (_m *MockRequester) GetServicesById(serverUrl string, id string) ([]string, error) {
	ret := _m.ctrl.Call(_m, "GetServicesById", serverUrl, id)
	ret0, _ := ret[0].([]string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockRequesterRecorder) GetServicesById(arg0, arg1 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "GetServicesById", arg0, arg1)
}
