package entities

type baseErrors struct {
	Message string   `json:"message"`
	Details []string `json:"details"`
}

func (e *baseErrors) Error() string {
	return e.Message
}

type InvalidEmailAndOrPasswordError struct {
	*baseErrors
}

func NewInvalidEmailAndOrPasswordError() error {
	return &InvalidEmailAndOrPasswordError{
		baseErrors: &baseErrors{
			Message: "invalid e-mail and/or password",
		},
	}
}
