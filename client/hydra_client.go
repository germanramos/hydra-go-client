package client

// import (
// 	"errors"
// 	"sync"
// 	"time"
// )

// type HydraClient interface {
// 	Get(appId string, forceCacheRefresh bool) ([]string, error)
// 	ReloadAppServers()
// 	ReloadHydraServers()
// 	SetAppsCacheMonitor(monitor *AppsCacheMonitor)
// 	SetHydraCacheMonitor(monitor *HydraCacheMonitor)
// 	SetMaxNumberOfRetriesPerHydraServer(numberOfRetries uint)
// 	SetWaitBetweenAllServersRetry(duration time.Duration)
// }

// const (
// 	AppRootPath string = "/app/"
// 	HydraAppId  string = "hydra"
// )

// type Client struct {
// 	sync.RWMutex
// 	appsCacheMonitor           *AppsCacheMonitor
// 	appServers                 map[string][]string
// 	hydraAvailable             bool
// 	hydraCacheMonitor          *HydraCacheMonitor
// 	hydraServers               []string
// 	hydraServersRequester      Requester
// 	maxNumberOfRetries         uint
// 	waitBetweenAllServersRetry time.Duration
// }

// func NewClient(hydraServers []string, requester Requester) *Client {
// 	return &Client{
// 		appServers:                 make(map[string][]string),
// 		hydraServers:               hydraServers,
// 		hydraServersRequester:      requester,
// 		maxNumberOfRetries:         0,
// 		waitBetweenAllServersRetry: time.Duration(0) * time.Millisecond,
// 	}
// }

// // Get returns to public api of one hydra server the urls
// // for the available servers for one application
// func (c *Client) Get(appId string, forceCacheRefresh bool) ([]string, error) {
// 	if len(appId) == 0 {
// 		return []string{}, errors.New("Invalid Argument: appId must be a single word")
// 	}

// 	if servers, ok := c.appServers[appId]; ok && !forceCacheRefresh {
// 		return servers, nil
// 	}

// 	return c.requestCandidateRefreshingCache(appId)
// }

// // requestCandidateRefreshingCache gets the candidate servers for one application,
// // refreshes the associated cache and returns the urls of servers
// func (c *Client) requestCandidateRefreshingCache(appId string) ([]string, error) {
// 	servers, err := c.requestCandidateServers(appId)
// 	if err != nil {
// 		return servers, err
// 	}
// 	c.refreshCache(appId, servers)

// 	return servers, nil
// }

// func (c *Client) refreshCache(appId string, servers []string) {
// 	c.Lock()
// 	c.appServers[appId] = servers
// 	c.Unlock()
// }

// func (c *Client) requestCandidateServers(appId string) ([]string, error) {
// 	var retries uint = 0
// 	var numberOfHydraServers uint = c.getNumberOfHydraServers()
// 	var totalNumberOfRetries uint = c.maxNumberOfRetries * numberOfHydraServers

// 	var currentHydraServerIndex uint = 0
// 	for c.maxNumberOfRetries == 0 || retries < totalNumberOfRetries {
// 		servers, err := c.hydraServersRequester.GetCandidateServers(c.hydraServers[currentHydraServerIndex]+AppRootPath, appId)
// 		if err == nil {
// 			return servers, nil
// 		} else {
// 			if currentHydraServerIndex+1 < c.getNumberOfHydraServers() {
// 				currentHydraServerIndex++
// 			} else {
// 				currentHydraServerIndex = 0
// 			}
// 			retries++
// 		}

// 		if retries%numberOfHydraServers == 0 {
// 			c.waitUntilTheNextRetry()
// 		}
// 	}

// 	return []string{}, errors.New("None Servers Accessible")
// }

// func (c *Client) waitUntilTheNextRetry() {
// 	time.Sleep(c.waitBetweenAllServersRetry)
// }

// func (c *Client) getNumberOfHydraServers() uint {
// 	c.RLock()
// 	defer c.RUnlock()
// 	return uint(len(c.hydraServers))
// }

// func (c *Client) getApplicationIds() []string {
// 	c.Lock()
// 	defer c.Unlock()

// 	var apps = []string{}
// 	for key := range c.appServers {
// 		apps = append(apps, key)
// 	}
// 	return apps
// }

// // ReloadAppServers reloads the application cache for all applications registered
// func (c *Client) ReloadAppServers() {
// 	c.refreshAppCache(c.retrieveNewServerConfiguration())
// }

// func (c *Client) refreshAppCache(newAppServers map[string][]string) {
// 	c.Lock()
// 	defer c.Unlock()

// 	c.appServers = newAppServers
// }

// func (c *Client) retrieveNewServerConfiguration() map[string][]string {
// 	var apps []string = c.getApplicationIds()
// 	var appsServersCache map[string][]string = make(map[string][]string)

// 	for _, appId := range apps {
// 		servers, err := c.requestCandidateServers(appId)
// 		if err == nil {
// 			appsServersCache[appId] = servers
// 		} else {
// 			appsServersCache[appId] = []string{}
// 		}
// 	}

// 	return appsServersCache
// }

// // ReloadHydraServers reloads the cache of hydra servers
// func (c *Client) ReloadHydraServers() {
// 	servers, err := c.requestCandidateServers(HydraAppId)
// 	c.Lock()
// 	if err == nil {
// 		c.hydraServers = servers
// 		if len(servers) > 0 {
// 			c.hydraAvailable = true
// 		} else {
// 			c.hydraAvailable = false
// 		}
// 	} else {
// 		c.hydraAvailable = false
// 	}
// 	c.Unlock()
// }

// func (c *Client) SetHydraCacheMonitor(monitor *HydraCacheMonitor) {
// 	c.hydraCacheMonitor = monitor
// }

// func (c *Client) SetAppsCacheMonitor(monitor *AppsCacheMonitor) {
// 	c.appsCacheMonitor = monitor
// }

// // SetMax NumberOfRetriesPerHydraServer configures the maximum number of retries
// // per hydra server. If it is configured to zero the retries will be infinity.
// func (c *Client) SetMaxNumberOfRetriesPerHydraServer(numberOfRetries uint) {
// 	c.maxNumberOfRetries = numberOfRetries
// }

// // SetWaitBetweenAllServersRetry configures the waiting time until the next cycle
// // of retries throught all hydra servers registered
// func (c *Client) SetWaitBetweenAllServersRetry(duration time.Duration) {
// 	c.waitBetweenAllServersRetry = duration
// }

// // IsHydraAvailable shows if reachable hydra servers
// func (c *Client) IsHydraAvailable() bool {
// 	c.Lock()
// 	defer c.Unlock()
// 	return c.hydraAvailable
// }
