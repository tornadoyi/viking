package reflect

import "github.com/tornadoyi/viking/goplus/runtime"


func CallValue(v Value, in []Value) (out []Value, err error) {
	defer runtime.CatchCallback(func(info *runtime.PanicInfo) {
		out = nil
		err = info.Error()
	})

	out = v.Call(in)
	return out, nil
}