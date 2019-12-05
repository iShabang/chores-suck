package tools

import (
	"github.com/dgrijalva/jwt-go"
	"time"
)

var (
	SECRET_KEY = []byte("cnBkyv93jqZ1DMWkDxHqCbfb@II*bq8!IUJnf#859VBz&n80$WQ9kIUEn5zOGz5M")
)

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

func GenToken(u string) (string, error, time.Time) {
	expireTime := time.Now().Add(5 * time.Minute)
	claims := &Claims{
		Username: u,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(SECRET_KEY)
	return tokenString, err, expireTime
}

func AuthToken(t string) (bool, error, Claims) {
	claims := Claims{}
	tkn, err := jwt.ParseWithClaims(t, claims, func(token *jwt.Token) (interface{}, error) {
		return SECRET_KEY, nil
	})
	return tkn.Valid, err, claims
}
