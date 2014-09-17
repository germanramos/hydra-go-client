package client

import (
	"time"
)

type HydraCacheMonitor struct {
	AbstractCacheMonitor
}

func NewHydraCacheMonitor(hydraClient HydraClient, refreshInterval time.Duration) *HydraCacheMonitor {
	h := new(HydraCacheMonitor)
	h.hydraClient = hydraClient
	h.running = false
	h.timeInterval = refreshInterval
	return h
}

// Run executes a coroutine that reload periodically the hydra cache of the Hydra Client
func (h *HydraCacheMonitor) Run() {
	h.controller = make(chan string, 1)
	h.running = true
	h.hydraClient.ReloadHydraServers()
	go func() {
	OuterLoop:
		for {
			select {
			case <-h.controller:
				break OuterLoop
			case <-time.After(h.timeInterval):
				h.hydraClient.ReloadHydraServers()
			}
		}
		h.running = false
	}()
}
