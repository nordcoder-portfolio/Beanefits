package jwtverifier

import (
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"

	"Beanefits/internal/http/middleware"

	"github.com/golang-jwt/jwt/v5"
)

type Verifier struct {
	secret []byte
	issuer string
}

func New(secret, issuer string) *Verifier {
	secret = strings.TrimSpace(secret)
	issuer = strings.TrimSpace(issuer)
	sum := sha256.Sum256([]byte(secret))
	log.Print("jwt secret fp", "sha256_8", fmt.Sprintf("%x", sum[:8]), "len", len(secret))
	return &Verifier{secret: []byte(secret), issuer: issuer}
}

type claims struct {
	jwt.RegisteredClaims
	Roles []string `json:"roles"`
}

func (v *Verifier) Verify(ctx context.Context, token string) (middleware.Claims, error) {
	parsed, err := jwt.ParseWithClaims(token, &claims{}, func(t *jwt.Token) (any, error) {
		if t.Method.Alg() != jwt.SigningMethodHS256.Alg() {
			return nil, errors.New("unexpected signing method")
		}
		return v.secret, nil
	}, jwt.WithIssuer(v.issuer), jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))
	if err != nil {
		return middleware.Claims{}, err
	}
	c, ok := parsed.Claims.(*claims)
	if !ok || !parsed.Valid {
		return middleware.Claims{}, errors.New("invalid token")
	}
	uid, err := strconv.ParseInt(c.Subject, 10, 64)
	if err != nil {
		return middleware.Claims{}, errors.New("invalid subject")
	}
	return middleware.Claims{UserID: uid, Roles: c.Roles}, nil
}
