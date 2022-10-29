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
// It returns  interface to use it and has Do and hooks which are WhenBeforeDo
// and WhenAfterDo, as well as options retry and context.
// When using HookWhenBeforeDo, it can modify a http.Request.
// When using HookWhenAfterDo, it can manipulate for a response data typed before
// RequestContext[T].Do.
type RequestContext[T any] struct {
	HttpClient  *http.Client
	HttpRequest *http.Request
	Context     context.Context
	// if set the CustomerEncoding, the DefaultEncoding is ignored
	CustomEncoding  Encoding
	DefaultEncoding HttpEncoding
	Header          http.Header
	Method          string
	UrlBuilder      UrlBuilder
	Body            interface{}

	HookWhenBeforeDo func(*RequestContext[T]) error
	HookWhenAfterDo  func(*ResponseContext[T]) error

	Retry Retry

	OriginalBody []byte
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

	return bytes.NewReader(buf), nil
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

	rsp, err := r.Retry.Do(r.HttpClient, r.HttpRequest, r.OriginalBody)
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

func (r *RequestContext[T]) WithRetry(retry Retry) *RequestContext[T] {
	r.Retry = retry

	return r
}

type RequestInterface[T any] interface {
	WithContext(context.Context) *RequestContext[T]
	WithRetry(retry Retry) *RequestContext[T]
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
