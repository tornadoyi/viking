package resolver

import (
	"fmt"
	_consul "github.com/hashicorp/consul/api"
	"github.com/tornadoyi/viking/goplus/runtime"
	"github.com/tornadoyi/viking/log"
	_resolver "github.com/tornadoyi/viking/network/grpc/resolver"
	"sync"
	"time"
)





func RegisterAgentResolver(agent IAgent, scheme string, interval time.Duration, filters... IAgentServiceFilter) error{
	if _resolver.Get(scheme) != nil { return fmt.Errorf("Repeated resolver %v", scheme)}
	r := NewAgentResolverBuilder(agent, scheme, interval, filters...)
	_resolver.Register(r)
	return nil
}



type AgentResolverBuilder struct {
	agent					IAgent
	scheme					string
	filters					[]IAgentServiceFilter
	interval				time.Duration
}

func NewAgentResolverBuilder(agent IAgent, scheme string,  interval time.Duration, filters... IAgentServiceFilter,) *AgentResolverBuilder{
	return &AgentResolverBuilder{
		agent:    agent,
		scheme:   scheme,
		filters:   filters,
		interval: interval,
	}
}

func (h *AgentResolverBuilder) Build(target _resolver.Target, cc _resolver.ClientConn, opts _resolver.BuildOptions) (_resolver.IResolver, error) {
	return NewAgentResolver(target, cc, h.agent, h.interval, h.filters...)
}
func (h *AgentResolverBuilder) Scheme() string { return h.scheme }

type AgentResolver struct {
	target     				_resolver.Target
	connection         		_resolver.ClientConn
	agent					IAgent
	filters					[]IAgentServiceFilter
	interval				time.Duration
	timer					*time.Timer
	mutex					sync.RWMutex
}

func NewAgentResolver(target _resolver.Target,  conn _resolver.ClientConn, agent IAgent, interval time.Duration, filters... IAgentServiceFilter) (_resolver.IResolver, error) {
	// create and update state
	r := &AgentResolver{
		target:     target,
		connection: conn,
		agent:     	agent,
		filters:	filters,
		interval:   interval,
		timer:      nil,
		mutex:      sync.RWMutex{},
	}
	err := r.init()
	if err != nil { return nil, err}
	return r, nil
}


func (*AgentResolver) ResolveNow(o _resolver.ResolveNowOptions) {}

func (*AgentResolver) Close() {}


func (h *AgentResolver) init() error{
	if err := h.updateState(); err != nil { return err}
	h.timer = time.AfterFunc(h.interval, func() {
		defer h.timer.Reset(h.interval)
		if err := h.updateState(); err != nil { log.Error(err) }
	})
	return nil
}


func (h *AgentResolver) updateState() (reterr error) {
	h.mutex.Lock()
	defer h.mutex.Unlock()
	defer runtime.CatchCallback(func(err error) { reterr = err })

	address := make([]_resolver.Address, 0)
	for _, s := range h.agent.Services() {
		for _, f := range h.filters {
			if f.Filter(h.target, s) { continue }
			address = append(address, _resolver.Address{Addr: fmt.Sprintf("%v:%v", s.Address, s.Port)})
		}
	}
	h.connection.UpdateState(_resolver.State{Addresses: address})

	return nil
}


type IAgentServiceFilter interface {
	Filter(_resolver.Target, *AgentService) bool
}


type AgentServiceNameFilter struct {
	services		map[string]bool
}

func NewAgentServiceNameFilter(services []string) *AgentServiceNameFilter {
	f := &AgentServiceNameFilter{ make(map[string]bool, len(services))}
	for _, name := range services { f.services[name] = true }
	return f
}

func (h *AgentServiceNameFilter) Filter(target _resolver.Target, s *AgentService) bool { if _, ok := h.services[s.Service]; ok { return false } else { return true } }


type IAgent interface {
	Services()		map[string]*AgentService
}


type AgentService	=	_consul.AgentService