package client_test

import (
	. "github.com/innotech/hydra-go-client/client"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"encoding/json"
	"net/http"
	"net/http/httptest"
)

var _ = Describe("HydraServersRequester", func() {
	var (
		requester *HydraServersRequester
	)

	BeforeEach(func() {
		requester = NewHydraServersRequester()
	})

	Describe("GetCandidateServers", func() {
		Context("when response status code is different from 200 OK", func() {
			It("should throw an error", func() {
				ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					// http.Error(rw, err.Error(), http.StatusInternalServerError)
					w.WriteHeader(http.StatusInternalServerError)
					// rw.Header().Set("Content-Type", "application/json")
					// rw.Write(jsonOutput)
				}))
				defer ts.Close()

				servers, err := requester.GetCandidateServers(ts.URL, "")

				Expect(servers).To(BeEmpty())
				Expect(err).To(HaveOccurred())
			})
		})

		Context("when response status code is 200 OK", func() {
			It("should return a server list", func() {
				var outputServers []string = []string{"http://www.server1.com", "http://www.server2.com", "http://www.server3.com"}
				ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					jsonOutput, _ := json.Marshal(outputServers)
					w.WriteHeader(http.StatusOK)
					w.Header().Set("Content-Type", "application/json")
					w.Write(jsonOutput)
				}))
				defer ts.Close()

				servers, err := requester.GetCandidateServers(ts.URL, "")

				Expect(err).ToNot(HaveOccurred())
				Expect(servers).To(HaveLen(len(outputServers)))
				Expect(servers).To(ConsistOf(outputServers))
			})
		})
	})
})
