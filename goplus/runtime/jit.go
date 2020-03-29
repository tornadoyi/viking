package runtime

import (
	"fmt"
	"reflect"
)

type JIT struct {
	funcs				map[string]*JITFunc
}

func NewJIT() *JIT{
	return &JIT{make(map[string]*JITFunc)}
}

func (h *JIT) AddFunc(name string, function interface{}, contexts... interface{}) error{
	if _, ok := h.funcs[name]; ok { return fmt.Errorf("Repeated jit function %v", name)}
	f :=  NewJITFunc(function, contexts...)
	f.SetName(name)
	h.funcs[name] = f
	return nil
}

func (h *JIT) HasFunc(name string) bool {
	_, ok := h.funcs[name]
	return ok
}

func (h *JIT) CallFunc(name string, args... interface{}) (interface{}, error) {
	f, ok := h.funcs[name]
	if !ok { return nil, fmt.Errorf("Can not found jit function  %v", name)}
	return f.Call(args...)
}



type JITFunc struct {
	// inputs
	name				string
	function			interface{}
	contexts			[]interface{}
	stack				StackInfo

	// compile
	tfunc				reflect.Type
	vfunc				reflect.Value
	vctxs				[]reflect.Value
	numIn				int
	numOut				int
	numArgs				int
	inputs				[]reflect.Value		// vctxs + vargs

	// state
	compiled			bool
}

func NewJITFunc(function interface{}, contexts... interface{}) *JITFunc {
	return &JITFunc{function:function, contexts:contexts, stack: Trace(1)}
}

func (h *JITFunc) Name() string { return h.name }

func (h *JITFunc) SetName(name string) { h.name = name }

func (h *JITFunc) Compile() error {
	if h.compiled { return nil }

	h.vfunc = reflect.ValueOf(h.function)
	h.tfunc = reflect.TypeOf(h.function)
	h.vctxs = make([]reflect.Value, len(h.contexts))
	if len(h.name) == 0 { h.name = h.tfunc.Name() }
	for i, ctx := range h.contexts { h.vctxs[i] = reflect.ValueOf(ctx)}
	if h.vfunc.Kind() != reflect.Func { return fmt.Errorf("JIT function %v with invalid function type %v", h.name, reflect.TypeOf(h.function)) }

	h.numIn = h.tfunc.NumIn()
	h.numOut = h.tfunc.NumOut()
	h.numArgs = h.numIn - len(h.vctxs)
	if h.numArgs < 0 { return fmt.Errorf("JIT function %v with invalid number of inputs %v, expect >= %v", h.name, h.numIn, len(h.vctxs))}
	h.inputs = make([]reflect.Value, h.numIn)
	for i, vctx := range h.vctxs { h.inputs[i] = vctx }

	h.compiled = true
	return nil
}

func (h *JITFunc) Call(args... interface{}) (ret interface{}, reterr error) {

	defer CatchCallback(func(info *PanicInfo){
		reterr = info.Error()
	})

	if err := h.Compile(); err != nil { return nil, err}

	// check input
	if h.numArgs != len(args) { return nil, fmt.Errorf("JIT function %v call error with invalid input number %v, expect %v.", h.name, len(args), h.numArgs) }

	// compile args
	st := len(h.vctxs)
	for i, arg := range args { h.inputs[st+i] = reflect.ValueOf(arg) }

	// call
	outs := h.vfunc.Call(h.inputs)
	switch h.numOut {
	case 0: return nil, nil
	case 1: return outs[0].Interface(), nil
	default:
		iouts := make([]interface{}, 0, h.numOut)
		for _, o := range outs { iouts = append(iouts, o.Interface()) }
		return iouts, nil
	}
}