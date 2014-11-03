package client_test

import (
	. "github.com/innotech/hydra-go-client/client"
	. "github.com/innotech/hydra-go-client/client/error"
	mock "github.com/innotech/hydra-go-client/client/mock"

	"github.com/innotech/hydra-go-client/vendors/code.google.com/p/gomock/gomock"
	. "github.com/innotech/hydra-go-client/vendors/github.com/onsi/ginkgo"
	. "github.com/innotech/hydra-go-client/vendors/github.com/onsi/gomega"

	"errors"
)

var _ = Describe("HydraClient", func() {
	const (
		hydra string = "hydra"
		// connection_timeout = 1000
		test_hydra_server_url         string = "http://localhost:8080"
		another_test_hydra_server_url string = "http://localhost:8081"
		test_app_server               string = "http://localhost:8080/app-server-first"
		another_test_app_server       string = "http://localhost:8081/app-server-second"
		service_id                    string = "testAppId"
	)

	var (
		test_hydra_servers []string = []string{test_hydra_server_url, another_test_hydra_server_url}
		test_services      []string = []string{test_app_server, another_test_app_server}

		hydraClient           *Client
		mockCtrl              *gomock.Controller
		mockHydraServiceCache *mock.MockHydraServiceCache
		mockServiceCache      *mock.MockServiceCache
		mockServiceRepository *mock.MockServiceRepository
	)

	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
		mockHydraServiceCache = mock.NewMockHydraServiceCache(mockCtrl)
		mockServiceCache = mock.NewMockServiceCache(mockCtrl)
		mockServiceRepository = mock.NewMockServiceRepository(mockCtrl)
		hydraClient = NewHydraClient(test_hydra_servers)
	})

	AfterEach(func() {
		mockCtrl.Finish()
	})

	Describe("Get", func() {
		Context("when cache is not invoked", func() {
			Context("when no services are cached", func() {
				It("should return the list of balanced services from Hydra", func() {
					c1 := mockServiceCache.EXPECT().Exists(gomock.Eq(service_id)).
						Return(false)
					c2 := mockHydraServiceCache.EXPECT().GetHydraServers().
						Return(test_hydra_servers).After(c1)
					c3 := mockServiceRepository.EXPECT().FindById(gomock.Eq(services_id), gomock.Eq(test_hydra_servers)).
						Return(test_services).After(c2)

					candidateServers, err := hydraClient.Get(service_id, false)

					mockServiceCache.EXPECT().PutService(gomock.Eq(services_id), gomock.Eq(candidateServers)).
						After(c3)
					Expect(err).ToNot(HaveOccurred(), "Must not return an error")
					Expect(candidateServers).ToNot(BeEmpty(), "Must not return an empty list of servers")
					Expect(candidateServers).ToNot(Equal(test_services), "Must return teh expected list of servers")
				})
			})
		})
		Context("when cache is not invoked", func() {
			Context("when no services are cached", func() {
				It("should call shortcutting the cache", func() {
					c1 := mockServiceCache.EXPECT().Exists(gomock.Eq(service_id)).
						Return(false)
					c2 := mockHydraServiceCache.EXPECT().GetHydraServers().
						Return(test_hydra_servers).Times(1).After(c1)
					c3 := mockServiceRepository.EXPECT().FindById(gomock.Eq(services_id), gomock.Eq(test_hydra_servers)).
						Return(test_services).Times(1).After(c2)
					c4 := mockServiceCache.EXPECT().PutService(gomock.Eq(services_id), gomock.Eq(test_services)).
						Times(1).After(c3)
					mockServiceCache.EXPECT().FindById(gomock.Eq(service_id)).Times(1).Return(test_services)

					_, err := hydraClient.Get(service_id, true)
					Expect(err).ToNot(HaveOccurred(), "Must not return an error")
					// Call twice to ensure that the second call hit the cache.
					candidateServers, err := hydraClient.Get(service_id, true)

					Expect(err).ToNot(HaveOccurred(), "Must not return an error")
					Expect(candidateServers).ToNot(BeEmpty(), "Must not return an empty list of servers")
					Expect(candidateServers).ToNot(Equal(test_services), "Must return teh expected list of servers")
				})
			})
		})
		It("should not accept an empty service id", func() {
			// TODO: If crash Finish mockController
			servers, err := hydraClient.Get("", true)

			Expect(err).To(HaveOccurred(), "Must return an error")
			Expect(servers).To(BeEmpty(), "Must return an empty list of servers")
		})
		It("should not accept an empty service id", func() {
			// TODO: If crash Finish mockController
			servers, err := hydraClient.Get("      ", true)

			Expect(err).To(HaveOccurred(), "Must return an error")
			Expect(servers).To(BeEmpty(), "Must return an empty list of servers")
		})
	})

	// Async test

	Describe("ReloadHydraServiceCache", func() {
		It("should reload Hydra servers", func() {
			c1 := mockHydraServiceCache.EXPECT().GetHydraServers().
				Return(test_hydra_servers)
			c2 := mockServiceRepository.EXPECT().FindById(gomock.Eq(services_id), gomock.Eq(test_hydra_servers)).
				Return(test_hydra_servers).After(c1)
			mockHydraServiceCache.EXPECT().Refresh(test_hydra_servers).After(c2)

			hydraClient.ReloadHydraServiceCache()
		})
		Context("when Hydra is not available", func() {
			It("should set Hydra as not available", func() {
				c1 := mockHydraServiceCache.EXPECT().GetHydraServers().
					Return(test_hydra_servers)
				c2 := mockServiceRepository.EXPECT().FindById(gomock.Eq(services_id), gomock.Eq(test_hydra_servers)).
					Return([]string{}, HydraNotAvailableError).After(c1)
				mockHydraServiceCache.EXPECT().Refresh(gomock.Eq(test_hydra_servers)).Times(0)

				hydraClient.ReloadHydraServiceCache()

				Expect(hydraClient.IsHydraAvailable()).To(BeFalse(), "Hydra is not available")
			})
		})
		Context("when no Hydra servers", func() {
			It("should reload Hydra servers and set Hydra as not available", func() {
				c1 := mockHydraServiceCache.EXPECT().GetHydraServers().
					Return(test_hydra_servers)
				c2 := mockServiceRepository.EXPECT().FindById(gomock.Eq(services_id), gomock.Eq(test_hydra_servers)).
					Return([]string{}, nil).After(c1)
				mockHydraServiceCache.EXPECT().Refresh(gomock.Eq([]string{}))

				hydraClient.ReloadHydraServiceCache()

				Expect(hydraClient.IsHydraAvailable()).To(BeFalse(), "Hydra is not available")
			})
		})
	})

	// Describe("InitHydraService", func() {
	// 	It("should init Hydra servers", func() {
	// 		hydraClient.InitHydraService()

	// 	})
	// })

	// Describe("SetConnectionTimeout", func() {

	// })

	Describe("ReloadServicesCache", func() {
		It("should reload the service cache", func() {
			appIds := []string{service_id}
			services := map[string][]string{
				service_id: test_services,
			}

			c1 := mockServicesCache.EXPECT().GetIds().Return(appIds)
			c2 := mockHydraServiceCache.EXPECT().GetHydraServers().
				Return(test_hydra_servers).After(c1)
			c3 := mockServiceRepository.EXPECT().FindByIds(gomock.Eq(appIds), gomock.Eq(test_hydra_servers)).
				Return(services, nil).After(c2)
			mockServicesCache.EXPECT().Refresh(gomock.Eq(services))

			hydraClient.SetHydraAvailable(true)
			hydraClient.ReloadServicesCache()
		})
		// TODO
		Context("when Hydra is not available", func() {
			It("should not reload the service cache", func() {
				appIds := []string{service_id}
				services := map[string][]string{
					service_id: test_services,
				}

				c1 := mockServicesCache.EXPECT().GetIds().Return(appIds)
				c2 := mockHydraServiceCache.EXPECT().GetHydraServers().
					Return(test_hydra_servers).After(c1)
				c3 := mockServiceRepository.EXPECT().FindByIds(gomock.Eq(appIds), gomock.Eq(test_hydra_servers)).
					Return(services, nil).After(c2)
				mockServicesCache.EXPECT().Refresh(gomock.Eq(services)).Times(0)

				hydraClient.SetHydraAvailable(true)
				hydraClient.ReloadServicesCache()
			})
		})
	})
})

