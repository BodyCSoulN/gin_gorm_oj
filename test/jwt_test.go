package test

import (
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"testing"
)

type UserClaims struct {
	Identity string `json:"identity"`
	Name     string `json:"name"`
	jwt.RegisteredClaims
}

var myKey = []byte("gin-gorm-oj-key")

// 生成token
func TestGenerateToken(t *testing.T) {
	UserClaim := &UserClaims{
		Identity:         "user_1",
		Name:             "bill",
		RegisteredClaims: jwt.RegisteredClaims{},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, UserClaim)
	tokenString, err := token.SignedString(myKey)
	if err != nil {
		t.Fatal(err)
	}
	// eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZGVudGl0eSI6InVzZXJfMSIsIm5hbWUiOiJiaWxsIn0.dwu4dyIVVxecaLukDqgC2TblOqaoskk8Q5Y-qH7sKAc
	fmt.Println(tokenString)
}

// 解析token
func TestParseToken(t *testing.T) {
	tokenString := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZGVudGl0eSI6InVzZXJfMSIsIm5hbWUiOiJiaWxsIn0.dwu4dyIVVxecaLukDqgC2TblOqaoskk8Q5Y-qH7sKAc"

	userClaim := new(UserClaims)
	claims, err := jwt.ParseWithClaims(tokenString, userClaim, func(token *jwt.Token) (interface{}, error) {
		return myKey, nil
	})
	if err != nil {
		t.Fatal(err)
	}

	if claims.Valid {
		fmt.Println(userClaim)
	} else {
		fmt.Println("parse err")
	}
}
