package failing_before_suite_test

import (
	. "github.com/innotech/hydra-go-client/vendors/github.com/onsi/ginkgo"
)

var _ = Describe("FailingBeforeSuite", func() {
	It("should run", func() {
		println("A TEST")
	})

	It("should run", func() {
		println("A TEST")
	})
})