///////////////////////////////////////////////////////////////

// import (
// 	. "github.com/innotech/hydra-go-client/client"
// 	mock "github.com/innotech/hydra-go-client/client/mock"

// 	"github.com/innotech/hydra-go-client/vendors/code.google.com/p/gomock/gomock"
// 	. "github.com/innotech/hydra-go-client/vendors/github.com/onsi/ginkgo"
// 	. "github.com/innotech/hydra-go-client/vendors/github.com/onsi/gomega"

// 	"errors"
// )

// var _ = Describe("HydraClient", func() {
// 	var (
// 		hydraClient   *Client
// 		mockCtrl      *gomock.Controller
// 		mockRequester *mock.MockRequester
// 	)

// 	BeforeEach(func() {
// 		mockCtrl = gomock.NewController(GinkgoT())
// 		mockRequester = mock.NewMockRequester(mockCtrl)
// 		hydraClient = NewClient([]string{"http://localhost:8080"}, mockRequester)
// 	})

// 	AfterEach(func() {
// 		mockCtrl.Finish()
// 	})

// 	Describe("Get", func() {
// 		Context("when an illegal application ID is passed as an argument", func() {
// 			It("should throw an error", func() {
// 				servers, err := hydraClient.Get("", false)
// 				Expect(servers).To(BeEmpty())
// 				Expect(err).To(HaveOccurred())
// 			})
// 		})
// 		Context("when the cache should not be refreshed", func() {
// 			Context("when the application ID doesn't exist", func() {
// 				It("should request servers from hydra server", func() {
// 					mockRequester.EXPECT().GetCandidateServers(gomock.Any(), gomock.Eq("app1"))
// 					_, _ = hydraClient.Get("app1", false)
// 				})
// 			})
// 			Context("when the application ID exists", func() {
// 				It("should not request servers from hydra server", func() {
// 					mockRequester.EXPECT().GetCandidateServers(gomock.Any(), gomock.Eq("app1"))
// 					_, _ = hydraClient.Get("app1", false)

