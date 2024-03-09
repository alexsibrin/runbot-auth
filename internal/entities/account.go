package entities

const (
	Active uint8 = iota
	Suspended
	Blocked
)

type Account struct {
	UUID      string
	Email     string
	Password  string
	Name      string
	Status    uint8
	CreatedAt int64
	UpdatedAt int64
}

func (e *Account) IsActive() bool {
	switch e.Status {
	case Active:
		return true
	default:
		return false
	}
}
