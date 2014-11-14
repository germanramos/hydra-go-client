package client_test

import (
	. "github.com/innotech/hydra-go-client/client"
	mock "github.com/innotech/hydra-go-client/client/mock"

	"github.com/innotech/hydra-go-client/vendors/code.google.com/p/gomock/gomock"
	. "github.com/innotech/hydra-go-client/vendors/github.com/onsi/ginkgo"
	. "github.com/innotech/hydra-go-client/vendors/github.com/onsi/gomega"

	// "time"
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

		test_hydra_servers []string = []string{seed_server}
	)

	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
		mockAppsMonitorMaker = mock.NewMockappsMonitorMaker(mockCtrl)
		mockClientMaker = mock.NewMockclientMaker(mockCtrl)
		mockHydraMonitorMaker = mock.NewMockhydraMonitorMaker(mockCtrl)
		mockHydraClient = mock.NewMockClient(mockCtrl)
		// TODO: hydraClientFactoryTimersFixture
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

	// Describe("Config", func() {
	// 	Context("when none seed servers is passed", func() {
	// 		It("should not create a client", func() {
	// 			client, err := Config([]string{}).Build()
	// 			Expect(err).To(HaveOccurred(), "Must return an error")
	// 			// TODO: Match error
	// 			Expect(client).To(BeNil(), "Must not return an Hydra client")
	// 		})
	// 	})
	// })

	Describe("Build", func() {
		// Covert up
		// It("should reload the hydra service cache", func() {
		// 	mockClientMaker.EXPECT().MakeClient(gomock.Eq(test_hydra_servers)).Return(mockHydraClient)
		// 	mockHydraClient.EXPECT().ReloadHydraServiceCache()

		// 	factory := Config(test_hydra_servers)
		// 	factory.ClientInstantiator = mockClientMaker
		// 	_ = factory.Build()
		// })
		// 	// With default timeout
		// 	It("should add a Hydra service cache monitor with default timeout and run it", func() {
		// 		// mockClientMaker.EXPECT().MakeClient(gomock.Eq(test_hydra_servers)).Return(mockHydraClient)

		// 		factory := Config(test_hydra_servers)
		// 		_ = factory.Build()
		// 	})
		// 	// It("should add a Hydra service cache monitor and run it", func() {
		// 	// 	mockClientMaker.EXPECT().MakeClient(gomock.Eq(test_hydra_servers)).Return(mockHydraClient)
		// 	// 	mockHydraClient.EXPECT().ReloadHydraServiceCache()

		// 	// 	factory := Config(test_hydra_servers)
		// 	// 	factory.ClientInstantiator = mockClientMaker
		// 	// 	_ = factory.Build()
		// 	// })
		// 	// Context("when the default configuration is not overwritten", func() {
		// 	// 	It("should set the hydra cache refresh time", func() {
		// 	// 		_ = Config(test_hydra_servers).Build()
		// 	// 	})
		// 	// })
		// 	// Context("text", body)
	})

	// Describe("SetWaitBetweenAllServersRetry", func() {
	// 	_ = Config(test_hydra_servers).WaitBetweenAllServerRetry(30).Build()
	// })
})

// 	// TODO: add descriptions in asserts

// 	// Describe("GetFactoryWithHydraServers", func() {
// 	// 	Context("when it is is called repeatedly", func() {
// 	// 		It("should return an unique hydra client", func() {
// 	// 			hydraClient := GetFactoryWithHydraServers(test_hydra_servers).Build()
// 	// 			anotherHydraClient := GetFactoryWithHydraServers(test_hydra_servers).Build()

// 	// 			Expect(hydraClient).ToNot(BeNil())
// 	// 			Expect(anotherHydraClietn).ToNot(BeNil())
// 	// 			Expect(hydraClient).To(Equal(anotherHydraClient))
// 	// 		})
// 	// 	})
// 	// })

// 	// Describe("GetHydraClient", func() {
// 	// 	It("should return an unique hydra client", func() {
// 	// 		hydraClient := GetFactoryWithHydraServers(test_hydra_servers).Build()
// 	// 		anotherHydraClient := GetHydraClient()

// 	// 		Expect(hydraClient).ToNot(BeNil())
// 	// 		Expect(anotherHydraClietn).ToNot(BeNil())
// 	// 		Expect(hydraClient).To(Equal(anotherHydraClient))
// 	// 	})
// 	// })

