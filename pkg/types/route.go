package types

type NoResponse struct{}

type Response[T any] struct {
	Payload T   `json:"payload,omitempty"`
	Status  int `json:"-"`
}
