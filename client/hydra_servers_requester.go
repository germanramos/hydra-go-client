package client

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type Requester interface {
	GetCandidateServers(hydraServerUrl string, appId string) ([]string, error)
}

type HydraServersRequester struct {
}

func NewHydraServersRequester() *HydraServersRequester {
	return new(HydraServersRequester)
}

// GetCandidateServers requests to public api of one hydra server the urls
// for the available servers for one application
func (h *HydraServersRequester) GetCandidateServers(hydraServerUrl string, appId string) ([]string, error) {
	res, errResponse := http.Get(hydraServerUrl + appId)
	if errResponse != nil {
		return []string{}, errResponse
	}

	defer res.Body.Close()
	body, errBody := ioutil.ReadAll(res.Body)
	if errBody != nil {
		return []string{}, errBody
	}

	var servers []string
	errJson := json.Unmarshal(body, &servers)
	if errJson != nil {
		return []string{}, errJson
	}
	return servers, nil
}
