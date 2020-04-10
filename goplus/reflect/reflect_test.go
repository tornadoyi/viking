package reflect

import (
	"strings"
	"testing"
)

type TestStruct struct {
	sInt32			int32
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

	InnerStruct struct{
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


func TestDFS(t *testing.T) {

	ts := &struct {
		Test1 struct{
			A1	*int
			b1	float64
		}

		Test2 struct{
			Test21 struct{
				A21	*int
				b21	float64
			}

			Test22 struct{
				A22	*int
				b22	float64
			}
		}
	}{}

	checks := map[string]string{
		"Test1": "Test1",
		"Test2": "Test2",
		"Test21": "Test2/Test21",
		"Test22": "Test2/Test22",
		"A1":	"Test1/A1",
		"b1":	"Test1/b1",
		"A21":	"Test2/Test21/A21",
		"b21":	"Test2/Test21/b21",
		"A22":	"Test2/Test22/A22",
		"b22":	"Test2/Test22/b22",
	}

	paths := make(map[string]string)

	DFS(ts, func(node Node) {
		sf := node.StructField()
		if sf == nil { return }
		ps := make([]string, 0)
		for _, p := range node.FieldPath() { ps = append(ps, p.Name) }
		paths[sf.Name] = strings.Join(ps, "/")
	})

	if len(paths) != len(checks) { t.Error("Error count of paths") }

	for k, v := range checks {
		p, ok := paths[k]
		if !ok {
			t.Errorf("No path %v", k)
			continue
		}
		if p != v { t.Errorf("Error path %v, expect %v", p, v) }
	}
}


func TestGetAndSetPrivateValue(t *testing.T){

	s := &struct {
		private			int
	}{0}

	v := ValueOf(s).Elem().FieldByName("private")
	if v.CanInterface() { t.Fatal("why private member can be accessed ?") }
	if v.CanSet() { t.Fatal("why private member can be set ?") }
	newPrivate := ValueOf(int(1))
	SetValue(v, newPrivate)
	if GetValue(v).Interface().(int) != 1 { t.Fail() }
}


func TestReconstruct(t *testing.T) {
	ts := &TestStruct{
		Slice: []*TestStruct{&TestStruct{Int8: 8}, &TestStruct{Bool: false}},
		Map: map[string]int{"Key1": 1, "Key2": 2},
	}
	if _, err := Reconstruct(ts); err != nil { t.Fatal(err) }
}