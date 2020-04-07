package redis

import (
	_redis "github.com/gomodule/redigo/redis"
	"github.com/tornadoyi/viking/log"
	"runtime"
	"time"
)

type Pool struct {
	_redis.Pool
}

func (h *Pool) init() {
	runtime.SetFinalizer(h, func(p *Pool) {
		if err := p.Close(); err != nil {
			log.Errorw("Redis pool closed with error", "error", err)
		}
	})
}

func (h *Pool) Do(commandName string, args ...interface{}) (*Result)  {
	conn := h.Get()
	defer conn.Close()
	reply, err := conn.Do(commandName, args...)
	return &Result{reply, err}
}

func (h *Pool) DoWithTimeout (timeout time.Duration, commandName string, args ...interface{}) (*Result) {
	conn := h.Get()
	defer conn.Close()
	reply, err := DoWithTimeout(conn, timeout, commandName, args)
	return &Result{reply, err}
}
