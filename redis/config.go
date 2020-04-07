package redis

import (
	"github.com/tornadoyi/viking/log"
	"time"
)



type PoolConfig struct {
	MaxIdle 					*int						`yaml:"max_idle"`
	MaxActive 					*int						`yaml:"max_active"`
	IdleTimeout 				*string						`yaml:"idle_timeout"`
	Wait 						*bool						`yaml:"wait"`
	MaxConnLifetime 			*string						`yaml:"max_conn_lifetime"`
}

func (h *PoolConfig) PoolOptions() []PoolOption {
	if h == nil { return nil }
	options := make([]PoolOption, 0)
	if h.MaxIdle != nil { options = append(options, PoolMaxIdle(*h.MaxIdle)) }
	if h.MaxActive != nil { options = append(options, PoolMaxActive(*h.MaxActive)) }
	if h.IdleTimeout != nil {
		if d, err := time.ParseDuration(*h.IdleTimeout); err == nil {
			options = append(options, PoolIdleTimeout(d))
		} else { log.Warnw("IdleTimeout parse failed",
			"error", err)
		}
	}
	if h.Wait != nil { options = append(options, PoolWait(*h.Wait)) }
	if h.MaxConnLifetime != nil {
		if d, err := time.ParseDuration(*h.MaxConnLifetime); err == nil {
			options = append(options, PoolMaxConnLifetime(d))
		} else { log.Warnw("MaxConnLifetime parse failed",
			"error", err)
		}
	}
	return options
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
	ConnectTimeout				*string						`yaml:"connect_timeout"`
	Database					*int						`yaml:"database"`
	KeepAlive					*string						`yaml:"keep_alive"`
	Password					*string						`yaml:"password"`
	ReadTimeout					*string						`yaml:"read_timeout"`
	//TLSConfig
	TLSSkipVerify				*bool						`yaml:"tls_skip_verify"`
	UseTLS						*bool						`yaml:"use_tls"`
	WriteTimeout				*string						`yaml:"write_timeout"`
}

func (h *DialConfig) DialOptions() []DialOption {
	if h == nil { return  nil}
	options := make([]DialOption, 0, 20)
	if h.ClientName != nil { options = append(options, DialClientName(*h.ClientName)) }
	if h.ConnectTimeout != nil {
		if d, err := time.ParseDuration(*h.ConnectTimeout); err == nil {
			options = append(options, DialConnectTimeout(d))
		} else { log.Warnw("ConnectTimeout parse failed", "error", err) }
	}
	if h.Database != nil { options = append(options, DialDatabase(*h.Database)) }
	if h.KeepAlive != nil {
		if d, err := time.ParseDuration(*h.KeepAlive); err == nil {
			options = append(options, DialKeepAlive(d))
		} else { log.Warnw("KeepAlive parse failed", "error", err) }
	}
	if h.Password != nil { options = append(options, DialPassword(*h.Password)) }
	if h.ReadTimeout != nil {
		if d, err := time.ParseDuration(*h.ReadTimeout); err == nil {
			options = append(options, DialReadTimeout(d))
		} else { log.Warnw("ReadTimeout parse failed", "error", err) }
	}
	if h.TLSSkipVerify != nil { options = append(options, DialTLSSkipVerify(*h.TLSSkipVerify)) }
	if h.UseTLS != nil { options = append(options, DialUseTLS(*h.UseTLS)) }
	if h.WriteTimeout != nil {
		if d, err := time.ParseDuration(*h.WriteTimeout); err == nil {
			options = append(options, DialWriteTimeout(d))
		} else { log.Warnw("WriteTimeout parse failed", "error", err) }
	}
	return options
}