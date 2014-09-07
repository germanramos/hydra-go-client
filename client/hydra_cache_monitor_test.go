package client_test

import (
	. "github.com/innotech/hydra-go-client/client"
	mock "github.com/innotech/hydra-go-client/client/mock"

	"code.google.com/p/gomock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("HydraCacheMonitor", func() {
	var (
		mockCtrl        *gomock.Controller
		mockHydraClient *mock.MockHydraClient
	)

	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
		mockHydraClient = mock.NewMockHydraClient(mockCtrl)
	})

	AfterEach(func() {
		mockCtrl.Finish()
	})

	Context("when new HydraCacheMonitor is instantiated", func() {
		It("should not be running", func() {
			hydraCacheMonitor := NewHydraCacheMonitor(mockHydraClient, 30000)
			Expect(hydraCacheMonitor.IsRunning()).To(BeFalse())
		})
	})

	Describe("Run", func() {
		It("should run successfully", func() {
			mockHydraClient.EXPECT().ReloadHydraServers()
			hydraCacheMonitor := NewHydraCacheMonitor(mockHydraClient, 30000)
			hydraCacheMonitor.Run()
			Eventually(func() bool {
				return hydraCacheMonitor.IsRunning()
			}).Should(BeTrue())
		})
	})

	Describe("Stop", func() {
		It("should stop the monitor", func() {
			mockHydraClient.EXPECT().ReloadHydraServers()
			hydraCacheMonitor := NewHydraCacheMonitor(mockHydraClient, 30000)
			hydraCacheMonitor.Run()
			Eventually(func() bool {
				return hydraCacheMonitor.IsRunning()
			}).Should(BeTrue())
			hydraCacheMonitor.Stop()
			Eventually(func() bool {
				return hydraCacheMonitor.IsRunning()
			}).Should(BeFalse())
		})
	})
})
