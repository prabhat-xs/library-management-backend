package utils

import (
	"math"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateJWT(id uint, email, role string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": email,
		"role":  role,
		"id":    id,
		"exp":   time.Now().Add(time.Hour * 72).Unix(),
	})
	return token.SignedString([]byte(os.Getenv("JWT_SECRET")))
}

func ValidateJWT(tokenString string) (id uint, email, role string, err error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if err != nil {
		return math.MaxUint64, "", "", err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		email = claims["email"].(string)
		role = claims["role"].(string)
		floatId := claims["id"].(float64)
		id = uint(floatId)
		return id, email, role, nil
	}

	return math.MaxUint64, "", "", jwt.ErrTokenInvalidClaims
}
