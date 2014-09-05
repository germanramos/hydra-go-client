package passing_before_suite_test

import (
	. "github.com/innotech/hydra-go-client/vendors/github.com/onsi/ginkgo"
	. "github.com/innotech/hydra-go-client/vendors/github.com/onsi/gomega"
)

var _ = Describe("PassingSuiteSetup", func() {
	It("should pass", func() {
		Ω(a).Should(Equal("ran before suite"))
		Ω(b).Should(BeEmpty())
	})

	It("should pass", func() {
		Ω(a).Should(Equal("ran before suite"))
		Ω(b).Should(BeEmpty())
	})

	It("should pass", func() {
		Ω(a).Should(Equal("ran before suite"))
		Ω(b).Should(BeEmpty())
	})

	It("should pass", func() {
		Ω(a).Should(Equal("ran before suite"))
		Ω(b).Should(BeEmpty())
	})
})
