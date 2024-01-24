package models

// The creating of an account

// AccountCreate the input model for the creating an accunt
type AccountCreate struct {
	Email    string
	Password string
	Name     string
}

// AccountCreateResponse the output models with a result for a created account
type AccountCreateResponse struct {
	Account *Account
	Token   *Token
}

// Account general account model
type Account struct {
	UUID      string
	Email     string
	Password  string
	Name      string
	CreatedAt int64
	UpdatedAt int64
}

// AccountGet input model for a getting an account by UUID
type AccountGet string

// AccountGetModel output model for a response to get account request
type AccountGetModel struct {
	UUID      string
	Email     string
	Name      string
	CreatedAt int64
	UpdatedAt int64
}

// UpdateAccount input model for an updating an account
type UpdateAccount struct {
	UUID  string
	Email string
	Phone string
}

// SignIn is a model for validating the user
type SignIn struct {
	Email    string
	Password string
}

// Token the general token model
type Token struct {
	Access  string
	Refresh string
}
