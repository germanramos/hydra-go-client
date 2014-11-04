package client

import (
	"time"
)

type CacheMonitor interface {
	Run()
}

type AbstractHydraCacheMonitor struct {
	hydraClient Client
	refreshTime time.Duration
}

func (a *AbstractHydraCacheMonitor) GetInterval() time.Duration {
	return a.refreshTime
}

func (a *AbstractHydraCacheMonitor) Run() {}
