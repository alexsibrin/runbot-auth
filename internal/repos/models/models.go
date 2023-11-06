package repos

type Auth struct {
	GUID     string
	Name     string
	Email    string
	Password string
	Level    uint8
	AddedAt  int64
}
