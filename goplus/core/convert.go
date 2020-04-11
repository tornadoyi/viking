package core

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"strconv"
	"strings"
)


var ErrNil = errors.New("convert failed, nil returned")

func Bool(v interface{}) (bool, error) {
	switch v := v.(type) {
	case bool : return v, nil
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64:
		return v != 0, nil
	case nil: return false, ErrNil
	case string:
		switch strings.ToLower(v) {
		case "true", "1", "yes": return true, nil
		case "false", "0", "no": return false, nil
		default: return false, fmt.Errorf("can not convert string %v to bool", v)
		}
	case []byte:
		n, err := strconv.ParseInt(string(v), 10, 0)
		if err != nil { return false, err}
		return n != 0, nil
	}
	return false, fmt.Errorf("unexpected type %T", v)
}

func Int(v interface{}) (int, error) {
	switch v := v.(type) {
	case bool : if v { return 1, nil } else { return 0, nil}
	case int : return v, nil
	case int8 : return int(v), nil
	case int16 : return int(v), nil
	case int32 : return int(v), nil
	case int64 : return int(v), nil
	case uint: return int(v), nil
	case uint8 : return int(v), nil
	case uint16 : return int(v), nil
	case uint32 : return int(v), nil
	case uint64 : return int(v), nil
	case float32 : return int(v), nil
	case float64 : return int(v), nil
	case nil: return 0, ErrNil
	case string:
		n, err := strconv.ParseInt(string(v), 10, 0)
		return int(n), err
	case []byte:
		n, err := strconv.ParseInt(string(v), 10, 0)
		return int(n), err
	}
	return 0, fmt.Errorf("unexpected type %T", v)
}

func Int8(v interface{}) (int8, error) {
	switch v := v.(type) {
	case bool : if v { return 1, nil } else { return 0, nil}
	case int : return int8(v), nil
	case int8 : return int8(v), nil
	case int16 : return int8(v), nil
	case int32 : return int8(v), nil
	case int64 : return int8(v), nil
	case uint : return int8(v), nil
	case uint8 : return int8(v), nil
	case uint16 : return int8(v), nil
	case uint32 : return int8(v), nil
	case uint64 : return int8(v), nil
	case float32 : return int8(v), nil
	case float64 : return int8(v), nil
	case nil: return 0, ErrNil
	case string:
		n, err := strconv.ParseInt(string(v), 10, 0)
		return int8(n), err
	case []byte:
		n, err := strconv.ParseInt(string(v), 10, 0)
		return int8(n), err
	}
	return 0, fmt.Errorf("unexpected type %T", v)
}

func Int16(v interface{}) (int16, error) {
	switch v := v.(type) {
	case bool : if v { return 1, nil } else { return 0, nil}
	case int : return int16(v), nil
	case int8 : return int16(v), nil
	case int16 : return int16(v), nil
	case int32 : return int16(v), nil
	case int64 : return int16(v), nil
	case uint: return int16(v), nil
	case uint8 : return int16(v), nil
	case uint16 : return int16(v), nil
	case uint32 : return int16(v), nil
	case uint64 : return int16(v), nil
	case float32 : return int16(v), nil
	case float64 : return int16(v), nil
	case nil: return 0, ErrNil
	case string:
		n, err := strconv.ParseInt(string(v), 10, 0)
		return int16(n), err
	case []byte:
		n, err := strconv.ParseInt(string(v), 10, 0)
		return int16(n), err
	}
	return 0, fmt.Errorf("unexpected type %T", v)
}

func Int32(v interface{}) (int32, error) {
	switch v := v.(type) {
	case bool : if v { return 1, nil } else { return 0, nil}
	case int : return int32(v), nil
	case int8 : return int32(v), nil
	case int16 : return int32(v), nil
	case int32 : return int32(v), nil
	case int64 : return int32(v), nil
	case uint: return int32(v), nil
	case uint8 : return int32(v), nil
	case uint16 : return int32(v), nil
	case uint32 : return int32(v), nil
	case uint64 : return int32(v), nil
	case float32 : return int32(v), nil
	case float64 : return int32(v), nil
	case nil: return 0, ErrNil
	case string:
		n, err := strconv.ParseInt(string(v), 10, 0)
		return int32(n), err
	case []byte:
		n, err := strconv.ParseInt(string(v), 10, 0)
		return int32(n), err
	}
	return 0, fmt.Errorf("unexpected type %T", v)
}

