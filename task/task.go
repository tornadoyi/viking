package task

import (
	"errors"
	"fmt"
	. "github.com/tornadoyi/viking/goplus/core"
	"github.com/tornadoyi/viking/goplus/runtime"
	"reflect"
	"sync"
	"time"
)

const (
	Init		int = iota
	Running
	Finished
	Canceled
)



func Create(f interface{}, args... interface{}) *Task {
	wg := &sync.WaitGroup{}
	return createTask(f, wg, args...)
}



type Task struct {
	function			reflect.Value
	arguments			[]reflect.Value
	state				int
	result				interface{}
	error				error
	stack				runtime.StackInfo
	wg					*sync.WaitGroup
	mutex				sync.RWMutex
}

func (h *Task) State() int {
	h.mutex.RLock()
	defer h.mutex.RUnlock()
	return h.state
}

func (h *Task) Error() error {
	h.mutex.RLock()
	defer h.mutex.RUnlock()
	return h.error
}

func (h *Task) Result() interface{} {
	h.mutex.RLock()
	defer h.mutex.RUnlock()
	return h.result
}

func (h *Task) Finished() bool {
	h.mutex.RLock()
	defer h.mutex.RUnlock()
	return h.state == Finished
}

func (h *Task) Canceled() bool {
	h.mutex.RLock()
	defer h.mutex.RUnlock()
	return h.state == Canceled
}

func (h *Task) Terminated() bool {
	h.mutex.RLock()
	defer h.mutex.RUnlock()
	return h.state == Finished || h.state == Canceled
}



func (h *Task) StateDesc() string {
	h.mutex.RLock()
	defer h.mutex.RUnlock()
	switch h.state {
	case Init: return "Init"
	case Running: return "Running"
	case Finished: return "Finished"
	case Canceled: return "Canceled"
	default: return "Invalid"
	}
}

func (h *Task) Start(){
	var stateErr error = nil
	h.mutex.Lock()
	if h.state != Init {
		stateErr = fmt.Errorf("Task can not start, current state is %v", h.StateDesc())
	} else {
		h.state = Running
		h.wg.Add(1)
	}
	h.mutex.Unlock()
	if stateErr != nil { panic(stateErr) }

	go func() {
		defer func(){
			h.mutex.Lock()
			if h.state == Running {
				h.wg.Done()
				h.state = Finished
			}
			h.mutex.Unlock()
		}()
		defer CatchCallback(func(info *PanicInfo){
			h.mutex.Lock()
			defer h.mutex.Unlock()
			h.error = info.Error()
			/*
			log.Critical(strings.Join([]string{
				"A task error occurred as below",
				fmt.Sprintf("%v", h.stack),
			}, "\n"))
			 */
		})

		// collect results
		vres := h.function.Call(h.arguments)
		var result interface{}
		if len(vres) == 1 { result = vres[0].Interface() } else {
			res := make([]interface{}, 0, len(vres))
			for _, v := range vres { res = append(res, v) }
			result = res
		}

		// save result
		h.mutex.Lock()
		if h.state != Canceled { h.result = result }
		h.mutex.Unlock()
	}()
}

func (h *Task) Cancel(){
	switch h.State() {
	case Init:
		h.mutex.Lock()
		h.state = Canceled
		h.mutex.Unlock()
		return
	case Canceled, Finished: return
	}

	h.mutex.Lock()
	h.wg.Done()
	h.state = Canceled
	h.mutex.Unlock()
}



func (h *Task) Wait(){
	switch h.State() {
	case Init: panic(errors.New("Can't wait for a task that hasn't started"))
	case Finished: return
	}
	h.wg.Wait()
}

func (h *Task) WaitTimeout(timeout time.Duration) {
	c := make(chan struct{})
	go func() {
		defer close(c)
		h.Wait()
	}()
	select {
	case <-c:  return
	case <-time.After(timeout):
		h.Cancel()
		return
	}
}






func createTask(f interface{}, wg *sync.WaitGroup, args... interface{}) *Task {
	if f == nil { panic(errors.New("Can not create task with nil function")) }
	vf := reflect.ValueOf(f)
	if vf.Kind() != reflect.Func { panic(fmt.Errorf("Can not create task with invalid function type %v", vf.Kind())) }
	vargs := make([]reflect.Value, len(args))
	for i, a := range (args){ vargs[i] = reflect.ValueOf(a) }
	return &Task{vf, vargs,Init, nil, nil, runtime.Trace(3),wg, sync.RWMutex{}}
}