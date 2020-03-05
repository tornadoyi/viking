package goplus

import "runtime"

type CodeStack struct {
	pc			uintptr
	file		string
	line		int
	function	*runtime.Func
}

func (h *CodeStack) PC() uintptr { return h.pc}

func (h *CodeStack) File() string { return h.file}

func (h *CodeStack) Line() int { return h.line}

func (h *CodeStack) Funtion() *runtime.Func { return h.function}

func TraceStack() []*CodeStack {

	stacks := make([]*CodeStack, 0)

	for skip := 0; ; skip++ {
		pc, file, line, ok := runtime.Caller(skip);
		if !ok { break }
		stacks = append(stacks, &CodeStack{pc, file, line, runtime.FuncForPC(pc)})
	}

	return stacks
}


