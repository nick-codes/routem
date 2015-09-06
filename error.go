package routem

type (
	httpError struct {
		code int
		err  error
	}
)

func (e *httpError) Code() int {
	return e.code
}

func (e *httpError) Error() string {
	return e.err.Error()
}

// NewHTTPError constructs a new HTTPError with the given code and err.
func NewHTTPError(code int, err error) HTTPError {
	return &httpError{
		code: code,
		err:  err,
	}
}
