package config

import (
	_consul "github.com/hashicorp/consul/api"
	"time"
)

type CatalogRegistrationConfig struct {
	ID              					*string									`yaml:"id"`
	Node            					*string									`yaml:"node"`
	Address         					*string									`yaml:"address"`
	TaggedAddresses 					map[string]string						`yaml:"tagged_addresses"`
	NodeMeta        					map[string]string						`yaml:"node_meta"`
	Datacenter      					*string									`yaml:"datacenter"`
	Service         					*AgentServiceConfig						`yaml:"service"`
	Check           					*AgentCheckConfig						`yaml:"check"`
	Checks          					[]*HealthCheckConfig					`yaml:"checks"`
	SkipNodeUpdate  					*bool									`yaml:"skip_node_update"`
}

func (h *CatalogRegistrationConfig) CatalogRegistration() *CatalogRegistration {
	c := &CatalogRegistration{}
	if h.ID != nil { c.ID = *h.ID }
	if h.Node != nil { c.Node = *h.Node }
	if h.Address != nil { c.Address = *h.Address }
	if h.TaggedAddresses != nil { c.TaggedAddresses = h.TaggedAddresses }
	if h.NodeMeta != nil { c.NodeMeta = h.NodeMeta }
	if h.Datacenter != nil { c.Datacenter = *h.Datacenter }
	if h.Service != nil { c.Service = h.Service.AgentService() }
	if h.Check != nil { c.Check = h.Check.AgentCheck() }
	if h.Checks != nil {
		checks := make([]*HealthCheck, 0, len(h.Checks))
		for _, c := range h.Checks {
			checks = append(checks, c.HealthCheck())
		}
		c.Checks = checks
	}
	if h.SkipNodeUpdate != nil { c.SkipNodeUpdate = *h.SkipNodeUpdate }

	return c
}


type AgentServiceConfig struct {
	Kind              					*ServiceKind 							`yaml:"kind"`
	ID                					*string 								`yaml:"id"`
	Service           					*string 								`yaml:"service"`
	Tags              					[]string 								`yaml:"tags"`
	Meta              					map[string]string 						`yaml:"meta"`
	Port              					*int 									`yaml:"port"`
	Address           					*string									`yaml:"address"`
	TaggedAddresses   					map[string]*ServiceAddressConfig 		`yaml:"tagged_addresses"`
	Weights           					*AgentWeightsConfig						`yaml:"weights"`
	EnableTagOverride 					*bool									`yaml:"enable_tag_override"`
	CreateIndex       					*uint64 								`yaml:"create_index"`
	ModifyIndex       					*uint64 								`yaml:"modify_index"`
	ContentHash       					*string    								`yaml:"content_hash"`
	Namespace 							*string 								`yaml:"namespace"`
}

func (h *AgentServiceConfig) AgentService() *AgentService {
	c := &AgentService{}
	if h.Kind != nil { c.Kind = *h.Kind }
	if h.ID != nil { c.ID = *h.ID }
	if h.Service != nil { c.Service = *h.Service }
	if h.Tags != nil { c.Tags = h.Tags }
	if h.Meta != nil { c.Meta = h.Meta }
	if h.Port != nil { c.Port = *h.Port }
	if h.Address != nil { c.Address = *h.Address }
	if h.TaggedAddresses != nil {
		c.TaggedAddresses = make(map[string]ServiceAddress, 0)
		for k, v := range h.TaggedAddresses {
			c.TaggedAddresses[k] = *v.ServiceAddress()
		}
	}
	if h.Weights != nil { c.Weights = *h.Weights.AgentWeights() }
	if h.EnableTagOverride != nil { c.EnableTagOverride = *h.EnableTagOverride }
	if h.CreateIndex != nil { c.CreateIndex = *h.CreateIndex }
	if h.ModifyIndex != nil { c.ModifyIndex = *h.ModifyIndex }
	if h.ContentHash != nil { c.ContentHash = *h.ContentHash }
	if h.Namespace != nil { c.Namespace = *h.Namespace }

	return c
}


type AgentCheckConfig struct {
	Node        						*string 								`yaml:"node"`
	CheckID     						*string 								`yaml:"check_id"`
	Name        						*string 								`yaml:"name"`
	Status      						*string 								`yaml:"status"`
	Notes       						*string 								`yaml:"notes"`
	Output      						*string 								`yaml:"output"`
	ServiceID   						*string 								`yaml:"serviceID"`
	ServiceName 						*string 								`yaml:"serviceName"`
	Type        						*string 								`yaml:"type"`
	Definition  						*HealthCheckDefinitionConfig 			`yaml:"definition"`
	Namespace   						*string 								`yaml:"namespace"`
}

func (h *AgentCheckConfig) AgentCheck() *AgentCheck {
	c := &AgentCheck{}
	if h.Node != nil { c.Node = *h.Node }
	if h.CheckID != nil { c.CheckID = *h.CheckID }
	if h.Name != nil { c.Name = *h.Name }
	if h.Status != nil { c.Status = *h.Status }
	if h.Notes != nil { c.Notes = *h.Notes }
	if h.Output != nil { c.Output = *h.Output }
	if h.ServiceID != nil { c.ServiceID = *h.ServiceID }
	if h.ServiceName != nil { c.ServiceName = *h.ServiceName }
	if h.Type != nil { c.Type = *h.Type }
	if h.Definition != nil { c.Definition = *h.Definition.HealthCheckDefinition() }
	if h.Namespace != nil { c.Namespace = *h.Namespace }
	return c
}


