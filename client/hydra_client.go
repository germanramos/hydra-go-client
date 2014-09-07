package client

import ()

type HydraClient interface {
	Get()
	ReloadHydraServers()
}

const (
	AppRootPath string = "/app/"
	HydraAppId  string = "hydra"
)

type Client struct {
}
