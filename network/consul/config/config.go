package config

import (
	_consul "github.com/hashicorp/consul/api"
	"time"
)

type ConnectConfig struct {
	Address 				*string									`yaml:"address"`
	Scheme 					*string									`yaml:"scheme"`
	Datacenter 				*string									`yaml:"data_center"`
	WaitTime 				*string									`yaml:"wait_time"`
	Token 					*string									`yaml:"token"`
	TokenFile 				*string									`yaml:"token_file"`
	Namespace 				*string									`yaml:"namespace"`
}

func (h *ConnectConfig) Config() *_consul.Config{
	c := _consul.DefaultConfig()
	if h.Address != nil { c.Address = *h.Address }
	if h.Scheme != nil { c.Scheme = *h.Scheme }
	if h.Datacenter != nil { c.Datacenter = *h.Datacenter }
	if h.WaitTime != nil { c.WaitTime, _ = time.ParseDuration(*h.WaitTime)}
	if h.Token != nil { c.Token = *h.Token }
	if h.TokenFile != nil { c.TokenFile = *h.TokenFile }
	if h.Namespace != nil { c.Namespace = *h.Namespace }

	return c
}



type ServiceAddressConfig struct {
	Address 			string                          		`yaml:"address"`
	Port    		 	int                          			`yaml:"port"`
}

func (h* ServiceAddressConfig) ServiceAddress() *ServiceAddress {
	return &ServiceAddress{
		h.Address,
		h.Port,
	}
}



type AgentWeightsConfig struct {
	Passing 		 	int										`yaml:"passing"`
	Warning 		 	int										`yaml:"warning"`
}

func (h* AgentWeightsConfig) AgentWeights() *_consul.AgentWeights {
	return &_consul.AgentWeights{
		h.Passing,
		h.Warning,
	}
}


type ServiceKind = _consul.ServiceKind
type ServiceAddress = _consul.ServiceAddress