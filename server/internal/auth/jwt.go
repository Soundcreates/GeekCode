package auth

import (
    "errors"
    "time"
    "strconv"
    "fmt"
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
        Subject: strconv.FormatUint(uint64(userId), 10), //converting uint to string
        ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)), // 1day
        IssuedAt : jwt.NewNumericDate(time.Now()),
    }

    fmt.Printf("Generating token for user ID: %d, Subject: %s\n", userId, claims.Subject) // Debug log

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
        fmt.Printf("Token parsing failed: %v\n", err) // Debug log
        return 0, errors.New("Invalid token")
    }

    //process of returning the user id
    fmt.Printf("Token Subject from claims: '%s'\n", claims.Subject) // Debug log
    
    userID , err  := strconv.ParseUint(claims.Subject,10,64)
    if err != nil {
        fmt.Printf("Failed to parse user ID from subject '%s': %v\n", claims.Subject, err) // Debug log
        return 0, errors.New("Invalid user ID in token")
    }
    
    fmt.Printf("Successfully parsed user ID: %d\n", userID) // Debug log
    return uint(userID) , nil 
}