package controllers

import "fmt"

type ErrDependenciesAreNil struct {
	controller string
	dependency string
}

func (err ErrDependenciesAreNil) Error() string {
	return fmt.Sprintf("controller %s got some dependencies are nil: %s", err.controller, err.dependency)
}

func NewErrDependenciesAreNil(c, d string) error {
	return ErrDependenciesAreNil{c, d}
}
