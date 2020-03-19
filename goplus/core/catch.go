package core

import (
	"errors"
	"fmt"
	"github.com/tornadoyi/viking/goplus/runtime"
	"github.com/tornadoyi/viking/log"
	"strings"
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
		log.Critical(strings.Join(msgs, "\n"))
	}
)

type PanicInfo struct {
	error		error
	stack		runtime.Stack
}

func (h *PanicInfo) Error() error { return h.error}

func (h *PanicInfo) Stack() runtime.Stack { return h.stack}

func Catch(){
	err := recover()
	if err == nil { return }
	info := &PanicInfo{errors.New(fmt.Sprintf("%v", err)), runtime.Trace(4)}
	if catchErrCallback != nil { catchErrCallback(info) }
}


func CatchCallback(cb func(*PanicInfo)) {
	err := recover()
	if err == nil { return }
	info := &PanicInfo{errors.New(fmt.Sprintf("%v", err)), runtime.Trace(4)}
	if catchErrCallback != nil { catchErrCallback(info) }
	if cb != nil { cb(info) }
}