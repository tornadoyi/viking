package redis

import (
	"time"
)



type PoolConfig struct {
	MaxIdle 					*int						`yaml:"max_idle"`
	MaxActive 					*int						`yaml:"max_active"`
	IdleTimeout 				*time.Duration				`yaml:"idle_timeout"`
	Wait 						*bool						`yaml:"wait"`
	MaxConnLifetime 			*time.Duration				`yaml:"max_conn_life_time"`
}


type PoolOption struct {
	f func(pool *Pool)
}

func PoolMaxIdle(v int) PoolOption { return PoolOption{func(p *Pool) {p.MaxIdle = v }} }
func PoolMaxActive(v int) PoolOption { return PoolOption{func(p *Pool) {p.MaxActive = v }}}
func PoolIdleTimeout(v time.Duration) PoolOption { return PoolOption{func(p *Pool) {p.IdleTimeout = v }}}
func PoolWait(v bool) PoolOption { return PoolOption{func(p *Pool) {p.Wait = v}} }
func PoolMaxConnLifetime(v time.Duration) PoolOption { return PoolOption{func(p *Pool) {p.MaxConnLifetime = v }}}



type DialConfig struct {
	ClientName					*string						`yaml:"client_name"`
	ConnectTimeout				*time.Duration				`yaml:"connect_timeout"`
	Database					*int						`yaml:"database"`
	KeepAlive					*time.Duration				`yaml:"keep_alive"`
	Password					*string						`yaml:"password"`
	ReadTimeout					*time.Duration				`yaml:"read_timeout"`
	//TLSConfig
	TLSSkipVerify				*bool						`yaml:"tls_skip_verify"`
	UseTLS						*bool						`yaml:"use_tls"`
	WriteTimeout				*time.Duration				`yaml:"write_timeout"`
}

func (h *DialConfig) DialOptions() []DialOption {
	options := make([]DialOption, 0, 20)
	if h.ClientName != nil { options = append(options, DialClientName(*h.ClientName)) }
	if h.ConnectTimeout != nil { options = append(options, DialConnectTimeout(*h.ConnectTimeout)) }
	if h.Database != nil { options = append(options, DialDatabase(*h.Database)) }
	if h.KeepAlive != nil { options = append(options, DialKeepAlive(*h.KeepAlive)) }
	if h.Password != nil { options = append(options, DialPassword(*h.Password)) }
	if h.ReadTimeout != nil { options = append(options, DialReadTimeout(*h.ReadTimeout)) }
	if h.TLSSkipVerify != nil { options = append(options, DialTLSSkipVerify(*h.TLSSkipVerify)) }
	if h.UseTLS != nil { options = append(options, DialUseTLS(*h.UseTLS)) }
	if h.WriteTimeout != nil { options = append(options, DialWriteTimeout(*h.WriteTimeout)) }
	return options
}