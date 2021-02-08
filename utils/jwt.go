package utils

import (
	"github.com/dgrijalva/jwt-go"
	"time"
)

var jwtSecret []byte

type Namespace struct {
	Domain   string `json:"domain"`
	Username string `json:"username"`
}

type Claims struct {
	Namespace Namespace `json:"https://www.teckstacks.cn/jwt/claims"`
	// Password string `json:"password"`
	jwt.StandardClaims
}

// GenerateToken generate tokens used for auth
func GenerateToken(username, domain string) (string, error) {
	nowTime := time.Now()
	expireTime := nowTime.Add(3 * time.Hour)

	jwtSecret = []byte("123456")

	claims := Claims{
		Namespace{
			domain,
			username,
		},
		//EncodeMD5(password),
		jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(),
			Issuer:    "gin-blog",
		},
	}

	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := tokenClaims.SignedString(jwtSecret)

	return token, err
}

// ParseToken parsing token
func ParseToken(token string) (*Claims, error) {
	tokenClaims, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if tokenClaims != nil {
		if claims, ok := tokenClaims.Claims.(*Claims); ok && tokenClaims.Valid {
			return claims, nil
		}
	}

	return nil, err
}
