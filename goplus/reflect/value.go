package reflect

import (
	"fmt"
	"github.com/tornadoyi/viking/goplus/runtime"
	"unsafe"
)

var (
	InvalidValue					= Value{}
)

func CallValue(v Value, in []Value) (out []Value, reterr error) {
	defer runtime.CatchCallback(func(err error) { out, reterr = nil, err })

	out = v.Call(in)
	return out, nil
}

// can get unexported value
func Readable(v Value) Value {
	if !v.IsValid() || v.CanInterface() { return v}
	if v.CanAddr() { return NewAt(v.Type(), unsafe.Pointer(v.UnsafeAddr())).Elem() }
	switch v.Kind() {
	case Bool: return ValueOf(v.Bool())
	case String: return ValueOf(v.String())
	case Int, Int8, Int16, Int32, Int64: return ValueOf(v.Int())
	case Uint, Uintptr, Uint8, Uint16, Uint32, Uint64: return ValueOf(v.Uint())
	case Float32, Float64: return ValueOf(v.Float())
	case Complex64, Complex128: return ValueOf(v.Complex())
	default: return InvalidValue
	}
}

// can set unexported value
func SetValue(dst Value, src Value) error {
	if !dst.IsValid() { return fmt.Errorf("invalid destination value") }
	src = Readable(src)
	if !src.IsValid() { return fmt.Errorf("source value %v can not readable", src) }
	if !dst.CanSet() {
		if !dst.CanAddr() { return fmt.Errorf("unddressed destionation value %v", dst)}
		dst = NewAt(dst.Type(), unsafe.Pointer(dst.UnsafeAddr())).Elem()
	}
	dst.Set(src)
	return nil
}