package client

// import (
// 	"errors"
// 	"time"
// )

// type HydraClientBuilder interface {
// 	Build() HydraClient
// 	Config([]string) error
// }

// const (
// 	DefaultAppsCacheDuration              time.Duration = time.Duration(20000) * time.Millisecond
// 	DefaultDurationBetweenAllServersRetry time.Duration = time.Duration(0) * time.Millisecond
// 	DefaultHydraServersCacheDuration      time.Duration = time.Duration(60000) * time.Millisecond
// 	DefaultNumberOfRetries                uint          = 10
// )

// type hydraClientFactory struct {
// 	appsCacheDuration              time.Duration
// 	hydraServers                   []string
// 	hydraServersCacheDuration      time.Duration
// 	maxNumberOfRetries             uint
// 	durationBetweenAllServersRetry time.Duration
// }

// var HydraClientFactory *hydraClientFactory = &hydraClientFactory{
// 	appsCacheDuration:              DefaultAppsCacheDuration,
// 	hydraServers:                   []string{},
// 	hydraServersCacheDuration:      DefaultHydraServersCacheDuration,
// 	maxNumberOfRetries:             DefaultNumberOfRetries,
// 	durationBetweenAllServersRetry: DefaultDurationBetweenAllServersRetry,
// }

// func (h *hydraClientFactory) Config(hydraServers []string) error {
// 	if hydraServers == nil {
// 		return errors.New("Invalid Argument: hydraServers can not be nil")
// 	}
// 	if len(hydraServers) == 0 {
// 		return errors.New("Invalid Argument: hydraServers can not be empty")
// 	}
// 	h.hydraServers = hydraServers
// 	return nil
// }

// func (h *hydraClientFactory) Build() *Client {
// 	hydraClient := NewClient(h.hydraServers, NewHydraServersRequester())
// 	hydraClient.SetMaxNumberOfRetriesPerHydraServer(h.maxNumberOfRetries)
// 	hydraClient.SetWaitBetweenAllServersRetry(h.durationBetweenAllServersRetry)

// 	hydraClient.ReloadHydraServers()
// 	h.configureCacheMonitors(hydraClient)

// 	return hydraClient
// }

// func (h *hydraClientFactory) configureCacheMonitors(hydraClient *Client) {
// 	hydraCacheMonitor := NewHydraCacheMonitor(hydraClient, h.hydraServersCacheDuration)
// 	hydraCacheMonitor.Run()
// 	appsCacheMonitor := NewAppsCacheMonitor(hydraClient, h.appsCacheDuration)
// 	appsCacheMonitor.Run()

// 	hydraClient.SetHydraCacheMonitor(hydraCacheMonitor)
// 	hydraClient.SetAppsCacheMonitor(appsCacheMonitor)
// }

// func (h *hydraClientFactory) GetAppsCacheDuration() time.Duration {
// 	return h.appsCacheDuration
// }

// func (h *hydraClientFactory) GetHydraServersCacheDuration() time.Duration {
// 	return h.hydraServersCacheDuration
// }

// func (h *hydraClientFactory) GetMaxNumberOfRetriesPerHydraServer() uint {
// 	return h.maxNumberOfRetries
// }

// func (h *hydraClientFactory) GetDurationBetweenAllServersRetry() time.Duration {
// 	return h.durationBetweenAllServersRetry
// }

// func (h *hydraClientFactory) WithAppsCacheDuration(duration time.Duration) *hydraClientFactory {
// 	h.appsCacheDuration = duration
// 	return h
// }

// func (h *hydraClientFactory) WithHydraServersCacheDuration(duration time.Duration) *hydraClientFactory {
// 	h.hydraServersCacheDuration = duration
// 	return h
// }

// func (h *hydraClientFactory) WithMaxNumberOfRetriesPerHydraServer(retries uint) *hydraClientFactory {
// 	h.maxNumberOfRetries = retries
// 	return h
// }

// func (h *hydraClientFactory) WaitBetweenAllServersRetry(duration time.Duration) *hydraClientFactory {
// 	h.durationBetweenAllServersRetry = duration
// 	return h
// }
