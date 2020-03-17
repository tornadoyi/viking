package config

import (
	"errors"
	"fmt"
	"reflect"
	"sort"
	"time"
	"github.com/tornadoyi/viking/log"
	"github.com/tornadoyi/viking/task"
)

var (
	configs 		= make(map[string]*Config)
	timers			= make(map[time.Duration]*UpdateTimer)
)


type Config struct {
	name			string
	executor		reflect.Value
	arguments		[]reflect.Value
	priority		int
	schedule		time.Duration
	data			interface{}

}

func (h *Config) Name() string{ return h.name }

func (h *Config) Priority() int { return h.priority }

func (h *Config) Schedule() time.Duration{ return h.schedule }

func (h *Config) Execute() interface{}{ return h.executor.Call(h.arguments)[0].Interface() }

func (h* Config) SetPriority(priority int){ h.priority = priority }

func (h* Config) SetSchedule(delay time.Duration){
	if _, ok := timers[delay]; ok { return }
	if delay <= 0 { panic(errors.New(fmt.Sprintf("Can not set timer's schedule with delay %v", delay))) }
	h.schedule = delay
}


type UpdateTimer struct {
	configs 	[]*Config
	timer		*time.Timer
}



func Add(name string, f interface{}, args... interface{}) *Config{
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


func GetContent(name string) interface{}{
	c, ok := configs[name]
	if !ok { panic(errors.New(fmt.Sprintf("Config %v is non-existent", name))) }
	return c.data
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

		ts := make(task.TaskGroup, 0, len(list))
		for _, c := range list{
			ts.Add(func(c *Config) {
				data := c.Execute()
				c.data = data
			}, c)
		}
		ts.Start()
		ts.Wait()

		errs = append(errs, ts.Errors()...)
	}

	for _, e := range errs {
		log.Error(e)
	}

}

