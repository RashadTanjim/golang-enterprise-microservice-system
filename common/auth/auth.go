package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Config defines JWT configuration used for signing and validation.
type Config struct {
	Secret   string
	Issuer   string
	Audience string
	TokenTTL time.Duration
}

// Claims represents JWT claims with roles.
type Claims struct {
	Roles []string `json:"roles,omitempty"`
	jwt.RegisteredClaims
}

// GenerateToken creates a signed JWT for the provided subject and roles.
func GenerateToken(cfg Config, subject string, roles []string) (string, error) {
	if cfg.Secret == "" {
		return "", errors.New("auth secret is empty")
	}

	issuedAt := time.Now().UTC()
	expiresAt := issuedAt.Add(cfg.TokenTTL)

	audience := jwt.ClaimStrings{}
	if cfg.Audience != "" {
		audience = jwt.ClaimStrings{cfg.Audience}
	}

	claims := Claims{
		Roles: roles,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    cfg.Issuer,
			Subject:   subject,
			Audience:  audience,
			IssuedAt:  jwt.NewNumericDate(issuedAt),
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			NotBefore: jwt.NewNumericDate(issuedAt),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(cfg.Secret))
}

// ParseToken validates and parses a JWT into claims.
func ParseToken(cfg Config, tokenString string) (*Claims, error) {
	if cfg.Secret == "" {
		return nil, errors.New("auth secret is empty")
	}

	options := []jwt.ParserOption{jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()})}
	if cfg.Issuer != "" {
		options = append(options, jwt.WithIssuer(cfg.Issuer))
	}
	if cfg.Audience != "" {
		options = append(options, jwt.WithAudience(cfg.Audience))
	}

	parser := jwt.NewParser(options...)
	parsedToken, err := parser.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(cfg.Secret), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := parsedToken.Claims.(*Claims)
	if !ok || !parsedToken.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}
