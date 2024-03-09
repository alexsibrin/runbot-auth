package models

import "encoding/json"

// The creating of an account

// SignUp the input model for the signing up
type SignUp struct {
	Email    string
	Password string
	Name     string
}

// SignUpResponse the response for the successful signing up
type SignUpResponse struct {
	Account *Account
	Token   *Token
}

// SignIn is a model for signing in
type SignIn struct {
	Email    string
	Password string
}

// SignIn the response for the successful signing in
type SignInResponse struct {
	Account *Account
	Token   *Token
}

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
	UpdatedAt int64 `json:"UpdatedAt,omitempty"`
}

// AccountGet input model for a getting an account by UUID
type AccountGet string

// AccountGetModel output model for a response to get account request
type AccountGetModel struct {
	UUID      string
	Email     string
	Name      string
	CreatedAt int64
	UpdatedAt int64 `json:"UpdatedAt,omitempty"`
}

// UpdateAccount input model for an updating an account
type UpdateAccount struct {
	UUID  string
	Email string
	Phone string
}

// Token the general token model
type Token struct {
	Access  string
	Refresh string `json:"-"`
}

func (t Token) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.Access)
}

type ChangeAccountStatus struct {
	UUID   string
	Status uint8
}

type ChangeAccountStatusResponse struct {
	UUID      string
	Status    uint8
	UpdatedAt int64
}
