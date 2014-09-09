package client

import (
	"time"
)

type CacheMonitor interface {
	GetInterval() time.Duration
	IsRunning() bool
	Run()
	Stop()
}
