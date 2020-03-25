package runtime

import "runtime"

var BlockProfile = runtime.BlockProfile
var Breakpoint = runtime.Breakpoint
var CPUProfile = runtime.CPUProfile
var Caller = runtime.Caller
var Callers = runtime.Callers
var GC = runtime.GC
var GOMAXPROCS = runtime.GOMAXPROCS
var GOROOT = runtime.GOROOT
var Goexit = runtime.Goexit
var GoroutineProfile = runtime.GoroutineProfile
var Gosched = runtime.Gosched
var KeepAlive = runtime.KeepAlive
var LockOSThread = runtime.LockOSThread
var MemProfile = runtime.MemProfile
var MutexProfile = runtime.MutexProfile
var NumCPU = runtime.NumCPU
var NumCgoCall = runtime.NumCgoCall
var NumGoroutine = runtime.NumGoroutine
var ReadMemStats = runtime.ReadMemStats
var ReadTrace = runtime.ReadTrace
var SetBlockProfileRate = runtime.SetBlockProfileRate
var SetCPUProfileRate = runtime.SetCPUProfileRate
var SetCgoTraceback = runtime.SetCgoTraceback
var SetFinalizer = runtime.SetFinalizer
var SetMutexProfileFraction = runtime.SetMutexProfileFraction
var Stack = runtime.Stack
var StartTrace = runtime.StartTrace
var StopTrace = runtime.StopTrace
var ThreadCreateProfile = runtime.ThreadCreateProfile
var UnlockOSThread = runtime.UnlockOSThread
var Version = runtime.Version


type BlockProfileRecord = runtime.BlockProfileRecord
type Error = runtime.Error
type Frame = runtime.Frame
type Frames = runtime.Frames
type Func = runtime.Func
type MemProfileRecord = runtime.MemProfileRecord
type MemStats = runtime.MemStats
type StackRecord = runtime.StackRecord
type TypeAssertionError = runtime.TypeAssertionError
