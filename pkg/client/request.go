package client

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"net/url"
)

var (
	DefaultContentType string = "application/json"
)

// A RequestContext is for a HTTP request. It returns response data typed.
// Its support to json format for a request and a response by default settings.
// If want other encoding and decoding, use CustomEncoding.
//
// The RequestContext[T] includes a http.Client standard library, so
// this RequestContext[T] has a Do function. It requires http.Client.
//
// To create this, use NewRequestContext that is in this package.
// It returns interface to use it and has Do and hooks which are WhenBeforeDo
// and WhenAfterDo, as well as options retry and context.
// When using HookWhenBeforeDo, it can modify a http.Request.
// When using HookWhenAfterDo, it can manipulate for a response data typed before
// RequestContext[T].Do.
//
// When using it, the T should be the response data that you expect data and
// when using Do function, it returns ResponseContext[T] with error,
// the ResponseContext[T] includes http.Response and ContextData of T.
// The ContextData of T is the actual data you want.
type RequestContext[T any] struct {
	// It is required for using the Do function.
	HttpClient *http.Client

	// When call the Do function, the HttpRequest will be generated by using
	// Body, Header, Method and UrlBuilder.
	HttpRequest *http.Request

	// If wants specific context on a request level.
	Context context.Context
	// If set the CustomerEncoding, the DefaultEncoding is ignored.
	//
	// The Encoding interface should require implement of Marshal and UnMarshal.
	// Marshal(data interface{}) ([]byte, error)
	// UnMarshal(reader io.ReadCloser, dest interface{}) error
	CustomEncoding Encoding

	// If you are using json format, this will support it.
	DefaultEncoding HttpEncoding

	// It is a http.Header
	Header http.Header

	// The Method should be http's method like GET, POST, PUT and etc..
	Method string

	// This interface builds url needed for the request.
	// To build url, it needs BaseURL, OperationPath, QueryParams and PathParams.
	// If custom UrlBuilder, it should implement Build function that returns url string
	// Without custom UrlBuilder, it requires the Url that is in this package.
	// UrlBuilder: &Url{
	// 	BaseUrl: "http://127.0.0.1:8080",
	// 	OperationsPath: "/todo/{id}",
	// 	PathParams: map[string]string{
	// 		"id": "iam",
	// 	},
	// 	QueryParams: url.Values{
	// 		"some": []string{"hello"},
	// 	},
	// }
	// Returns string encoded from "http://127.0.0.1:8080/todo/iam?some=hello"
	UrlBuilder UrlBuilder

	// When using Body, it will be used by DefaultEncoding and CustomerEncoding.
	Body interface{}

	// When using HookWhenBeforeDo, it can modify a http.Request.
	HookWhenBeforeDo func(*RequestContext[T]) error

	// When using HookWhenAfterDo, it can manipulate for a response data typed before
	// RequestContext[T].Do.
	HookWhenAfterDo func(*ResponseContext[T]) error

	// If not set Retry, it will be ignored.
	// When using Retry, it requires both RetryInterval and
	// RetryMax of Retry that is in this pacakage.
	//
	// If RetryInterval is 300ms and RetryMax is 3,
	// the total waiting time is 300ms * 3.
	// RetryInterval is internally calculated in milliseconds.
	//
	// Example
	// If RetryInterval is 300ms and RetryMax is 3,
	// Retry: Retry{
	// RetryInterval: 300
	// RetryMax: 3
	// }
	Retry *Retry

	// It is related to Retry for reusing a request.
	originalBody []byte
}

func (r *RequestContext[T]) buildBody() (io.Reader, error) {
	if r.Body == nil {
		return nil, nil
	}

	if r.CustomEncoding != nil {
		buf, err := r.CustomEncoding.Marshal(r.Body)
		if err != nil {
			return nil, err
		}

		r.originalBody = buf
		return bytes.NewReader(buf), nil
	}

	contentType := r.Header.Get("Content-Type")
	if contentType == "" {
		contentType = DefaultContentType
	}
	buf, err := r.DefaultEncoding.Marshal(contentType, r.Body)
	if err != nil {
		return nil, err
	}

	r.originalBody = buf
	return bytes.NewReader(buf), nil
}

func (r *RequestContext[T]) newRequest() (*RequestContext[T], error) {
	url, err := r.UrlBuilder.Build()
	if err != nil {
		return nil, err
	}

	// If has Body, it returns io.Reader
	reader, err := r.buildBody()
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(r.Method, url, reader)
	if err != nil {
		return nil, err
	}

	r.HttpRequest = req
	r.HttpRequest.Header = r.Header

	if r.Context != nil {
		r.HttpRequest = r.HttpRequest.WithContext(r.Context)
	}

	return r, err
}

func (r *RequestContext[T]) Do() (*ResponseContext[T], error) {
	req, err := r.newRequest()
	if err != nil {
		return nil, err
	}

	if r.HookWhenBeforeDo != nil {
		err = r.HookWhenBeforeDo(r)
		if err != nil {
			return nil, err
		}
	}

	rsp, err := r.Retry.Do(r.HttpClient, r.HttpRequest, r.originalBody)
	// rsp, err := r.HttpClient.Do(req.HttpRequest)
	if err != nil {
		return nil, err
	}

	var rspData T
	if 0 < rsp.ContentLength {
		if r.CustomEncoding != nil {
			err = req.CustomEncoding.UnMarshal(rsp.Body, &rspData)
			if err != nil {
				return nil, err
			}
		} else {
			err = req.DefaultEncoding.UnMarshal(rsp, rsp.Body, &rspData)
			if err != nil {
				return nil, err
			}
		}
	}
	defer rsp.Body.Close()

	rspContext := ResponseContext[T]{}

	rspContext.HttpResponse = rsp
	rspContext.ContextData = rspData

	if r.HookWhenAfterDo != nil {
		err = r.HookWhenAfterDo(&rspContext)
		if err != nil {
			return &rspContext, err
		}
	}

	return &rspContext, nil
}