func Int64(v interface{}) (int64, error) {
	switch v := v.(type) {
	case bool : if v { return 1, nil } else { return 0, nil}
	case int : return int64(v), nil
	case int8 : return int64(v), nil
	case int16 : return int64(v), nil
	case int32 : return int64(v), nil
	case int64 : return int64(v), nil
	case uint: return int64(v), nil
	case uint8 : return int64(v), nil
	case uint16 : return int64(v), nil
	case uint32 : return int64(v), nil
	case uint64 : return int64(v), nil
	case float32 : return int64(v), nil
	case float64 : return int64(v), nil
	case nil: return 0, ErrNil
	case string:
		n, err := strconv.ParseInt(string(v), 10, 0)
		return int64(n), err
	case []byte:
		n, err := strconv.ParseInt(string(v), 10, 0)
		return int64(n), err
	}
	return 0, fmt.Errorf("unexpected type %T", v)
}

func UInt(v interface{}) (uint, error) {
	switch v := v.(type) {
	case bool : if v { return 1, nil } else { return 0, nil}
	case int : return uint(v), nil
	case int8 : return uint(v), nil
	case int16 : return uint(v), nil
	case int32 : return uint(v), nil
	case int64 : return uint(v), nil
	case uint: return uint(v), nil
	case uint8 : return uint(v), nil
	case uint16 : return uint(v), nil
	case uint32 : return uint(v), nil
	case uint64 : return uint(v), nil
	case float32 : return uint(v), nil
	case float64 : return uint(v), nil
	case nil: return 0, ErrNil
	case string:
		n, err := strconv.ParseInt(string(v), 10, 0)
		return uint(n), err
	case []byte:
		n, err := strconv.ParseInt(string(v), 10, 0)
		return uint(n), err
	}
	return 0, fmt.Errorf("unexpected type %T", v)
}

func UInt8(v interface{}) (uint8, error) {
	switch v := v.(type) {
	case bool : if v { return 1, nil } else { return 0, nil}
	case int : return uint8(v), nil
	case int8 : return uint8(v), nil
	case int16 : return uint8(v), nil
	case int32 : return uint8(v), nil
	case int64 : return uint8(v), nil
	case uint: return uint8(v), nil
	case uint8 : return uint8(v), nil
	case uint16 : return uint8(v), nil
	case uint32 : return uint8(v), nil
	case uint64 : return uint8(v), nil
	case float32 : return uint8(v), nil
	case float64 : return uint8(v), nil
	case nil: return 0, ErrNil
	case string:
		n, err := strconv.ParseInt(string(v), 10, 0)
		return uint8(n), err
	case []byte:
		n, err := strconv.ParseInt(string(v), 10, 0)
		return uint8(n), err
	}
	return 0, fmt.Errorf("unexpected type %T", v)
}

func UInt16(v interface{}) (uint16, error) {
	switch v := v.(type) {
	case bool : if v { return 1, nil } else { return 0, nil}
	case int : return uint16(v), nil
	case int8 : return uint16(v), nil
	case int16 : return uint16(v), nil
	case int32 : return uint16(v), nil
	case int64 : return uint16(v), nil
	case uint: return uint16(v), nil
	case uint8 : return uint16(v), nil
	case uint16 : return uint16(v), nil
	case uint32 : return uint16(v), nil
	case uint64 : return uint16(v), nil
	case float32 : return uint16(v), nil
	case float64 : return uint16(v), nil
	case nil: return 0, ErrNil
	case string:
		n, err := strconv.ParseInt(string(v), 10, 0)
		return uint16(n), err
	case []byte:
		n, err := strconv.ParseInt(string(v), 10, 0)
		return uint16(n), err
	}
	return 0, fmt.Errorf("unexpected type %T", v)
}

func UInt32(v interface{}) (uint32, error) {
	switch v := v.(type) {
	case bool : if v { return 1, nil } else { return 0, nil}
	case int : return uint32(v), nil
	case int8 : return uint32(v), nil
	case int16 : return uint32(v), nil
	case int32 : return uint32(v), nil
	case int64 : return uint32(v), nil
	case uint: return uint32(v), nil
	case uint8 : return uint32(v), nil
	case uint16 : return uint32(v), nil
	case uint32 : return uint32(v), nil
	case uint64 : return uint32(v), nil
	case float32 : return uint32(v), nil
	case float64 : return uint32(v), nil
	case nil: return 0, ErrNil
	case string:
		n, err := strconv.ParseInt(string(v), 10, 0)
		return uint32(n), err
	case []byte:
		n, err := strconv.ParseInt(string(v), 10, 0)
		return uint32(n), err
	}
	return 0, fmt.Errorf("unexpected type %T", v)
}

