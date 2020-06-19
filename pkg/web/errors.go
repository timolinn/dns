package web

// Error represents errors that happens on the web layer
type Error struct {
	Err    error
	Status int
	Fields []FieldError
}

// FieldError is used to indicate an error with a specific request field.
type FieldError struct {
	Field string `json:"field"`
	Error string `json:"error"`
}

type ErrorResponse struct {
	Error  string `json:"error"`
	Fields []FieldError
}

func NewRequestError(err error, statusCode int) error {
	return &Error{err, statusCode, nil}
}

func (e Error) Error() string {
	return e.Err.Error()
}
