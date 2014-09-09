package client

import (
	"time"
)

type HydraCacheMonitor struct {
	controller   chan string
	hydraClient  HydraClient
	running      bool
	timeInterval time.Duration
}

func NewHydraCacheMonitor(hydraClient HydraClient, refreshInterval time.Duration) *HydraCacheMonitor {
	return &HydraCacheMonitor{
		hydraClient:  hydraClient,
		running:      false,
		timeInterval: refreshInterval,
	}
}

func (h *HydraCacheMonitor) GetInterval() time.Duration {
	return h.timeInterval
}

func (h *HydraCacheMonitor) IsRunning() bool {
	return h.running
}

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

func (h *HydraCacheMonitor) Stop() {
	h.controller <- "stop"
}
