package task

import (
	"sync"
	"time"
)


func CreateGroup() *TaskGroup{ return &TaskGroup{make([]*Task, 0), &sync.WaitGroup{}} }


type TaskGroup struct {
	tasks			[]*Task
	wg				*sync.WaitGroup
}

func (h *TaskGroup) Terminated() bool{
	for _, t := range h.tasks{
		if !t.Terminated() { return false}
	}
	return true
}

func (h *TaskGroup) Errors() []error {
	errs := make([]error, 0)
	for _, t := range h.tasks{
		err := t.Error()
		if err == nil { continue }
		errs = append(errs, err)
	}
	return errs
}

func (h *TaskGroup) Add(f interface{}, args... interface{}) *Task{
	t := createTask(f, h.wg, args...)
	h.tasks = append(h.tasks, t)
	return t
}


func (h *TaskGroup) Range(f func(int, *Task) bool)  {
	for i, t := range h.tasks {
		if f(i, t) { continue }
		break
	}
}



func (h *TaskGroup) Start(){
	for _, t := range h.tasks{
		t.Start()
	}
}


func (h *TaskGroup) Cancel(){
	for _, t := range h.tasks{
		t.Cancel()
	}
}


func (h *TaskGroup) Wait(){
	h.wg.Wait()
}

func (h *TaskGroup) WaitTimeout(timeout time.Duration){
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