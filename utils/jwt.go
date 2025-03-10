package utils

import (
    "os"
    "time"

    "github.com/golang-jwt/jwt/v5"
)

func GenerateJWT(email, role string,id uint) (string, error) {
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
        "email": email,
        "role":  role,
        "id": id,
        "exp":   time.Now().Add(time.Hour * 72).Unix(),
    })
    return token.SignedString([]byte(os.Getenv("JWT_SECRET")))
}

func ValidateJWT(tokenString string) (email, role string, err error) {
    token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
        return []byte(os.Getenv("JWT_SECRET")), nil
    })

    if err != nil {
        return "", "", err
    }

    if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
        email = claims["email"].(string)
        role = claims["role"].(string)
        return email, role, nil
    }

    return "", "", jwt.ErrTokenInvalidClaims
}