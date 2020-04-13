package runtime

import (
	"fmt"
	"github.com/tornadoyi/viking/log"
)

var (
	catchErrCallback = func(err error){
		log.Error("A panic occurred as below\n" + ErrorWithStack(err))
	}
)


func Catch(){
	err := recover()
	if err == nil { return }
	perr := &PanicError{fmt.Errorf("%v", err), Trace(3)}
	if catchErrCallback != nil { catchErrCallback(perr) }
}


func CatchCallback(cb func(error)) {
	err := recover()
	if err == nil { return }
	perr := &PanicError{fmt.Errorf("%v", err), Trace(3)}
	//if catchErrCallback != nil { catchErrCallback(perr) }
	if cb != nil { cb(perr) }
}


func ErrorWithStack(err error) string{
	perr, ok := err.(*PanicError)
	if !ok { return err.Error()}
	return err.Error() + "\n" + fmt.Sprintf("%v\n", perr.stack)
}


type PanicError struct {
	error error
	stack StackInfo
}

func (h *PanicError) Error() string { return h.error.Error()}

func (h *PanicError) Stack() StackInfo { return h.stack}

