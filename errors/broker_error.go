package errors

import (
	"fmt"
)

type BrokerError struct {
	wrapped_err error
}

func NewBrokerError(err error) *BrokerError {
	return &BrokerError{
		wrapped_err: err,
	}
}

func (e *BrokerError) Error() string {
	return e.wrapped_err.Error()
}

func (e *BrokerError) ToJson() string {
	return fmt.Sprintf("{ \"description\": \"%s\"}", e.wrapped_err.Error())
}
