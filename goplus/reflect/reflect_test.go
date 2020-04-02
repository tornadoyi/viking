package reflect

import (
	"testing"
)

type TestStruct struct {
	Bool			bool
	Int				int
	Int8			int8
	Int16			int16
	Int32			int16
	Int64			int64
	Uint			uint
	Uint8			uint8
	Uint16			uint16
	Uint32			uint32
	Uint64			uint64
	Uintptr			uintptr
	Float32			float32
	Float64			float64
	Complex64		complex64
	Complex128		complex128
	Array			[]string
	Chan			chan int
	Func			func(int, string) float32
	Interface		interface{}
	Map				map[string]int
	Slice			[]*TestStruct
	String			string

	InStruct struct{
		A 		*int
		B		int
	}
}


func TestInstantiate(t *testing.T) {
	var Bool			*bool
	var ts 				TestStruct

	if err := Instantiate(&ts); err != nil { t.Error(err) }
	if err := Instantiate(&Bool); err != nil { t.Error(err) }
}