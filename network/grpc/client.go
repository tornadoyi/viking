package grpc

import (
	"errors"
	"fmt"
	"github.com/tornadoyi/viking/network/grpc/backoff"
	"github.com/tornadoyi/viking/network/grpc/keepalive"
	_grpc "google.golang.org/grpc"
	"sync"
	"time"
)

var (
	clients = make(map[string]*Client)
	clientMutex sync.Mutex
)


type Client = _grpc.ClientConn



func CreateClient(name string, address string,  opt ...DialOption) (*Client, error) {
	defer clientMutex.Unlock()
	clientMutex.Lock()

	if _, ok := clients[name]; ok { return nil, errors.New(fmt.Sprintf("Repteated client %v", name))}
	client, err := _grpc.Dial(address, opt...)
	if err != nil { return nil, err}
	clients[name] = client

	return client, nil
}


func GetClinet(name string) (*Client, bool) {
	defer clientMutex.Unlock()
	clientMutex.Lock()
	client, ok := clients[name]
	return client, ok
}

func RemoveClient(name string) {
	defer clientMutex.Unlock()
	clientMutex.Lock()
	if _, ok := clients[name]; !ok { return }
	delete(clients, name)
}


type DialConfig struct {
	FailOnNonTempDialError		*bool										`yaml:"fail_on_non_temp_dial_error"`
	Authority					*string										`yaml:"authority"`
	Block						*bool										`yaml:"block"`
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
	MinConnectTimeout 			time.Duration							`yaml:"min_connect_timeout"`
}

func (h *ConnectParamsConfig) DialOption() DialOption {
	return WithConnectParams(_grpc.ConnectParams{
		h.Backoff.Config(),
		h.MinConnectTimeout,
	})
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
