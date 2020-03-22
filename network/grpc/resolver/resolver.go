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
var Register = _resolver.Register
var Get = _resolver.Get
var SetDefaultScheme = _resolver.SetDefaultScheme
var GetDefaultScheme = _resolver.GetDefaultScheme