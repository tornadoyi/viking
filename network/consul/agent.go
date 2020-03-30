package consul

import (
	"fmt"
	_consul "github.com/hashicorp/consul/api"
	"github.com/tornadoyi/viking/log"
	"github.com/tornadoyi/viking/network/consul/config"
	"sync"
	"time"
)


type Agent struct {
	*_consul.Agent
	client	*Client
	timer				*time.Timer
	focus				map[string]bool
	services			map[string][]*AgentService
	mutex				*sync.RWMutex
}

func NewAgent(client *Client) *Agent {
	return &Agent{
		Agent:    client.client.Agent(),
		client:   client,
		focus:    make(map[string]bool),
		services: make(map[string][]*AgentService),
		mutex:    &sync.RWMutex{},
	}
}


func (h *Agent) Services() map[string][]*AgentService {
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
		h.updateServices(services)
	})
}


func (h *Agent) RegisterResolver(scheme string, services []string,  interval time.Duration) error {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	// add service to focus
	for _, s := range services{
		h.focus[s] = true
	}

	// register
	return RegisterAgentResolver(h, scheme, services, interval)
}



func (h *Agent) updateServices(services map[string]*AgentService) {
	focusServices := make(map[string][]*AgentService)
	for _, s := range services {
		if _, ok := h.focus[s.Service]; !ok { continue }
		svrs, _ := focusServices[s.Service]
		focusServices[s.Service] = append(svrs, s)
	}
	h.services = focusServices
}









type AgentService = _consul.AgentService