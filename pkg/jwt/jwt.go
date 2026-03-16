package jwt

import (
	"errors"
	"time"

	"github.com/Prizze/TaskScheduler/internal/config"
	"github.com/dgrijalva/jwt-go"
)

type Claims struct {
	UserID int64
	jwt.StandardClaims
}

func GenerateJWT(userID int64, cfg *config.Config) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)

	claims := &Claims{
		UserID: userID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodES256, claims)

	return token.SignedString(cfg.JwtSecret)
}

func ParseJWT(tokenString string, cfg *config.Config) (*Claims, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return cfg.JwtSecret, nil
	})
	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}
