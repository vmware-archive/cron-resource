package main_test

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"time"

	"github.com/gorhill/cronexpr"
	"github.com/pivotal-cf-experimental/cron-resource/models"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"
)

var _ = Describe("Check", func() {
	var checkCmd *exec.Cmd

	BeforeEach(func() {
		checkCmd = exec.Command(checkPath)
	})

	Context("MustParse", func() {
		It("can parse crontab expressions", func() {
			expr := cronexpr.MustParse("* * * * 0-5")
			Expect(expr).ToNot(BeNil())
		})
	})

	Context("when a crontab expression is specified", func() {
		var request models.CheckRequest
		var response models.CheckResponse
		var session *gexec.Session

		BeforeEach(func() {
			request = models.CheckRequest{
				Source: models.Source{
					Location: "America/New_York",
				},
			}
			response = models.CheckResponse{}
		})

		JustBeforeEach(func() {
			stdin, err := checkCmd.StdinPipe()
			Expect(err).ShouldNot(HaveOccurred())

			session, err = gexec.Start(checkCmd, GinkgoWriter, GinkgoWriter)
			Expect(err).ShouldNot(HaveOccurred())

			err = json.NewEncoder(stdin).Encode(request)
			Expect(err).ShouldNot(HaveOccurred())
		})

		Context("the expression is invalid", func() {
			It("exits with status code 1", func() {
				request.Source.Expression = "invalid"
				Eventually(session.Err).Should(gbytes.Say("invalid crontab expression"))
				Eventually(session).Should(gexec.Exit(1))
			})
		})

		Context("expression is valid", func() {
			JustBeforeEach(func() {
				Eventually(session).Should(gexec.Exit(0))

				err := json.Unmarshal(session.Out.Contents(), &response)
				Expect(err).ShouldNot(HaveOccurred())
			})

			Context("wildcard expression", func() {
				BeforeEach(func() {
					request.Source.Expression = "* * * * *"
				})

				Context("when no version is given", func() {
					It("outputs time.Now()", func() {
						Expect(response).To(HaveLen(1))
						Expect(response[0].Time.Unix()).To(BeNumerically("~", time.Now().Unix(), 1))
					})
				})

				Context("when a version is given", func() {
					Context("when the resource has already triggered in that minute", func() {
						BeforeEach(func() {
							request.Version.Time = time.Now()
						})

						It("doesn't print anything", func() {
							Expect(response).To(BeEmpty())
						})
					})

					Context("when the resource hasn't triggered in that minute", func() {
						BeforeEach(func() {
							request.Version.Time = time.Now().Add(-2 * time.Minute)
						})

						It("outputs time.Now()", func() {
							Expect(response).To(HaveLen(1))
							Expect(response[0].Time.Unix()).To(BeNumerically("~", time.Now().Unix(), 1))
						})
					})
				})
			})

			Context("when given a crontab expression that triggers 30 minutes ago", func() {
				BeforeEach(func() {
					halfHourAgo := time.Now().Add(-30 * time.Minute)
					request.Source.Expression = fmt.Sprintf("%d * * * *", halfHourAgo.Minute())
				})

				Context("when no version is given", func() {
					It("outputs time.Now()", func() {
						Expect(response).To(HaveLen(1))
						Expect(response[0].Time.Unix()).To(BeNumerically("~", time.Now().Unix(), 1))
					})
				})
			})

			Context("when given a crontab expression that triggers in the previous hour", func() {
				BeforeEach(func() {
					request.Source.Expression = fmt.Sprintf("0 %d * * *", time.Now().Hour()-1)
				})

				Context("when no version is given", func() {
					It("doesn't output any versions", func() {
						Expect(response).To(BeEmpty())
					})

					Context("when FireImmediately is true", func() {
						BeforeEach(func() {
							request.Source.FireImmediately = true
						})

						It("outputs time.Now()", func() {
							Expect(response).To(HaveLen(1))
							Expect(response[0].Time.Unix()).To(BeNumerically("~", time.Now().Unix(), 1))
						})
					})
				})

				Context("when a version that is a day old is given", func() {
					BeforeEach(func() {
						request.Version.Time = time.Now().Add(-25 * time.Hour)
					})

					It("outputs time.Now()", func() {
						Expect(response).To(HaveLen(1))
						Expect(response[0].Time.Unix()).To(BeNumerically("~", time.Now().Unix(), 1))
					})
				})

				Context("when a version is given", func() {
					BeforeEach(func() {
						request.Version.Time = time.Now().Add(-1 * time.Hour)
					})

					It("doesn't output any versions", func() {
						Expect(response).To(BeEmpty())
					})
				})
			})

			Context("Timezones", func() {
				var loc *time.Location
				var err error

				Context("when a different timezone is specified", func() {
					BeforeEach(func() {
						request.Source.Location = "America/Vancouver"
						loc, err = time.LoadLocation("America/Vancouver")
						Expect(err).NotTo(HaveOccurred())
					})

					Context("when given a crontab expression that triggers in the current hour in the given timezone", func() {

						BeforeEach(func() {
							request.Source.Expression = fmt.Sprintf("* %d * * *", time.Now().In(loc).Hour())
						})

						It("outputs time.Now()", func() {
							Expect(response).To(HaveLen(1))
							Expect(response[0].Time.Unix()).To(BeNumerically("~", time.Now().Unix(), 1))
							_, offset := response[0].Time.Zone()
							_, expectedOffset := time.Now().In(loc).Zone()
							Expect(offset).To(Equal(expectedOffset))
						})
					})
				})

				Context("when no timezone is given", func() {
					BeforeEach(func() {
						request.Source.Location = ""
						loc, err = time.LoadLocation("UTC")
						Expect(err).NotTo(HaveOccurred())
					})

					Context("when given a crontab expression that triggers in the current hour", func() {
						BeforeEach(func() {
							request.Source.Expression = fmt.Sprintf("* %d * * *", time.Now().In(loc).Hour())
						})

						It("outputs time.Now()", func() {
							Expect(response).To(HaveLen(1))
							Expect(response[0].Time.Unix()).To(BeNumerically("~", time.Now().Unix(), 1))
							_, offset := response[0].Time.Zone()
							_, expectedOffset := time.Now().In(loc).Zone()
							Expect(offset).To(Equal(expectedOffset))
						})
					})
				})
			})
		})
	})
})
