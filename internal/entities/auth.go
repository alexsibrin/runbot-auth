package entities

type Auth struct {
	Email     string
	Password  string
	Name      string
	CreatedAt string
}

type Token struct {
	Access  string
	Refresh string
}
