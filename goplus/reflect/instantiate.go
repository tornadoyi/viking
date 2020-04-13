package reflect

import "github.com/tornadoyi/viking/goplus/runtime"

func Instantiate(obj interface{}) (reterr error){

	defer runtime.CatchCallback(func(err error) { reterr = err })

	var dispatch func(Type) Value
	var newStruct func(Type) Value

	assgin := func(dst, src Value) {
		if dst.Kind() != Ptr && src.Kind() == Ptr { src = src.Elem() }
		//fmt.Println(src.Type(), "=", dst.Type())
		dst.Set(src)
	}

	dispatch = func(t Type) Value {
		switch t.Kind() {
		case Map:
			return MakeMap(t)
		case Slice:
			return MakeSlice(t, 0, 0)
		case Chan:
			return MakeChan(t, 0)
		case Struct:
			return newStruct(t)
		case Ptr:
			return dispatch(t.Elem())
		default:
			return New(t)
		}
	}

	newStruct = func(t Type) Value {
		v := New(t)
		ve := v.Elem()
		for i:=0; i<ve.NumField(); i++{
			f := ve.Field(i)
			if !f.CanSet() { continue }
			assgin(f, dispatch(f.Type()))
		}
		return v
	}

	v := ValueOf(obj).Elem()
	assgin(v, dispatch(v.Type()))

	return nil
}


