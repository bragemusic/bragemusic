package server

type Response struct {
	Payload any `json:"payload,omitempty"`
	Status  int `json:"-"`
}
