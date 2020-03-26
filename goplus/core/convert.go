package core

import (
	"errors"
	"fmt"
	"strconv"
)


var ErrNil = errors.New("convert failed, nil returned")

func Int(v interface{}) (int, error) {
	switch v := v.(type) {
	case int64:
		x := int(v)
		if int64(x) != v {
			return 0, strconv.ErrRange
		}
		return x, nil
	case []byte:
		n, err := strconv.ParseInt(string(v), 10, 0)
		return int(n), err
	case nil:
		return 0, ErrNil
	case string:
		return strconv.Atoi(v)
	}
	return 0, fmt.Errorf("unexpected type for Int, got type %T", v)
}

func Int64(v interface{}) (int64, error) {
	switch v := v.(type) {
	case int64:
		return v, nil
	case []byte:
		n, err := strconv.ParseInt(string(v), 10, 64)
		return n, err
	case nil:
		return 0, ErrNil
	case string:
		return strconv.ParseInt(v, 10, 64)
	}
	return 0, fmt.Errorf("unexpected type for Int64, got type %T", v)
}

var errNegativeInt = errors.New("unexpected value for Uint64")


func Uint64(v interface{}) (uint64, error) {
	switch v := v.(type) {
	case int64:
		if v < 0 {
			return 0, errNegativeInt
		}
		return uint64(v), nil
	case []byte:
		n, err := strconv.ParseUint(string(v), 10, 64)
		return n, err
	case nil:
		return 0, ErrNil
	case string:
		return strconv.ParseUint(v, 10,64)
	}
	return 0, fmt.Errorf("unexpected type for Uint64, got type %T", v)
}


func Float64(v interface{}) (float64, error) {
	switch v := v.(type) {
	case []byte:
		n, err := strconv.ParseFloat(string(v), 64)
		return n, err
	case nil:
		return 0, ErrNil
	case string:
		return strconv.ParseFloat(v, 64)
	}
	return 0, fmt.Errorf("unexpected type for Float64, got type %T", v)
}


func String(v interface{}) (string, error) {
	switch v := v.(type) {
	case []byte:
		return string(v), nil
	case string:
		return v, nil
	case nil:
		return "", ErrNil
	case int:
		return strconv.Itoa(v), nil
	}
	return "", fmt.Errorf("unexpected type for String, got type %T", v)
}


func Bytes(v interface{}) ([]byte, error) {
	switch v := v.(type) {
	case []byte:
		return v, nil
	case string:
		return []byte(v), nil
	case nil:
		return nil, ErrNil
	}
	return nil, fmt.Errorf("unexpected type for Bytes, got type %T", v)
}


func Bool(v interface{}) (bool, error) {
	switch v := v.(type) {
	case int64:
		return v != 0, nil
	case []byte:
		return strconv.ParseBool(string(v))
	case nil:
		return false, ErrNil
	case string:
		return strconv.ParseBool(v)
	}
	return false, fmt.Errorf("unexpected type for Bool, got type %T", v)
}


func Values(v interface{}) ([]interface{}, error) {
	switch v := v.(type) {
	case []interface{}:
		return v, nil
	case nil:
		return nil, ErrNil
	}
	return nil, fmt.Errorf("unexpected type for Values, got type %T", v)
}


