
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
	MaxConnectionIdle 			time.Duration 			`yaml:"max_connection_idle"`
	MaxConnectionAge 			time.Duration			`yaml:"max_connection_age"`
	MaxConnectionAgeGrace 		time.Duration 			`yaml:"max_connection_age_grace"`
	Time 						time.Duration 			`yaml:"time"`
	Timeout 					time.Duration 			`yaml:"timeout"`
}

func (h* ServerParametersConfig) ServerOption() _grpc.ServerOption {
	if h == nil { return  nil}
	return _grpc.KeepaliveParams(ServerParameters{
		h.MaxConnectionIdle,
		h.MaxConnectionAge,
		h.MaxConnectionAgeGrace,
		h.Time,
		h.Timeout,
	})
}


type ClientParametersConfig struct {
	Time 						time.Duration 			`yaml:"time"`
	Timeout 					time.Duration 			`yaml:"timeout"`
	PermitWithoutStream			bool					`yaml:"permit_without_stream"`
}

func (h* ClientParametersConfig) DialOption() _grpc.DialOption {
	if h == nil { return  nil}
	return _grpc.WithKeepaliveParams(ClientParameters{h.Time, h.Timeout, h.PermitWithoutStream})
}


type EnforcementPolicyConfig struct {
	MinTime 					time.Duration 			`yaml:"min_time"`
	PermitWithoutStream 		bool					`yaml:"permit_without_stream"`
}

func (h* EnforcementPolicyConfig) ServerOption() _grpc.ServerOption {
	if h == nil { return  nil}
	return _grpc.KeepaliveEnforcementPolicy(EnforcementPolicy{h.MinTime, h.PermitWithoutStream})
}