// 					mockRequester.EXPECT().GetCandidateServers(gomock.Any(), gomock.Any()).Times(0)
// 					_, _ = hydraClient.Get("app1", false)
// 				})
// 			})
// 		})
// 		Context("when the cache should be refreshed", func() {
// 			Context("when the application ID doesn't exist", func() {
// 				It("should request servers from hydra server", func() {
// 					mockRequester.EXPECT().GetCandidateServers(gomock.Any(), gomock.Eq("app1"))
// 					_, _ = hydraClient.Get("app1", true)
// 				})
// 			})
// 			Context("when the application ID exists", func() {
// 				It("should request servers from hydra server", func() {
// 					mockRequester.EXPECT().GetCandidateServers(gomock.Any(), gomock.Eq("app1"))
// 					_, _ = hydraClient.Get("app1", false)

// 					mockRequester.EXPECT().GetCandidateServers(gomock.Any(), gomock.Eq("app1"))
// 					_, _ = hydraClient.Get("app1", true)
// 				})
// 			})
// 		})
// 	})

// 	Describe("ReloadHydraServers", func() {
// 		Context("when hydra server is not accessible", func() {
// 			It("should consider that Hydra is not available", func() {
// 				mockRequester.EXPECT().GetCandidateServers(gomock.Any(), gomock.Eq("hydra")).Return([]string{}, errors.New("Not Found"))
// 				hydraClient.SetMaxNumberOfRetriesPerHydraServer(1)
// 				hydraClient.ReloadHydraServers()
// 				Expect(hydraClient.IsHydraAvailable()).To(BeFalse())
// 			})
// 		})
// 		Context("when hydra server is not accessible", func() {
// 			Context("when hydra server responses with an empty list of servers", func() {
// 				It("should consider that Hydra is not available", func() {
// 					mockRequester.EXPECT().GetCandidateServers(gomock.Any(), gomock.Eq("hydra")).Return([]string{}, nil)
// 					hydraClient.SetMaxNumberOfRetriesPerHydraServer(1)
// 					hydraClient.ReloadHydraServers()
// 					Expect(hydraClient.IsHydraAvailable()).To(BeFalse())
// 				})
// 			})
// 			Context("when hydra server responses with a list of servers", func() {
// 				It("should consider that Hydra is available", func() {
// 					mockRequester.EXPECT().GetCandidateServers(gomock.Any(), gomock.Eq("hydra")).Return([]string{"http://localhost:8081"}, nil)
// 					hydraClient.SetMaxNumberOfRetriesPerHydraServer(1)
// 					hydraClient.ReloadHydraServers()
// 					Expect(hydraClient.IsHydraAvailable()).To(BeTrue())
// 				})
// 				It("should update hydra cache", func() {
// 					mockRequester.EXPECT().GetCandidateServers(gomock.Any(), gomock.Eq("hydra")).Return([]string{"http://localhost:8082"}, nil)
// 					hydraClient.SetMaxNumberOfRetriesPerHydraServer(1)
// 					hydraClient.ReloadHydraServers()
// 					mockRequester.EXPECT().GetCandidateServers(gomock.Eq("http://localhost:8082"+AppRootPath), gomock.Eq("app1"))
// 					_, _ = hydraClient.Get("app1", false)
// 				})
// 			})
// 		})
// 	})