func (r *RequestContext[T]) WhenAfterDo(hook func(*ResponseContext[T]) error) RequestInterface[T] {
	r.HookWhenAfterDo = hook

	return r
}

func (r *RequestContext[T]) WhenBeforeDo(hook func(*RequestContext[T]) error) RequestInterface[T] {
	r.HookWhenBeforeDo = hook

	return r
}

func (r *RequestContext[T]) WithContext(ctx context.Context) RequestInterface[T] {
	r.Context = ctx

	return r
}

func (r *RequestContext[T]) WithRetry(opts ...RetryPolicyOpt) RequestInterface[T] {
	for _, opt := range opts {
		opt(r.Retry)
	}

	return r
}

type RequestInterface[T any] interface {
	// If uses this WithContext, HttpRequest of RequestContext[T] will be applied
	// and replaced to NewHttpRequest when call Do function.
	WithContext(context.Context) RequestInterface[T]

	// When uses this WithRetry, the first request depends
	// on the client's timeout or context.
	// When failed at the first time, it will retry the request as much as
	// the RetryMax value of Retry when call Do function.
	// If uses this WithRetry without options, it will be set by DefaultSetting.
	// DefaultRetry is the default implementation of Retry and is used by RetryPolicyNoBackOff.
	WithRetry(opts ...RetryPolicyOpt) RequestInterface[T]

	// When using WhenBeforeDo, it can modify a http.Request.
	WhenBeforeDo(func(*RequestContext[T]) error) RequestInterface[T]

	// When call this Do funcation, returns ResponseContext[T] and error. In addition,
	// ContextData of ResponseContext[T] is actual data that you expect data.
	Do() (*ResponseContext[T], error)

	// When using WhenAfterDo, it can manipulate for a response data typed before
	// RequestInterface.Do returns ResponseContext[T]
	WhenAfterDo(func(*ResponseContext[T]) error) RequestInterface[T]
}

// It returns interface to use it and has Do and hooks which are WhenBeforeDo
// and WhenAfterDo, as well as options retry and context.
// When using HookWhenBeforeDo, it can modify a http.Request.
// When using HookWhenAfterDo, it can manipulate for a response data typed before
// Do function
//
// To send an HTTP request and return an HTTP response, call Do function.
func newRequest[T any](httpClient *Client, r *RequestContext[T]) RequestInterface[T] {
	if httpClient == nil || r == nil {
		return nil
	}

	r.HttpClient = httpClient.HttpClient
	r.CustomEncoding = httpClient.Encoding
	return r
}

func NewRequestContext[T any](client *Client, contextModel *RequestContextModel) RequestInterface[T] {
	return newRequest(
		client,
		&RequestContext[T]{
			Context: contextModel.Context,
			Method:  contextModel.Method,
			UrlBuilder: &Url{
				BaseUrl:       contextModel.BaseUrl,
				OperationPath: contextModel.OperationPath,
				QueryParams:   contextModel.QueryParams,
				PathParams:    contextModel.PathParams,
			},
			Header:         contextModel.Header,
			Body:           contextModel.Body,
			CustomEncoding: contextModel.Encoding,
			Retry: &Retry{
				Policy: &RetryPolicy{
					RetryMax: 0,
				},
			},
		},
	)
}

type RequestContextModelOpt func(*RequestContextModel)

type RequestContextModel struct {
	Context       context.Context
	Method        string
	BaseUrl       string
	OperationPath string
	QueryParams   url.Values
	PathParams    map[string]string
	Header        http.Header
	Body          interface{}
	Encoding      Encoding
}

func NewRequestContextModel(opts ...RequestContextModelOpt) *RequestContextModel {
	model := &RequestContextModel{}

	for _, opt := range opts {
		opt(model)
	}

	return model
}

func WithHttpMethod(method string) RequestContextModelOpt {
	return func(requestContextModel *RequestContextModel) {
		requestContextModel.Method = method
	}
}
func WithBody(body interface{}) RequestContextModelOpt {
	return func(requestContextModel *RequestContextModel) {
		requestContextModel.Body = body
	}
}

func WithUrl(baseUrl string, operationPath string) RequestContextModelOpt {
	return func(requestContextModel *RequestContextModel) {
		requestContextModel.BaseUrl = baseUrl
		requestContextModel.OperationPath = operationPath
	}
}

type PathOpt func(map[string]string)

func WithPathParam(key string, value string) PathOpt {
	return func(m map[string]string) {
		m[key] = value
	}
}

func WithPathParams(opts ...PathOpt) RequestContextModelOpt {
	params := make(map[string]string)

	for _, opt := range opts {
		opt(params)
	}

	return func(requestContextModel *RequestContextModel) {
		requestContextModel.PathParams = params
	}
}

type QueryOpt func(*url.Values)

func WithQueryParam(key string, value string) QueryOpt {

	return func(v *url.Values) {
		v.Add(key, value)
	}
}

func WithQueryParams(opts ...QueryOpt) RequestContextModelOpt {
	queries := url.Values{}
	for _, opt := range opts {
		opt(&queries)

	}
	return func(rcm *RequestContextModel) {
		rcm.QueryParams = queries
	}
}

func WithQueryValues(urlValues *url.Values) RequestContextModelOpt {
	return func(rcm *RequestContextModel) {
		rcm.QueryParams = *urlValues
	}
}

func WithRequestContextModel(requestContextModel *RequestContextModel) RequestContextModelOpt {
	return func(rcm *RequestContextModel) {
		*rcm = *requestContextModel
	}
}
