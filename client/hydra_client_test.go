package client_test

import (
	. "github.com/innotech/hydra-go-client/client"
	. "github.com/innotech/hydra-go-client/client/error"
	mock "github.com/innotech/hydra-go-client/client/mock"

	"github.com/innotech/hydra-go-client/vendors/code.google.com/p/gomock/gomock"
	. "github.com/innotech/hydra-go-client/vendors/github.com/onsi/ginkgo"
	. "github.com/innotech/hydra-go-client/vendors/github.com/onsi/gomega"
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

		hydraClient           *HydraClient
		mockCtrl              *gomock.Controller
		mockHydraServiceCache *mock.MockHydraCache
		mockServiceCache      *mock.MockServiceCache
		mockServiceRepository *mock.MockServiceRepository
	)

	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
		mockHydraServiceCache = mock.NewMockHydraCache(mockCtrl)
		mockServiceCache = mock.NewMockServiceCache(mockCtrl)
		mockServiceRepository = mock.NewMockServiceRepository(mockCtrl)
		hydraClient = NewHydraClient(test_hydra_servers)
		hydraClient.HydraServiceCache = mockHydraServiceCache
		hydraClient.ServicesCache = mockServiceCache
		hydraClient.ServicesRepository = mockServiceRepository
	})

	AfterEach(func() {
		mockCtrl.Finish()
	})

	Describe("Get", func() {
		Context("when no services are cached", func() {
			It("should return the list of balanced services from Hydra", func() {
				c1 := mockServiceCache.EXPECT().Exists(gomock.Eq(service_id)).
					Return(false)
				c2 := mockHydraServiceCache.EXPECT().GetHydraServers().
					Return(test_hydra_servers).After(c1)
				c3 := mockServiceRepository.EXPECT().FindById(gomock.Eq(service_id), gomock.Eq(test_hydra_servers)).
					Return(test_services, nil).After(c2)
				mockServiceCache.EXPECT().PutService(gomock.Eq(service_id), gomock.Eq(test_services)).
					After(c3)

				candidateServers, err := hydraClient.Get(service_id)

				Expect(err).ToNot(HaveOccurred(), "Must not return an error")
				Expect(candidateServers).ToNot(BeEmpty(), "Must not return an empty list of servers")
				Expect(candidateServers).To(Equal(test_services), "Must return the expected list of servers")
			})
		})
		It("should not accept an empty service id", func() {
			servers, err := hydraClient.Get("")

			Expect(err).To(HaveOccurred(), "Must return an error")
			Expect(servers).To(BeEmpty(), "Must return an empty list of servers")
		})
		It("should not accept an empty service id", func() {
			servers, err := hydraClient.Get("      ")

			Expect(err).To(HaveOccurred(), "Must return an error")
			Expect(servers).To(BeEmpty(), "Must return an empty list of servers")
		})
	})
	Describe("GetShortcuttingTheCache", func() {
		It("should call shortcutting the cache", func() {
			c1 := mockHydraServiceCache.EXPECT().GetHydraServers().
				Return(test_hydra_servers)
			c2 := mockServiceRepository.EXPECT().FindById(gomock.Eq(service_id), gomock.Eq(test_hydra_servers)).
				Return(test_services, nil).After(c1)
			c3 := mockServiceCache.EXPECT().PutService(gomock.Eq(service_id), gomock.Eq(test_services)).
				After(c2)
			c4 := mockHydraServiceCache.EXPECT().GetHydraServers().
				Return(test_hydra_servers).After(c3)
			c5 := mockServiceRepository.EXPECT().FindById(gomock.Eq(service_id), gomock.Eq(test_hydra_servers)).
				Return(test_services, nil).After(c4)
			mockServiceCache.EXPECT().PutService(gomock.Eq(service_id), gomock.Eq(test_services)).
				After(c5)

			_, err := hydraClient.GetShortcuttingTheCache(service_id)
			Expect(err).ToNot(HaveOccurred(), "Must not return an error")
			// Call twice to ensure that the second call hit the cache.
			candidateServers, err := hydraClient.GetShortcuttingTheCache(service_id)

			Expect(err).ToNot(HaveOccurred(), "Must not return an error")
			Expect(candidateServers).ToNot(BeEmpty(), "Must not return an empty list of servers")
			Expect(candidateServers).To(Equal(test_services), "Must return the expected list of servers")
		})
		Context("when services are cached", func() {
			It("should response using the cache", func() {
				c1 := mockHydraServiceCache.EXPECT().GetHydraServers().
					Return(test_hydra_servers)
				c2 := mockServiceRepository.EXPECT().FindById(gomock.Eq(service_id), gomock.Eq(test_hydra_servers)).
					Return(test_services, nil).After(c1)
				mockServiceCache.EXPECT().PutService(gomock.Eq(service_id), gomock.Eq(test_services)).
					After(c2)

				candidateServers, err := hydraClient.GetShortcuttingTheCache(service_id)

				Expect(err).ToNot(HaveOccurred(), "Must not return an error")
				Expect(candidateServers).ToNot(BeEmpty(), "Must not return an empty list of servers")
				Expect(candidateServers).To(Equal(test_services), "Must return teh expected list of servers")
			})
		})
	})

	Describe("ReloadHydraServiceCache", func() {
		It("should reload Hydra servers", func() {
			c1 := mockHydraServiceCache.EXPECT().GetHydraServers().
				Return(test_hydra_servers)
			c2 := mockServiceRepository.EXPECT().FindById(gomock.Eq(hydra), gomock.Eq(test_hydra_servers)).
				Return(test_hydra_servers, nil).After(c1)
			mockHydraServiceCache.EXPECT().Refresh(test_hydra_servers).After(c2)

			hydraClient.ReloadHydraServiceCache()
		})
		Context("when Hydra is not available", func() {
			It("should set Hydra as not available", func() {
				c1 := mockHydraServiceCache.EXPECT().GetHydraServers().
					Return(test_hydra_servers)
				mockServiceRepository.EXPECT().FindById(gomock.Eq(hydra), gomock.Eq(test_hydra_servers)).
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
				mockServiceRepository.EXPECT().FindById(gomock.Eq(hydra), gomock.Eq(test_hydra_servers)).
					Return([]string{}, nil).After(c1)
				mockHydraServiceCache.EXPECT().Refresh(gomock.Eq([]string{}))

				hydraClient.ReloadHydraServiceCache()

				Expect(hydraClient.IsHydraAvailable()).To(BeFalse(), "Hydra is not available")
			})
		})
	})

	// TODO
	// Describe("SetConnectionTimeout", func() {
	// })

	Describe("ReloadServicesCache", func() {
		It("should reload the service cache", func() {
			appIds := []string{service_id}
			services := map[string][]string{
				service_id: test_services,
			}

			c1 := mockServiceCache.EXPECT().GetIds().Return(appIds)
			c2 := mockHydraServiceCache.EXPECT().GetHydraServers().
				Return(test_hydra_servers).After(c1)
			c3 := mockServiceRepository.EXPECT().FindByIds(gomock.Eq(appIds), gomock.Eq(test_hydra_servers)).
				Return(services, nil).After(c2)
			mockServiceCache.EXPECT().Refresh(gomock.Eq(services)).After(c3)

			hydraClient.SetHydraAvailable(true)
			hydraClient.ReloadServicesCache()
		})
		Context("when Hydra availability is set to false", func() {
			It("should not reload the service cache", func() {
				mockServiceRepository.EXPECT().FindByIds(gomock.Any(), gomock.Any()).Times(0)
				mockServiceCache.EXPECT().Refresh(gomock.Any()).Times(0)

				hydraClient.SetHydraAvailable(false)
				hydraClient.ReloadServicesCache()
			})
		})
		Context("when no Hydra server responds", func() {
			It("should not reload the service cache", func() {
				appIds := []string{service_id}

				c1 := mockServiceCache.EXPECT().GetIds().Return(appIds)
				c2 := mockHydraServiceCache.EXPECT().GetHydraServers().
					Return(test_hydra_servers).After(c1)
				mockServiceRepository.EXPECT().FindByIds(gomock.Eq(appIds), gomock.Any()).
					Return(nil, HydraNotAvailableError).After(c2)
				mockServiceCache.EXPECT().Refresh(gomock.Any()).Times(0)

				hydraClient.SetHydraAvailable(true)
				hydraClient.ReloadServicesCache()

				Expect(hydraClient.IsHydraAvailable()).To(BeFalse(), "Hydra is not available")
			})
		})
	})
})
