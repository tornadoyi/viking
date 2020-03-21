package core

import (
	"sync"
)

type Dict	map[interface{}]interface{}


type AtomDict struct {
	sync.Map
}


func (h *AtomDict) Set(key interface{}, value interface{}) { h.Map.Store(key, value) }

func (h *AtomDict) Get(key interface{}) (interface{}, bool) { return h.Map.Load(key) }

func (h *AtomDict) Exists(key interface{}) bool {
	_, ok := h.Get(key)
	return ok
}
