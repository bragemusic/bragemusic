package serverclient

type ImageUpload struct {
	Name string `json:"name"`
	Type string `json:"type"`
	Data []byte `json:"data"`
}
