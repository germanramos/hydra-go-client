package client

type CacheMonitor interface {
	IsRunning() bool
	Run()
	Stop()
}
