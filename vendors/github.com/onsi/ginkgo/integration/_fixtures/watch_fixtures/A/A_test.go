package A_test

import (
	. "github.com/innotech/hydra-go-client/vendors/github.com/onsi/ginkgo/integration/_fixtures/watch_fixtures/A"

	. "github.com/innotech/hydra-go-client/vendors/github.com/onsi/ginkgo"
	. "github.com/innotech/hydra-go-client/vendors/github.com/onsi/gomega"
)

var _ = Describe("A", func() {
	It("should do it", func() {
		Ω(DoIt()).Should(Equal("done!"))
	})
})
