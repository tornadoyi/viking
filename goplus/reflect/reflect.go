package reflect

import (
	_reflect "reflect"
)


const (
	Invalid Kind = iota
	Bool
	Int
	Int8
	Int16
	Int32
	Int64
	Uint
	Uint8
	Uint16
	Uint32
	Uint64
	Uintptr
	Float32
	Float64
	Complex64
	Complex128
	Array
	Chan
	Func
	Interface
	Map
	Ptr
	Slice
	String
	Struct
	UnsafePointer
)

type ChanDir = _reflect.ChanDir
type Kind = _reflect.Kind
type MapIter = _reflect.MapIter
type Method = _reflect.Method
type SelectCase = _reflect.SelectCase
type SelectDir = _reflect.SelectDir
type SliceHeader = _reflect.SliceHeader
type StringHeader = _reflect.StringHeader
type StructField = _reflect.StructField
type StructTag = _reflect.StructTag
type Type = _reflect.Type
type Value = _reflect.Value
type ValueError = _reflect.ValueError



var Copy = _reflect.Copy
var DeepEqual = _reflect.DeepEqual
var Swapper = _reflect.Swapper

var ArrayOf = _reflect.ArrayOf
var ChanOf = _reflect.ChanOf
var FuncOf = _reflect.FuncOf
var MapOf = _reflect.MapOf
var PtrTo = _reflect.PtrTo
var SliceOf = _reflect.SliceOf
var StructOf = _reflect.StructOf
var TypeOf = _reflect.TypeOf
var Append = _reflect.Append
var AppendSlice = _reflect.AppendSlice
var Indirect = _reflect.Indirect
var MakeChan = _reflect.MakeChan
var MakeFunc = _reflect.MakeFunc
var MakeMap = _reflect.MakeMap
var MakeMapWithSize = _reflect.MakeMapWithSize
var MakeSlice = _reflect.MakeSlice
var New = _reflect.New
var NewAt = _reflect.NewAt
var Select = _reflect.Select
var ValueOf = _reflect.ValueOf
var Zero = _reflect.Zero

