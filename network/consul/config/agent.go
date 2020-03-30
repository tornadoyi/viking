package config

import _consul "github.com/hashicorp/consul/api"

type AgentServiceRegistrationConfig struct {
	Kind				*ServiceKind               		`yaml:"kind"`
	ID					*string                    				`yaml:"id"`
	Name   				*string                    				`yaml:"name"`
	Tags   				[]string                  				`yaml:"tags"`
	Port   				*int                       				`yaml:"port"`
	Address				*string                    				`yaml:"address"`
	TaggedAddresses   	map[string]ServiceAddressConfig 		`yaml:"tagged_addresses"`
	EnableTagOverride 	*bool                     				`yaml:"enable_tag_override"`
	Meta              	map[string]string         				`yaml:"meta"`
	Weights           	*AgentWeightsConfig             		`yaml:"weights"`
	Check             	*AgentServiceCheckConfig				`yaml:"check"`
	Checks            	[]*AgentServiceCheckConfig				`yaml:"checks"`
	//Proxy             	*_consul.AgentServiceConnectProxyConfig `yaml:"proxy"`
	//Connect           	*_consul.AgentServiceConnect            `yaml:"Connect"`
	Namespace         	*string                          		`yaml:"namespace"`
}

func (h *AgentServiceRegistrationConfig) AgentServiceRegistration() *AgentServiceRegistration {
	c := &AgentServiceRegistration{}

	if h.Kind != nil { c.Kind = *h.Kind }
	if h.ID != nil { c.ID = *h.ID }
	if h.Name != nil { c.Name = *h.Name }
	if len(h.Tags) != 0 { c.Tags = h.Tags }
	if h.Port != nil { c.Port = *h.Port }
	if h.Address != nil { c.Address = *h.Address }
	if h.TaggedAddresses != nil {
		c.TaggedAddresses = make(map[string]ServiceAddress, 0)
		for k, v := range h.TaggedAddresses {
			c.TaggedAddresses[k] = *v.ServiceAddress()
		}
	}
	if h.EnableTagOverride != nil { c.EnableTagOverride = *h.EnableTagOverride }
	if h.Meta != nil { c.Meta = h.Meta }
	if h.Weights != nil { c.Weights = h.Weights.AgentWeights() }
	if h.Check != nil { c.Check = h.Check.AgentServiceCheck() }
	if h.Checks != nil {
		checks := make([]*AgentServiceCheck, 0, len(h.Checks))
		for _, c := range h.Checks {
			checks = append(checks, c.AgentServiceCheck())
		}
		c.Checks = checks
	}
	if h.Namespace != nil { c.Namespace = *h.Namespace }


	return c
}


type AgentServiceCheckConfig struct {
	CheckID           									*string              					`yaml:"check_id"`
	Name              									*string              					`yaml:"name"`
	Args              									[]string            					`yaml:"args"`
	DockerContainerID 									*string              					`yaml:"docker_container_id"`
	Shell             									*string              					`yaml:"shell"` // Only supported for Docker.
	Interval          									*string              					`yaml:"interval"`
	Timeout           									*string              					`yaml:"timeout"`
	TTL               									*string              					`yaml:"ttl"`
	HTTP              									*string              					`yaml:"http"`
	Header            									map[string][]string 					`yaml:"header"`
	Method            									*string              					`yaml:"method"`
	Body              									*string              					`yaml:"body"`
	TCP               									*string              					`yaml:"tcp"`
	Status            									*string              					`yaml:"status"`
	Notes             									*string              					`yaml:"notes"`
	TLSSkipVerify     									*bool                					`yaml:"tls_skip_verify"`
	GRPC              									*string              					`yaml:"grpc"`
	GRPCUseTLS        									*bool                					`yaml:"grpc_use_tls"`
	AliasNode         									*string              					`yaml:"alias_node"`
	AliasService      									*string              					`yaml:"alias_service"`
	DeregisterCriticalServiceAfter 						*string 						`yaml:"deregister_critical_service_after"`
}

func (h *AgentServiceCheckConfig) AgentServiceCheck() *AgentServiceCheck {
	c := &AgentServiceCheck{}
	if h.CheckID != nil { c.CheckID = *h.CheckID }
	if h.Name != nil { c.Name = *h.Name }
	if h.Args != nil { c.Args = h.Args }
	if h.DockerContainerID != nil { c.DockerContainerID = *h.DockerContainerID }
	if h.Shell != nil { c.Shell = *h.Shell }
	if h.Interval != nil { c.Interval = *h.Interval }
	if h.Timeout != nil { c.Timeout = *h.Timeout }
	if h.TTL != nil { c.TTL = *h.TTL }
	if h.HTTP != nil { c.HTTP = *h.HTTP }
	if h.Header != nil { c.Header = h.Header }
	if h.Method != nil { c.Method = *h.Method }
	if h.Body != nil { c.Body = *h.Body }
	if h.TCP != nil { c.TCP = *h.TCP }
	if h.Status != nil { c.Status = *h.Status }
	if h.Notes != nil { c.Notes = *h.Notes }
	if h.TLSSkipVerify != nil { c.TLSSkipVerify = *h.TLSSkipVerify }
	if h.GRPC != nil { c.GRPC = *h.GRPC }
	if h.GRPCUseTLS != nil { c.GRPCUseTLS = *h.GRPCUseTLS }
	if h.AliasNode != nil { c.AliasNode = *h.AliasNode }
	if h.AliasService != nil { c.AliasService = *h.AliasService }
	if h.DeregisterCriticalServiceAfter != nil { c.DeregisterCriticalServiceAfter = *h.DeregisterCriticalServiceAfter }

	return c
}



type AgentServiceRegistration = _consul.AgentServiceRegistration
type AgentServiceCheck = _consul.AgentServiceCheck