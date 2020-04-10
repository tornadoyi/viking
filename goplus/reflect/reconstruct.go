package reflect

import (
	"fmt"
)

func Reconstruct(obj interface{}) (ret interface{}, err error){
	o := ValueOf(obj)
	if o.Kind() != Ptr { return nil, fmt.Errorf("Reconstructed object must be a pointer")}
	v := reconstruct(o)
	return v.Interface(), nil
}

func reconstruct(o Value) Value {
	switch o.Kind() {
	case Map: return reconstructMap(o)
	case Slice, Array: return reconstructSlice(o)
	case Chan: return reconstructChan(o)
	case Struct: return reconstructStruct(o)
	case Ptr: return reconstructPtr(o)
	}
	// default
	ret := New(o.Type()).Elem()
	SetValue(ret, o)
	return ret
}

func reconstructPtr(o Value) Value {
	if o.IsNil() { return New(o.Type()).Elem() }
	return reconstruct(o.Elem()).Addr()
}

func reconstructSlice(o Value) Value {
	if o.Len() == 0 { return MakeSlice(o.Type(), 0, o.Cap()) }
	s := make([]Value, o.Len())
	for i:=0; i<o.Len(); i++{
		v := o.Index(i)
		s[i] = reconstruct(v)
	}
	var tp Type
	for _, v := range s {
		if tp == nil { tp = v.Type(); continue }
		if v.Type() == tp { continue }
		var t interface{}
		tp = TypeOf(t)
		break
	}
	ret := MakeSlice(SliceOf(tp), 0, o.Cap())
	ret = Append(ret, s...)
	return ret
}

func reconstructMap(o Value) Value {
	if o.Len() == 0 { return MakeMap(o.Type()) }
	kvs := make(map[Value]Value, o.Len())
	it := o.MapRange()
	for it.Next() {
		k, v := it.Key(), it.Value()
		nk := reconstruct(k)
		nv := reconstruct(v)
		kvs[nk] = nv
	}
	var ktype, vtype Type
	for k, _ := range kvs {
		if ktype == nil {
			ktype = k.Type();
			continue
		}
		if ktype == k.Type() {
			continue
		}
		var t interface{}
		ktype = TypeOf(t)
		break
	}
	for _, v := range kvs {
		if vtype == nil {
			vtype = v.Type();
			continue
		}
		if vtype == v.Type() {
			continue
		}
		var t interface{}
		vtype = TypeOf(t)
		break
	}

	ret := MakeMapWithSize(MapOf(ktype, vtype), o.Len())
	for k, v := range kvs { ret.SetMapIndex(k, v) }
	return ret
}

func reconstructChan(o Value) Value {
	return MakeChan(o.Type(), o.Cap())
}

func reconstructStruct(o Value) Value {
	fields := make([]StructField, o.NumField())
	nvs := make([]Value, o.NumField())
	tp := o.Type()
	for i:=0; i<tp.NumField(); i++{
		f := tp.Field(i)
		v := o.Field(i)
		nvs[i] = reconstruct(v)
		fields[i] = StructField{
			f.Name,//strings.Title(f.Name),
			f.PkgPath,
			nvs[i].Type(),
			"",
			0,
			nil,
			f.Anonymous,
		}
	}
	ret := New(StructOf(fields)).Elem()
	for i, _ := range fields {
		SetValue(ret.Field(i), nvs[i])
	}
	return ret
}

