package goplus

import (
	"fmt"
	"os"
	"strings"
)

var (
	catchErrCallback = func(info *PanicInfo){
		if info == nil { return }
		msgs := make([]string, 0, 2*len(info.stacks))
		msgs = append(msgs, fmt.Sprintf("A panic occured as below"))

		stacks := info.Stacks()
		for i:=len(stacks)-1; i>2; i--{
			s := stacks[i]
			fname := "<nil>"
			if s.Funtion() != nil { fname = s.Funtion().Name() }
			msgs = append(msgs, fmt.Sprintf(
				"%v [0x%x]",
				 fname, s.PC(),
			))
			msgs = append(msgs, fmt.Sprintf(
				"\t%v:%v",
				s.File(), s.Line(),
			))
		}
		msgs = append(msgs, fmt.Sprintf("error: %v\n", info.Error()))
		fmt.Fprint(os.Stderr, strings.Join(msgs, "\n"))
	}
)

type PanicInfo struct {
	error		interface{}
	stacks		[]*CodeStack
}

func (h *PanicInfo) Error() interface{} { return h.error}

func (h *PanicInfo) Stacks() []*CodeStack { return h.stacks}

func Catch(){
	err := recover()
	fmt.Println(err)
	if err == nil { return }
	if catchErrCallback != nil {
		catchErrCallback(&PanicInfo{err, TraceStack()})
	}
}


