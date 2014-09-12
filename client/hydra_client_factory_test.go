package client_test

import (
	. "github.com/innotech/hydra-go-client/client"
	mock "github.com/innotech/hydra-go-client/client/mock"

	"github.com/innotech/hydra-go-client/vendors/code.google.com/p/gomock/gomock"
	. "github.com/innotech/hydra-go-client/vendors/github.com/onsi/ginkgo"
	. "github.com/innotech/hydra-go-client/vendors/github.com/onsi/gomega"

	"time"
)

var _ = Describe("HydraClientFactory", func() {
	var (
		mockCtrl        *gomock.Controller
		mockHydraClient *mock.MockHydraClient
	)

	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
		mockHydraClient = mock.NewMockHydraClient(mockCtrl)
	})

	AfterEach(func() {
		mockCtrl.Finish()
	})

	Describe("Build", func() {
		It("should build a HydraClient", func() {
			var client *Client
			Expect(HydraClientFactory.Build()).To(BeAssignableToTypeOf(client))
		})
	})

	Describe("Config", func() {
		Context("when hydra server list argument is nil", func() {
			It("should throw an error", func() {
				err := HydraClientFactory.Config(nil)
				Expect(err).Should(HaveOccurred())
			})
		})
		Context("when hydra server list argument is a empty list", func() {
			It("should throw an error", func() {
				err := HydraClientFactory.Config([]string{})
				Expect(err).Should(HaveOccurred())
			})
		})
		Context("when hydra server list argument is a valid list of servers", func() {
			It("should set hydra server list successfully", func() {
				err := HydraClientFactory.Config([]string{"http://localhost:8080"})
				Expect(err).ShouldNot(HaveOccurred())
				// TODO: check
			})
		})
	})

	Describe("WithHydraServersCacheDuration", func() {
		Context("when duration argument is a valid uint number", func() {
			It("should set hydra servers cache duration successfully", func() {
				const hydraServersCacheDuration time.Duration = time.Duration(30000) * time.Millisecond
				h := HydraClientFactory.WithHydraServersCacheDuration(hydraServersCacheDuration)
				Expect(h).To(Equal(HydraClientFactory))
				Expect(hydraServersCacheDuration).To(Equal(HydraClientFactory.GetHydraServersCacheDuration()))
			})
		})
	})

	Describe("WithAppsCacheDuration", func() {
		Context("when duration argument is a valid uint number", func() {
			It("should set apps cache duration successfully", func() {
				const appsCacheDuration time.Duration = time.Duration(30000) * time.Millisecond
				h := HydraClientFactory.WithAppsCacheDuration(appsCacheDuration)
				Expect(h).To(Equal(HydraClientFactory))
				Expect(appsCacheDuration).To(Equal(HydraClientFactory.GetAppsCacheDuration()))
			})
		})
	})

	Describe("WithMaxNumberOfRetriesPerHydraServer", func() {
		Context("when duration argument is a valid uint number", func() {
			It("should set apps cache duration successfully", func() {
				const retries uint = 3
				h := HydraClientFactory.WithMaxNumberOfRetriesPerHydraServer(retries)
				Expect(h).To(Equal(HydraClientFactory))
				Expect(retries).To(Equal(HydraClientFactory.GetMaxNumberOfRetriesPerHydraServer()))
			})
		})
	})

	Describe("WaitBetweenAllServersRetry", func() {
		Context("when duration argument is a valid uint number", func() {
			It("should set wait between all servers retry successfully", func() {
				const appsCacheDuration time.Duration = time.Duration(30000) * time.Millisecond
				h := HydraClientFactory.WaitBetweenAllServersRetry(appsCacheDuration)
				Expect(h).To(Equal(HydraClientFactory))
				Expect(appsCacheDuration).To(Equal(HydraClientFactory.GetDurationBetweenAllServersRetry()))
			})
		})
	})
})
