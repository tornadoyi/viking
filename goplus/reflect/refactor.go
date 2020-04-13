package reflect

import (
	"encoding/json"
	"fmt"
	"gopkg.in/yaml.v2"
	"strings"
)


func RefactorJson(obj interface{}, opt... RefactorOption) ([]byte, error){
	opt = append(opt, RefactorTitle(true), RefactorMarshallKinds())
	st, err := Refactor(obj, opt...)
	if err != nil { return nil, err}
	data, err := json.Marshal(st)
	if err != nil { return nil, err}
	return data, nil
}

func RefactorYaml(obj interface{}, opt... RefactorOption) ([]byte, error){
	opt = append(opt, RefactorTitle(true), RefactorMarshallKinds())
	st, err := Refactor(obj, opt...)
	if err != nil { return nil, err}
	data, err := yaml.Marshal(st)
	if err != nil { return nil, err}
	return data, nil
}

func Refactor(obj interface{}, opt... RefactorOption) (ret interface{}, err error){
	// check type
	o := ValueOf(obj)
	if o.Kind() != Ptr { return nil, fmt.Errorf("refactoring object must be a pointer")}

	// init config
	cfg := &RefactorConfig{}
	for _, o := range opt { o.apply(cfg) }

	// refactor
	if v := refactor(o, cfg); !v.IsValid() { return nil, nil } else { return v.Interface(), nil }
}

func refactor(o Value, cfg *RefactorConfig) Value {
	// refactor function
	if v := callMethod(o, "Refactor"); v.IsValid() { return v}

	// check valid kind
	if !validKind(o, cfg) { return InvalidValue}

	switch o.Kind() {
	case Map: return refactorMap(o, cfg)
	case Slice, Array: return refactorSlice(o, cfg)
	case Chan: return refactorChan(o, cfg)
	case Struct: return refactorStruct(o, cfg)
	case Interface: return refactorInterface(o, cfg)
	case Ptr: return refactorPtr(o, cfg)
	}
	// default
	ret := New(o.Type()).Elem()
	SetValue(ret, o)
	return ret
}

func refactorPtr(o Value, cfg *RefactorConfig) Value {
	if o.IsNil() { return New(o.Type()).Elem() }
	if v := refactor(o.Elem(), cfg); !v.IsValid() { return v} else { return v.Addr()}
}

func refactorInterface(o Value, cfg *RefactorConfig) Value {
	if o.IsNil() { return New(o.Type()).Elem() }
	return refactor(o.Elem(), cfg)
}

func refactorSlice(o Value, cfg *RefactorConfig) Value {
	s := make([]Value, 0, o.Len())
	for i:=0; i<o.Len(); i++{
		v := o.Index(i)
		if nv := refactor(v, cfg); !nv.IsValid() { continue } else { s = append(s, nv) }
	}
	if len(s) == 0 { return MakeSlice(o.Type(), 0, o.Cap()) }

	var tp Type
	for _, v := range s {
		if tp == nil { tp = v.Type(); continue }
		if v.Type() == tp { continue }
		tp = nil
		break
	}

	var ret Value
	if tp == nil {
		var a []interface{}
		ret = MakeSlice(TypeOf(a), 0, o.Cap())
	} else {
		ret = MakeSlice(SliceOf(tp), 0, o.Cap())
	}

	ret = Append(ret, s...)
	return ret
}

func refactorMap(o Value, cfg *RefactorConfig) Value {
	kvs := make(map[Value]Value, o.Len())
	it := o.MapRange()
	for it.Next() {
		k, v := it.Key(), it.Value()
		nk := refactor(k, cfg)
		if !nk.IsValid() { continue }
		nv := refactor(v, cfg)
		if !nv.IsValid() { continue }
		kvs[nk] = nv
	}
	if len(kvs) == 0 { return MakeMap(o.Type()) }

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
	// format: 1. ignore/name 	2. func
	parseTag := func(stag string) (bool, string, string) {
		ignore, name, fname := false, "", ""
		for i, t := range strings.Split(stag, ",") {
			if len(t) == 0 { continue }
			switch i {
			case 0: if t == "-" { ignore = true } else { name = t }
			case 1: fname = t
			}
		}
		return ignore, name, fname
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
		if !cfg.WithoutTag { tag = f.Tag }

		// parse tag
		method := ""
		if t, ok := f.Tag.Lookup("refactor"); ok {
			ignore, name, fname := parseTag(t)
			if ignore { continue }
			if len(name) != 0 { fieldName = name }
			method = fname
		}

		// pkgpath
		pkgpath := f.PkgPath
		c := int(fieldName[0])
		if 65 <= c && c <= 90 { pkgpath = "" }

		// new value
		v := o.Field(i)
		var nv Value
		if len(method) > 0 { nv = callMethod(Access(v), method) }
		if !nv.IsValid() { nv = refactor(v, cfg) }
		if !nv.IsValid() { continue }

		// add field
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


func callMethod(o Value, name string) Value {
	if !o.IsValid() { return InvalidValue}
	f := o.MethodByName(name)
	if !f.IsValid() || f.IsZero() || !f.CanInterface() { return InvalidValue}

	outs := f.Call(nil)
	if len(outs) != 1 { return InvalidValue}
	return outs[0]
}


func validKind(o Value, cfg *RefactorConfig) bool {
	if !cfg.ContainKind(o.Kind()) { return false}
	if o.Kind() == Ptr {
		if o.IsNil() { if cfg.SkipNil { return false } else { return true} }
		return validKind(o.Elem(), cfg)
	}
	return true
}




type RefactorConfig struct {
	Title					bool						`json:"title" yaml:"title"`
	WithoutTag				bool						`json:"without_tag" yaml:"without_tag"`
	Kinds					[]Kind						`json:"kinds" yaml:"kinds"`
	kindDict				map[Kind]bool
	SkipNil					bool						`json:"skip_nil" yaml:"skip_nil"`
}

func (h *RefactorConfig) ContainKind(k Kind) bool {
	if len(h.Kinds) == 0 { return true}
	if h.kindDict == nil {
		h.kindDict = make(map[Kind]bool)
		for _, k := range h.Kinds { h.kindDict[k] = true }
	}
	_, ok := h.kindDict[k]
	return ok
}



type RefactorOption struct {
	apply	func(*RefactorConfig)
}


func RefactorTitle(title bool) RefactorOption{
	return RefactorOption{func(c *RefactorConfig){
		c.Title = title
	}}
}

func RefactorWithoutTag() RefactorOption{
	return RefactorOption{func(c *RefactorConfig){
		c.WithoutTag = true
	}}
}

func RefactorKinds(kinds []Kind) RefactorOption {
	return RefactorOption{func(c *RefactorConfig){
		c.Kinds = append(c.Kinds, kinds...)
	}}
}


func RefactorMarshallKinds() RefactorOption {
	return RefactorOption{func(c *RefactorConfig){
		c.Kinds = append(c.Kinds,
			Bool,
			Int, Int8, Int16, Int32, Int64,
			Uint, Uint8, Uint16, Uint32, Uint64,
			Float32, Float64,
			Array, Slice, Map,
			String, Interface,  Struct,
			Ptr,
		)
	}}
}

func RefactorSkipNil() RefactorOption {
	return RefactorOption{func(c *RefactorConfig){
		c.SkipNil = true
	}}
}
