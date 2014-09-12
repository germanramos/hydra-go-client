package client_test

import (
	. "github.com/innotech/hydra-go-client/client"
	mock "github.com/innotech/hydra-go-client/client/mock"

	"github.com/innotech/hydra-go-client/vendors/code.google.com/p/gomock/gomock"
	. "github.com/innotech/hydra-go-client/vendors/github.com/onsi/ginkgo"
	. "github.com/innotech/hydra-go-client/vendors/github.com/onsi/gomega"

	"time"
)

var _ = Describe("AppsCacheMonitor", func() {
	var (
		mockCtrl         *gomock.Controller
		mockHydraClient  *mock.MockHydraClient
		appsCacheMonitor *AppsCacheMonitor
	)

	var refreshInterval time.Duration = time.Duration(3000) * time.Millisecond

	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
		mockHydraClient = mock.NewMockHydraClient(mockCtrl)
		appsCacheMonitor = NewAppsCacheMonitor(mockHydraClient, refreshInterval)
	})

	AfterEach(func() {
		mockCtrl.Finish()
	})

	Context("when new AppsCacheMonitor is instantiated", func() {
		It("should not be running", func() {
			Expect(appsCacheMonitor.IsRunning()).To(BeFalse())
		})
	})

	// TODO: Refactor to abstract class
	Describe("Get", func() {
		It("should return the refresh interval", func() {
			Expect(appsCacheMonitor.GetInterval()).To(Equal(refreshInterval))
		})
	})

	Describe("Run", func() {
		It("should run successfully", func() {
			mockHydraClient.EXPECT().ReloadAppServers()
			appsCacheMonitor.Run()
			Eventually(func() bool {
				return appsCacheMonitor.IsRunning()
			}).Should(BeTrue())
			appsCacheMonitor.Stop()
		})
	})

	Describe("Stop", func() {
		It("should stop the monitor", func() {
			mockHydraClient.EXPECT().ReloadAppServers()
			appsCacheMonitor.Run()
			Eventually(func() bool {
				return appsCacheMonitor.IsRunning()
			}).Should(BeTrue())
			appsCacheMonitor.Stop()
			Eventually(func() bool {
				return appsCacheMonitor.IsRunning()
			}).Should(BeFalse())
		})
	})
})
