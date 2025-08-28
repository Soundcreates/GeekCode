package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwtKey []byte

//initialize jwt secret key

func Init(secret string){
	jwtKey = []byte(secret)

}

//this func generates token for the user
func GenerateToken(userId uint, secret string) (string,error){
	claims := &jwt.RegisteredClaims{
		Subject: string(rune(userId)),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)), // 1day
		IssuedAt : jwt.NewNumericDate(time.Now()),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims) 
	return token.SignedString([]byte(secret))

}

//this func checks if the user has a token

func ValidateToken(tokenStr , secret string) (uint , error) {
	claims := &jwt.RegisteredClaims{}

	token, err := jwt.ParseWithClaims(tokenStr, claims , func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})

	if err  != nil || !token.Valid{
		return 0, errors.New("Invalid token")
	}

	//process of returning the user id
	runes := []rune(claims.Subject)
	
	return uint(runes[0]) , nil 
}
