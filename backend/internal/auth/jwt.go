package auth

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type claims struct {
	UserID string `json:"userId"`
	jwt.RegisteredClaims
}

func SignToken(userID string) (string, error) {
	secret := os.Getenv("JWT_SECRET")
	c := claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, c).SignedString([]byte(secret))
}

func VerifyToken(tokenStr string) (string, error) {
	secret := os.Getenv("JWT_SECRET")
	token, err := jwt.ParseWithClaims(tokenStr, &claims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(secret), nil
	})
	if err != nil {
		return "", err
	}
	c, ok := token.Claims.(*claims)
	if !ok || !token.Valid {
		return "", errors.New("invalid token")
	}
	return c.UserID, nil
}

func CheckCredentials(username, password string) (string, bool) {
	pairs := [][2]string{
		{os.Getenv("USER1_NAME"), os.Getenv("USER1_PASSWORD")},
		{os.Getenv("USER2_NAME"), os.Getenv("USER2_PASSWORD")},
	}
	for _, pair := range pairs {
		if pair[0] != "" && pair[0] == username && pair[1] == password {
			return username, true
		}
	}
	return "", false
}
