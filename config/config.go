package config

import (
	"errors"
	"fmt"
	"reflect"
	"sort"
	"time"
	"viking/task"
)

var (
	configs 		= make(map[string]*Config)
	timers			= make(map[time.Duration]*UpdateTimer)
)


type Config struct {
	name			string
	executor		reflect.Value
	arguments		[]reflect.Value
	priority		time.Duration
	schedule		time.Duration
	data			interface{}

}

func (h *Config) Name() string{ return h.name }

func (h *Config) Priority() time.Duration{ return h.priority }

func (h *Config) Schedule() time.Duration{ return h.schedule }

func (h *Config) Execute() interface{}{ return h.executor.Call(h.arguments) }

func (h* Config) SetPriority(priority time.Duration){ h.priority = priority }

func (h* Config) SetSchedule(delay time.Duration){
	if _, ok := timers[delay]; ok { return }
	if delay <= 0 { panic(errors.New(fmt.Sprintf("Can not set timer's schedule with delay %v", delay))) }
	h.schedule = delay
}


type UpdateTimer struct {
	configs 	[]*Config
	timer		*time.Timer
}



func AddConfig(name string, f interface{}, args... interface{}) *Config{
	if _, ok := configs[name]; ok { panic(errors.New(fmt.Sprintf("Repeated config %k", name))) }
	vargs := make([]reflect.Value, len(args))
	for i, a := range (args){ vargs[i] = reflect.ValueOf(a) }
	config := &Config{
		name,
		reflect.ValueOf(f),
		vargs,
		1,
		0,
		nil,
	}
	configs[name] = config
	return config
}


func Start(){

	addSchedule := func(c *Config){
		if c.schedule <= 0 {return }
		updater, ok := timers[c.schedule]
		if ok {
			updater.configs = append(updater.configs, c)
			return
		}
		timer := time.AfterFunc(c.schedule, func(){
			defer updater.timer.Reset(c.schedule)

			updater, _ := timers[c.schedule]
			t := task.Create(updateConfigs, updater.configs)
			t.Wait()
		})
		timers[c.schedule] = &UpdateTimer{[]*Config{c}, timer}
	}

	// execute all configs
	cfgs := make([]*Config, 0, len(configs))
	for _, c := range configs{ cfgs = append(cfgs, c) }
	updateConfigs(cfgs)

	// add schedule
	for _, c := range configs{ addSchedule(c) }
}


func updateConfigs(configs []*Config){
	priors := make([]time.Duration, 0, len(configs))
	priorConfigs := make(map[time.Duration][]*Config)

	for _, c := range configs{
		if _, ok := priorConfigs[c.priority]; !ok {
			priors = append(priors, c.priority)
		}
		priorConfigs[c.priority] = append(priorConfigs[c.priority], c)
	}

	sort.Slice(priors, func(i, j int) bool {
		return priors[i] < priors[j]
	})

	for _, p := range priors {
		list, _ := priorConfigs[p]
		if len(list) == 0 { continue }

		ts := make(task.TaskGroup, 0, len(list))
		for _, c := range list{
			ts.Add(func(c *Config) {
				data := c.Execute()
				c.data = data
			}, c)
		}
		ts.Start()
		ts.Wait()
	}

}

