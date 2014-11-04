package client_test

import (
	. "github.com/innotech/hydra-go-client/client"
	mock "github.com/innotech/hydra-go-client/client/mock"

	"github.com/innotech/hydra-go-client/vendors/code.google.com/p/gomock/gomock"
	. "github.com/innotech/hydra-go-client/vendors/github.com/onsi/ginkgo"

	"time"
)

var _ = Describe("ServicesCacheMonitor", func() {
	var (
		mockCtrl             *gomock.Controller
		mockHydraClient      *mock.MockClient
		servicesCacheMonitor *ServicesCacheMonitor

		refreshTime time.Duration = time.Duration(3000) * time.Millisecond
	)

	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
		mockHydraClient = mock.NewMockClient(mockCtrl)
		servicesCacheMonitor = NewServicesCacheMonitor(mockHydraClient, refreshTime)
	})

	AfterEach(func() {
		mockCtrl.Finish()
	})

	Describe("Run", func() {
		It("should call to reload services cache", func() {
			mockHydraClient.EXPECT().ReloadServicesCache()

			servicesCacheMonitor.Run()
		})
	})
})
