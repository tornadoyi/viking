package reflect

import (
	"fmt"
	"strings"
)

func Refactor(obj interface{}, opt... RefactorOption) (ret interface{}, err error){
	cfg := &RefactorConfig{
		false,
		true,
	}
	for _, o := range opt { o.apply(cfg) }
	o := ValueOf(obj)
	if o.Kind() != Ptr { return nil, fmt.Errorf("refactoring object must be a pointer")}
	v := refactor(o, cfg)
	return v.Interface(), nil
}

func refactor(o Value, cfg *RefactorConfig) Value {
	switch o.Kind() {
	case Map: return refactorMap(o, cfg)
	case Slice, Array: return refactorSlice(o, cfg)
	case Chan: return refactorChan(o, cfg)
	case Struct: return refactorStruct(o, cfg)
	case Ptr: return refactorPtr(o, cfg)
	}
	// default
	ret := New(o.Type()).Elem()
	SetValue(ret, o)
	return ret
}

func refactorPtr(o Value, cfg *RefactorConfig) Value {
	if o.IsNil() { return New(o.Type()).Elem() }
	return refactor(o.Elem(), cfg).Addr()
}

func refactorSlice(o Value, cfg *RefactorConfig) Value {
	if o.Len() == 0 { return MakeSlice(o.Type(), 0, o.Cap()) }
	s := make([]Value, o.Len())
	for i:=0; i<o.Len(); i++{
		v := o.Index(i)
		s[i] = refactor(v, cfg)
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

func refactorMap(o Value, cfg *RefactorConfig) Value {
	if o.Len() == 0 { return MakeMap(o.Type()) }
	kvs := make(map[Value]Value, o.Len())
	it := o.MapRange()
	for it.Next() {
		k, v := it.Key(), it.Value()
		nk := refactor(k, cfg)
		nv := refactor(v, cfg)
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

func refactorChan(o Value, cfg *RefactorConfig) Value {
	return MakeChan(o.Type(), o.Cap())
}

func refactorStruct(o Value, cfg *RefactorConfig) Value {
	// format: 1. ignore/name
	parseTag := func(stag string) (bool, string) {
		ignore, name := false, ""
		for i, t := range strings.Split(stag, ",") {
			if len(t) == 0 { continue }
			switch i {
			case 0: if t == "-" { ignore = true } else { name = t }
			}
		}
		return ignore, name
	}

	fields := make([]StructField, 0, o.NumField())
	nvs := make([]Value, 0, o.NumField())
	tp := o.Type()
	for i:=0; i<tp.NumField(); i++{
		f := tp.Field(i)

		// field name
		fieldName := f.Name
		if cfg.Title { fieldName = strings.Title(fieldName) }

		// field tag
		tag := StructTag("")
		if cfg.WithTag { tag = f.Tag }

		// parse tag
		if t, ok := f.Tag.Lookup("refactor"); ok {
			ignore, name := parseTag(t)
			if ignore { continue }
			if len(name) != 0 { fieldName = name }
		}

		// pkgpath
		pkgpath := f.PkgPath
		c := int(fieldName[0])
		if 65 <= c && c <= 80 { pkgpath = "" }

		// add field
		v := o.Field(i)
		nv := refactor(v, cfg)
		nvs = append(nvs, nv)
		fields = append(fields, StructField{
			fieldName,
			pkgpath,
			nv.Type(),
			tag,
			0,
			nil,
			f.Anonymous,
		})
	}
	ret := New(StructOf(fields)).Elem()
	for i, _ := range fields {
		SetValue(ret.Field(i), nvs[i])
	}
	return ret
}




type RefactorConfig struct {
	Title					bool						`json:"title" yaml:"title"`
	WithTag					bool						`json:"with_tag" yaml:"with_tag"`
}


type RefactorOption struct {
	apply	func(*RefactorConfig)
}


func RefactorTitle(title bool) RefactorOption{
	return RefactorOption{func(c *RefactorConfig){
		c.Title = title
	}}
}

func RefactorWithTag(withtag bool) RefactorOption{
	return RefactorOption{func(c *RefactorConfig){
		c.WithTag = withtag
	}}
}