package utils

import (
	"log"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var jwtKey = []byte("r3&$@^#e2f@%frg54w3%#^#yh5*(*$#grb5")

type Claims struct {
	UserID int64
	Phone  string
	jwt.StandardClaims
}

// GeneratedToken ...
func GeneratedToken(userID int64, phone string, expire int) (string, error) {
	expireTime := time.Now().Add(time.Duration(expire) * time.Hour)
	claims := &Claims{
		UserID: userID,
		Phone:  phone,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(), // 过期时间
			IssuedAt:  time.Now().Unix(), // 签名时间
			Issuer:    "Feiyizhou",       // 签名颁发者
			Subject:   "user token",      // 签名主题
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}
	return tokenString, err
}

// ParseToken ...
func ParseToken(tokenString string) (int64, string, error) {
	claims := new(Claims)
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (i interface{}, err error) {
		return jwtKey, nil
	})
	if err != nil {
		log.Printf("auth toke err, err: %v \n", err)
		return 0, "", err
	}
	if !token.Valid {
		log.Println("token is invalid")
		return 0, "", err
	}
	return claims.UserID, claims.Phone, err
}
