package gwork

import (
	"fmt"
)

type ServiceError struct {
	Msg string
}

func (e *ServiceError) Error() string {
	return fmt.Sprintf("%s", e.Msg)
}

func Error(msg string) error {
	return &ServiceError{msg}
}

func CheckErr(err error) {
	if err != nil {
		panic(err)
	}
}
