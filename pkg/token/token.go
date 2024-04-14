package token

import (
	"fmt"
	"slices"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type (
	BannerAPIClaims struct {
		Roles   []Role `json:"roles,omitempty"`
		Comment string `json:"comment,omitempty"`

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

func IsAdmin(roles []string) bool {
	return slices.Contains(roles, string(RoleAdmin))
}

func IsUser(roles []string) bool {
	return slices.Contains(roles, string(RoleUser))
}

func Create(r Role) (string, error) {
	if hmacSecret == "" {
		return "", errors.New("hmac secret signing missing")
	}

	tokenString, err := jwt.NewWithClaims(jwt.SigningMethodHS256,
		BannerAPIClaims{
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

func Validate(token string) (jwt.MapClaims, error) {
	if hmacSecret == "" {
		return nil, errors.New("hmac secret signing missing")
	}

	token = strings.TrimPrefix(token, "Bearer ")

	parse, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(hmacSecret), nil
	})
	if err != nil {
		return nil, errors.WithMessage(err, "parse parsing failed")
	}

	if claims, ok := parse.Claims.(jwt.MapClaims); ok {
		return claims, nil
	}

	return nil, errors.New("unexpected empty claims")
}

// GetRoles implements the custom Claims getter
func GetRoles(m jwt.MapClaims) ([]string, error) {
	return parseString(m, "roles")
}

// parseString tries to parse a key in the map claims type as a [string] type.
// If the key does not exist, an empty string is returned. If the key has the
// wrong type, an error is returned.
func parseString(m jwt.MapClaims, key string) ([]string, error) {
	var (
		ok    bool
		raw   interface{}
		roles []string
	)
	raw, ok = m[key]
	if !ok {
		return nil, errors.WithMessage(fmt.Errorf("%s key is missing", key), jwt.ErrInvalidType.Error())
	}

	rawRoles, ok := raw.([]interface{})
	if !ok {
		return nil, errors.WithMessage(fmt.Errorf("%s is invalid", key), jwt.ErrInvalidType.Error())
	}

	for _, v := range rawRoles {
		if val, ok := v.(string); ok {
			roles = append(roles, val)
		}
	}

	return roles, nil
}
