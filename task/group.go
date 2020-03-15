package task

import "time"

type TaskGroup	[]*Task

func (h *TaskGroup) Finish() bool{
	for _, t := range *h{
		if !t.Finish() { return false}
	}
	return true
}

func (h *TaskGroup) Errors() []error {
	errs := make([]error, 0)
	for _, t := range *h{
		err := t.Error()
		if err == nil { continue }
		errs = append(errs, err)
	}
	return errs
}

func (h *TaskGroup) Start(){
	for _, t := range *h{
		t.Start()
	}
}




func (h *TaskGroup) Add(f interface{}, args... interface{}) *Task{
	t := Create(f, args...)
	*h = append(*h, t)
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