package client_test

import (
	. "github.com/innotech/hydra-go-client/client"
	mock "github.com/innotech/hydra-go-client/client/mock"

	"code.google.com/p/gomock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"errors"
)

var _ = Describe("HydraClient", func() {
	var (
		hydraClient   *Client
		mockCtrl      *gomock.Controller
		mockRequester *mock.MockRequester
	)

	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
		mockRequester = mock.NewMockRequester(mockCtrl)
		hydraClient = NewClient([]string{"http://localhost:8080"}, mockRequester)
	})

	AfterEach(func() {
		mockCtrl.Finish()
	})

	Describe("Get", func() {
		Context("when an illegal application ID is passed as an argument", func() {
			It("should throw an error", func() {
				servers, err := hydraClient.Get("", false)
				Expect(servers).To(BeEmpty())
				Expect(err).To(HaveOccurred())
			})
		})
		Context("when the cache should not be refreshed", func() {
			Context("when the application ID doesn't exist", func() {
				It("should request servers from hydra server", func() {
					mockRequester.EXPECT().GetCandidateServers(gomock.Any(), gomock.Eq("app1"))
					_, _ = hydraClient.Get("app1", false)
				})
			})
			Context("when the application ID exists", func() {
				It("should not request servers from hydra server", func() {
					mockRequester.EXPECT().GetCandidateServers(gomock.Any(), gomock.Eq("app1"))
					_, _ = hydraClient.Get("app1", false)

					mockRequester.EXPECT().GetCandidateServers(gomock.Any(), gomock.Any()).Times(0)
					_, _ = hydraClient.Get("app1", false)
				})
			})
		})
		Context("when the cache should be refreshed", func() {
			Context("when the application ID doesn't exist", func() {
				It("should request servers from hydra server", func() {
					mockRequester.EXPECT().GetCandidateServers(gomock.Any(), gomock.Eq("app1"))
					_, _ = hydraClient.Get("app1", true)
				})
			})
			Context("when the application ID exists", func() {
				It("should request servers from hydra server", func() {
					mockRequester.EXPECT().GetCandidateServers(gomock.Any(), gomock.Eq("app1"))
					_, _ = hydraClient.Get("app1", false)

					mockRequester.EXPECT().GetCandidateServers(gomock.Any(), gomock.Eq("app1"))
					_, _ = hydraClient.Get("app1", true)
				})
			})
		})
	})

	Describe("ReloadHydraServers", func() {
		Context("when hydra server is not accessible", func() {
			It("should consider that Hydra is not available", func() {
				mockRequester.EXPECT().GetCandidateServers(gomock.Any(), gomock.Eq("hydra")).Return([]string{}, errors.New("Not Found"))
				hydraClient.SetMaxNumberOfRetriesPerHydraServer(1)
				hydraClient.ReloadHydraServers()
				Expect(hydraClient.IsHydraAvailable()).To(BeFalse())
			})
		})
		Context("when hydra server responses with a list of servers", func() {
			It("should consider that Hydra is available", func() {
				mockRequester.EXPECT().GetCandidateServers(gomock.Any(), gomock.Eq("hydra")).Return([]string{"http://localhost:8080"}, nil)
				hydraClient.SetMaxNumberOfRetriesPerHydraServer(1)
				hydraClient.ReloadHydraServers()
				Expect(hydraClient.IsHydraAvailable()).To(BeTrue())
			})
		})
	})

	Describe("ReloadAppServers", func() {
		Context("when no application registered", func() {
			It("should not send any request to hydra servers", func() {
				mockRequester.EXPECT().GetCandidateServers(gomock.Any(), gomock.Any()).Times(0)
				hydraClient.ReloadAppServers()
			})
		})
		// Context("when some applications are registered", func() {
		// 	It("should require update the application cache", func() {
		// 		mockRequester.EXPECT().GetCandidateServers(gomock.Any(), gomock.Eq("app1")).Return([]string{"http://localhost:8080"}, nil)
		// 		_, _ = hydraClient.Get("app1", false)
		// 		// Expect(se, ...)

		// 		// mockRequester.EXPECT().GetCandidateServers(gomock.Any(), gomock.Any())
		// 		mockRequester.EXPECT().GetCandidateServers(gomock.Any(), gomock.Eq("app1"))
		// 		hydraClient.ReloadAppServers()
		// 	})
		// })
	})
})
