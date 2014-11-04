package client

import (
	"errors"
	"strings"
)

type Client interface {
	Get(serviceId string, shortcutCache bool) ([]string, error)
	GetHydraServers() []string
	IsHydraAvailable() bool
	ReloadHydraServiceCache()
	ReloadServicesCache()
	SetHydraAvailable(available bool)
	SetMaxNumberOfRetries(numberOfRetries int)
	SetWaitBetweenAllServersRetry(millisecondsToRetry int)
}

const (
	HydraAppId string = "hydra"
)

type HydraClient struct {
	hydraAvailable     bool // TODO: should be atomic
	HydraServiceCache  HydraCache
	ServicesCache      ServiceCache
	ServicesRepository ServiceRepository
}

func NewHydraClient(seedHydraServers []string) *HydraClient {
	return &HydraClient{
		hydraAvailable:    false,
		HydraServiceCache: NewHydraServiceCache(seedHydraServers),
	}
}

// Retrieve a list of servers sorted by hydra available for a concrete
// application. This method can shortcut the cache.
// In android this method must be called in a async task to avoid the interaction of the main thread
// with the network, use getAsync instead.
func (h *HydraClient) get(serviceId string, shortcutCache bool) ([]string, error) {
	if strings.Trim(serviceId, " ") == "" {
		return []string{}, errors.New("Illegal Argument: serviceId must be a single word")
	}

	if !shortcutCache && h.ServicesCache.Exists(serviceId) {
		return h.ServicesCache.FindById(serviceId), nil
	}

	servers, err := h.ServicesRepository.FindById(serviceId, h.HydraServiceCache.GetHydraServers())
	if err == nil {
		h.ServicesCache.PutService(serviceId, servers)
	} else {
		// TODO: Lock
		h.hydraAvailable = false
	}

	return servers, nil
}

func (h *HydraClient) Get(serviceId string) ([]string, error) {
	return h.get(serviceId, false)
}

func (h *HydraClient) GetShortcuttingTheCache(serviceId string) ([]string, error) {
	return h.get(serviceId, true)
}

// Return a future with the server request. Avoid the interaction of the calling thread with the network.
func (h *HydraClient) IsHydraAvailable() bool {
	// TODO: Lock
	return h.hydraAvailable
}

func (h *HydraClient) ReloadHydraServiceCache() {
	// TODO: Lock
	newHydraServers, err := h.ServicesRepository.FindById(HydraAppId, h.HydraServiceCache.GetHydraServers())
	if err == nil {
		h.HydraServiceCache.Refresh(newHydraServers)
		if len(newHydraServers) > 0 {
			h.hydraAvailable = true
		} else {
			h.hydraAvailable = false
		}
	} else {
		h.hydraAvailable = false
	}
}

func (h *HydraClient) ReloadServicesCache() {
	if !h.hydraAvailable {
		return
	}

	servers, err := h.ServicesRepository.FindByIds(h.ServicesCache.GetIds(), h.HydraServiceCache.GetHydraServers())
	if err != nil {
		h.hydraAvailable = false
		return
	}
	h.ServicesCache.Refresh(servers)
}

func (h *HydraClient) SetMaxNumberOfRetries(numberOfRetries int) {
	h.ServicesRepository.SetMaxNumberOfRetries(numberOfRetries)
}

func (h *HydraClient) SetWaitBetweenAllServersRetry(millisecondsToRetry int) {
	h.ServicesRepository.SetWaitBetweenAllServersRetry(millisecondsToRetry)
}

// TODO
// func (h *HydraClient) SetConnectionTimeout(connection int) {

// }

func (h *HydraClient) GetHydraServers() []string {
	return h.HydraServiceCache.GetHydraServers()
}

func (h *HydraClient) SetHydraAvailable(available bool) {
	h.hydraAvailable = available
}
