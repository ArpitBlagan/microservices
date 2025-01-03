package utils

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type Claims struct {
    UserID uint `json:"user_id"`
    jwt.RegisteredClaims
}
var jwtSecret = []byte("your-secret-key")
func CreateJwt(userId uint)(string,error){
	claims := &Claims{
        UserID: userId,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)), // 24-hour expiry
        },
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString(jwtSecret)
}

func ValidateJwt(tokenString string)(*Claims,error){
	
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
        return jwtSecret, nil
    })

    if err != nil {
        return nil, err
    }

    claims, ok := token.Claims.(*Claims)
    if !ok || !token.Valid {
        return nil, err
    }

    return claims, nil
}

