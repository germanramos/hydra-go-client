package client_test

import (
	. "github.com/innotech/hydra-go-client/client"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("HydraClientFactory", func() {
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
			It("should throw an error", func() {
				err := HydraClientFactory.Config([]string{"http://localhost:8080"})
				Expect(err).ShouldNot(HaveOccurred())
			})
		})
	})
})
