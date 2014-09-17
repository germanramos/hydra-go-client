package client

import (
	"time"
)

type AppsCacheMonitor struct {
	AbstractCacheMonitor
}

func NewAppsCacheMonitor(hydraClient HydraClient, refreshInterval time.Duration) *AppsCacheMonitor {
	a := new(AppsCacheMonitor)
	a.hydraClient = hydraClient
	a.running = false
	a.timeInterval = refreshInterval
	return a
}

func (a *AppsCacheMonitor) Run() {
	a.controller = make(chan string, 1)
	a.running = true
	a.hydraClient.ReloadAppServers()
	go func() {
	OuterLoop:
		for {
			select {
			case <-a.controller:
				break OuterLoop
			case <-time.After(a.timeInterval):
				a.hydraClient.ReloadAppServers()
			}
		}
		a.running = false
	}()
}
