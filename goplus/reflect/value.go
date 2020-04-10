package reflect

import (
	"github.com/tornadoyi/viking/goplus/runtime"
	"unsafe"
)


func CallValue(v Value, in []Value) (out []Value, err error) {
	defer runtime.CatchCallback(func(info *runtime.PanicInfo) {
		out = nil
		err = info.Error()
	})

	out = v.Call(in)
	return out, nil
}

// can get unexported value
func GetValue(v Value) Value {
	if v.CanInterface() { return v}
	return NewAt(v.Type(), unsafe.Pointer(v.UnsafeAddr())).Elem()
}

// can set unexported value
func SetValue(dst Value, src Value) {
	src = GetValue(src)
	if !dst.CanSet() { dst = NewAt(dst.Type(), unsafe.Pointer(dst.UnsafeAddr())).Elem() }
	dst.Set(src)
}