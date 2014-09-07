package client_test

import (
	. "github.com/innotech/hydra-go-client/client"
	mock "github.com/innotech/hydra-go-client/client/mock"

	"code.google.com/p/gomock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("AppsCacheMonitor", func() {
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

	Context("when new AppsCacheMonitor is instantiated", func() {
		It("should not be running", func() {
			appCacheMonitor := NewAppsCacheMonitor(mockHydraClient, 30000)
			Expect(appCacheMonitor.IsRunning()).To(BeFalse())
		})
	})

	Describe("Run", func() {
		It("should run successfully", func() {
			mockHydraClient.EXPECT().ReloadAppServers()
			appCacheMonitor := NewAppsCacheMonitor(mockHydraClient, 30000)
			appCacheMonitor.Run()
			Eventually(func() bool {
				return appCacheMonitor.IsRunning()
			}).Should(BeTrue())
		})
	})

	Describe("Stop", func() {
		It("should stop the monitor", func() {
			mockHydraClient.EXPECT().ReloadAppServers()
			appCacheMonitor := NewAppsCacheMonitor(mockHydraClient, 30000)
			appCacheMonitor.Run()
			Eventually(func() bool {
				return appCacheMonitor.IsRunning()
			}).Should(BeTrue())
			appCacheMonitor.Stop()
			Eventually(func() bool {
				return appCacheMonitor.IsRunning()
			}).Should(BeFalse())
		})
	})
})
