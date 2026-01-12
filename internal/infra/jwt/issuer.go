package jwt

import (
	"context"
	"crypto/sha256"
	"fmt"
	"log"
	"strings"
	"time"

	pgdto "Beanefits/internal/repository/postgres/dto"

	jwt "github.com/golang-jwt/jwt/v5"
)

type Issuer struct {
	secret []byte
	issuer string
	ttl    time.Duration
}

func NewIssuer(secret string, issuer string) *Issuer {
	secret = strings.TrimSpace(secret)
	issuer = strings.TrimSpace(issuer)
	sum := sha256.Sum256([]byte(secret))
	log.Print("jwt secret fp", "sha256_8", fmt.Sprintf("%x", sum[:8]), "len", len(secret))
	return &Issuer{
		secret: []byte(secret),
		issuer: issuer,
		ttl:    24 * time.Hour,
	}
}

type claims struct {
	Roles []string `json:"roles"`
	jwt.RegisteredClaims
}

func (i *Issuer) Issue(ctx context.Context, user pgdto.UserWithRoles) (string, error) {
	rs := make([]string, 0, len(user.Roles))
	for _, r := range user.Roles {
		rs = append(rs, string(r))
	}

	now := time.Now()
	c := claims{
		Roles: rs,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   fmt.Sprintf("%d", user.ID),
			Issuer:    i.issuer,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(i.ttl)),
		},
	}

	t := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	return t.SignedString(i.secret)
}
