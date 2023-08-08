package utils

type unwrapErrors interface {
	Unwrap() []error
}

// UnwrapErrors unwraps error from errors.Join(), or any error that implements `Unwrap() []error`.
func UnwrapErrors(err error) []error {
	e, ok := err.(unwrapErrors)
	if !ok {
		return nil
	}
	return e.Unwrap()
}
