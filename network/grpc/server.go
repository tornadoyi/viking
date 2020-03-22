package grpc

import (
	"errors"
	"fmt"
	"github.com/tornadoyi/viking/goplus/core"
	"github.com/tornadoyi/viking/log"
	"net"
	"github.com/tornadoyi/viking/network/grpc/keepalive"
	"reflect"
	"runtime"
	_grpc "google.golang.org/grpc"
)

var (
	servers = core.AtomicDict{}
)


func CreateServer(name string, network string, address string, opt ...ServerOption) (*Server, error) {

	if servers.Exists(name) { return nil, errors.New(fmt.Sprintf("Repteated server %v", name))}

	// create listener
	listener, err := net.Listen(network, address)
	if err != nil { return nil, err}

	// create grpc server
	server := &Server{_grpc.NewServer(opt...),name, listener,network, address}
	server.init()

	// save
	servers.Set(name, server)
	return server, nil
}


func GetServer(name string) (*Server, bool) {
	server, ok := servers.Get(name)
	if !ok { return nil, false }
	return server.(*Server), true
}

func RemoveServer(name string) {
	servers.Delete(name)
}


type Server struct {
	*_grpc.Server
	name				string
	listener			net.Listener
	network				string
	address				string
}

func (h* Server) init() {
	runtime.SetFinalizer(h, func (server *Server){
		if err := server.listener.Close(); err != nil {
			log.Error(err)
		}
	})
}

func (h *Server) Name() string { return h.name }

func (h *Server) Network() string { return h.network }

func (h *Server) Address() string { return h.address }

func (h *Server) Serve() error{ return h.Server.Serve(h.listener) }




func (h *Server) RegisterService(service interface{}, register interface{}) error {
	vf := reflect.ValueOf(register)
	if vf.Kind() != reflect.Func {
		return errors.New(fmt.Sprintf("Register failed for server %v with invalid function type %v", h.name, vf.Kind()))
	}

	paramCount := reflect.TypeOf(register).NumIn()
	if reflect.TypeOf(register).NumIn() != 2 {
		return errors.New(fmt.Sprintf("Invalid register parameters for server %v, parameter count is %v, expect 2", h.name, paramCount))
	}

	args := []reflect.Value{reflect.ValueOf(h.Server), reflect.ValueOf(service)}
	vf.Call(args)

	return nil
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