package client

import (
	"bytes"
	"encoding/json"
	"io"
)

type Encoding interface {
	Marshal(data interface{}) (io.Reader, error)
	UnMarshal(reader io.ReadCloser, dest interface{}) error
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
