package commons

type RequestData[T any] struct {
	Data *T `json:"data,omitempty"`
}
