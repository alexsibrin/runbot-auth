package controllers

import "fmt"

type ErrUnitIsNil struct {
	controller string
	dependency string
}

func (err ErrUnitIsNil) Error() string {
	return fmt.Sprintf("controller %s got some dependencies are nil: %s", err.controller, err.dependency)
}

func NewErrUnitIsNil(c, d string) error {
	return ErrUnitIsNil{c, d}
}
