package consul

import (
	"errors"
	"fmt"
	"github.com/tornadoyi/viking/log"
	_resolver "github.com/tornadoyi/viking/network/grpc/resolver"
	"strings"
	"sync"
	"time"
)



func RegisterAgentResolver(agent *Agent, scheme string, services []string,  interval time.Duration) error{
	if _resolver.Get(scheme) != nil { return fmt.Errorf("Repeated resolver %v", scheme)}
	r := NewAgentResolverBuilder(agent, scheme, services, interval)
	_resolver.Register(r)
	return nil
}



type AgentResolverBuilder struct {
	agent					*Agent
	scheme					string
	services				[]string
	interval				time.Duration
}

func NewAgentResolverBuilder(agent *Agent, scheme string, services []string,  interval time.Duration) *AgentResolverBuilder{
	return &AgentResolverBuilder{
		agent:    agent,
		scheme:   scheme,
		services: services,
		interval: interval,
	}
}

func (h *AgentResolverBuilder) Build(target _resolver.Target, cc _resolver.ClientConn, opts _resolver.BuildOptions) (_resolver.IResolver, error) {
	return NewAgentResolver(target, cc, h.agent, h.services, h.interval)
}
func (h *AgentResolverBuilder) Scheme() string { return h.scheme }

type AgentResolver struct {
	target     				_resolver.Target
	connection         		_resolver.ClientConn
	agent					*Agent
	services				[]string
	interval				time.Duration
	timer					*time.Timer
	mutex					sync.RWMutex
}

func NewAgentResolver(target _resolver.Target,  conn _resolver.ClientConn, agent *Agent, services []string, interval time.Duration) (_resolver.IResolver, error) {
	// create and update state
	r := &AgentResolver{
		target:     target,
		connection: conn,
		agent:     	agent,
		services:   services,
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


func (h *AgentResolver) updateState() error {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	address := make([]_resolver.Address, 0)
	services := h.agent.Services()

	errs := []string{}
	for _, name := range h.services {
		slist, ok := services[name]
		if !ok {
			errs = append(errs, fmt.Sprintf("Can not find agent service &v", name))
			continue
		}
		for _, s := range slist {
			address = append(address, _resolver.Address{
				Addr: fmt.Sprintf("%v:%v", s.Address, s.Port)})
		}
	}
	h.connection.UpdateState(_resolver.State{Addresses: address})

	if errs == nil { return nil}
	return errors.New(strings.Join(errs, "\n"))
}