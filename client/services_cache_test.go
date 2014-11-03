package client_test

import (
	. "github.com/innotech/hydra-go-client/client"

	. "github.com/innotech/hydra-go-client/vendors/github.com/onsi/ginkgo"
	. "github.com/innotech/hydra-go-client/vendors/github.com/onsi/gomega"
)

var _ = Describe("ServicesCache", func() {
	const (
		service_id              string = "service"
		test_app_server         string = "http://localhost:8080/app-server-one"
		another_test_app_server string = "http://localhost:8081/app-server-second"
	)

	var (
		test_services []string = []string{test_app_server, another_test_app_server}

		servicesCache *ServicesCache
	)

	BeforeEach(func() {
		servicesCache = NewServicesCache()
	})

	Describe("FindById", func() {
		It("should return the services stored by id", func() {
			servicesCache.PutService(service_id, test_services)
			services := servicesCache.FindById(service_id)

			Expect(services).ToNot(BeEmpty(), "Must return a not empty list of services")
			Expect(services).To(Equal(test_services), "The returned services must be the expected")
		})
		Context("when the service (service id) doesn't exist", func() {
			It("should return the services stored by id", func() {
				services := servicesCache.FindById(service_id)

				Expect(services).To(BeEmpty(), "Must return an  empty list of services")
				Expect(services).To(Equal([]string{}), "The returned services must be the expected")
			})
		})
	})

	Describe("Exists", func() {
		Context("when service exists", func() {
			It("should return true", func() {
				servicesCache.PutService(service_id, test_services)
				exists := servicesCache.Exists(service_id)

				Expect(exists).To(BeTrue(), "Must return true")
			})
		})
	})

	Describe("GetIds", func() {
		Context("when services are registered", func() {
			It("should return the service ids", func() {
				servicesCache.PutService(service_id, test_services)
				ids := servicesCache.GetIds()

				Expect(ids).ToNot(BeEmpty(), "Must return an  empty list of services")
				Expect(ids).To(ContainElement(service_id), "Must return the expected ids")
			})
		})
	})

	Describe("Refresh", func() {
		It("should refresh the whole cache", func() {
			s := map[string][]string{
				service_id: test_services,
			}
			servicesCache.Refresh(s)
			services := servicesCache.FindById(service_id)

			Expect(services).ToNot(BeEmpty(), "Must not return an empty list of services")
			Expect(services).To(Equal(test_services), "The returned services must be the expected")
		})
	})
})
