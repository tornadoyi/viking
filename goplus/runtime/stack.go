package runtime

import (
	"fmt"
	"runtime"
	"strings"
)

const (
	MAX_STACK_DEPTH		= 64
)

var (
	stackStringMethod	func(h Stack) string
)

func SetStackStringMethod(f func(h Stack) string){ stackStringMethod = f }


type Stack []uintptr

func Trace(skip int) Stack{
	var pcs [MAX_STACK_DEPTH]uintptr
	n := runtime.Callers(skip, pcs[:])
	var s Stack = pcs[0:n]
	return s
}

func (h Stack) Frames() []runtime.Frame{
	frames := runtime.CallersFrames(h)
	frameList := make([]runtime.Frame, 0, len(h))
	for ; ; {
		f, ok := frames.Next()
		if !ok { break }
		frameList = append(frameList, f)
	}
	return frameList
}

func (h Stack) String() string{
	if stackStringMethod != nil { return stackStringMethod(h)}
	frames := h.Frames()
	msgs := make([]string, 0, len(frames))
	for i:=len(frames)-1; i>=0; i--{
		pc := (h)[i]
		frame := frames[i]

		msgs = append(msgs, fmt.Sprintf(
			"%v [0x%x]",
			frame.Function, pc,
		))
		msgs = append(msgs, fmt.Sprintf(
			"\t%v:%v",
			frame.File, frame.Line,
		))
	}
	return strings.Join(msgs, "\n")
}
