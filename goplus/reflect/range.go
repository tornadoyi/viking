package reflect


type Node interface {
	Value() *Value
	StructField() *StructField
	FieldPath() []*StructField
	ValuePath() []*Value
	Interrupt()
	Stop()
}

type _Node struct {
	fieldPath				[]*StructField
	valuePath				[]*Value
	value					*Value
	structField				*StructField
	interrupt				bool
	stop					bool
}

func (h *_Node) Value() *Value { return h.value}
func (h *_Node) StructField() *StructField { return h.structField}
func (h *_Node) FieldPath() []*StructField { return h.fieldPath}
func (h *_Node) ValuePath() []*Value { return h.valuePath}
func (h *_Node) Interrupt() { h.interrupt = true }
func (h *_Node) Stop() { h.stop = true }



func DFS(obj interface{}, fn func(Node)) {

	node := &_Node{
		fieldPath:   make([]*StructField, 0),
		valuePath:   make([]*Value, 0),
	}

	var dfs func(Value)
	dfs = func(v Value) {
		if !v.IsValid() { return }
		// call
		node.value = &v
		fn(node)

		// check state
		if node.stop { return }
		if node.interrupt {
			node.interrupt = false
			return
		}

		switch v.Kind() {
		case Struct:
			t := v.Type()
			for i:=0; i<t.NumField(); i++{
				sf := t.Field(i)
				vf := v.Field(i)
				node.structField = &sf
				node.fieldPath = append(node.fieldPath, &sf)
				node.valuePath = append(node.valuePath, &vf)

				dfs(vf)

				node.structField = nil
				node.fieldPath = node.fieldPath[:len(node.fieldPath)-1]
				node.valuePath = node.valuePath[:len(node.valuePath)-1]

				if node.stop { break }
			}
		case Ptr:
			dfs(v.Elem())
		}
	}

	dfs(ValueOf(obj))
}
