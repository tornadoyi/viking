package http

import (
	"fmt"
	"net"
	"time"
)

var (
	servers		=  	map[string]*Server{}
)


func CreateServer(name string, args...ServerOption) (*Server, error) {
	// check
	if _, ok := servers[name]; ok { return nil, fmt.Errorf("Repeated server %v", name)}

	// apply arguments
	s := &Server{}
	for _, a := range args { a.apply(s) }

	// save
	servers[name] = s
	return s, nil
}



func HandlerOption(Handler RequestHandler) ServerOption {
	return newFuncServerOption(func(o *Server) {
		o.Handler = Handler
	})
}


func ErrorHandlerOption(ErrorHandler func(ctx *RequestCtx, err error)) ServerOption {
	return newFuncServerOption(func(o *Server) {
		o.ErrorHandler = ErrorHandler
	})
}


func HeaderReceivedOption(HeaderReceived func(header *RequestHeader) RequestConfig) ServerOption {
	return newFuncServerOption(func(o *Server) {
		o.HeaderReceived = HeaderReceived
	})
}

func NameOption(name string) ServerOption {
	return newFuncServerOption(func(o *Server) {
		o.Name = name
	})
}

func ConcurrencyOption(Concurrency int) ServerOption {
	return newFuncServerOption(func(o *Server) {
		o.Concurrency = Concurrency
	})
}

func DisableKeepaliveOption(DisableKeepalive bool) ServerOption {
	return newFuncServerOption(func(o *Server) {
		o.DisableKeepalive = DisableKeepalive
	})
}

func ReadBufferSizeOption(ReadBufferSize int) ServerOption {
	return newFuncServerOption(func(o *Server) {
		o.ReadBufferSize = ReadBufferSize
	})
}

func WriteBufferSizeOption(WriteBufferSize int) ServerOption {
	return newFuncServerOption(func(o *Server) {
		o.WriteBufferSize = WriteBufferSize
	})
}

func ReadTimeoutOption(ReadTimeout time.Duration) ServerOption {
	return newFuncServerOption(func(o *Server) {
		o.ReadTimeout = ReadTimeout
	})
}

func WriteTimeoutOption(WriteTimeout time.Duration) ServerOption {
	return newFuncServerOption(func(o *Server) {
		o.WriteTimeout = WriteTimeout
	})
}

func IdleTimeoutOption(IdleTimeout time.Duration) ServerOption {
	return newFuncServerOption(func(o *Server) {
		o.IdleTimeout = IdleTimeout
	})
}

func MaxConnsPerIPOption(MaxConnsPerIP int) ServerOption {
	return newFuncServerOption(func(o *Server) {
		o.MaxConnsPerIP = MaxConnsPerIP
	})
}

func MaxRequestsPerConnOption(MaxRequestsPerConn int) ServerOption {
	return newFuncServerOption(func(o *Server) {
		o.MaxRequestsPerConn = MaxRequestsPerConn
	})
}

func MaxKeepaliveDurationOption(MaxKeepaliveDuration time.Duration) ServerOption {
	return newFuncServerOption(func(o *Server) {
		o.MaxKeepaliveDuration = MaxKeepaliveDuration
	})
}

func TCPKeepaliveOption(TCPKeepalive bool) ServerOption {
	return newFuncServerOption(func(o *Server) {
		o.TCPKeepalive = TCPKeepalive
	})
}

func TCPKeepalivePeriodOption(TCPKeepalivePeriod time.Duration) ServerOption {
	return newFuncServerOption(func(o *Server) {
		o.TCPKeepalivePeriod = TCPKeepalivePeriod
	})
}

func MaxRequestBodySizeOption(MaxRequestBodySize int) ServerOption {
	return newFuncServerOption(func(o *Server) {
		o.MaxRequestBodySize = MaxRequestBodySize
	})
}

func ReduceMemoryUsageOption(ReduceMemoryUsage bool) ServerOption {
	return newFuncServerOption(func(o *Server) {
		o.ReduceMemoryUsage = ReduceMemoryUsage
	})
}

func GetOnlyOption(GetOnly bool) ServerOption {
	return newFuncServerOption(func(o *Server) {
		o.GetOnly = GetOnly
	})
}

func LogAllErrorsOption(LogAllErrors bool) ServerOption {
	return newFuncServerOption(func(o *Server) {
		o.LogAllErrors = LogAllErrors
	})
}

func DisableHeaderNamesNormalizingOption(DisableHeaderNamesNormalizing bool) ServerOption {
	return newFuncServerOption(func(o *Server) {
		o.DisableHeaderNamesNormalizing = DisableHeaderNamesNormalizing
	})
}

func SleepWhenConcurrencyLimitsExceededOption(SleepWhenConcurrencyLimitsExceeded time.Duration) ServerOption {
	return newFuncServerOption(func(o *Server) {
		o.SleepWhenConcurrencyLimitsExceeded = SleepWhenConcurrencyLimitsExceeded
	})
}

func NoDefaultServerHeaderOption(NoDefaultServerHeader bool) ServerOption {
	return newFuncServerOption(func(o *Server) {
		o.NoDefaultServerHeader = NoDefaultServerHeader
	})
}

func NoDefaultDateOption(NoDefaultDate bool) ServerOption {
	return newFuncServerOption(func(o *Server) {
		o.NoDefaultDate = NoDefaultDate
	})
}

func NoDefaultContentTypeOption(NoDefaultContentType bool) ServerOption {
	return newFuncServerOption(func(o *Server) {
		o.NoDefaultContentType = NoDefaultContentType
	})
}

func ConnStateOption(ConnState func(net.Conn, ConnState)) ServerOption {
	return newFuncServerOption(func(o *Server) {
		o.ConnState = ConnState
	})
}

func LoggerOption(Logger Logger) ServerOption {
	return newFuncServerOption(func(o *Server) {
		o.Logger = Logger
	})
}

func KeepHijackedConnsOption(KeepHijackedConns bool) ServerOption {
	return newFuncServerOption(func(o *Server) {
		o.KeepHijackedConns = KeepHijackedConns
	})
}


type ServerOption interface {
	apply(*Server)
}

type funcServerOption struct {
	f func(*Server)
}

func (h *funcServerOption) apply(do *Server) {
	h.f(do)
}

func newFuncServerOption(f func(*Server)) *funcServerOption {
	return &funcServerOption{
		f: f,
	}
}

