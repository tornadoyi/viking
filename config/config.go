package config

import (
	"fmt"
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
	mutex			= sync.RWMutex{}
)


type Config struct {
	name			string
	executor		reflect.Value
	arguments		[]reflect.Value
	priority		int
	schedule		time.Duration
	data			atomic.Value

}

func (h *Config) Name() string{ return h.name }

func (h *Config) Priority() int { return h.priority }

func (h *Config) Schedule() time.Duration{ return h.schedule }

func (h *Config) Execute() interface{}{ return h.executor.Call(h.arguments)[0].Interface() }

func (h* Config) SetPriority(priority int){ h.priority = priority }

func (h* Config) SetSchedule(delay time.Duration){
	h.schedule = delay
	addTimer(delay)
}



func Add(name string, f interface{}, args... interface{}) (*Config, error) {
	mutex.Lock()
	defer mutex.Unlock()
	_, ok := configs[name]
	if ok { return nil, fmt.Errorf("Repeated config %k", name)}
	vargs := make([]reflect.Value, len(args))
	for i, a := range (args){ vargs[i] = reflect.ValueOf(a) }
	config := &Config{
		name,
		reflect.ValueOf(f),
		vargs,
		1,
		0,
		atomic.Value{},
	}
	configs[name] = config
	return config, nil
}


func GetContent(name string) (interface{}, bool) {
	mutex.RLock()
	defer mutex.RUnlock()
	c, ok := configs[name]
	if !ok { return nil, false }
	return c.data.Load(), true
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

		ts := task.CreateGroup()
		for _, c := range list{
			ts.Add(func(c *Config) {
				data := c.Execute()
				c.data.Store(data)
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
		t := task.Create(updateConfigs, cfgs)
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
