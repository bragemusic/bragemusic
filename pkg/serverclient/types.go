package serverclient

type FileUpload struct {
	Name string `json:"name"`
	Type string `json:"type"`
	Data []byte `json:"data"`
}
