package client

import (
	"context"
	"io"
	"net/http"
	"net/url"
	"time"
)

var (
	DefaultContentType string = "application/json"
)

type RequestContext[T any] struct {
	HttpClient      *http.Client
	HttpRequest     *http.Request
	Context         context.Context
	CustomEncoding  Encoding
	DefaultEncoding HttpEncoding
	Header          http.Header
	Method          string
	UrlBuilder      UrlBuilder
	Body            interface{}

	HookWhenBeforeDo func(*RequestContext[T]) error
	HookWhenAfterDo  func(*ResponseContext[T]) error

	Retry         bool
	RetryInterval time.Duration
	RetryMax      int
}

func (r *RequestContext[T]) buildBody() (io.Reader, error) {
	if r.Body == nil {
		return nil, nil
	}

	if r.CustomEncoding != nil {
		return r.CustomEncoding.Marshal(r.Body)
	}

	contentType := r.Header.Get("Content-Type")
	if contentType == "" {
		contentType = DefaultContentType
	}
	return r.DefaultEncoding.Marshal(contentType, r.Body)
}

func (r *RequestContext[T]) newRequest() (*RequestContext[T], error) {
	url, err := r.UrlBuilder.Build()
	if err != nil {
		return nil, err
	}

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

	rsp, err := r.HttpClient.Do(req.HttpRequest)
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

func (r *RequestContext[T]) WhenAfterDo(hook func(*ResponseContext[T]) error) *RequestContext[T] {

	r.HookWhenAfterDo = hook

	return r
}

func (r *RequestContext[T]) WhenBeforeDo(hook func(*RequestContext[T]) error) *RequestContext[T] {

	r.HookWhenBeforeDo = hook

	return r
}

func (r *RequestContext[T]) WithContext(ctx context.Context) *RequestContext[T] {
	r.Context = ctx

	return r
}

type RequestInterface[T any] interface {
	WithContext(context.Context) *RequestContext[T]
	WhenBeforeDo(func(*RequestContext[T]) error) *RequestContext[T]
	Do() (*ResponseContext[T], error)
	WhenAfterDo(func(*ResponseContext[T]) error) *RequestContext[T]
}

func NewRequest[T any](httpClient *Client, r *RequestContext[T]) RequestInterface[T] {
	if httpClient == nil || r == nil {
		return nil
	}

	r.HttpClient = httpClient.HttpClient
	r.CustomEncoding = httpClient.Encoding
	return r
}

type RequestContextModel struct {
	Context       context.Context
	Method        string
	BaseUrl       string
	OperationPath string
	QueryParams   url.Values
	PathParams    map[string]string
	Header        http.Header
	Body          interface{}
}

func NewRequestContext[T any](contextModel RequestContextModel) *RequestContext[T] {
	return &RequestContext[T]{
		Context: contextModel.Context,
		Method:  contextModel.Method,
		UrlBuilder: &Url{
			BaseUrl:       contextModel.BaseUrl,
			OperationPath: contextModel.OperationPath,
			QueryParams:   contextModel.QueryParams,
			PathParams:    contextModel.PathParams,
		},
		Header: contextModel.Header,
		Body:   contextModel.Body,
	}
}
