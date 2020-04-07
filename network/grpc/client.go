package grpc

import (
	"errors"
	"fmt"
	"github.com/tornadoyi/viking/log"
	"github.com/tornadoyi/viking/network/grpc/backoff"
	"github.com/tornadoyi/viking/network/grpc/keepalive"
	_grpc "google.golang.org/grpc"
	"reflect"
	"runtime"
	"sync"
	"time"
)

var (
	clients		 = map[string]*Client{}
	cmutex		 = sync.RWMutex{}
)


func CreateClient(name string, target string, register interface{}, options... DialOption) (*Client, error) {
	cmutex.Lock()
	defer cmutex.Unlock()
	if _, ok := clients[name]; ok { return nil, fmt.Errorf("Repteated client %v", name)}

	// dial
	conn, err := _grpc.Dial(target, options...)
	if err != nil { return nil, err}

	// register service
	vf := reflect.ValueOf(register)
	if vf.Kind() != reflect.Func {
		return nil, errors.New(fmt.Sprintf("Register failed for client %v with invalid function type %v", name, vf.Kind()))
	}
	if parmInCount := reflect.TypeOf(register).NumIn(); parmInCount != 1 {
		return nil, errors.New(fmt.Sprintf("Invalid register parameters for client %v, input parameter count is %v, expect 1", name, parmInCount))
	}
	if parmOutCount := reflect.TypeOf(register).NumOut(); parmOutCount != 1 {
		return nil, errors.New(fmt.Sprintf("Invalid register parameters for client %v, output parameter count is %v, expect 1", name, parmOutCount))
	}
	args := []reflect.Value{reflect.ValueOf(conn)}
	service := vf.Call(args)[0].Interface()

	// create client and save
	client := &Client{name, conn, service}
	client.init()
	clients[name] = client
	return client, nil
}

func GetClient(name string) (*Client, bool) {
	cmutex.RLock()
	defer cmutex.RUnlock()
	c, ok := clients[name]
	if !ok { return nil, false }
	return c, true
}

func RemoveClient(name string){
	cmutex.Lock()
	defer cmutex.Unlock()
	delete(clients, name)
}



type ClientConn = _grpc.ClientConn

type Client struct {
	name						string
	connection					*ClientConn
	service						interface{}
}

func (h* Client) init() {
	runtime.SetFinalizer(h, func (client *Client){
		if err := h.connection.Close(); err != nil {
			log.Errorw("Grpc client connection close", "error", err)
		}
	})
}

func (h *Client) Name() string { return h.name }

func (h *Client) Service() interface{} { return h.service }





type DialConfig struct {
	FailOnNonTempDialError		*bool										`yaml:"fail_on_non_temp_dial_error"`
	Authority					*string										`yaml:"authority"`
	Block						*bool										`yaml:"block"`
	BalancerName				*string										`yaml:"balance_name"`
	ChannelzParentID			*int64										`yaml:"channelz_parent_id"`
	DefaultServiceConfig		*string										`yaml:"default_service_config"`
	DisableHealthCheck			*bool										`yaml:"disable_health_check"`
	DisableRetry				*bool										`yaml:"disable_retry"`
	DisableServiceConfig		*bool										`yaml:"disable_service_config"`
	InitialConnWindowSize		*int32										`yaml:"initial_connWindow_size"`
	Insecure					*bool										`yaml:"insecure"`
	MaxHeaderListSize			*uint32										`yaml:"maxHeader_list_size"`
	ReadBufferSize				*int										`yaml:"read_buffer_size"`
	UserAgent					*string										`yaml:"user_agent"`
	WriteBufferSize				*int										`yaml:"Write_buffer_size"`
	ConnectParams				*ConnectParamsConfig						`yaml:"connect_params"`
	KeepaliveParams				*keepalive.ClientParametersConfig			`yaml:"keepalive"`
}


