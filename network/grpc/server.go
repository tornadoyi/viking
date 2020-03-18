package grpc

import (
	"errors"
	"fmt"
	"github.com/tornadoyi/viking/network/grpc/keepalive"
	"sync"
	_grpc "google.golang.org/grpc"
)

var (
	servers = make(map[string]*Server, 0)
	mutex sync.Mutex
)


type Server = _grpc.Server

func CreateServer(name string, opt ...ServerOption) (*Server, error) {
	defer mutex.Unlock()
	mutex.Lock()

	if _, ok := servers[name]; ok { return nil, errors.New(fmt.Sprintf("Repteated server %v", name))}

	server := _grpc.NewServer(opt...)
	servers[name] = server

	return server, nil
}


func GetServer(name string) (*Server, bool) {
	defer mutex.Unlock()
	mutex.Lock()
	server, ok := servers[name]
	return server, ok
}

func RemoveServer(name string) {
	defer mutex.Unlock()
	mutex.Lock()
	if _, ok := servers[name]; !ok { return }
	delete(servers, name)
}




type ServerConfig struct {
	InitialConnWindowSize		*int32									`yaml:"initial_connWindow_size"`
	InitialWindowSize			*int32									`yaml:"initial_window_size"`
	MaxConcurrentStreams		*uint32									`yaml:"max_concurrent_streams"`
	MaxHeaderListSize			*uint32									`yaml:"max_header_list_size"`
	MaxRecvMsgSize				*int									`yaml:"max_recv_msg_size"`
	MaxSendMsgSize				*int									`yaml:"max_send_msg_size"`
	ReadBufferSize				*int									`yaml:"read_buffer_size"`
	WriteBufferSize				*int									`yaml:"write_buffer_size"`
	Keepalive					*keepalive.ServerParametersConfig		`yaml:"keepalive"`
}


func (h* ServerConfig) ServerOptions() []ServerOption {
	if h == nil { return nil }
	options := make([]ServerOption, 0, 10)

	if h.InitialConnWindowSize != nil { options = append(options, InitialConnWindowSize(*h.InitialConnWindowSize)) }
	if h.InitialWindowSize != nil { options = append(options, InitialWindowSize(*h.InitialWindowSize)) }
	if h.MaxConcurrentStreams != nil { options = append(options, MaxConcurrentStreams(*h.MaxConcurrentStreams)) }
	if h.MaxHeaderListSize != nil { options = append(options, MaxHeaderListSize(*h.MaxHeaderListSize)) }
	if h.MaxRecvMsgSize != nil { options = append(options, MaxRecvMsgSize(*h.MaxRecvMsgSize)) }
	if h.MaxSendMsgSize != nil { options = append(options, MaxSendMsgSize(*h.MaxSendMsgSize)) }
	if h.ReadBufferSize != nil { options = append(options, ReadBufferSize(*h.ReadBufferSize)) }
	if h.WriteBufferSize != nil { options = append(options, WriteBufferSize(*h.WriteBufferSize)) }
	if h.Keepalive != nil { options = append(options, h.Keepalive.ServerOption()) }

	return options
}



// export ServerOption
type ServerOption = _grpc.ServerOption
var ChainStreamInterceptor = _grpc.ChainStreamInterceptor
var ChainUnaryInterceptor = _grpc.ChainUnaryInterceptor
var ConnectionTimeout = _grpc.ConnectionTimeout
var Creds = _grpc.Creds
var CustomCodec = _grpc.CustomCodec
var HeaderTableSize = _grpc.HeaderTableSize
var InTapHandle = _grpc.InTapHandle
var InitialConnWindowSize = _grpc.InitialConnWindowSize
var InitialWindowSize = _grpc.InitialWindowSize
var KeepaliveEnforcementPolicy = _grpc.KeepaliveEnforcementPolicy
var KeepaliveParams = _grpc.KeepaliveParams
var MaxConcurrentStreams = _grpc.MaxConcurrentStreams
var MaxHeaderListSize = _grpc.MaxHeaderListSize
var MaxMsgSize = _grpc.MaxMsgSize
var MaxRecvMsgSize = _grpc.MaxRecvMsgSize
var MaxSendMsgSize = _grpc.MaxSendMsgSize
var RPCCompressor = _grpc.RPCCompressor
var RPCDecompressor = _grpc.RPCDecompressor
var ReadBufferSize = _grpc.ReadBufferSize
var StatsHandler = _grpc.StatsHandler
var StreamInterceptor = _grpc.StreamInterceptor
var UnaryInterceptor = _grpc.UnaryInterceptor
var UnknownServiceHandler = _grpc.UnknownServiceHandler
var WriteBufferSize = _grpc.WriteBufferSize