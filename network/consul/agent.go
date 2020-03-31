package consul

import (
	"fmt"
	_consul "github.com/hashicorp/consul/api"
	"github.com/tornadoyi/viking/log"
	"github.com/tornadoyi/viking/network/consul/config"
	"github.com/tornadoyi/viking/network/consul/resolver"
	"sync"
	"time"
)


type Agent struct {
	*_consul.Agent
	client	*Client
	timer				*time.Timer
	services			map[string]*AgentService
	mutex				*sync.RWMutex
}

func NewAgent(client *Client) *Agent {
	return &Agent{
		Agent:    client.client.Agent(),
		client:   client,
		services: make(map[string]*AgentService),
		mutex:    &sync.RWMutex{},
	}
}


func (h *Agent) Services() map[string]*AgentService {
	h.mutex.RLock()
	defer h.mutex.RUnlock()
	return h.services
}


func (h *Agent) RegisterService(cfg *config.AgentServiceRegistration) error {
	// check
	if cfg == nil { return fmt.Errorf("Empty registration config")}

	// deregister
	if err := h.Agent.ServiceDeregister(cfg.ID); err != nil { return err}

	// register
	err := h.Agent.ServiceRegister(cfg)
	if err != nil { return err }

	// start health checking server
	return listenHealthCheck(cfg.Check.HTTP, nil)

}

func (h *Agent) SetInterval(interval time.Duration) {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	if h.timer != nil { h.timer.Stop() }
	h.timer = time.AfterFunc(interval, func() {
		h.mutex.Lock()
		defer h.mutex.Unlock()
		defer h.timer.Reset(interval)
		services, err := h.Agent.Services()
		if err != nil { log.Error(err); return }
		h.services = services
	})
}


func (h *Agent) RegisterResolver(scheme string, interval time.Duration, filters... resolver.IAgentServiceFilter) error {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	// register
	return resolver.RegisterAgentResolver(h, scheme, interval, filters...)
}




type AgentService = _consul.AgentService