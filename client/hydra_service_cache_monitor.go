package client

import (
	"time"
)

type HydraServiceCacheMonitor struct {
	AbstractHydraCacheMonitor
}

func NewHydraServiceCacheMonitor(hydraClient Client, refreshTime time.Duration) *HydraServiceCacheMonitor {
	a := new(HydraServiceCacheMonitor)
	a.hydraClient = hydraClient
	a.refreshTime = refreshTime
	return a
}

// Run executes a goroutine that reload periodically the hydra service cache of the Hydra Client
func (a *HydraServiceCacheMonitor) Run() {
	a.hydraClient.ReloadHydraServiceCache()
	go func() {
		for {
			select {
			case <-time.After(a.refreshTime):
				a.hydraClient.ReloadHydraServiceCache()
			}
		}
	}()
}
