package hasher

import "golang.org/x/crypto/bcrypt"

type StringHasher struct{}

func NewStringHasher() *StringHasher {
	return &StringHasher{}
}

func (sh *StringHasher) Hash(str string) (string, error) {
	h, err := bcrypt.GenerateFromPassword([]byte(str), bcrypt.DefaultCost)
	return string(h), err
}

func (sh *StringHasher) Compare(str, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(str))
}
