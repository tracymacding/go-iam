package gerror

import (
	"fmt"
)

type IAMError struct {
	Status int
	Cause  error
}

func NewIAMError(status int, err error) error {
	return &IAMError{
		Status: status,
		Cause:  err,
	}
}

func (err *IAMError) Error() string {
	return fmt.Sprintf("%d-%s", err.Status, err.Cause)
}
