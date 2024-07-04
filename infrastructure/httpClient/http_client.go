package httpClient

import "github.com/valyala/fasthttp"

func GetRequest(token string) (*fasthttp.Request, *fasthttp.Response) {
	req := fasthttp.AcquireRequest()
	req.Header.SetMethod(fasthttp.MethodGet)
	if token != "" {
		req.Header.Add("Authorization", "Bearer "+token)
	}
	resp := fasthttp.AcquireResponse()
	return req, resp
}

func ReleaseRR(req *fasthttp.Request, resp *fasthttp.Response) {
	fasthttp.ReleaseRequest(req)
	fasthttp.ReleaseResponse(resp)
}
