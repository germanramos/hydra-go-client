package client_test

import (
	. "github.com/innotech/hydra-go-client/client"
	mock "github.com/innotech/hydra-go-client/client/mock"

	"code.google.com/p/gomock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"time"
)

var _ = Describe("HydraCacheMonitor", func() {
	var (
		mockCtrl          *gomock.Controller
		mockHydraClient   *mock.MockHydraClient
		hydraCacheMonitor *HydraCacheMonitor
	)

	var refreshInterval time.Duration = time.Duration(3000) * time.Millisecond

	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
		mockHydraClient = mock.NewMockHydraClient(mockCtrl)
		hydraCacheMonitor = NewHydraCacheMonitor(mockHydraClient, refreshInterval)
	})

	AfterEach(func() {
		mockCtrl.Finish()
	})

	Context("when new HydraCacheMonitor is instantiated", func() {
		It("should not be running", func() {
			Expect(hydraCacheMonitor.IsRunning()).To(BeFalse())
		})
	})

	// TODO: Refactor to abstract class
	Describe("Get", func() {
		It("should return the refresh interval", func() {
			Expect(hydraCacheMonitor.GetInterval()).To(Equal(refreshInterval))
		})
	})

	Describe("Run", func() {
		It("should run successfully", func() {
			mockHydraClient.EXPECT().ReloadHydraServers()
			hydraCacheMonitor.Run()
			Eventually(func() bool {
				return hydraCacheMonitor.IsRunning()
			}).Should(BeTrue())
			hydraCacheMonitor.Stop()
		})
	})

	Describe("Stop", func() {
		It("should stop the monitor", func() {
			mockHydraClient.EXPECT().ReloadHydraServers()
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
