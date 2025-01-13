package swarmgo

import (
	"net/http"
	"net/url"

	"github.com/go-resty/resty/v2"
)

const (
	// MethodGet HTTP method
	MethodGet = "GET"

	// MethodPost HTTP method
	MethodPost = "POST"

	// MethodPut HTTP method
	MethodPut = "PUT"

	// MethodDelete HTTP method
	MethodDelete = "DELETE"

	// MethodPatch HTTP method
	MethodPatch = "PATCH"

	// MethodHead HTTP method
	MethodHead = "HEAD"

	// MethodOptions HTTP method
	MethodOptions = "OPTIONS"
)

type (
	HttpClientFunction struct {
		BaseFunction
		url     string
		method  string
		client  *resty.Client
		request *resty.Request
	}
)

func NewHttpClientFunction() *HttpClientFunction {
	fn := &HttpClientFunction{}
	baseFn, err := NewCustomFunction(fn)
	if err != nil {
		return nil
	}
	fn.BaseFunction = *baseFn
	fn.BaseFunction.SetFunction(fn.work)
	return fn
}

func (fn *HttpClientFunction) work(args map[string]interface{}, contextVariables map[string]interface{}) Result {
	result := Result{
		Success: false,
		Data:    "",
	}
	if fn.method == "" {
		fn.method = MethodGet
	}

	var resp *resty.Response
	var err error
	switch fn.method {
	case MethodGet:
		resp, err = fn.request.Get(fn.url)
	case MethodPost:
		resp, err = fn.request.Post(fn.url)
	case MethodPut:
		resp, err = fn.request.Put(fn.url)
	case MethodDelete:
		resp, err = fn.request.Delete(fn.url)
	case MethodPatch:
		resp, err = fn.request.Patch(fn.url)
	case MethodHead:
		resp, err = fn.request.Head(fn.url)
	}

	if err != nil {
		result.Data = err.Error()
		return result
	}

	body_bytes := resp.Body()
	result.Data = string(body_bytes)

	return result

}

var _ AgentFunction = (*HttpClientFunction)(nil)

func (fn *HttpClientFunction) GetName() string {
	return "http_client_request_data"
}

func (fn *HttpClientFunction) GetDescription() string {
	return "get data from a url by http request, support GET, POST, PUT, DELETE, PATCH, OPTIONS, HEAD, TRACE, CONNECT"
}

func (fn *HttpClientFunction) SetProxy(proxy string) *HttpClientFunction {
	fn.client.SetProxy(proxy)
	return fn
}

func (fn *HttpClientFunction) SetUrl(url string) *HttpClientFunction {
	fn.url = url
	return fn
}

func (model *HttpClientFunction) SetHeader(header, value string) *HttpClientFunction {
	model.request.SetHeader(header, value)
	return model
}

func (model *HttpClientFunction) SetFormDataFromValues(data url.Values) *HttpClientFunction {
	model.request.SetFormDataFromValues(data)
	return model
}

func (model *HttpClientFunction) SetBody(body interface{}) *HttpClientFunction {
	model.request.SetBody(body)
	return model
}

func (model *HttpClientFunction) SetFile(param, filePath string) *HttpClientFunction {
	model.request.SetFile(param, filePath)
	return model
}

func (model *HttpClientFunction) SetPathParam(param, value string) *HttpClientFunction {
	model.request.SetPathParam(param, value)
	return model
}

func (model *HttpClientFunction) SetCookie(hc *http.Cookie) *HttpClientFunction {
	model.request.SetCookie(hc)
	return model
}

func (model *HttpClientFunction) SetQueryParam(param, value string) *HttpClientFunction {
	model.request.SetQueryParam(param, value)
	return model
}

func (model *HttpClientFunction) SetCookies(rs []*http.Cookie) *HttpClientFunction {
	model.request.SetCookies(rs)
	return model
}

func (model *HttpClientFunction) SetMethod(method string) *HttpClientFunction {
	model.method = method
	return model
}
