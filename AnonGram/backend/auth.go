package main

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"os"
	"time"
)

type Claims struct {
	Role string `json:"role"`
	jwt.RegisteredClaims
}

func hashPassword(p string) (string, error) {
	b, e := bcrypt.GenerateFromPassword([]byte(p), bcrypt.DefaultCost)
	return string(b), e
}
func checkPassword(h, p string) bool {
	return bcrypt.CompareHashAndPassword([]byte(h), []byte(p)) == nil
}
func makeToken(s, r string, d time.Duration) (string, error) {
	c := Claims{Role: r, RegisteredClaims: jwt.RegisteredClaims{Subject: s, ExpiresAt: jwt.NewNumericDate(time.Now().Add(d)), IssuedAt: jwt.NewNumericDate(time.Now())}}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, c).SignedString([]byte(os.Getenv("JWT_SECRET")))
}
func ParseAccessToken(v string) (*Claims, error) {
	if v == "" {
		return nil, errors.New("empty")
	}
	t, e := jwt.ParseWithClaims(v, &Claims{}, func(t *jwt.Token) (any, error) {
		if t.Method.Alg() != "HS256" {
			return nil, errors.New("alg")
		}
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
	if e != nil {
		return nil, e
	}
	c, ok := t.Claims.(*Claims)
	if !ok || !t.Valid {
		return nil, errors.New("invalid")
	}
	return c, nil
}
