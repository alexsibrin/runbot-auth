package rest

type SignUpModel struct {
	Email    string
	Password string
	Name     string
}

type LogInModel struct {
	Email    string
	Password string
}
