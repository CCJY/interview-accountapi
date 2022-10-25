package commons

type ResponseData[T any] struct {
	Data         T      `json:"data,omitempty"`
	ErrorMessage string `json:"error_message,omitempty"`
}

type ResponseDataArray[T any] struct {
	Data         *[]T   `json:"data,omitempty"`
	ErrorMessage string `json:"error_message,omitempty"`
}
