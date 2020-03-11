package task

import "time"

type TaskGroup	[]*Task

func (h *TaskGroup) Start(){
	for _, t := range *h{
		t.Start()
	}
}

func (h *TaskGroup) Add(f interface{}, args... interface{}) *Task{
	t := Create(f, args...)
	return t
}

func (h *TaskGroup) Wait(){
	for _, t := range *h{
		t.Wait()
	}
}

func (h *TaskGroup) WaitTimeout(timeout time.Duration) bool{
	c := make(chan struct{})
	go func() {
		defer close(c)
		h.Wait()
	}()
	select {
	case <-c:
		return true
	case <-time.After(timeout):
		return false
	}
}