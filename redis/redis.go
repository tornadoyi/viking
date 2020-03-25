package redis

import (
	"fmt"
	_redis "github.com/gomodule/redigo/redis"
	"sync"
)

var (
	pools 			=	map[string]*Pool{}
	mutex			= 	sync.RWMutex{}
)

func CreatePool(name string, network, host string, poolOptions []PoolOption, dialOptions []DialOption) (*Pool, error) {
	mutex.Lock()
	defer mutex.Unlock()
	// check
	if _, ok := pools[name]; ok { return nil, fmt.Errorf("Repeated redis pool %v", name) }

	// create pool
	pool := &Pool{}
	pool.init()
	pool.Dial = func() (Conn, error) {
		con, err := _redis.Dial(network, host, dialOptions...)
		if err != nil {
			return nil, err
		}
		return con, nil
	}
	for _, opt := range poolOptions { opt.f(pool) }

	// save
	pools[name] = pool

	return pool, nil
}


func GetPool(name string) (*Pool, bool) {
	mutex.RLock()
	defer mutex.RUnlock()
	pool, ok := pools[name]
	if !ok { return nil, false }
	return pool, true
}

func RemovePool(name string) {
	mutex.Lock()
	defer mutex.Unlock()
	delete(pools, name)
}



// export
type DialOption = _redis.DialOption
type Conn = _redis.Conn
type Args = _redis.Args
type Argument = _redis.Argument
type ConnWithTimeout = _redis.ConnWithTimeout
type Message = _redis.Message
type Pong = _redis.Pong
type PoolStats = _redis.PoolStats
type PubSubConn = _redis.PubSubConn
type Scanner = _redis.Scanner
type Script = _redis.Script
type SlowLog = _redis.SlowLog
type Subscription = _redis.Subscription


var Bool = _redis.Bool
var ByteSlices = _redis.ByteSlices
var Bytes = _redis.Bytes
var DoWithTimeout = _redis.DoWithTimeout
var Float64 = _redis.Float64
var Float64s = _redis.Float64s
var Int = _redis.Int
var Int64 = _redis.Int64
var Int64Map = _redis.Int64Map
var Int64s = _redis.Int64s
var IntMap = _redis.IntMap
var Ints = _redis.Ints
var MultiBulk = _redis.MultiBulk
var Positions = _redis.Positions
var ReceiveWithTimeout = _redis.ReceiveWithTimeout
var Scan = _redis.Scan
var ScanSlice = _redis.ScanSlice
var ScanStruct = _redis.ScanStruct
var String = _redis.String
var StringMap = _redis.StringMap
var Strings = _redis.Strings
var Uint64 = _redis.Uint64
var Values = _redis.Values

var DialClientName = _redis.DialClientName
var DialConnectTimeout = _redis.DialConnectTimeout
var DialDatabase = _redis.DialDatabase
var DialKeepAlive = _redis.DialKeepAlive
var DialNetDial = _redis.DialNetDial
var DialPassword = _redis.DialPassword
var DialReadTimeout = _redis.DialReadTimeout
var DialTLSConfig = _redis.DialTLSConfig
var DialTLSSkipVerify = _redis.DialTLSSkipVerify
var DialUseTLS = _redis.DialUseTLS
var DialWriteTimeout = _redis.DialWriteTimeout

