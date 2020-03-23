
package keepalive

import (
	_grpc "google.golang.org/grpc"
	_keepalive "google.golang.org/grpc/keepalive"
	"time"
)


type ClientParameters = _keepalive.ClientParameters

type ServerParameters = _keepalive.ServerParameters

type EnforcementPolicy = _keepalive.EnforcementPolicy




type ServerParametersConfig struct {
	MaxConnectionIdle 			string 					`yaml:"max_connection_idle"`
	MaxConnectionAge 			string					`yaml:"max_connection_age"`
	MaxConnectionAgeGrace 		string 					`yaml:"max_connection_age_grace"`
	Time 						string 					`yaml:"time"`
	Timeout 					string 					`yaml:"timeout"`
}

func (h* ServerParametersConfig) ServerOption() _grpc.ServerOption {
	if h == nil { return  nil}
	p := ServerParameters{}
	p.MaxConnectionIdle, _ = time.ParseDuration(h.MaxConnectionIdle)
	p.MaxConnectionAge, _ = time.ParseDuration(h.MaxConnectionAge)
	p.MaxConnectionAgeGrace, _ = time.ParseDuration(h.MaxConnectionAgeGrace)
	p.Time, _ = time.ParseDuration(h.Time)
	p.Timeout, _ = time.ParseDuration(h.Timeout)
	return _grpc.KeepaliveParams(p)
}


type ClientParametersConfig struct {
	Time 						string		 			`yaml:"time"`
	Timeout 					string		 			`yaml:"timeout"`
	PermitWithoutStream			bool					`yaml:"permit_without_stream"`
}

func (h* ClientParametersConfig) DialOption() _grpc.DialOption {
	if h == nil { return  nil}
	p := ClientParameters{}
	p.Time, _ = time.ParseDuration(h.Time)
	p.Timeout, _ = time.ParseDuration(h.Timeout)
	p.PermitWithoutStream = h.PermitWithoutStream
	return _grpc.WithKeepaliveParams(p)
}


type EnforcementPolicyConfig struct {
	MinTime 					string					`yaml:"min_time"`
	PermitWithoutStream 		bool					`yaml:"permit_without_stream"`
}

func (h* EnforcementPolicyConfig) ServerOption() _grpc.ServerOption {
	if h == nil { return  nil}
	p := EnforcementPolicy{}
	p.MinTime, _ = time.ParseDuration(h.MinTime)
	p.PermitWithoutStream = h.PermitWithoutStream
	return _grpc.KeepaliveEnforcementPolicy(p)
}
