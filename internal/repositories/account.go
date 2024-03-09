package repositories

type Account struct {
	UUID      string
	Name      string
	Email     string
	Password  string
	Status    uint8
	CreatedAt int64
	UpdatedAt int64
}
