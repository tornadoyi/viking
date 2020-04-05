package config

import (
	"fmt"
	"github.com/tornadoyi/viking/goplus/runtime"
	"github.com/tornadoyi/viking/log"
	"github.com/tornadoyi/viking/task"
	"reflect"
	"sort"
	"sync"
	"sync/atomic"
	"time"
)

var (
	configs 		= map[string]*Config{}
	timers			= map[time.Duration]*time.Timer{}

	// event
	onLoadStart	 	= (*runtime.JITFunc)(nil)
	onLoadEnd	 	= (*runtime.JITFunc)(nil)

	mutex			= sync.RWMutex{}

)



func AddConfig(name string, f interface{}, args... interface{}) (*Config, error) {
	mutex.Lock()
	defer mutex.Unlock()
	_, ok := configs[name]
	if ok { return nil, fmt.Errorf("Repeated config %k", name)}
	vargs := make([]reflect.Value, len(args))
	for i, a := range (args){ vargs[i] = reflect.ValueOf(a) }
	config := NewConfig(name, f, args...)
	configs[name] = config
	return config, nil
}

func AddParser(parser IParser) (*Config, error) {
	cfg, err := AddConfig(parser.Name(), parser.ParseFunc(), parser.Arguments()...)
	if err != nil { return nil, err}
	err = cfg.AddFunc("Parser", func() IParser { return parser })
	if err != nil { return nil, err}
	return cfg, nil
}


func GetConfig(name string) (*Config, bool) {
	mutex.RLock()
	defer mutex.RUnlock()
	c, ok := configs[name]
	return c, ok
}


func GetConfigData(name string) (interface{}, bool) {
	mutex.RLock()
	defer mutex.RUnlock()
	c, ok := configs[name]
	if !ok { return nil, false }
	return c.Data(), true
}

func Configs() []*Config {
	mutex.RLock()
	defer mutex.RUnlock()
	cfgs := make([]*Config, 0, len(configs))
	for _, cfg := range configs { cfgs = append(cfgs, cfg) }
	return cfgs
}



func RegisterLoadStartEvent(f func(*Config)) {
	mutex.Lock()
	defer mutex.Unlock()
	onLoadStart = runtime.NewJITFunc(f)
}

func RegisterLoadEndEvent(f func(*Config)) {
	mutex.Lock()
	defer mutex.Unlock()
	onLoadEnd = runtime.NewJITFunc(f)
}



func Start() []error{
	cfgs := make([]*Config, 0)
	mutex.RLock()
	for _, cfg := range configs { cfgs = append(cfgs, cfg) }
	mutex.RUnlock()

	// execute all configs
	return updateConfigs(cfgs)

}


func updateConfigs(configs []*Config) []error{
	priors := make([]int, 0, len(configs))
	priorConfigs := make(map[int][]*Config)

	for _, c := range configs{
		if _, ok := priorConfigs[c.priority]; !ok {
			priors = append(priors, c.priority)
		}
		priorConfigs[c.priority] = append(priorConfigs[c.priority], c)
	}

	sort.Slice(priors, func(i, j int) bool {
		return priors[i] < priors[j]
	})

	errs := make([]error, 0)

	for _, p := range priors {
		list, _ := priorConfigs[p]
		if len(list) == 0 { continue }

		ts := task.NewGroup()
		for _, c := range list{
			ts.Add(func(c *Config) {
				if onLoadStart != nil { onLoadStart.Call(c) }
				data := c.Execute()
				c.data.Store(data)
				if onLoadEnd != nil { onLoadEnd.Call(c) }
			}, c)
		}
		ts.Start()
		ts.Wait()

		errs = append(errs, ts.Errors()...)
	}

	return errs
}


func addTimer(delay time.Duration) {
	if delay <= 0 {return }
	mutex.Lock()
	defer mutex.Unlock()

	if _, ok := timers[delay]; ok { return }
	timers[delay] = time.AfterFunc(delay, func() {
		// collect configs
		cfgs := make([]*Config, 0)
		var timer *time.Timer
		mutex.Lock()
		timer, _ = timers[delay]
		for _, cfg := range configs {
			if cfg.schedule != delay { continue }
			cfgs = append(cfgs, cfg)
		}
		// clean unused timer
		if len(cfgs) == 0 || timer == nil  { delete(timers, delay) }
		mutex.Unlock()

		// exit when timer is unused
		if len(cfgs) == 0 || timer == nil { return }

		// update config
		t := task.NewTask(updateConfigs, cfgs)
		t.Start()
		t.Wait()
		if t.Error() != nil {
			log.Error(t.Error())
		} else {
			errs := t.Result().([]error)
			for _, err := range errs { log.Error(err) }
		}

		// reset timer
		timer.Reset(delay)
	})
}





type Config struct {
	*runtime.JIT
	name									string
	executor								*runtime.JITFunc
	priority								int
	schedule								time.Duration
	data									atomic.Value
	lastExecuteStartTime					int64
	lastExecuteEndTime						int64
}

func NewConfig(name string, f interface{}, args... interface{}) *Config {
	return &Config{
		JIT:	  runtime.NewJIT(),
		name:     name,
		executor: runtime.NewJITFunc(f, args...),
		data:     atomic.Value{},
	}
}

func (h *Config) Name() string{ return h.name }

func (h *Config) Priority() int { return h.priority }

func (h *Config) Schedule() time.Duration{ return h.schedule }

func (h *Config) LastExecuteStartTime() int64 { return h.lastExecuteStartTime}

func (h *Config) LastExecuteEndTime() int64 { return h.lastExecuteEndTime }

func (h *Config) Execute() interface{} {
	h.lastExecuteStartTime = time.Now().UnixNano()
	result, err := h.executor.Call()
	h.lastExecuteEndTime = time.Now().UnixNano()
	if err != nil { panic(err) }
	return result
}

func (h* Config) SetPriority(priority int){ h.priority = priority }

func (h *Config) Data() interface{} { return h.data.Load() }

func (h* Config) SetSchedule(delay time.Duration){
	h.schedule = delay
	addTimer(delay)
}




type IParser interface {
	Name()					string
	DefaultConfig()			interface{}
	ParseFunc()				interface{}
	Arguments()				[]interface{}
}


type BaseParser struct {
	name					string
	defaultConfig			interface{}
	parseFunc				interface{}
	arguments				[]interface{}
}


func NewBaseParser(name string, defaultConfig, parseFunc interface{}, args... interface{}) *BaseParser {
	return &BaseParser{
		name:          name,
		defaultConfig: defaultConfig,
		parseFunc:     parseFunc,
		arguments:     args,
	}
}

func (h *BaseParser) Name() string { return h.name}
func (h *BaseParser) DefaultConfig() interface{} { return h.defaultConfig}
func (h *BaseParser) ParseFunc() interface{} { return h.parseFunc }
func (h *BaseParser) Arguments() []interface{} { return h.arguments }
