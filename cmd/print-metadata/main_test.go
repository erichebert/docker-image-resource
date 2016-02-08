package main_test

import (
	"encoding/json"
	"os/exec"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/onsi/gomega/gexec"
)

type imageMetadata struct {
	Env []string `json:"env"`
}

var _ = Describe("print-metadata", func() {
	var (
		cmd     *exec.Cmd
		session *gexec.Session

		metadata imageMetadata
	)

	BeforeEach(func() {
		cmd = exec.Command(printMetadataPath)
	})

	JustBeforeEach(func() {
		var err error
		session, err = gexec.Start(cmd, GinkgoWriter, GinkgoWriter)
		Expect(err).NotTo(HaveOccurred())
		Eventually(session).Should(gexec.Exit(0))

		err = json.Unmarshal(session.Out.Contents(), &metadata)
		Expect(err).NotTo(HaveOccurred())
	})

	Context("when it is running in an environment with environment variables", func() {
		BeforeEach(func() {
			cmd.Env = []string{
				"SOME=foo",
				"AMAZING=bar",
				"ENV=baz",
			}
		})

		It("outputs them on stdout", func() {
			Expect(metadata.Env).To(ConsistOf([]string{
				"SOME=foo",
				"AMAZING=bar",
				"ENV=baz",
			}))
		})
	})

	Context("when it is running in an environment with environment variables in the blacklist", func() {
		BeforeEach(func() {
			cmd.Env = []string{
				"SOME=foo",
				"HOSTNAME=bar",
				"ENV=baz",
			}
		})

		It("outputs everything but them on stdout", func() {
			Expect(metadata.Env).To(ConsistOf([]string{
				"SOME=foo",
				"ENV=baz",
			}))
		})
	})
})
