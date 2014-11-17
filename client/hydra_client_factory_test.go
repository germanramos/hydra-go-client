package client_test

import (
	. "github.com/innotech/hydra-go-client/client"
	mock "github.com/innotech/hydra-go-client/client/mock"

	"github.com/innotech/hydra-go-client/vendors/code.google.com/p/gomock/gomock"
	. "github.com/innotech/hydra-go-client/vendors/github.com/onsi/ginkgo"
	. "github.com/innotech/hydra-go-client/vendors/github.com/onsi/gomega"

	"time"
)

var _ = FDescribe("HydraClientFactory", func() {
	const (
		seed_server string = "http://localhost:8080"
	)

	var (
		mockCtrl              *gomock.Controller
		mockAppsMonitorMaker  *mock.MockappsMonitorMaker
		mockClientMaker       *mock.MockclientMaker
		mockHydraMonitorMaker *mock.MockhydraMonitorMaker
		mockHydraClient       *mock.MockClient
		mockHydraMonitor      *mock.MockCacheMonitor
		mockServicesMonitor   *mock.MockCacheMonitor

		test_hydra_servers []string = []string{seed_server}
	)

	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
		mockAppsMonitorMaker = mock.NewMockappsMonitorMaker(mockCtrl)
		mockClientMaker = mock.NewMockclientMaker(mockCtrl)
		mockHydraMonitorMaker = mock.NewMockhydraMonitorMaker(mockCtrl)
		mockHydraClient = mock.NewMockClient(mockCtrl)
		mockHydraMonitor = mock.NewMockCacheMonitor(mockCtrl)
		mockServicesMonitor = mock.NewMockCacheMonitor(mockCtrl)
	})

	AfterEach(func() {
		mockCtrl.Finish()
		Reset()
	})

	It("should get an unique Hydra client", func() {
		mockClientMaker.EXPECT().MakeClient(gomock.Eq(test_hydra_servers)).Return(mockHydraClient)
		mockHydraClient.EXPECT().SetMaxNumberOfRetries(gomock.Any()).Times(1)
		mockHydraClient.EXPECT().SetWaitBetweenAllServersRetry(gomock.Any()).Times(1)
		c1 := mockHydraClient.EXPECT().ReloadHydraServiceCache().Times(0)
		mockHydraClient.EXPECT().ReloadHydraServiceCache().AnyTimes().After(c1)
		c2 := mockHydraClient.EXPECT().ReloadServicesCache().Times(0)
		mockHydraClient.EXPECT().ReloadServicesCache().AnyTimes().After(c2)

		factory, _ := Config(test_hydra_servers)
		factory.ClientInstantiator = mockClientMaker
		hydraClient := factory.Build()
		anotherHydraClient := GetHydraClient()

		Expect(hydraClient).ToNot(BeNil(), "Client must not be nil")
		Expect(anotherHydraClient).ToNot(BeNil(), "The second client must not be nil")
		Expect(hydraClient).To(Equal(anotherHydraClient), "The clients must be the same")
	})

	Context("when calls to Config many times", func() {
		It("should get an unique Hydra client", func() {
			mockClientMaker.EXPECT().MakeClient(gomock.Eq(test_hydra_servers)).Return(mockHydraClient)
			mockHydraClient.EXPECT().SetMaxNumberOfRetries(gomock.Any()).Times(1)
			mockHydraClient.EXPECT().SetWaitBetweenAllServersRetry(gomock.Any()).Times(1)
			c1 := mockHydraClient.EXPECT().ReloadHydraServiceCache().Times(0)
			mockHydraClient.EXPECT().ReloadHydraServiceCache().AnyTimes().After(c1)
			c2 := mockHydraClient.EXPECT().ReloadServicesCache().Times(0)
			mockHydraClient.EXPECT().ReloadServicesCache().AnyTimes().After(c2)

			factory, _ := Config(test_hydra_servers)
			factory.ClientInstantiator = mockClientMaker
			hydraClient := factory.Build()
			factory2, _ := Config(test_hydra_servers)
			anotherHydraClient := factory2.Build()

			Expect(hydraClient).NotTo(BeNil(), "Client must not be nil")
			Expect(anotherHydraClient).NotTo(BeNil(), "The second client must not be nil")
			Expect(hydraClient).To(Equal(anotherHydraClient), "The clients must be the same")
		})
	})

	Describe("Config", func() {
		Context("when nil seed servers is passed", func() {
			It("should not create a client", func() {
				factory, err := Config(nil)
				Expect(err).To(HaveOccurred(), "Must return an error")
				// TODO: Match error
				Expect(factory).To(BeNil(), "Must not return an HydraClientFactory")
			})
		})
		Context("when none seed servers is passed", func() {
			It("should not create a client", func() {
				factory, err := Config([]string{})
				Expect(err).To(HaveOccurred(), "Must return an error")
				// TODO: Match error
				Expect(factory).To(BeNil(), "Must not return an HydraClientFactory")
			})
		})
	})

	Describe("Build", func() {
		It("should add a hydra service cache monitor with default timeout and run it", func() {
			mockClientMaker.EXPECT().MakeClient(gomock.Eq(test_hydra_servers)).Return(mockHydraClient)
			mockHydraClient.EXPECT().SetMaxNumberOfRetries(gomock.Any()).Times(1)
			mockHydraClient.EXPECT().SetWaitBetweenAllServersRetry(gomock.Any()).Times(1)
			c1 := mockHydraClient.EXPECT().ReloadHydraServiceCache().Times(0)
			mockHydraClient.EXPECT().ReloadHydraServiceCache().AnyTimes().After(c1)
			c2 := mockHydraClient.EXPECT().ReloadServicesCache().Times(0)
			mockHydraClient.EXPECT().ReloadServicesCache().AnyTimes().After(c2)

			mockHydraMonitorMaker.EXPECT().MakeHydraMonitor(gomock.Any(), gomock.Eq(time.Duration(60)*time.Second)).
				Return(mockHydraMonitor)
			mockHydraMonitor.EXPECT().Run()

			factory, _ := Config(test_hydra_servers)
			factory.ClientInstantiator = mockClientMaker
			factory.HydraMonitorInstantiator = mockHydraMonitorMaker
			_ = factory.Build()
		})
		Context("when hydra cache refresh time is configured with", func() {
			It("should add a Hydra service cache monitor with default timeout and run it", func() {
				mockClientMaker.EXPECT().MakeClient(gomock.Eq(test_hydra_servers)).Return(mockHydraClient)
				mockHydraClient.EXPECT().SetMaxNumberOfRetries(gomock.Any()).Times(1)
				mockHydraClient.EXPECT().SetWaitBetweenAllServersRetry(gomock.Any()).Times(1)
				c1 := mockHydraClient.EXPECT().ReloadHydraServiceCache().Times(0)
				mockHydraClient.EXPECT().ReloadHydraServiceCache().AnyTimes().After(c1)
				c2 := mockHydraClient.EXPECT().ReloadServicesCache().Times(0)
				mockHydraClient.EXPECT().ReloadServicesCache().AnyTimes().After(c2)

				const timeout int = 10
				mockHydraMonitorMaker.EXPECT().MakeHydraMonitor(gomock.Any(), gomock.Eq(time.Duration(timeout)*time.Second)).
					Return(mockHydraMonitor)
				mockHydraMonitor.EXPECT().Run()

				factory, _ := Config(test_hydra_servers)
				factory.ClientInstantiator = mockClientMaker
				factory.HydraMonitorInstantiator = mockHydraMonitorMaker
				_ = factory.WithHydraCacheRefreshTime(timeout).Build()
			})
		})
		Context("when hydra cache refresh time is configured and", func() {
			It("should add a Hydra service cache monitor with default timeout and run it", func() {
				mockClientMaker.EXPECT().MakeClient(gomock.Eq(test_hydra_servers)).Return(mockHydraClient)
				mockHydraClient.EXPECT().SetMaxNumberOfRetries(gomock.Any()).Times(1)
				mockHydraClient.EXPECT().SetWaitBetweenAllServersRetry(gomock.Any()).Times(1)
				c1 := mockHydraClient.EXPECT().ReloadHydraServiceCache().Times(0)
				mockHydraClient.EXPECT().ReloadHydraServiceCache().AnyTimes().After(c1)
				c2 := mockHydraClient.EXPECT().ReloadServicesCache().Times(0)
				mockHydraClient.EXPECT().ReloadServicesCache().AnyTimes().After(c2)

				const timeout int = 10
				mockHydraMonitorMaker.EXPECT().MakeHydraMonitor(gomock.Any(), gomock.Eq(time.Duration(timeout)*time.Second)).
					Return(mockHydraMonitor)
				mockHydraMonitor.EXPECT().Run()

				factory, _ := Config(test_hydra_servers)
				factory.ClientInstantiator = mockClientMaker
				factory.HydraMonitorInstantiator = mockHydraMonitorMaker
				_ = factory.AndHydraRefreshTime(timeout).Build()
			})
		})
		Context("when hydra cache refresh time is disabled with", func() {
			It("should not add a Hydra service cache monitor", func() {
				mockClientMaker.EXPECT().MakeClient(gomock.Eq(test_hydra_servers)).Return(mockHydraClient)
				mockHydraClient.EXPECT().SetMaxNumberOfRetries(gomock.Any()).Times(1)
				mockHydraClient.EXPECT().SetWaitBetweenAllServersRetry(gomock.Any()).Times(1)
				c1 := mockHydraClient.EXPECT().ReloadHydraServiceCache().Times(0)
				mockHydraClient.EXPECT().ReloadHydraServiceCache().AnyTimes().After(c1)
				c2 := mockHydraClient.EXPECT().ReloadServicesCache().Times(0)
				mockHydraClient.EXPECT().ReloadServicesCache().AnyTimes().After(c2)

				mockHydraMonitorMaker.EXPECT().MakeHydraMonitor(gomock.Any(), gomock.Any()).
					Times(0)

				factory, _ := Config(test_hydra_servers)
				factory.ClientInstantiator = mockClientMaker
				factory.HydraMonitorInstantiator = mockHydraMonitorMaker
				_ = factory.WithoutHydraServerRefresh().Build()
			})
		})
		Context("when hydra cache refresh time is disabled and", func() {
			It("should not add a Hydra service cache monitor", func() {
				mockClientMaker.EXPECT().MakeClient(gomock.Eq(test_hydra_servers)).Return(mockHydraClient)
				mockHydraClient.EXPECT().SetMaxNumberOfRetries(gomock.Any()).Times(1)
				mockHydraClient.EXPECT().SetWaitBetweenAllServersRetry(gomock.Any()).Times(1)
				c1 := mockHydraClient.EXPECT().ReloadHydraServiceCache().Times(0)
				mockHydraClient.EXPECT().ReloadHydraServiceCache().AnyTimes().After(c1)
				c2 := mockHydraClient.EXPECT().ReloadServicesCache().Times(0)
				mockHydraClient.EXPECT().ReloadServicesCache().AnyTimes().After(c2)

				mockHydraMonitorMaker.EXPECT().MakeHydraMonitor(gomock.Any(), gomock.Any()).
					Times(0)

				factory, _ := Config(test_hydra_servers)
				factory.ClientInstantiator = mockClientMaker
				factory.HydraMonitorInstantiator = mockHydraMonitorMaker
				_ = factory.AndWithoutHydraServerRefresh().Build()
			})
		})

		It("should add an apps service cache monitor with default timeout and run it", func() {
			mockClientMaker.EXPECT().MakeClient(gomock.Eq(test_hydra_servers)).Return(mockHydraClient)
			mockHydraClient.EXPECT().SetMaxNumberOfRetries(gomock.Any()).Times(1)
			mockHydraClient.EXPECT().SetWaitBetweenAllServersRetry(gomock.Any()).Times(1)
			c1 := mockHydraClient.EXPECT().ReloadHydraServiceCache().Times(0)
			mockHydraClient.EXPECT().ReloadHydraServiceCache().AnyTimes().After(c1)
			c2 := mockHydraClient.EXPECT().ReloadServicesCache().Times(0)
			mockHydraClient.EXPECT().ReloadServicesCache().AnyTimes().After(c2)

			mockAppsMonitorMaker.EXPECT().MakeAppsMonitor(gomock.Any(), gomock.Eq(time.Duration(20)*time.Second)).
				Return(mockServicesMonitor)
			mockServicesMonitor.EXPECT().Run()

			factory, _ := Config(test_hydra_servers)
			factory.ClientInstantiator = mockClientMaker
			factory.AppsMonitorInstantiator = mockAppsMonitorMaker
			_ = factory.Build()
		})
		Context("when apps cache refresh time is configured with", func() {
			It("should add an apps service cache monitor with default timeout and run it", func() {
				mockClientMaker.EXPECT().MakeClient(gomock.Eq(test_hydra_servers)).Return(mockHydraClient)
				mockHydraClient.EXPECT().SetMaxNumberOfRetries(gomock.Any()).Times(1)
				mockHydraClient.EXPECT().SetWaitBetweenAllServersRetry(gomock.Any()).Times(1)
				c1 := mockHydraClient.EXPECT().ReloadHydraServiceCache().Times(0)
				mockHydraClient.EXPECT().ReloadHydraServiceCache().AnyTimes().After(c1)
				c2 := mockHydraClient.EXPECT().ReloadServicesCache().Times(0)
				mockHydraClient.EXPECT().ReloadServicesCache().AnyTimes().After(c2)

				const timeout int = 90
				mockAppsMonitorMaker.EXPECT().MakeAppsMonitor(gomock.Any(), gomock.Eq(time.Duration(timeout)*time.Second)).
					Return(mockServicesMonitor)
				mockServicesMonitor.EXPECT().Run()

				factory, _ := Config(test_hydra_servers)
				factory.ClientInstantiator = mockClientMaker
				factory.AppsMonitorInstantiator = mockAppsMonitorMaker
				_ = factory.WithAppsCacheRefreshTime(timeout).Build()
			})
		})
		Context("when apps cache refresh time is configured and", func() {
			It("should add an apps service cache monitor with default timeout and run it", func() {
				mockClientMaker.EXPECT().MakeClient(gomock.Eq(test_hydra_servers)).Return(mockHydraClient)
				mockHydraClient.EXPECT().SetMaxNumberOfRetries(gomock.Any()).Times(1)
				mockHydraClient.EXPECT().SetWaitBetweenAllServersRetry(gomock.Any()).Times(1)
				c1 := mockHydraClient.EXPECT().ReloadHydraServiceCache().Times(0)
				mockHydraClient.EXPECT().ReloadHydraServiceCache().AnyTimes().After(c1)
				c2 := mockHydraClient.EXPECT().ReloadServicesCache().Times(0)
				mockHydraClient.EXPECT().ReloadServicesCache().AnyTimes().After(c2)

				const timeout int = 90
				mockAppsMonitorMaker.EXPECT().MakeAppsMonitor(gomock.Any(), gomock.Eq(time.Duration(timeout)*time.Second)).
					Return(mockServicesMonitor)
				mockServicesMonitor.EXPECT().Run()

				factory, _ := Config(test_hydra_servers)
				factory.ClientInstantiator = mockClientMaker
				factory.AppsMonitorInstantiator = mockAppsMonitorMaker
				_ = factory.AndAppsCacheRefreshTime(timeout).Build()
			})
		})
		Context("when apps cache refresh time is disabled with", func() {
			It("should not add an apps service cache monitor", func() {
				mockClientMaker.EXPECT().MakeClient(gomock.Eq(test_hydra_servers)).Return(mockHydraClient)
				mockHydraClient.EXPECT().SetMaxNumberOfRetries(gomock.Any()).Times(1)
				mockHydraClient.EXPECT().SetWaitBetweenAllServersRetry(gomock.Any()).Times(1)
				c1 := mockHydraClient.EXPECT().ReloadHydraServiceCache().Times(0)
				mockHydraClient.EXPECT().ReloadHydraServiceCache().AnyTimes().After(c1)
				c2 := mockHydraClient.EXPECT().ReloadServicesCache().Times(0)
				mockHydraClient.EXPECT().ReloadServicesCache().AnyTimes().After(c2)

				mockAppsMonitorMaker.EXPECT().MakeAppsMonitor(gomock.Any(), gomock.Any()).
					Times(0)

				factory, _ := Config(test_hydra_servers)
				factory.ClientInstantiator = mockClientMaker
				factory.AppsMonitorInstantiator = mockAppsMonitorMaker
				_ = factory.WithoutAppsRefresh().Build()
			})
		})
		Context("when apps cache refresh time is disabled with", func() {
			It("should not add an apps service cache monitor", func() {
				mockClientMaker.EXPECT().MakeClient(gomock.Eq(test_hydra_servers)).Return(mockHydraClient)
				mockHydraClient.EXPECT().SetMaxNumberOfRetries(gomock.Any()).Times(1)
				mockHydraClient.EXPECT().SetWaitBetweenAllServersRetry(gomock.Any()).Times(1)
				c1 := mockHydraClient.EXPECT().ReloadHydraServiceCache().Times(0)
				mockHydraClient.EXPECT().ReloadHydraServiceCache().AnyTimes().After(c1)
				c2 := mockHydraClient.EXPECT().ReloadServicesCache().Times(0)
				mockHydraClient.EXPECT().ReloadServicesCache().AnyTimes().After(c2)

				mockAppsMonitorMaker.EXPECT().MakeAppsMonitor(gomock.Any(), gomock.Any()).
					Times(0)

				factory, _ := Config(test_hydra_servers)
				factory.ClientInstantiator = mockClientMaker
				factory.AppsMonitorInstantiator = mockAppsMonitorMaker
				_ = factory.AndWithoutAppsRefresh().Build()
			})
		})

		Context("when number of retries are configured with", func() {
			It("should set the number of retries", func() {
				const numberOfRetries int = 90
				mockClientMaker.EXPECT().MakeClient(gomock.Eq(test_hydra_servers)).Return(mockHydraClient)
				mockHydraClient.EXPECT().SetMaxNumberOfRetries(gomock.Eq(numberOfRetries)).Times(1)
				mockHydraClient.EXPECT().SetWaitBetweenAllServersRetry(gomock.Any()).Times(1)
				c1 := mockHydraClient.EXPECT().ReloadHydraServiceCache().Times(0)
				mockHydraClient.EXPECT().ReloadHydraServiceCache().AnyTimes().After(c1)
				c2 := mockHydraClient.EXPECT().ReloadServicesCache().Times(0)
				mockHydraClient.EXPECT().ReloadServicesCache().AnyTimes().After(c2)

				factory, _ := Config(test_hydra_servers)
				factory.ClientInstantiator = mockClientMaker
				_ = factory.WithNumberOfRetries(numberOfRetries).Build()
			})
		})
		Context("when number of retries are configured and", func() {
			It("should set the number of retries", func() {
				const numberOfRetries int = 90
				mockClientMaker.EXPECT().MakeClient(gomock.Eq(test_hydra_servers)).Return(mockHydraClient)
				mockHydraClient.EXPECT().SetMaxNumberOfRetries(gomock.Eq(numberOfRetries)).Times(1)
				mockHydraClient.EXPECT().SetWaitBetweenAllServersRetry(gomock.Any()).Times(1)
				c1 := mockHydraClient.EXPECT().ReloadHydraServiceCache().Times(0)
				mockHydraClient.EXPECT().ReloadHydraServiceCache().AnyTimes().After(c1)
				c2 := mockHydraClient.EXPECT().ReloadServicesCache().Times(0)
				mockHydraClient.EXPECT().ReloadServicesCache().AnyTimes().After(c2)

				factory, _ := Config(test_hydra_servers)
				factory.ClientInstantiator = mockClientMaker
				_ = factory.AndNumberOfRetries(numberOfRetries).Build()
			})
		})

		Context("when milliseconds to retry are configured", func() {
			It("should set wait between all servers retry", func() {
				const millisecondsToRetry int = 30
				mockClientMaker.EXPECT().MakeClient(gomock.Eq(test_hydra_servers)).Return(mockHydraClient)
				mockHydraClient.EXPECT().SetMaxNumberOfRetries(gomock.Any()).Times(1)
				mockHydraClient.EXPECT().SetWaitBetweenAllServersRetry(gomock.Eq(millisecondsToRetry)).Times(1)
				c1 := mockHydraClient.EXPECT().ReloadHydraServiceCache().Times(0)
				mockHydraClient.EXPECT().ReloadHydraServiceCache().AnyTimes().After(c1)
				c2 := mockHydraClient.EXPECT().ReloadServicesCache().Times(0)
				mockHydraClient.EXPECT().ReloadServicesCache().AnyTimes().After(c2)

				factory, _ := Config(test_hydra_servers)
				factory.ClientInstantiator = mockClientMaker
				_ = factory.WaitBetweenAllServersRetry(millisecondsToRetry).Build()
			})
		})
	})
})
