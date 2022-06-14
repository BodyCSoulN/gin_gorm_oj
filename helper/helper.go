package helper

import (
	"crypto/md5"
	"crypto/tls"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"github.com/jordan-wright/email"
	uuid "github.com/satori/go.uuid"
	"math/rand"
	"net/smtp"
	"strings"
	"time"
)

type UserClaims struct {
	Identity string `json:"identity"`
	Name     string `json:"name"`
	IsAdmin  int    `json:"is_admin"'`
	jwt.RegisteredClaims
}

// GetMd5 生成Md5
func GetMd5(s string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(s)))
}

var myKey = []byte("gin-gorm-oj-key")

// GenerateToken 生成token
func GenerateToken(identity, name string, isAdmin int) (string, error) {
	UserClaim := &UserClaims{
		Identity:         identity,
		Name:             name,
		IsAdmin:          isAdmin,
		RegisteredClaims: jwt.RegisteredClaims{},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, UserClaim)
	tokenString, err := token.SignedString(myKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

// ParseToken 解析token
func ParseToken(tokenString string) (*UserClaims, error) {

	userClaim := new(UserClaims)
	claims, err := jwt.ParseWithClaims(tokenString, userClaim, func(token *jwt.Token) (interface{}, error) {
		return myKey, nil
	})
	if err != nil {
		return nil, err
	}

	if !claims.Valid {
		return nil, fmt.Errorf("parse Token Error:%v", err)
	}
	return userClaim, nil
}

// SendVerifyCode 发送验证码
func SendVerifyCode(toUserEmail, verifyCode string) error {
	e := email.NewEmail()
	e.From = "bill <liuzexin98@163.com>"
	e.To = []string{toUserEmail}

	e.Subject = "验证码发送"
	e.HTML = []byte("您的验证码是：<b>" + verifyCode + "</b>，5分钟内有效。")
	err := e.SendWithTLS("smtp.163.com:465", smtp.PlainAuth("", "liuzexin98@163.com", "TCBIRLQHRZCNTKDF", "smtp.163.com"), &tls.Config{
		ServerName:         "smtp.163.com",
		InsecureSkipVerify: true,
	})
	return err
}

// GenerateUUid 生成UUi d
func GenerateUUid() string {
	return uuid.NewV4().String()
}

// GenerateVerifyCode 生成验证码
func GenerateVerifyCode() string {
	nums := [10]byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	rand.Seed(time.Now().Unix())
	builder := strings.Builder{}
	for i := 0; i < 6; i++ {
		fmt.Fprintf(&builder, "%d", nums[rand.Intn(len(nums))])
	}

	return builder.String()
}
