package reflect

import (
	"encoding/json"
	"fmt"
	"strings"
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

	InnerStruct struct{
		A 		*int
		B		int
	}


	privateString	string
	privateString2	string
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



type RefactorObject struct {
	Int64						int64
	String						string
}
func (h *RefactorObject) Refactor() interface{} {
	return struct {
		Int64_1					int64
		Int64_2					int64
		String1					string
		String2					string
	}{h.Int64, h.Int64+1, h.String, h.String}
}

func TestRefactor(t *testing.T) {

	t.Run("Types", func(t *testing.T) {
		ts := &TestStruct{
			Slice: []*TestStruct{&TestStruct{Int8: 8}, &TestStruct{Bool: false}},
			Map: map[string]int{"Key1": 1, "Key2": 2},
		}
		if _,  err := Refactor(ts); err != nil { t.Fatal(err) }
	})

	t.Run("Private", func(t *testing.T) {
		type Test struct {
			Public			string
			private			string
		}
		ts := &Test{"public value", "private value"}
		newts, err := Refactor(ts)
		if err != nil { t.Fatal(err) }
		if !strings.Contains(fmt.Sprintf("%v", newts), ts.Public) { t.Fail() }
		if !strings.Contains(fmt.Sprintf("%v", newts), ts.private) { t.Fail() }
	})


	t.Run("Option", func(t *testing.T) {
		type Test struct {
			Public			string
			private			string
		}
		ts := &Test{"public value", "private value"}
		newts, err := Refactor(ts, []RefactorOption{
			RefactorTitle(true),
			RefactorWithTag(false)}...)
		if err != nil { t.Fatal(err) }
		bs, err := json.Marshal(newts)
		if err != nil { t.Fatal(err) }
		if !strings.Contains(string(bs), "private") {t.Fail()}
	})

	t.Run("Tag", func(t *testing.T) {
		type Test struct {
			NoRefactor		string				`refactor:"-"`
			Public			string
			private			string				`refactor:"PrivateKey"`
		}
		ts := &Test{"NoRefactorValue", "public value", "private value"}
		newts, err := Refactor(ts)
		if err != nil { t.Fatal(err) }
		bs, err := json.Marshal(newts)
		if err != nil { t.Fatal(err) }
		if strings.Contains(string(bs), "NoRefactorValue") {t.Fatalf("tag \"-\" doesn't work")}
		if !strings.Contains(string(bs), "PrivateKey") {t.Fatalf("tag \"name\" doesn't work")}
	})

	t.Run("RefactorFunction", func(t *testing.T) {

		type Test struct {
			NoRefactor		string				`refactor:"-"`
			RefObject		*RefactorObject
		}
		ts := &Test{"NoRefactorValue", &RefactorObject{0, "hello"}}
		newts, err := Refactor(ts)
		if err != nil { t.Fatal(err) }
		bs, err := json.Marshal(newts)
		if err != nil { t.Fatal(err) }
		if !strings.Contains(string(bs), "Int64_1") {t.Fatalf("tag \"name\" doesn't work")}
	})
}
