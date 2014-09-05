package client

import (
	"errors"
)

type HydraClientBuilder interface {
	Config([]string) error
}

type hydraClientFactory struct {
	hydraServers []string
}

var HydraClientFactory *hydraClientFactory = new(hydraClientFactory)

func (h *hydraClientFactory) Config(hydraServers []string) error {
	if hydraServers == nil {
		return errors.New("Invalid Argument: hydraServers can not be nil")
	}
	if len(hydraServers) == 0 {
		return errors.New("Invalid Argument: hydraServers can not be empty")
	}
	h.hydraServers = hydraServers
	return nil
}
