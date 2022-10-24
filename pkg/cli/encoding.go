package cli

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"strings"
)

type Encoding interface {
	Marshal(data interface{}) (io.Reader, error)
	UnMarshal(reader io.ReadCloser, dest interface{}) error
}

type HttpEncoding struct {
}

func (e *HttpEncoding) Marshal(contentType string, data interface{}) (io.Reader, error) {
	switch {
	case strings.Contains(contentType, "json"):
		buf, err := json.Marshal(data)
		if err != nil {
			return nil, err
		}
		return bytes.NewReader(buf), nil
	default:
		return nil, fmt.Errorf("invalid marshal")
	}
}

func (e *HttpEncoding) UnMarshal(contentType string, reader io.ReadCloser, dest interface{}) error {
	switch {
	case strings.Contains(contentType, "json"):
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

type JSONEncoding struct {
}

func (e *JSONEncoding) Marshal(data interface{}) (io.Reader, error) {
	buf, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	return bytes.NewReader(buf), nil
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
