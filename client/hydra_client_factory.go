package client

import (
	"errors"
	"time"
)

const (
	default_hydra_server_refresh time.Duration = time.Duration(60) * time.Second
	default_hydra_apps_refresh   time.Duration = time.Duration(20) * time.Second
	default_retries_number       int           = 10
)

type clientMaker interface {
	MakeClient(seedHydraServers []string) Client
}

type clientInstantiator struct{}

func (c *clientInstantiator) MakeClient(seedHydraServers []string) Client {
	return NewHydraClient(seedHydraServers)
}

type hydraMonitorMaker interface {
	MakeHydraMonitor(hydraClient Client, refreshTime time.Duration) CacheMonitor
}

type hydraMonitorInstantiator struct{}

func (c *hydraMonitorInstantiator) MakeHydraMonitor(hydraClient Client, refreshTime time.Duration) CacheMonitor {
	return NewHydraServiceCacheMonitor(hydraClient, refreshTime)
}

type appsMonitorMaker interface {
	MakeAppsMonitor(hydraClient Client, refreshTime time.Duration) CacheMonitor
}

type appsMonitorInstantiator struct{}

func (c *appsMonitorInstantiator) MakeAppsMonitor(hydraClient Client, refreshTime time.Duration) CacheMonitor {
	return NewServicesCacheMonitor(hydraClient, refreshTime)
}

type hydraClientFactory struct {
	AppsMonitorInstantiator  appsMonitorMaker
	ClientInstantiator       clientMaker
	HydraMonitorInstantiator hydraMonitorMaker

	hydraClient            Client
	hydraServerRefreshTime time.Duration
	hydraAppsRefreshTime   time.Duration
	hydraMonitor           CacheMonitor
	appsMonitor            CacheMonitor
	numberOfRetries        int
	millisecondsToRetry    int
	hydraServers           []string
	enableAppRefresh       bool
	enableHydraRefresh     bool
	// connectionTimeout int
}

var factory *hydraClientFactory = &hydraClientFactory{
	AppsMonitorInstantiator:  new(appsMonitorInstantiator),
	ClientInstantiator:       new(clientInstantiator),
	HydraMonitorInstantiator: new(hydraMonitorInstantiator),

	hydraServerRefreshTime: default_hydra_server_refresh,
	hydraAppsRefreshTime:   default_hydra_apps_refresh,
	numberOfRetries:        default_retries_number,
	millisecondsToRetry:    0,
	hydraServers:           []string{},
	enableAppRefresh:       true,
	enableHydraRefresh:     true,
	// connectionTimeout:      1000,
}

// TODO: Maybe create an error type
func Config(hydraServers []string) (*hydraClientFactory, error) {
	if hydraServers == nil {
		return nil, errors.New("Illegal Argument: hydraServers can not be nil")
	}
	if len(hydraServers) == 0 {
		return nil, errors.New("Illegal Argument: hydraServers can not be empty")
	}
	factory.hydraServers = hydraServers
	return factory, nil
}

func (h *hydraClientFactory) Build() Client {
	if h.hydraClient != nil {
		return h.hydraClient
	}

	h.hydraClient = h.ClientInstantiator.MakeClient(h.hydraServers)
	h.hydraClient.SetMaxNumberOfRetries(h.numberOfRetries)
	h.hydraClient.SetWaitBetweenAllServersRetry(h.millisecondsToRetry)
	h.hydraClient.ReloadHydraServiceCache()
	// TODO: set connection timeout
	h.configureCacheRefreshMonitors()

	return h.hydraClient
}

func (h *hydraClientFactory) configureCacheRefreshMonitors() {
	if h.enableHydraRefresh {
		h.hydraMonitor = h.HydraMonitorInstantiator.MakeHydraMonitor(h.hydraClient, h.hydraServerRefreshTime)
		h.hydraMonitor.Run()
	}
	if h.enableAppRefresh {
		h.appsMonitor = h.AppsMonitorInstantiator.MakeAppsMonitor(h.hydraClient, h.hydraAppsRefreshTime)
		h.appsMonitor.Run()
	}
}

func GetHydraClient() Client {
	return factory.getHydraClient()
}

func (h *hydraClientFactory) getHydraClient() Client {
	return h.hydraClient
}

func Reset() {
	factory.AppsMonitorInstantiator = new(appsMonitorInstantiator)
	factory.ClientInstantiator = new(clientInstantiator)
	factory.HydraMonitorInstantiator = new(hydraMonitorInstantiator)

	factory.hydraClient = nil
	factory.hydraMonitor = nil
	factory.appsMonitor = nil
	factory.hydraServerRefreshTime = default_hydra_server_refresh
	factory.hydraAppsRefreshTime = default_hydra_apps_refresh
	factory.enableAppRefresh = true
	factory.enableHydraRefresh = true
}

func (h *hydraClientFactory) WithHydraCacheRefreshTime(timeoutSeconds int) *hydraClientFactory {
	h.hydraServerRefreshTime = time.Duration(timeoutSeconds) * time.Second
	return h
}

func (h *hydraClientFactory) WithAppsCacheRefreshTime(timeoutSeconds int) *hydraClientFactory {
	h.hydraAppsRefreshTime = time.Duration(timeoutSeconds) * time.Second
	return h
}

func (h *hydraClientFactory) AndAppsCacheRefreshTime(timeoutSeconds int) *hydraClientFactory {
	return h.WithAppsCacheRefreshTime(timeoutSeconds)
}

func (h *hydraClientFactory) AndHydraRefreshTime(timeoutSeconds int) *hydraClientFactory {
	return h.WithHydraCacheRefreshTime(timeoutSeconds)
}

func (h *hydraClientFactory) WithNumberOfRetries(numberOfRetries int) *hydraClientFactory {
	h.numberOfRetries = numberOfRetries
	return h
}

func (h *hydraClientFactory) AndNumberOfRetries(numberOfRetries int) *hydraClientFactory {
	return h.WithNumberOfRetries(numberOfRetries)
}

func (h *hydraClientFactory) WaitBetweenAllServersRetry(millisecondsToRetry int) *hydraClientFactory {
	h.millisecondsToRetry = millisecondsToRetry
	return h
}

func (h *hydraClientFactory) WithoutAppsRefresh() *hydraClientFactory {
	h.enableAppRefresh = false
	return h
}

func (h *hydraClientFactory) WithoutHydraServerRefresh() *hydraClientFactory {
	h.enableHydraRefresh = false
	return h
}

func (h *hydraClientFactory) AndWithoutHydraServerRefresh() *hydraClientFactory {
	return h.WithoutHydraServerRefresh()
}

func (h *hydraClientFactory) AndWithoutAppsRefresh() *hydraClientFactory {
	return h.WithoutAppsRefresh()
}