// 	// Describe("GetFactoryWithHydraServers", func() {
// 	// 	Context("when it is is called repeatedly", func() {
// 	// 		It("should return an unique hydra client", func() {
// 	// 			hydraClient := GetFactoryWithHydraServers(test_hydra_servers).Build()
// 	// 			anotherHydraClient := GetFactoryWithHydraServers(test_hydra_servers).Build()

// 	// 			Expect(hydraClient).ToNot(BeNil())
// 	// 			Expect(anotherHydraClietn).ToNot(BeNil())
// 	// 			Expect(hydraClient).To(Equal(anotherHydraClient))
// 	// 		})
// 	// 	})
// 	// })

// 	//////////////////////////////////////////////////////////////////////

// 	It("should be instantiated with default configuration", func() {
// 		Expect(HydraClientFactory.GetAppsCacheDuration()).To(Equal(DefaultAppsCacheDuration))
// 		Expect(HydraClientFactory.GetDurationBetweenAllServersRetry()).To(Equal(DefaultDurationBetweenAllServersRetry))
// 		Expect(HydraClientFactory.GetHydraServersCacheDuration()).To(Equal(DefaultHydraServersCacheDuration))
// 		Expect(HydraClientFactory.GetMaxNumberOfRetriesPerHydraServer()).To(Equal(DefaultNumberOfRetries))
// 	})

// 	Describe("Build", func() {
// 		It("should build a HydraClient", func() {
// 			var client *Client
// 			Expect(HydraClientFactory.Build()).To(BeAssignableToTypeOf(client))
// 		})
// 	})

// 	Describe("Config", func() {
// 		Context("when hydra server list argument is nil", func() {
// 			It("should throw an error", func() {
// 				err := HydraClientFactory.Config(nil)
// 				Expect(err).Should(HaveOccurred())
// 			})
// 		})
// 		Context("when hydra server list argument is a empty list", func() {
// 			It("should throw an error", func() {
// 				err := HydraClientFactory.Config([]string{})
// 				Expect(err).Should(HaveOccurred())
// 			})
// 		})
// 		Context("when hydra server list argument is a valid list of servers", func() {
// 			It("should set hydra server list successfully", func() {
// 				err := HydraClientFactory.Config([]string{"http://localhost:8080"})
// 				Expect(err).ShouldNot(HaveOccurred())
// 			})
// 		})
// 	})

// 	Describe("WithHydraServersCacheDuration", func() {
// 		Context("when duration argument is a valid uint number", func() {
// 			It("should set hydra servers cache duration successfully", func() {
// 				const hydraServersCacheDuration time.Duration = time.Duration(30000) * time.Millisecond
// 				h := HydraClientFactory.WithHydraServersCacheDuration(hydraServersCacheDuration)
// 				Expect(h).To(Equal(HydraClientFactory))
// 				Expect(hydraServersCacheDuration).To(Equal(HydraClientFactory.GetHydraServersCacheDuration()))
// 			})
// 		})
// 	})

// 	Describe("WithAppsCacheDuration", func() {
// 		Context("when duration argument is a valid uint number", func() {
// 			It("should set apps cache duration successfully", func() {
// 				const appsCacheDuration time.Duration = time.Duration(30000) * time.Millisecond
// 				h := HydraClientFactory.WithAppsCacheDuration(appsCacheDuration)
// 				Expect(h).To(Equal(HydraClientFactory))
// 				Expect(appsCacheDuration).To(Equal(HydraClientFactory.GetAppsCacheDuration()))
// 			})
// 		})
// 	})

// 	Describe("WithMaxNumberOfRetriesPerHydraServer", func() {
// 		Context("when duration argument is a valid uint number", func() {
// 			It("should set apps cache duration successfully", func() {
// 				const retries uint = 3
// 				h := HydraClientFactory.WithMaxNumberOfRetriesPerHydraServer(retries)
// 				Expect(h).To(Equal(HydraClientFactory))
// 				Expect(retries).To(Equal(HydraClientFactory.GetMaxNumberOfRetriesPerHydraServer()))
// 			})
// 		})
// 	})

// 	Describe("WaitBetweenAllServersRetry", func() {
// 		Context("when duration argument is a valid uint number", func() {
// 			It("should set wait between all servers retry successfully", func() {
// 				const appsCacheDuration time.Duration = time.Duration(30000) * time.Millisecond
// 				h := HydraClientFactory.WaitBetweenAllServersRetry(appsCacheDuration)
// 				Expect(h).To(Equal(HydraClientFactory))
// 				Expect(appsCacheDuration).To(Equal(HydraClientFactory.GetDurationBetweenAllServersRetry()))
// 			})
// 		})
// 	})
// })
