package test

import (
	"crypto/tls"
	"fmt"
	"github.com/jordan-wright/email"
	"math/rand"
	"net/smtp"
	"strings"
	"testing"
	"time"
)

func TestSendEmail(t *testing.T) {
	e := email.NewEmail()
	e.From = "bill <liuzexin98@163.com>"
	e.To = []string{"1264968325@qq.com"}

	e.Subject = "验证码发送测试"
	e.HTML = []byte("您的验证码是<b>123456</b>")
	//err := e.Send("smtp.163.com:456", smtp.PlainAuth("", "liuzexin98@163.com", "TCBIRLQHRZCNTKDF", "smtp.163.com"))
	err := e.SendWithTLS("smtp.163.com:465", smtp.PlainAuth("", "liuzexin98@163.com", "TCBIRLQHRZCNTKDF", "smtp.163.com"), &tls.Config{
		ServerName:         "smtp.163.com",
		InsecureSkipVerify: true,
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestGenerateVerifyCode(t *testing.T) {
	nums := [10]byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	rand.Seed(time.Now().Unix())
	builder := strings.Builder{}
	for i := 0; i < 6; i++ {
		fmt.Fprintf(&builder, "%d", nums[rand.Intn(len(nums))])
	}
	fmt.Println(builder.String())
}
