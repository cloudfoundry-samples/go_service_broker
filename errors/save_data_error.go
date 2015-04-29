package errors

import (
	"fmt"
)

type SaveDataError struct {
	Data   string
	Reason error
}

func NewSaveDataError(data string, reason error) *SaveDataError {
	return &SaveDataError{
		Data:   data,
		Reason: reason,
	}
}

func (e *SaveDataError) Error() string {
	return fmt.Sprintf("Can not save %s due to %s", e.Data, e.Reason.Error())
}
