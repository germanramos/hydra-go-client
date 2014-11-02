package client

// import (
// 	"time"
// )

// type CacheMonitor interface {
// 	GetInterval() time.Duration
// 	IsRunning() bool
// 	Run()
// 	Stop()
// }

// type AbstractCacheMonitor struct {
// 	controller   chan string
// 	hydraClient  HydraClient
// 	running      bool
// 	timeInterval time.Duration
// }

// func (a *AbstractCacheMonitor) GetInterval() time.Duration {
// 	return a.timeInterval
// }

// func (a *AbstractCacheMonitor) IsRunning() bool {
// 	return a.running
// }

// func (a *AbstractCacheMonitor) Run() {
// }

// func (a *AbstractCacheMonitor) Stop() {
// 	a.controller <- "stop"
// }
