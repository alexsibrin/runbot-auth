package jwtapp

import (
	"errors"
	"fmt"
	"github.com/alexsibrin/runbot-auth/internal/entities"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"time"
)

var (
	ErrTokenIsNonValid = errors.New("token is not valid")
)

type myClaims struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	UUID  string `json:"uuid"`
	jwt.RegisteredClaims
}

type Config struct {
	Salt      string
	Issuer    string
	Subject   string
	Audience  []string
	ExpiresIn time.Duration
}

type JwtWrapper struct {
	config *Config
}

func New(c *Config) *JwtWrapper {
	return &JwtWrapper{
		config: c,
	}
}

func (j *JwtWrapper) AccessToken(a *entities.Account) (string, error) {
	return j.createToken(a, j.config.ExpiresIn, jwt.SigningMethodHS256)
}

func (j *JwtWrapper) RefreshToken(a *entities.Account) (string, error) {
	return j.createToken(a, j.config.ExpiresIn*10, jwt.SigningMethodHS512)
}

func (j *JwtWrapper) createToken(a *entities.Account, expiresin time.Duration, method jwt.SigningMethod) (string, error) {
	claims := j.convertEntity2Claims(a, expiresin)
	token := jwt.NewWithClaims(method, claims)

	signedString, err := token.SignedString([]byte(j.config.Salt))
	if err != nil {
		return "", err
	}

	return signedString, nil
}

func (j *JwtWrapper) Decrypt(t string) (*entities.Account, error) {
	token, err := jwt.Parse(t, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(j.config.Salt), nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(myClaims); ok && token.Valid {
		return j.convertClaims2Entity(&claims), nil
	}
	return nil, ErrTokenIsNonValid
}

func (j *JwtWrapper) convertEntity2Claims(a *entities.Account, expiresat time.Duration) *myClaims {
	now := time.Now()
	claims := jwt.RegisteredClaims{
		Issuer:    j.config.Issuer,
		Subject:   j.config.Subject,
		Audience:  j.config.Audience,
		ExpiresAt: jwt.NewNumericDate(now.Add(expiresat)),
		NotBefore: jwt.NewNumericDate(now),
		IssuedAt:  jwt.NewNumericDate(now),
		ID:        uuid.NewString(),
	}

	return &myClaims{
		Name:             a.Name,
		Email:            a.Email,
		UUID:             a.UUID,
		RegisteredClaims: claims,
	}
}

func (j *JwtWrapper) convertClaims2Entity(claims *myClaims) *entities.Account {
	return &entities.Account{
		UUID:  claims.UUID,
		Email: claims.Email,
		Name:  claims.Name,
	}
}
