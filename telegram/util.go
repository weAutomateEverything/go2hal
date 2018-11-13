package telegram

import (
	"github.com/dgrijalva/jwt-go"
	"os"
	"time"
)

func makeToken(roomid uint32) (string, error) {
	claims := CustomClaims{
		roomid,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 120).Unix(),
			IssuedAt:  jwt.TimeFunc().Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(os.Getenv("JWT_KEY")))
}
