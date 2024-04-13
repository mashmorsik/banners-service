package token

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"time"
)

type (
	BannerApiClaims struct {
		Roles   []Role
		Comment string

		jwt.RegisteredClaims
	}

	Role string
)

var (
	hmacSecret = ""
)

const (
	RoleAdmin Role = "admin"
	RoleUser  Role = "user"
)

func NewTokenManager(secret string) {
	hmacSecret = secret
}

func Create(r Role) (string, error) {
	if hmacSecret == "" {
		return "", errors.New("hmac secret signing missing")
	}

	tokenString, err := jwt.NewWithClaims(jwt.SigningMethodHS256,
		BannerApiClaims{
			RegisteredClaims: jwt.RegisteredClaims{
				Issuer:    "ca.banners-service",
				Subject:   string(r),
				Audience:  []string{"banners-service-api"},
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
				IssuedAt:  jwt.NewNumericDate(time.Now()),
				NotBefore: jwt.NewNumericDate(time.Now()),
				ID:        uuid.Must(uuid.NewV7()).String(),
			},
			Roles:   []Role{r},
			Comment: "avito-top",
		}).SignedString([]byte(hmacSecret))
	if err != nil {
		return "", errors.WithMessage(err, "token signing failed")
	}

	return tokenString, nil
}
