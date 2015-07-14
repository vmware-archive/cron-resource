package main_test

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
	"github.com/pivotal-cf-experimental/concourse-cron-resource/models"
)

var _ = Describe("In", func() {
	var tmpdir string
	var destination string

	var inCmd *exec.Cmd

	BeforeEach(func() {
		var err error

		tmpdir, err = ioutil.TempDir("", "in-destination")
		Expect(err).ShouldNot(HaveOccurred())

		destination = path.Join(tmpdir, "in-dir")

		inCmd = exec.Command(inPath, destination)
	})

	AfterEach(func() {
		os.RemoveAll(tmpdir)
	})

	Context("when executed", func() {
		var request models.InRequest
		var response models.InResponse

		BeforeEach(func() {
			request = models.InRequest{
				Version: models.Version{Time: time.Now()},
				Source: models.Source{
					Expression: "* * * * *",
					Location:   "America/New_York",
				},
			}

			response = models.InResponse{}
		})

		JustBeforeEach(func() {
			stdin, err := inCmd.StdinPipe()
			Expect(err).ShouldNot(HaveOccurred())

			session, err := gexec.Start(inCmd, GinkgoWriter, GinkgoWriter)
			Expect(err).ShouldNot(HaveOccurred())

			err = json.NewEncoder(stdin).Encode(request)
			Expect(err).ShouldNot(HaveOccurred())

			Eventually(session).Should(gexec.Exit(0))

			err = json.Unmarshal(session.Out.Contents(), &response)
			Expect(err).ShouldNot(HaveOccurred())
		})

		It("reports the version's time as the version", func() {
			Expect(response.Version.Time.UnixNano()).Should(Equal(request.Version.Time.UnixNano()))
		})

		It("writes the requested version and source to the destination", func() {
			input, err := os.Open(filepath.Join(destination, "input"))
			Expect(err).ShouldNot(HaveOccurred())

			var requested models.InRequest
			err = json.NewDecoder(input).Decode(&requested)
			Expect(err).ShouldNot(HaveOccurred())

			Expect(requested.Version.Time.Unix()).Should(Equal(request.Version.Time.Unix()))
			Expect(requested.Source).Should(Equal(request.Source))
		})
	})
})
