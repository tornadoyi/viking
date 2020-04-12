package core

import "sync"

type DBuffer struct {
	buffer0						interface{}
	buffer1						interface{}
	curIdx						int8
	mutex						*sync.RWMutex
}

func NewDBuffer(current, backup interface{}) *DBuffer {
	return &DBuffer{
		buffer0: current,
		buffer1: backup,
		curIdx:  0,
		mutex:   &sync.RWMutex{},
	}
}

func (h *DBuffer) Current() interface{} {
	var buf interface{}
	h.mutex.RLock()
	if h.curIdx == 0 { buf = h.buffer0 } else { buf = h.buffer1 }
	h.mutex.RUnlock()
	return buf
}

func (h *DBuffer) Backup() interface{} {
	var buf interface{}
	h.mutex.RLock()
	if h.curIdx == 0 { buf = h.buffer1 } else { buf = h.buffer0 }
	h.mutex.RUnlock()
	return buf
}

func (h *DBuffer) Swap() {
	h.mutex.Lock()
	if h.curIdx == 0 { h.curIdx = 1 } else { h.curIdx = 0 }
	h.mutex.Unlock()
}
