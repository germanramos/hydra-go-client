package client

import (
	. "github.com/innotech/hydra-go-client/client/error"

	"encoding/json"
	"io/ioutil"
	"net/http"
)

type Requester interface {
	GetServicesById(serverUrl string, id string) ([]string, error)
}

// TODO: connection timeout
type HydraRequester struct {
}

func NewHydraRequester() *HydraRequester {
	return new(HydraRequester)
}

// Return the candidate url's of the servers sorted by the hydra active algorithm.
func (h *HydraRequester) GetServicesById(serverUrl string, id string) ([]string, error) {
	res, errResponse := http.Get(serverUrl + id)
	if errResponse != nil {
		return []string{}, InaccessibleHydraServerError
	}

	if res.StatusCode != http.StatusOK {
		return []string{}, IncorrectHydraServerResponseError
	}

	defer res.Body.Close()
	body, errBody := ioutil.ReadAll(res.Body)
	if errBody != nil {
		return []string{}, IncorrectHydraServerResponseError
	}

	var servers []string
	errJson := json.Unmarshal(body, &servers)
	if errJson != nil {
		return []string{}, IncorrectHydraServerResponseError
	}
	return servers, nil
}
