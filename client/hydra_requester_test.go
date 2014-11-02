package client_test

import (
	. "github.com/innotech/hydra-go-client/client"
	. "github.com/innotech/hydra-go-client/client/error"
	. "github.com/innotech/hydra-go-client/client/mock"
	. "github.com/innotech/hydra-go-client/vendors/github.com/onsi/ginkgo"
	. "github.com/innotech/hydra-go-client/vendors/github.com/onsi/gomega"

	"encoding/json"
	"net/http"
	"time"
)

var _ = FDescribe("HydraRequester", func() {
	const (
		app_id              string = "testAppId"
		test_hydra_server_0 string = "http://localhost:8080"
		test_hydra_server_1 string = "http://localhost:8081"
		test_hydra_server_2 string = "http://localhost:8082"
	)

	var (
		hydraRequester *HydraRequester
	)

	BeforeEach(func() {
		hydraRequester = NewHydraRequester()
	})

	Describe("GetServersById", func() {
		Context("when hydra server responds successfully", func() {
			It("should return a list of servers", func() {
				var appServers []string = []string{test_hydra_server_0, test_hydra_server_1, test_hydra_server_2}
				routes := []Route{
					Route{
						Pattern: AppRootPath + app_id,
						Handler: func(w http.ResponseWriter, r *http.Request) {
							jsonOutput, _ := json.Marshal(appServers)
							w.WriteHeader(http.StatusOK)
							w.Header().Set("Content-Type", "application/json")
							w.Write(jsonOutput)
						},
					},
				}
				ts := RunHydraServerMock(routes)
				defer ts.Close()
				time.Sleep(time.Duration(200) * time.Millisecond)

				candidateServers, err := hydraRequester.GetServicesById(ts.URL+AppRootPath, app_id)
				Expect(err).ToNot(HaveOccurred(), "Must not return an error")
				Expect(candidateServers).ToNot(BeEmpty(), "Must return a not empty list of servers")
				Expect(candidateServers).To(HaveLen(len(appServers)), "The number of candidate servers must be the expected")
				Expect(candidateServers).To(Equal(appServers), "The expected servers are returned")
			})
		})
		Context("when hydra server responds with bad request status", func() {
			It("should return an empty list of servers and an valid error", func() {
				routes := []Route{
					Route{
						Pattern: AppRootPath + app_id,
						Handler: func(w http.ResponseWriter, r *http.Request) {
							w.WriteHeader(http.StatusBadRequest)
						},
					},
				}
				ts := RunHydraServerMock(routes)
				defer ts.Close()
				time.Sleep(time.Duration(200) * time.Millisecond)

				candidateServers, err := hydraRequester.GetServicesById(ts.URL+AppRootPath, app_id)
				Expect(candidateServers).To(BeEmpty(), "Must return an empty list of servers")
				Expect(err).To(HaveOccurred(), "Must return an error")
				// TODO
				Expect(err).To(MatchError(IncorrectHydraServerResponse), "The expected error is returned")
			})
		})
		Context("when hydra server is not accesible", func() {
			It("should return an empty list of servers and an valid error", func() {
				candidateServers, err := hydraRequester.GetServicesById(test_hydra_server_0+AppRootPath, app_id)
				Expect(candidateServers).To(BeEmpty(), "Must return an empty list of servers")
				Expect(err).To(HaveOccurred(), "Must return an error")
				// TODO
				Expect(err).To(MatchError(InaccessibleHydraServer), "The expected error is returned")
			})
		})
	})
})
