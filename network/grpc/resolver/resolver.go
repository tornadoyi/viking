package resolver

import (
	_resolver "google.golang.org/grpc/resolver"
)



type ResolverBuilderConfig struct {
	Scheme			string						`yaml:"scheme"`
	Address			map[string][]string			`yaml:"address"`
}

func (h *ResolverBuilderConfig) ResolverBuilder() *ResolverBuilder{
	return &ResolverBuilder{h.Scheme, h.Address}
}



type ResolverBuilder struct {
	scheme			string
	address			map[string][]string
}

func (h *ResolverBuilder) Build(target _resolver.Target, cc _resolver.ClientConn, opts _resolver.BuildOptions) (_resolver.Resolver, error) {
	r := &Resolver{target, cc, h.address}
	r.init()
	return r, nil
}
func (h *ResolverBuilder) Scheme() string { return h.scheme }

type Resolver struct {
	target     				_resolver.Target
	connection         		_resolver.ClientConn
	address 				map[string][]string
}

func (h *Resolver) init() {
	addrStrs := h.address[h.target.Endpoint]
	addrs := make([]_resolver.Address, len(addrStrs))
	for i, s := range addrStrs {
		addrs[i] = _resolver.Address{Addr: s}
	}
	h.connection.UpdateState(_resolver.State{Addresses: addrs})
}
func (*Resolver) ResolveNow(o _resolver.ResolveNowOptions) {}

func (*Resolver) Close() {}




// export
type IResolver	= _resolver.Resolver
type Target	= _resolver.Target
type ClientConn	= _resolver.ClientConn
type ResolveNowOptions	= _resolver.ResolveNowOptions
type BuildOptions	= _resolver.BuildOptions
type Address	= _resolver.Address
type State	= _resolver.State


var Register = _resolver.Register
var Get = _resolver.Get
var SetDefaultScheme = _resolver.SetDefaultScheme
var GetDefaultScheme = _resolver.GetDefaultScheme