package task

import (
	"errors"
	"fmt"
	"github.com/tornadoyi/viking/goplus/runtime"
	"sync"
	"time"
)

const (
	Init		int = iota
	Running
	Finished
	Canceled
)


type Task struct {
	function			*runtime.JITFunc
	state				int
	result				interface{}
	error				error
	stack				runtime.StackInfo
	wg					*sync.WaitGroup
	mutex				sync.RWMutex
}

func NewTask(f interface{}, args... interface{}) *Task {
	return newTask(&sync.WaitGroup{}, f, args...)
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

		// collect results
		result, err := h.function.Call()

		// save result
		h.mutex.Lock()
		if h.state != Canceled {
			h.result = result
			h.error = err
		}
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



func newTask(wg *sync.WaitGroup, f interface{}, args... interface{}) *Task {
	return &Task{
		function: runtime.NewJITFunc(f, args...),
		state:    Init,
		result:   nil,
		error:    nil,
		stack:    runtime.Trace(2),
		wg:       wg,
		mutex:    sync.RWMutex{},
	}
}