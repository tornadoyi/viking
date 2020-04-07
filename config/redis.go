package config

import (
	"fmt"
	"github.com/tornadoyi/viking/goplus/runtime"
	"github.com/tornadoyi/viking/log"
	"github.com/tornadoyi/viking/redis"
	"strings"
)


func AddRedisConfig(name, pool, command string, cmdArgs []interface{}, defaultConfig interface{}, parseFunc interface{}, args... interface{}) (*Config, error) {
	return AddParser(NewRedisParser(name, pool, command, cmdArgs, defaultConfig, parseFunc, args...))
}


type RedisParser struct {
	*BaseParser
	pool					string
	command					string
	cmdArgs					[]interface{}
	parser					*runtime.JITFunc
}

func NewRedisParser(name, pool, command string, cmdArgs []interface{}, defaultConfig interface{}, parseFunc interface{}, args... interface{}) *RedisParser {
	p := &RedisParser{
		pool:			pool,
		command:		command,
		cmdArgs:		cmdArgs,
		parser:			runtime.NewJITFunc(parseFunc, args...),
	}
	p.BaseParser = NewBaseParser(name, defaultConfig, p.parse)
	return p
}

func (h *RedisParser) Pool() string { return h.pool}
func (h *RedisParser) Command() string { return h.command}
func (h *RedisParser) CmdArgs() []interface{} { return h.cmdArgs}
func (h *RedisParser) CmdWithArgDesc() string {
	cmds := make([]string, 0, len(h.cmdArgs)+1)
	cmds = append(cmds, h.command)
	for _, arg := range h.cmdArgs { cmds = append(cmds, fmt.Sprintf("%v", arg)) }
	return strings.Join(cmds, " ")
}

func (h *RedisParser) parse() interface{} {
	content := h.defaultConfig
	pool, ok := redis.GetPool(h.pool)
	if !ok {
		log.Warnw("Can not fetch redis config with invalid pool name", "config", h.name, "pool", h.pool)
	} else {
		r := pool.Do(h.command, h.cmdArgs...)
		if r.Error() != nil {
			log.Warnw("Redis config fetch failed", "config", h.name, "pool", h.pool, "cmd", h.CmdWithArgDesc())
		} else  {
			v, _ := r.Interface()
			if v == nil {
				log.Warnw("Redis config parse failed, a nil parsed result", "config", h.name)
			} else {
				content = v
			}
		}
	}

	data, err := h.parser.Call(content)
	if err != nil { panic(err) }
	return data
}