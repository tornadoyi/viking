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
func Access(v Value) Value {
	if !v.IsValid() || v.CanInterface() { return v}
	if v.CanAddr() { return NewAt(v.Type(), unsafe.Pointer(v.UnsafeAddr())).Elem() }
	switch v.Kind() {
	case Bool: return ValueOf(v.Bool())
	case String: return ValueOf(v.String())
	case Int:	return ValueOf(int(v.Int()))
	case Int8:  return ValueOf(int8(v.Int()))
	case Int16: return ValueOf(int16(v.Int()))
	case Int32: return ValueOf(int32(v.Int()))
	case Int64: return ValueOf(int64(v.Int()))
	case Uint:  return ValueOf(uint(v.Uint()))
	case Uint8:  return ValueOf(uint8(v.Uint()))
	case Uint16: return ValueOf(uint16(v.Uint()))
	case Uint32: return ValueOf(uint32(v.Uint()))
	case Uint64: return ValueOf(uint64(v.Uint()))
	case Float32: return ValueOf(float32(v.Float()))
	case Float64: return ValueOf(float64(v.Float()))
	case Complex64: return ValueOf(complex64(v.Complex()))
	case Complex128: return ValueOf(v.Complex())
	default: return InvalidValue
	}
}

// can set unexported value
func SetValue(dst Value, src Value) error {
	if !dst.IsValid() { return fmt.Errorf("invalid destination value") }
	src = Access(src)
	if !src.IsValid() { return fmt.Errorf("source value %v can not readable", src) }
	if !dst.CanSet() {
		if !dst.CanAddr() { return fmt.Errorf("unddressed destionation value %v", dst)}
		dst = NewAt(dst.Type(), unsafe.Pointer(dst.UnsafeAddr())).Elem()
	}
	dst.Set(src)
	return nil
}