func (h *DialConfig) DialOptions() []DialOption {
	if h == nil { return  nil}
	options := make([]DialOption, 0, 20)

	if h.FailOnNonTempDialError != nil { options = append(options, FailOnNonTempDialError(*h.FailOnNonTempDialError)) }
	if h.Authority != nil { options = append(options, WithAuthority(*h.Authority)) }
	if h.Block != nil && *h.Block { options = append(options, WithBlock()) }
	if h.BalancerName != nil { options = append(options, WithBalancerName(*h.BalancerName)) }
	if h.ChannelzParentID != nil { options = append(options, WithChannelzParentID(*h.ChannelzParentID)) }
	if h.DefaultServiceConfig != nil { options = append(options, WithDefaultServiceConfig(*h.DefaultServiceConfig)) }
	if h.DisableHealthCheck != nil && *h.DisableHealthCheck { options = append(options, WithDisableHealthCheck()) }
	if h.DisableRetry != nil && *h.DisableRetry { options = append(options, WithDisableRetry()) }
	if h.DisableServiceConfig != nil && *h.DisableServiceConfig { options = append(options, WithDisableServiceConfig()) }
	if h.InitialConnWindowSize != nil { options = append(options, WithInitialConnWindowSize(*h.InitialConnWindowSize)) }
	if h.Insecure != nil && *h.Insecure { options = append(options, WithInsecure()) }
	if h.MaxHeaderListSize != nil { options = append(options, WithMaxHeaderListSize(*h.MaxHeaderListSize)) }
	if h.ReadBufferSize != nil { options = append(options, WithReadBufferSize(*h.ReadBufferSize)) }
	if h.UserAgent != nil { options = append(options, WithUserAgent(*h.UserAgent)) }
	if h.ConnectParams != nil { options = append(options, h.ConnectParams.DialOption()) }
	if h.KeepaliveParams != nil { options = append(options, h.KeepaliveParams.DialOption()) }

	return options
}




type ConnectParamsConfig struct {
	Backoff 					backoff.Config							`yaml:"backoff"`
	MinConnectTimeout 			string									`yaml:"min_connect_timeout"`
}

func (h *ConnectParamsConfig) DialOption() DialOption {
	p := _grpc.ConnectParams{}
	p.Backoff = h.Backoff.Config()
	p.MinConnectTimeout, _ = time.ParseDuration(h.MinConnectTimeout)
	return WithConnectParams(p)
}




// export DialOption
type DialOption  = _grpc.DialOption
var FailOnNonTempDialError = _grpc.FailOnNonTempDialError
var WithAuthority = _grpc.WithAuthority
var WithBackoffConfig = _grpc.WithBackoffConfig
var WithBackoffMaxDelay = _grpc.WithBackoffMaxDelay
var WithBalancer = _grpc.WithBalancer
var WithBalancerName = _grpc.WithBalancerName
var WithBlock = _grpc.WithBlock
var WithChainStreamInterceptor = _grpc.WithChainStreamInterceptor
var WithChainUnaryInterceptor = _grpc.WithChainUnaryInterceptor
var WithChannelzParentID = _grpc.WithChannelzParentID
var WithCodec = _grpc.WithCodec
var WithCompressor = _grpc.WithCompressor
var WithConnectParams = _grpc.WithConnectParams
var WithContextDialer = _grpc.WithContextDialer
var WithCredentialsBundle = _grpc.WithCredentialsBundle
var WithDecompressor = _grpc.WithDecompressor
var WithDefaultCallOptions = _grpc.WithDefaultCallOptions
var WithDefaultServiceConfig = _grpc.WithDefaultServiceConfig
var WithDialer = _grpc.WithDialer
var WithDisableHealthCheck = _grpc.WithDisableHealthCheck
var WithDisableRetry = _grpc.WithDisableRetry
var WithDisableServiceConfig = _grpc.WithDisableServiceConfig
var WithInitialConnWindowSize = _grpc.WithInitialConnWindowSize
var WithInitialWindowSize = _grpc.WithInitialWindowSize
var WithInsecure = _grpc.WithInsecure
var WithKeepaliveParams = _grpc.WithKeepaliveParams
var WithMaxHeaderListSize = _grpc.WithMaxHeaderListSize
var WithMaxMsgSize = _grpc.WithMaxMsgSize
var WithPerRPCCredentials = _grpc.WithPerRPCCredentials
var WithReadBufferSize = _grpc.WithReadBufferSize
var WithResolvers = _grpc.WithResolvers
var WithServiceConfig = _grpc.WithServiceConfig
var WithStatsHandler = _grpc.WithStatsHandler
var WithStreamInterceptor = _grpc.WithStreamInterceptor
var WithTimeout = _grpc.WithTimeout
var WithTransportCredentials = _grpc.WithTransportCredentials
var WithUnaryInterceptor = _grpc.WithUnaryInterceptor
var WithUserAgent = _grpc.WithUserAgent
var WithWriteBufferSize = _grpc.WithWriteBufferSize
