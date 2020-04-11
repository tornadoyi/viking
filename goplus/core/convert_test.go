package core

import (
	"testing"
)

func TestConvert(t *testing.T) {
	cases := []interface{}{
		int(0), int8(0), int16(0), int32(0), int64(0),
		uint(0), uint8(0), uint16(0), uint32(0), uint64(0),
		float32(1.2), float64(2.1),
		"1234", []byte("1234")}
	for _, c := range cases {
		if _, err := Int(c); err != nil { t.Fatal(err) }
		if _, err := Int8(c); err != nil { t.Fatal(err) }
		if _, err := Int16(c); err != nil { t.Fatal(err) }
		if _, err := Int32(c); err != nil { t.Fatal(err) }
		if _, err := Int64(c); err != nil { t.Fatal(err) }
		if _, err := UInt(c); err != nil { t.Fatal(err) }
		if _, err := UInt8(c); err != nil { t.Fatal(err) }
		if _, err := UInt16(c); err != nil { t.Fatal(err) }
		if _, err := UInt32(c); err != nil { t.Fatal(err) }
		if _, err := UInt64(c); err != nil { t.Fatal(err) }
		if _, err := Float32(c); err != nil { t.Fatal(err) }
		if _, err := Float64(c); err != nil { t.Fatal(err) }
		if _, err := String(c); err != nil { t.Fatal(err) }
		if _, err := BigBytes(c); err != nil { t.Fatal(err) }
		if _, err := LittleBytes(c); err != nil { t.Fatal(err) }
	}
}
