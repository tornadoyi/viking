package core

import (
	"fmt"
	"os"
	"strings"
	"viking/goplus/runtime"
	"viking/log"
)

var (
	catchErrCallback = func(info *PanicInfo){
		if info == nil { return }
		msgs := make([]string, 0)
		msgs = append(
			msgs,
			fmt.Sprintf("A panic occured as below"),
			fmt.Sprintf("%v", info.stack),
			fmt.Sprintf("error: %v\n", info.Error()),
		)
		log.Critical(os.Stderr, strings.Join(msgs, "\n"))
	}
)

type PanicInfo struct {
	error		interface{}
	stack		runtime.Stack
}

func (h *PanicInfo) Error() interface{} { return h.error}

func (h *PanicInfo) Stack() runtime.Stack { return h.stack}

func Catch(){
	err := recover()
	if err == nil { return }
	info := &PanicInfo{err, runtime.Trace(4)}
	if catchErrCallback != nil { catchErrCallback(info) }
}


func CatchCallback(cb func(*PanicInfo)) {
	err := recover()
	if err == nil { return }
	info := &PanicInfo{err, runtime.Trace(4)}
	if cb != nil { cb(info) }
	if catchErrCallback != nil { catchErrCallback(info) }
}