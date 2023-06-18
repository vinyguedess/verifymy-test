package entities

type baseErrors struct {
	Message string   `json:"message"`
	Details []string `json:"details"`
}

func (e *baseErrors) Error() string {
	return e.Message
}

func NewUnexpectedError(err error) error {
	return &baseErrors{
		Message: "unexpected error",
		Details: []string{err.Error()},
	}
}

func NewError(message string, details []string) error {
	return &baseErrors{
		Message: message,
		Details: details,
	}
}

type EmailAlreadyInUseError struct {
	*baseErrors
}

func NewEmailAlreadyInUseError(email string) error {
	return &EmailAlreadyInUseError{
		baseErrors: &baseErrors{
			Message: "e-mail is already in use",
			Details: []string{email},
		},
	}
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
