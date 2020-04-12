package http

import "github.com/tornadoyi/viking/goplus/runtime"

type interactFunc					func(*Request, *Response) (interface{}, error)

func DoWithContext(f interactFunc) (ret interface{}, err error) {
	// catch error
	defer runtime.CatchCallback(func(info *runtime.PanicInfo) {
		err = info.Error()
	})

	// get from pool
	req := AcquireRequest()
	resp := AcquireResponse()
	defer ReleaseRequest(req)
	defer ReleaseResponse(resp)

	ret, err = f(req, resp)
	return ret, err
}