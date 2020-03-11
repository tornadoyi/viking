package core

type Pair struct {
	first	interface{}
	second	interface{}
}

func (h *Pair) First() interface{} {return h.first}

func (h *Pair) Second() interface{} {return h.second}

func (h *Pair) SetFirst(v interface{}) { h.first = v}

func (h *Pair) SetSecond(v interface{}) { h.second = v}