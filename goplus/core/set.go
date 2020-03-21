package core

import "sync"

type  Set	map[interface{}]bool

func (h Set) Add(key interface{}) {
	h[key] = true
}

func (h Set) Delete(key interface{}) {
	if _, ok := h[key]; ok {
		delete (h, key)
	}
}

func (h Set) Exists(key interface{}) bool {
	_, ok := h[key]
	return ok
}


type AtomSet struct {
	set				Set
	mutex			sync.RWMutex
}

func (h *AtomSet) Add(key interface{}) {
	defer h.mutex.Unlock()
	h.mutex.Lock()
	h.set.Add(key)
}

func (h *AtomSet) Delete(key interface{}) {
	defer h.mutex.Unlock()
	h.mutex.Lock()
	h.Delete(key)
}

func (h *AtomSet) Exists(key interface{}) bool {
	defer h.mutex.RUnlock()
	h.mutex.RLock()
	return h.Exists(key)
}