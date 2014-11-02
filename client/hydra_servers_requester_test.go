package client_test

// import (
// 	. "github.com/innotech/hydra-go-client/client"

// 	. "github.com/innotech/hydra-go-client/client/mock"
// 	. "github.com/innotech/hydra-go-client/vendors/github.com/onsi/ginkgo"
// 	. "github.com/innotech/hydra-go-client/vendors/github.com/onsi/gomega"

// 	"encoding/json"
// 	"net/http"
// 	"time"
// )

// var _ = Describe("HydraServersRequester", func() {
// 	var (
// 		requester *HydraServersRequester
// 		appId     string
// 	)

// 	BeforeEach(func() {
// 		requester = NewHydraServersRequester()
// 		appId = "app1"
// 	})

// 	Describe("GetCandidateServers", func() {
// 		Context("when hydra server is not accessible", func() {
// 			It("should throw an error", func() {
// 				servers, err := requester.GetCandidateServers("htttp://localhost:3537"+AppRootPath, appId)

// 				Expect(servers).To(BeEmpty())
// 				Expect(err).To(HaveOccurred())
// 			})
// 		})
// 		Context("when hydra server is accessible", func() {
// 			Context("when response status code is different from 200 OK", func() {
// 				It("should throw an error", func() {
// 					routes := []Route{
// 						Route{
// 							Pattern: AppRootPath + appId,
// 							Handler: func(w http.ResponseWriter, r *http.Request) {
// 								w.WriteHeader(http.StatusInternalServerError)
// 							},
// 						},
// 					}
// 					ts := RunHydraServerMock(routes)
// 					defer ts.Close()
// 					time.Sleep(time.Duration(200) * time.Millisecond)

// 					servers, err := requester.GetCandidateServers(ts.URL+AppRootPath, appId)

// 					Expect(servers).To(BeEmpty())
// 					Expect(err).To(HaveOccurred())
// 				})
// 			})
// 			Context("when response status code is 200 OK", func() {
// 				It("should return a server list", func() {
// 					var outputServers []string = []string{"http://www.server1.com", "http://www.server2.com", "http://www.server3.com"}
// 					routes := []Route{
// 						Route{
// 							Pattern: AppRootPath + appId,
// 							Handler: func(w http.ResponseWriter, r *http.Request) {
// 								jsonOutput, _ := json.Marshal(outputServers)
// 								w.WriteHeader(http.StatusOK)
// 								w.Header().Set("Content-Type", "application/json")
// 								w.Write(jsonOutput)
// 							},
// 						},
// 					}
// 					ts := RunHydraServerMock(routes)
// 					defer ts.Close()
// 					time.Sleep(time.Duration(200) * time.Millisecond)

// 					servers, err := requester.GetCandidateServers(ts.URL+AppRootPath, appId)

// 					Expect(err).ToNot(HaveOccurred())
// 					Expect(servers).To(HaveLen(len(outputServers)))
// 					Expect(servers).To(ConsistOf(outputServers))
// 				})
// 			})
// 		})
// 	})
// })
