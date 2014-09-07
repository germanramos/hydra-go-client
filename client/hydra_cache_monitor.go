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

func NewHydraCacheMonitor(hydraClient HydraClient, refreshInterval uint) *HydraCacheMonitor {
	return &HydraCacheMonitor{
		hydraClient:  hydraClient,
		running:      false,
		timeInterval: time.Duration(refreshInterval) * time.Millisecond,
	}
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

func (h *HydraCacheMonitor) IsRunning() bool {
	return h.running
}
