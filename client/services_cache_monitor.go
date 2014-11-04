package client

import (
	"time"
)

type ServicesCacheMonitor struct {
	AbstractHydraCacheMonitor
}

func NewServicesCacheMonitor(hydraClient Client, refreshTime time.Duration) *ServicesCacheMonitor {
	a := new(ServicesCacheMonitor)
	a.hydraClient = hydraClient
	a.refreshTime = refreshTime
	return a
}

// Run executes a goroutine that reload periodically the service cache of the Hydra Client
func (a *ServicesCacheMonitor) Run() {
	a.hydraClient.ReloadServicesCache()
	go func() {
		for {
			select {
			case <-time.After(a.refreshTime):
				a.hydraClient.ReloadServicesCache()
			}
		}
	}()
}
