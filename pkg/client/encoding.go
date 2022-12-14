package client

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type HttpEncoding struct {
}

func (e *HttpEncoding) Marshal(contentType string, data interface{}) ([]byte, error) {
	switch {
	case strings.Contains(contentType, "json"):
		buf, err := json.Marshal(data)
		if err != nil {
			return nil, err
		}
		return buf, nil
		// return bytes.NewReader(buf), nil
	default:
		return nil, fmt.Errorf("invalid marshal")
	}
}

func (e *HttpEncoding) UnMarshal(response *http.Response, reader io.ReadCloser, dest interface{}) error {
	switch {
	case strings.Contains(response.Header.Get("Content-Type"), "json"):
		bodyBytes, err := io.ReadAll(reader)

		if err != nil {
			return err
		}
		err = json.Unmarshal(bodyBytes, &dest)
		if err != nil {
			return err
		}

		return nil
	default:
		return fmt.Errorf("invalid unmarshal")
	}
}

// Custom Encoding

type Encoding interface {
	Marshal(data interface{}) ([]byte, error)
	UnMarshal(reader io.ReadCloser, dest interface{}) error
}

type JSONEncoding struct {
}

func (e *JSONEncoding) Marshal(data interface{}) ([]byte, error) {
	buf, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	return buf, err
}

func (e *JSONEncoding) UnMarshal(reader io.ReadCloser, dest interface{}) error {
	bodyBytes, err := io.ReadAll(reader)

	if err != nil {
		return err
	}
	err = json.Unmarshal(bodyBytes, &dest)
	if err != nil {
		return err
	}

	return nil
}