type HealthCheckConfig struct {
	Node        						*string 								`yaml:"node"`
	CheckID     						*string 								`yaml:"check_id"`
	Name        						*string 								`yaml:"name"`
	Status      						*string 								`yaml:"status"`
	Notes       						*string 								`yaml:"notes"`
	Output      						*string 								`yaml:"output"`
	ServiceID   						*string 								`yaml:"service_id"`
	ServiceName 						*string 								`yaml:"service_name"`
	ServiceTags 						[]string 								`yaml:"service_tags"`
	Type        						*string 								`yaml:"type"`
	Namespace   						*string 								`yaml:"namespace"`
	Definition 							*HealthCheckDefinitionConfig			`yaml:"definition"`
	CreateIndex 						*uint64									`yaml:"create_index"`
	ModifyIndex 						*uint64									`yaml:"modify_index"`
}

func (h *HealthCheckConfig) HealthCheck() *HealthCheck {
	c := &HealthCheck{}
	if h.Node != nil { c.Node = *h.Node }
	if h.CheckID != nil { c.CheckID = *h.CheckID }
	if h.Name != nil { c.Name = *h.Name }
	if h.Status != nil { c.Status = *h.Status }
	if h.Notes != nil { c.Notes = *h.Notes }
	if h.Output != nil { c.Output = *h.Output }
	if h.ServiceID != nil { c.ServiceID = *h.ServiceID }
	if h.ServiceName != nil { c.ServiceName = *h.ServiceName }
	if h.ServiceTags != nil { c.ServiceTags = h.ServiceTags }
	if h.Type != nil { c.Type = *h.Type }
	if h.Namespace != nil { c.Namespace = *h.Namespace }
	if h.Definition != nil { c.Definition = *h.Definition.HealthCheckDefinition() }
	if h.CreateIndex != nil { c.CreateIndex = *h.CreateIndex }
	if h.ModifyIndex != nil { c.ModifyIndex = *h.ModifyIndex }
	return c
}


type HealthCheckDefinitionConfig struct {
	HTTP                                   	*string								`yaml:"http"`
	Header                                 	*map[string][]string				`yaml:"header"`
	Method                                 	*string								`yaml:"method"`
	Body                                   	*string								`yaml:"body"`
	TLSSkipVerify                          	*bool								`yaml:"tls_skip_verify"`
	TCP                                    	*string								`yaml:"tcp"`
	IntervalDuration                       	*string								`yaml:"interval_duration"`
	TimeoutDuration                        	*string								`yaml:"timeout_duration"`
	DeregisterCriticalServiceAfterDuration 	*string								`yaml:"deregister_critical_service_after_duration"`
}


func (h *HealthCheckDefinitionConfig) HealthCheckDefinition() *HealthCheckDefinition {
	c := &HealthCheckDefinition{}
	if h.HTTP != nil { c.HTTP = *h.HTTP }
	if h.Header != nil { c.Header = *h.Header }
	if h.Method != nil { c.Method = *h.Method }
	if h.Body != nil { c.Body = *h.Body }
	if h.TLSSkipVerify != nil { c.TLSSkipVerify = *h.TLSSkipVerify }
	if h.TCP != nil { c.TCP = *h.TCP }
	if h.IntervalDuration != nil { c.IntervalDuration, _ = time.ParseDuration(*h.IntervalDuration) }
	if h.TimeoutDuration != nil { c.TimeoutDuration, _ = time.ParseDuration(*h.TimeoutDuration) }
	if h.DeregisterCriticalServiceAfterDuration != nil { c.DeregisterCriticalServiceAfterDuration, _ = time.ParseDuration(*h.DeregisterCriticalServiceAfterDuration) }
	return c
}


type CatalogDeregistrationConfig struct {
	Node       							*string								`yaml:"node"`
	Address    							*string								`yaml:"address"`
	Datacenter 							*string								`yaml:"datacenter"`
	ServiceID  							*string								`yaml:"service_id"`
	CheckID    							*string								`yaml:"check_id"`
	Namespace  							*string 							`yaml:"namespace"`
}

func (h *CatalogDeregistrationConfig) CatalogDeregistration() *CatalogDeregistration {
	c := &CatalogDeregistration{}
	if h.Node != nil { c.Node = *h.Node }
	if h.Address != nil { c.Address = *h.Address }
	if h.Datacenter != nil { c.Datacenter = *h.Datacenter }
	if h.ServiceID != nil { c.ServiceID = *h.ServiceID }
	if h.CheckID != nil { c.CheckID = *h.CheckID }
	if h.Namespace != nil { c.Namespace = *h.Namespace }
	return c
}



type CatalogRegistration = _consul.CatalogRegistration
type CatalogDeregistration = _consul.CatalogDeregistration
type HealthCheckDefinition = _consul.HealthCheckDefinition
type HealthCheck = _consul.HealthCheck
type AgentCheck = _consul.AgentCheck
type AgentService = _consul.AgentService