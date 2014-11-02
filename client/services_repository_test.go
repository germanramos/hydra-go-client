package client_test

import (
	. "github.com/innotech/hydra-go-client/client"
	. "github.com/innotech/hydra-go-client/client/error"
	mock "github.com/innotech/hydra-go-client/client/mock"

	"github.com/innotech/hydra-go-client/vendors/code.google.com/p/gomock/gomock"
	. "github.com/innotech/hydra-go-client/vendors/github.com/onsi/ginkgo"
	. "github.com/innotech/hydra-go-client/vendors/github.com/onsi/gomega"
)

var _ = Describe("ServicesRepository", func() {
	const (
		service_id                    string = "testAppId"
		test_hydra_server_url         string = "http://localhost:8080"
		another_test_hydra_server_url string = "http://localhost:8081"
		test_app_server               string = "http://localhost:8080/app-server-one"
		another_test_app_server       string = "http://localhost:8081/app-server-two"
	)

	var (
		test_services      []string = []string{test_app_server, another_test_app_server}
		test_hydra_servers []string = []string{test_hydra_server_url, another_test_hydra_server_url}

		mockCtrl           *gomock.Controller
		mockRequester      *mock.MockRequester
		servicesRepository *ServicesRepository
	)

	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
		mockRequester = mock.NewMockRequester(mockCtrl)
		servicesRepository = NewServicesRepository()
		servicesRepository.HydraRequester = mockRequester
	})

	AfterEach(func() {
		mockCtrl.Finish()
	})

	Describe("FindById", func() {
		Context("when the first hydra server responds succesfully", func() {
			It("should return the correct list of service servers from first hydra server", func() {
				mockRequester.EXPECT().GetServicesById(gomock.Eq(test_hydra_server_url+AppRootPath), gomock.Eq(service_id)).
					Return(test_services, nil)
				services, err := servicesRepository.FindById(service_id, test_hydra_servers)

				Expect(err).ToNot(HaveOccurred(), "Must not return an error")
				Expect(services).ToNot(BeEmpty(), "Must return a not empty set of services")
				Expect(services).To(Equal(test_services), "The expected services are returned")
			})
		})
		Context("when the response from first hydra server fails", func() {
			It("should return the correct list of service servers from second hydra server", func() {
				c1 := mockRequester.EXPECT().GetServicesById(gomock.Eq(test_hydra_server_url+AppRootPath), gomock.Eq(service_id)).
					Return([]string{}, IncorrectHydraServerResponseError)
				mockRequester.EXPECT().GetServicesById(gomock.Eq(another_test_hydra_server_url+AppRootPath), gomock.Eq(service_id)).
					Return(test_services, nil).After(c1)
				services, err := servicesRepository.FindById(service_id, test_hydra_servers)

				Expect(err).ToNot(HaveOccurred(), "Must not return an error")
				Expect(services).ToNot(BeEmpty(), "Must return a not empty set of services")
				Expect(services).To(Equal(test_services), "The expected services are returned")
			})
		})
		Context("when no hydra server responds successfully", func() {
			It("should return an empty list of servers", func() {
				servicesRepository.SetMaxNumberOfRetries(1)

				c1 := mockRequester.EXPECT().GetServicesById(gomock.Eq(test_hydra_server_url+AppRootPath), gomock.Eq(service_id)).
					Return([]string{}, IncorrectHydraServerResponseError)
				mockRequester.EXPECT().GetServicesById(gomock.Eq(another_test_hydra_server_url+AppRootPath), gomock.Eq(service_id)).
					Return([]string{}, IncorrectHydraServerResponseError).After(c1)
				services, err := servicesRepository.FindById(service_id, test_hydra_servers)

				Expect(err).ToNot(HaveOccurred(), "Must not return an error")
				Expect(services).To(BeEmpty(), "Must return an empty set of services")
			})
			Context("when max number of retries is greater than one", func() {
				It("should return an empty list of servers", func() {
					const retries int = 2
					servicesRepository.SetMaxNumberOfRetries(retries)

					mockRequester.EXPECT().GetServicesById(gomock.Eq(test_hydra_server_url+AppRootPath), gomock.Eq(service_id)).
						Return([]string{}, IncorrectHydraServerResponseError).Times(retries)
					mockRequester.EXPECT().GetServicesById(gomock.Eq(another_test_hydra_server_url+AppRootPath), gomock.Eq(service_id)).
						Return([]string{}, IncorrectHydraServerResponseError).Times(retries)
					services, err := servicesRepository.FindById(service_id, test_hydra_servers)

					Expect(err).ToNot(HaveOccurred(), "Must not return an error")
					Expect(services).To(BeEmpty(), "Must return an empty set of services")
				})
			})
		})

		Context("when no accesible hydra server", func() {
			It("should return an empty list of servers", func() {
				servicesRepository.SetMaxNumberOfRetries(1)

				c1 := mockRequester.EXPECT().GetServicesById(gomock.Eq(test_hydra_server_url+AppRootPath), gomock.Eq(service_id)).
					Return([]string{}, InaccessibleHydraServerError)
				mockRequester.EXPECT().GetServicesById(gomock.Eq(another_test_hydra_server_url+AppRootPath), gomock.Eq(service_id)).
					Return([]string{}, InaccessibleHydraServerError).After(c1)
				services, err := servicesRepository.FindById(service_id, test_hydra_servers)

				Expect(err).To(HaveOccurred(), "Must return an error")
				Expect(err).To(MatchError(HydraNotAvailableError), "Must return an HydraNotAvailableError")
				Expect(services).To(BeEmpty(), "Must return an empty set of services")
			})
		})
	})

	Describe("FindByIds", func() {
		It("should return a map of service ids and servers", func() {
			mockRequester.EXPECT().GetServicesById(gomock.Eq(test_hydra_server_url+AppRootPath), gomock.Eq(service_id)).
				Return(test_services, nil)
			services, err := servicesRepository.FindByIds([]string{service_id}, test_hydra_servers)

			Expect(err).ToNot(HaveOccurred(), "Must not return an error")
			Expect(services).NotTo(BeEmpty(), "Must return an empty map of services")
			expectedServices := map[string][]string{
				service_id: test_services,
			}
			Expect(services).To(Equal(expectedServices), "Must return teh expected map of services")
		})
	})
})
