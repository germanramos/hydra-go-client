package client_test

import (
	. "github.com/innotech/hydra-go-client/client"
	mock "github.com/innotech/hydra-go-client/client/mock"

	"github.com/innotech/hydra-go-client/vendors/code.google.com/p/gomock/gomock"
	. "github.com/innotech/hydra-go-client/vendors/github.com/onsi/ginkgo"

	"time"
)

var _ = Describe("HydraCacheMonitor", func() {
	var (
		mockCtrl            *gomock.Controller
		mockHydraClient     *mock.MockClient
		hydraServersMonitor *HydraServiceCacheMonitor

		refreshTime time.Duration = time.Duration(3000) * time.Millisecond
	)

	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
		mockHydraClient = mock.NewMockClient(mockCtrl)
		hydraServersMonitor = NewHydraServiceCacheMonitor(mockHydraClient, refreshTime)
	})

	AfterEach(func() {
		mockCtrl.Finish()
	})

	Describe("Run", func() {
		It("should call to init Hydra service", func() {
			mockHydraClient.EXPECT().ReloadHydraServiceCache()

			hydraServersMonitor.Run()
		})
	})
})
