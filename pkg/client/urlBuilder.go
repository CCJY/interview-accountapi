package client

import (
	"fmt"
	"net/url"
	"strings"
)

type UrlBuilder interface {
	Build() (string, error)
}

type Url struct {
	BaseUrl       string
	OperationPath string
	QueryParams   url.Values
	PathParams    map[string]string
}

func (u *Url) Build() (string, error) {
	baseUrl, err := url.Parse(u.BaseUrl)
	if err != nil {
		return "", err
	}

	for k, v := range u.PathParams {
		u.OperationPath = strings.ReplaceAll(u.OperationPath, fmt.Sprintf("{%s}", k), v)
	}

	url, err := baseUrl.Parse(u.OperationPath)
	if err != nil {
		return "", err
	}

	url.RawQuery = u.QueryParams.Encode()

	return url.String(), err
}
