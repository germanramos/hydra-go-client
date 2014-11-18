package client

import (
	"sync"
)

type ServiceCache interface {
	Exists(serviceId string) bool
	FindById(serviceId string) []string
	GetIds() []string
	PutService(serviceId string, servers []string)
	Refresh(newAppserverCache map[string][]string)
	RemoveServer(serviceId, serverToRemove string)
}

type ServicesCache struct {
	sync.RWMutex
	serviceCache map[string][]string
}

func NewServicesCache() *ServicesCache {
	return &ServicesCache{
		serviceCache: make(map[string][]string),
	}
}

func (s *ServicesCache) FindById(serviceId string) []string {
	s.RLock()
	defer s.RUnlock()
	if services, ok := s.serviceCache[serviceId]; ok {
		return services
	}
	return []string{}
}

func (s *ServicesCache) Exists(serviceId string) bool {
	s.RLock()
	defer s.RUnlock()
	_, exists := s.serviceCache[serviceId]
	return exists
}

func (s *ServicesCache) GetIds() []string {
	ids := []string{}
	s.RLock()
	defer s.RUnlock()
	for id, _ := range s.serviceCache {
		ids = append(ids, id)
	}
	return ids
}

func (s *ServicesCache) Refresh(newAppserverCache map[string][]string) {
	s.Lock()
	defer s.Unlock()
	s.serviceCache = newAppserverCache
}

func (s *ServicesCache) PutService(serviceId string, servers []string) {
	s.Lock()
	defer s.Unlock()
	s.serviceCache[serviceId] = servers
}

// Remove a server from a application, this method was called normally if the server fails.
func (s *ServicesCache) RemoveServer(serviceId, serverToRemove string) {
	s.Lock()
	defer s.Unlock()
	servers := s.serviceCache[serviceId]
	finalServers := []string{}
	for _, server := range servers {
		if server != serverToRemove {
			finalServers = append(finalServers, server)
		}
	}

	if len(finalServers) == 0 {
		delete(s.serviceCache, serviceId)
	} else {
		s.serviceCache[serviceId] = finalServers
	}
}
