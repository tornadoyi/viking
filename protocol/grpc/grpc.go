package grpc

import (
	_grpc "google.golang.org/grpc"
)

type DialOption  = _grpc.DialOption

type ServerOption = _grpc.ServerOption



// export DialOption
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


// export ServerOption
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