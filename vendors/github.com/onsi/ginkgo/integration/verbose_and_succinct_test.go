package integration_test

import (
	. "github.com/innotech/hydra-go-client/vendors/github.com/onsi/ginkgo"
	. "github.com/innotech/hydra-go-client/vendors/github.com/onsi/gomega"
	"github.com/innotech/hydra-go-client/vendors/github.com/onsi/gomega/gexec"
)

var _ = Describe("Verbose And Succinct Mode", func() {
	var pathToTest string
	var otherPathToTest string

	Context("when running one package", func() {
		BeforeEach(func() {
			pathToTest = tmpPath("ginkgo")
			copyIn("passing_ginkgo_tests", pathToTest)
		})

		It("should default to non-succinct mode", func() {
			session := startGinkgo(pathToTest, "--noColor")
			Eventually(session).Should(gexec.Exit(0))
			output := session.Out.Contents()

			Ω(output).Should(ContainSubstring("Running Suite: Passing_ginkgo_tests Suite"))
		})
	})

	Context("when running more than one package", func() {
		BeforeEach(func() {
			pathToTest = tmpPath("ginkgo")
			copyIn("passing_ginkgo_tests", pathToTest)
			otherPathToTest = tmpPath("more_ginkgo")
			copyIn("more_ginkgo_tests", otherPathToTest)
		})

		Context("with no flags set", func() {
			It("should default to succinct mode", func() {
				session := startGinkgo(pathToTest, "--noColor", pathToTest, otherPathToTest)
				Eventually(session).Should(gexec.Exit(0))
				output := session.Out.Contents()

				Ω(output).Should(ContainSubstring("] Passing_ginkgo_tests Suite - 3/3 specs ••• SUCCESS!"))
				Ω(output).Should(ContainSubstring("] More_ginkgo_tests Suite - 2/2 specs •• SUCCESS!"))
			})
		})

		Context("with --succinct=false", func() {
			It("should not be in succinct mode", func() {
				session := startGinkgo(pathToTest, "--noColor", "--succinct=false", pathToTest, otherPathToTest)
				Eventually(session).Should(gexec.Exit(0))
				output := session.Out.Contents()

				Ω(output).Should(ContainSubstring("Running Suite: Passing_ginkgo_tests Suite"))
				Ω(output).Should(ContainSubstring("Running Suite: More_ginkgo_tests Suite"))
			})
		})

		Context("with -v", func() {
			It("should not be in succinct mode, but should be verbose", func() {
				session := startGinkgo(pathToTest, "--noColor", "-v", pathToTest, otherPathToTest)
				Eventually(session).Should(gexec.Exit(0))
				output := session.Out.Contents()

				Ω(output).Should(ContainSubstring("Running Suite: Passing_ginkgo_tests Suite"))
				Ω(output).Should(ContainSubstring("Running Suite: More_ginkgo_tests Suite"))
				Ω(output).Should(ContainSubstring("should proxy strings"))
				Ω(output).Should(ContainSubstring("should always pass"))
			})
		})
	})
})
