package client_test

import (
	. "github.com/innotech/hydra-go-client/client"

	. "github.com/innotech/hydra-go-client/vendors/github.com/onsi/ginkgo"
	. "github.com/innotech/hydra-go-client/vendors/github.com/onsi/gomega"
)

var _ = FDescribe("HydraServiceCache", func() {
	const (
		test_hydra_server_url         string = "http://localhost:8080"
		another_test_hydra_server_url string = "http://localhost:8081"
	)

	var (
		test_hydra_servers     []string = []string{test_hydra_server_url}
		test_new_hydra_servers []string = []string{another_test_hydra_server_url}
		hydraServiceCache      *HydraServiceCache
	)

	BeforeEach(func() {
		hydraServiceCache = NewHydraServiceCache(test_hydra_servers)
	})

	Describe("GetHydraServers", func() {
		It("should return the cached hydra servers", func() {
			hydraServers := hydraServiceCache.GetHydraServers()

			Expect(hydraServers).ToNot(BeEmpty(), "Must return a not empty set of servers")
			Expect(hydraServers).To(Equal(test_hydra_servers), "The expected servers are returned")
		})
		Context("when the cache is refreshed", func() {
			It("should return the cached hydra servers", func() {
				hydraServiceCache.Refresh(test_new_hydra_servers)
				hydraServers := hydraServiceCache.GetHydraServers()

				Expect(hydraServers).ToNot(BeEmpty(), "Must return a not empty set of servers")
				Expect(hydraServers).To(Equal(test_new_hydra_servers), "The new servers are returned")
			})
		})
		Context("when the cache is refreshed with an empty list of servers", func() {
			It("should return the seed hydra servers", func() {
				hydraServiceCache.Refresh(test_new_hydra_servers)
				hydraServiceCache.Refresh([]string{})
				hydraServers := hydraServiceCache.GetHydraServers()

				Expect(hydraServers).ToNot(BeEmpty(), "Must return a not empty set of servers")
				Expect(hydraServers).To(Equal([]string{another_test_hydra_server_url, test_hydra_server_url}), "The known servers are returned")
			})
		})
	})
})
