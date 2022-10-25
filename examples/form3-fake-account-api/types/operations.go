package types

type RequestData[T any] struct {
	Data T `json:"data,omitempty"`
}

type ResponseData[T any] struct {
	Data         T      `json:"data,omitempty"`
	ErrorMessage string `json:"error_message,omitempty"`
}

type ResponseDataArray[T any] struct {
	Data         *[]T   `json:"data,omitempty"`
	ErrorMessage string `json:"error_message,omitempty"`
}

func NewRequestData[T any](data T) *RequestData[T] {
	return &RequestData[T]{Data: data}
}
