package client

import (
	"time"
)

type AppsCacheMonitor struct {
	controller   chan string
	hydraClient  HydraClient
	running      bool
	timeInterval time.Duration
}

func NewAppsCacheMonitor(hydraClient HydraClient, refreshInterval time.Duration) *AppsCacheMonitor {
	return &AppsCacheMonitor{
		hydraClient:  hydraClient,
		running:      false,
		timeInterval: refreshInterval,
	}
}

func (a *AppsCacheMonitor) GetInterval() time.Duration {
	return a.timeInterval
}

func (a *AppsCacheMonitor) IsRunning() bool {
	return a.running
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

func (a *AppsCacheMonitor) Stop() {
	a.controller <- "stop"
}
