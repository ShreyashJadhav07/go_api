package utils

import (

	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwtKey = []byte(os.Getenv("JWT_SECRET"))

type Claims struct{
	Email string `json:"email"`
	jwt.RegisteredClaims
}

func GenerateToken(email string) (string,error){
	expirationTime :=time.Now().Add(24 * time.Hour)
	
	claims :=&Claims{
		Email:email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt: jwt.NewNumericDate(time.Now()),
			Issuer: "auth-api",
		},

	}


	token:=jwt.NewWithClaims(jwt.SigningMethodHS256,claims)

	tokenString,err := token.SignedString(jwtKey)

	if err!=nil{
		log.Printf("Error signing JWT token: %V",err)
		return "" ,err
	}

	return tokenString,nil
}