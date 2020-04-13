package http

import "github.com/tornadoyi/viking/goplus/runtime"

type interactFunc					func(*Request, *Response) (interface{}, error)

func DoWithContext(f interactFunc) (ret interface{}, reterr error) {
	// catch error
	defer runtime.CatchCallback(func(err error) { reterr = err })

	// get from pool
	req := AcquireRequest()
	resp := AcquireResponse()
	defer ReleaseRequest(req)
	defer ReleaseResponse(resp)

	ret, reterr = f(req, resp)
	return ret, reterr
}