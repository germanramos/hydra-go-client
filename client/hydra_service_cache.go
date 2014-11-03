package client

type HydraCache interface {
	GetHydraServers() []string
	Refresh(newHydraServers []string)
}

type HydraServiceCache struct {
	hydraSeedServers []string
	hydraServers     []string
}

func NewHydraServiceCache(hydraServers []string) *HydraServiceCache {
	return &HydraServiceCache{
		hydraSeedServers: hydraServers,
		hydraServers:     hydraServers,
	}
}

func (h *HydraServiceCache) GetHydraServers() []string {
	return h.hydraServers
}

func (h *HydraServiceCache) Refresh(newHydraServers []string) {
	if len(newHydraServers) > 0 {
		// TODO: Lock
		h.hydraServers = newHydraServers
	} else if missingServers := h.getMissingSeedServers(); len(missingServers) > 0 {
		for _, server := range missingServers {
			h.hydraServers = append(h.hydraServers, server)
		}
	}
}

func (h *HydraServiceCache) getMissingSeedServers() []string {
	missingServers := []string{}
	for _, seedServer := range h.hydraSeedServers {
		for i := 0; i < len(h.hydraServers); i++ {
			if h.hydraServers[i] == seedServer {
				break
			} else if i == len(h.hydraServers)-1 {
				missingServers = append(missingServers, seedServer)

			}
		}
	}
	return missingServers
}