// 	Describe("ReloadAppServers", func() {
// 		Context("when no applications registered", func() {
// 			It("should not send any request to hydra servers", func() {
// 				mockRequester.EXPECT().GetCandidateServers(gomock.Any(), gomock.Any()).Times(0)
// 				hydraClient.ReloadAppServers()
// 			})
// 		})
// 		Context("when one application is registered", func() {
// 			It("should require update the application cache", func() {
// 				mockRequester.EXPECT().GetCandidateServers(gomock.Any(), gomock.Eq("app1")).Return([]string{"http://localhost:8080"}, nil)
// 				_, _ = hydraClient.Get("app1", false)

// 				mockRequester.EXPECT().GetCandidateServers(gomock.Any(), gomock.Eq("app1"))
// 				hydraClient.ReloadAppServers()
// 			})
// 		})
// 		Context("when multiple applications are registered", func() {
// 			It("should require update the application cache", func() {
// 				mockRequester.EXPECT().GetCandidateServers(gomock.Any(), gomock.Eq("app1")).Return([]string{"http://localhost:8080"}, nil)
// 				_, _ = hydraClient.Get("app1", false)
// 				mockRequester.EXPECT().GetCandidateServers(gomock.Any(), gomock.Eq("app2")).Return([]string{"http://localhost:8081"}, nil)
// 				_, _ = hydraClient.Get("app2", false)

// 				mockRequester.EXPECT().GetCandidateServers(gomock.Any(), gomock.Eq("app1"))
// 				mockRequester.EXPECT().GetCandidateServers(gomock.Any(), gomock.Eq("app2"))
// 				hydraClient.ReloadAppServers()
// 			})
// 		})
// 	})
// })
