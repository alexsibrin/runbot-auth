package services

type IAuth interface {
	SignUp()
}

type IServices interface {
	IAuth
}
