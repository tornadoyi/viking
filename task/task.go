package task

import (
	"errors"
	"fmt"
	"github.com/tornadoyi/viking/goplus/runtime"
	"github.com/tornadoyi/viking/log"
	"sync"
	"time"
)

const (
	Init				State = iota
	Running
	Finished
	Canceled
)


type Task struct {
	function			*runtime.JITFunc
	state				State
	result				interface{}
	error				error
	stack				runtime.StackInfo
	wg					*sync.WaitGroup
	terminateCallback	*runtime.JITFunc
	mutex				sync.RWMutex
}

func NewTask(f interface{}, args... interface{}) *Task {
	return newTask(&sync.WaitGroup{}, f, args...)
}

func (h *Task) State() State {
	h.mutex.RLock()
	s := h.state
	h.mutex.RUnlock()
	return s
}

func (h *Task) Error() error {
	h.mutex.RLock()
	err := h.error
	h.mutex.RUnlock()
	return err
}

func (h *Task) Result() interface{} {
	h.mutex.RLock()
	result := h.result
	h.mutex.RUnlock()
	return result
}

func (h *Task) Finished() bool {
	h.mutex.RLock()
	ret := h.state == Finished
	h.mutex.RUnlock()
	return ret
}

func (h *Task) Canceled() bool {
	h.mutex.RLock()
	ret := h.state == Canceled
	h.mutex.RUnlock()
	return ret
}

func (h *Task) Terminated() bool {
	h.mutex.RLock()
	ret := h.state == Finished || h.state == Canceled
	h.mutex.RUnlock()
	return ret
}

func (h *Task) SetTerminateCallback (f func(*Task)) {
	h.mutex.Lock()
	h.terminateCallback = runtime.NewJITFunc(f)
	h.mutex.Unlock()
}

func (h *Task) Start(){
	var stateErr error = nil
	h.mutex.Lock()
	if h.state != Init {
		stateErr = fmt.Errorf("task can not start, current state is %v", h.state)
	} else {
		h.state = Running
		h.wg.Add(1)
	}
	h.mutex.Unlock()
	if stateErr != nil { panic(stateErr) }

	go func() {
		result, err := h.function.Call()
		h.terminate(false, result, err)
	}()
}

func (h *Task) Cancel(){
	h.terminate(true, nil, errors.New("positive cancel"))
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

func (h *Task) terminate(cancel bool, result interface{}, err error) {
	// check state
	s := h.State()
	if s == Canceled || s == Finished { return }

	// terminate
	h.mutex.Lock()
	if s == Running { h.wg.Done() }
	if cancel { h.state = Canceled } else { h.state = Finished }
	h.error = err
	h.result = result
	h.mutex.Unlock()

	// callback
	if h.terminateCallback != nil {
		_, err := h.terminateCallback.Call(h)
		if err != nil { log.Error(err) }
	}
}





type 		State				int
func (h State) String() string {
	switch h {
	case Init: 		return "Init"
	case Running: 	return "Running"
	case Finished: 	return "Finished"
	case Canceled: 	return "Canceled"
	default: 		return "Invalid"
	}
}