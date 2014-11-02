package client

import (
	. "github.com/innotech/hydra-go-client/client/error"

	"time"
)

const (
	AppRootPath string = "/app/"
)

type ServicesRepository struct {
	HydraRequester             Requester
	maxNumberOfRetries         int
	waitBetweenAllServersRetry int
}

func NewServicesRepository() *ServicesRepository {
	return &ServicesRepository{
		HydraRequester:             NewHydraRequester(),
		maxNumberOfRetries:         0,
		waitBetweenAllServersRetry: 0,
	}
}

func (s *ServicesRepository) FindById(id string, servers []string) ([]string, error) {
	var (
		newCandidateServers []string = []string{}
		err                 error
		retries             int = 0
		failedServers       int = 0
	)
	// Infinite loop if maxNumberOfRetries is set to 0.
	// In this case retries can overflow it value, java automatically set to
	// the integer minimum value an the loop goes on
	for s.maxNumberOfRetries == 0 || retries < s.maxNumberOfRetries {
		for _, hydraServer := range servers {
			newCandidateServers, err = s.HydraRequester.GetServicesById(hydraServer+AppRootPath, id)
			if err == nil {
				return newCandidateServers, nil
			} else {
				switch err.(type) {
				case InaccessibleHydraServer:
					failedServers++
				case IncorrectHydraServerResponse:
					continue
				}
			}
		}
		retries++
		s.waitUntilTheNextRetry()
	}

	if len(servers)*s.maxNumberOfRetries == failedServers {
		return []string{}, HydraNotAvailableError
	}
	return newCandidateServers, nil
}

func (s *ServicesRepository) FindByIds(ids []string, servers []string) (map[string][]string, error) {
	newAppServerCache := make(map[string][]string)

	for _, applicationId := range ids {
		newAppServers, err := s.FindById(applicationId, servers)
		if err != nil {
			return newAppServerCache, err
		}
		newAppServerCache[applicationId] = newAppServers
	}

	return newAppServerCache, nil
}

func (s *ServicesRepository) waitUntilTheNextRetry() {
	time.Sleep(time.Duration(s.waitBetweenAllServersRetry) * time.Millisecond)
}

func (s *ServicesRepository) SetMaxNumberOfRetries(numberOfRetries int) {
	s.maxNumberOfRetries = numberOfRetries
}

func (s *ServicesRepository) SetWaitBetweenAllServersRetry(millisecondsToRetry int) {
	s.waitBetweenAllServersRetry = millisecondsToRetry
}

// TODO: set connection timeout
