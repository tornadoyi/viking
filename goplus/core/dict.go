package core

import (
	"sync"
)

type Dict	map[interface{}]interface{}

func (d Dict) Exists(key interface{}) bool {
	_, ok := d[key]
	return ok
}

func (h Dict) Keys() []interface{} {
	keys := make([]interface{}, 0, len(h))
	for k, _ := range h {
		keys = append(keys, k)
	}
	return keys
}

func (h Dict) Values() []interface{} {
	values := make([]interface{}, 0, len(h))
	for _, v := range h {
		values = append(values, v)
	}
	return values
}


type AtomicDict struct {
	sync.Map
}


func (h *AtomicDict) Set(key interface{}, value interface{}) { h.Map.Store(key, value) }

func (h *AtomicDict) Get(key interface{}) (interface{}, bool) { return h.Map.Load(key) }

func (h *AtomicDict) Exists(key interface{}) bool {
	_, ok := h.Get(key)
	return ok
}

func (h *AtomicDict) Keys() []interface{} {
	keys := make([]interface{}, 0)
	h.Range(func(key, value interface{}) bool {
		keys = append(keys, key)
		return true
	})
	return keys
}

func (h *AtomicDict) Values() []interface{} {
	values := make([]interface{}, 0)
	h.Range(func(key, value interface{}) bool {
		values = append(values, value)
		return true
	})
	return values
}
