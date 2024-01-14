package controllers

import "fmt"

type ErrEmptyValue struct {
	unit string
}

func (err ErrEmptyValue) Error() string {
	return fmt.Sprintf("%s is empty", err.unit)
}

func NewErrEmptyValue(u string) error {
	return ErrEmptyValue{u}
}
