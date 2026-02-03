package bragerr

type BragErr struct {
	Code    string `json:"code"`
	Title   string `json:"title"`
	Message string `json:"message"`

	Status  int    `json:"-"`
	Service string `json:"-"`
	Err     error  `json:"-"`
}

func (e *BragErr) Error() string {
	if e.Err != nil {
		return e.Code + ": " + e.Err.Error()
	}
	return e.Code
}

func (e *BragErr) Unwrap() error {
	return e.Err
}

type BragErrFactory struct {
	service string
}

func NewFactory(service string) BragErrFactory {
	return BragErrFactory{
		service: service,
	}
}