func UInt64(v interface{}) (uint64, error) {
	switch v := v.(type) {
	case bool : if v { return 1, nil } else { return 0, nil}
	case int : return uint64(v), nil
	case int8 : return uint64(v), nil
	case int16 : return uint64(v), nil
	case int32 : return uint64(v), nil
	case int64 : return uint64(v), nil
	case uint: return uint64(v), nil
	case uint8 : return uint64(v), nil
	case uint16 : return uint64(v), nil
	case uint32 : return uint64(v), nil
	case uint64 : return uint64(v), nil
	case float32 : return uint64(v), nil
	case float64 : return uint64(v), nil
	case nil: return 0, ErrNil
	case string:
		n, err := strconv.ParseInt(string(v), 10, 0)
		return uint64(n), err
	case []byte:
		n, err := strconv.ParseInt(string(v), 10, 0)
		return uint64(n), err
	}
	return 0, fmt.Errorf("unexpected type %T", v)
}

func Float32(v interface{}) (float32, error) {
	switch v := v.(type) {
	case bool : if v { return 1, nil } else { return 0, nil}
	case int : return float32(v), nil
	case int8 : return float32(v), nil
	case int16 : return float32(v), nil
	case int32 : return float32(v), nil
	case int64 : return float32(v), nil
	case uint: return float32(v), nil
	case uint8 : return float32(v), nil
	case uint16 : return float32(v), nil
	case uint32 : return float32(v), nil
	case uint64 : return float32(v), nil
	case float32 : return float32(v), nil
	case float64 : return float32(v), nil
	case nil: return 0, ErrNil
	case string:
		n, err := strconv.ParseInt(string(v), 10, 0)
		return float32(n), err
	case []byte:
		n, err := strconv.ParseInt(string(v), 10, 0)
		return float32(n), err
	}
	return 0, fmt.Errorf("unexpected type %T", v)
}

func Float64(v interface{}) (float64, error) {
	switch v := v.(type) {
	case bool : if v { return 1, nil } else { return 0, nil}
	case int : return float64(v), nil
	case int8 : return float64(v), nil
	case int16 : return float64(v), nil
	case int32 : return float64(v), nil
	case int64 : return float64(v), nil
	case uint: return float64(v), nil
	case uint8 : return float64(v), nil
	case uint16 : return float64(v), nil
	case uint32 : return float64(v), nil
	case uint64 : return float64(v), nil
	case float32 : return float64(v), nil
	case float64 : return float64(v), nil
	case nil: return 0, ErrNil
	case string:
		n, err := strconv.ParseFloat(string(v), 64)
		return float64(n), err
	case []byte:
		n, err := strconv.ParseFloat(string(v), 64)
		return float64(n), err
	}
	return 0, fmt.Errorf("unexpected type %T", v)
}

func String(v interface{}) (string, error) {
	switch v := v.(type) {
	case nil: return "", ErrNil
	case string: return v, nil
	case []byte: return string(v), nil
	case int, int8, int16, int32, int64:
		if n, err := Int64(v); err != nil { return "", err } else {
			return strconv.FormatInt(n, 10), nil
		}
	case uint, uint8, uint16, uint32, uint64:
		if n, err := UInt64(v); err != nil { return "", err } else {
			return strconv.FormatUint(n, 10), nil
		}
	case float32, float64:
		if n, err := Float64(v); err != nil { return "", err } else {
			return strconv.FormatFloat(n, 'f', -1, 64), nil
		}
	default:
		return fmt.Sprintf("%v", v), nil
	}
}

func Bytes(v interface{}, order binary.ByteOrder) ([]byte, error) {
	buffer := bytes.NewBuffer([]byte{})
	switch v.(type) {
	case int: v = int64(v.(int))
	case uint: v = uint64(v.(uint))
	case string: return []byte(v.(string)), nil
	case []byte: return v.([]byte), nil
	}
	if err := binary.Write(buffer, order, v); err != nil { return nil, err}
	return buffer.Bytes(), nil
}

func BigBytes(v interface{}) ([]byte, error) { return Bytes(v, binary.BigEndian) }

func LittleBytes(v interface{}) ([]byte, error) { return Bytes(v, binary.LittleEndian) }



