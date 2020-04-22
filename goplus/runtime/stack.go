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
	stackStringMethod	func(h StackInfo) string
)

func SetStackStringMethod(f func(h StackInfo) string){ stackStringMethod = f }


type StackInfo []uintptr

func Trace(skip int) StackInfo{
	var pcs [MAX_STACK_DEPTH]uintptr
	n := runtime.Callers(1+skip, pcs[:])
	var s StackInfo = pcs[0:n]
	return s
}

func (h StackInfo) Frames() []runtime.Frame{
	frames := runtime.CallersFrames(h)
	frameList := make([]runtime.Frame, 0, len(h))
	for {
		f, ok := frames.Next()
		if !ok { break }
		frameList = append(frameList, f)
	}
	return frameList
}

func (h StackInfo) String() string{
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
