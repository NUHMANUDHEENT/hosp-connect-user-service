package utils

import (
	"encoding/base64"
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/go-redis/redis"
	"gopkg.in/gomail.v2"
)

var Rdb = redis.NewClient(&redis.Options{
	Addr:     os.Getenv("REDIS_PORT"),
	Password: "",
	DB:       0,
})

func GenerateToken() (string, error) {
	token := make([]byte, 32)
	if _, err := rand.Read(token); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(token), nil
}

func SignUpverify(email string) (string, error) {

	token, err := GenerateToken()
	if err != nil {
		return "Failed to generate token", err
	}

	err = Rdb.Set(token, email, 15*time.Minute).Err()
	if err != nil {
		return "Failed to store token in Redis", err
	}

	verificationURL := fmt.Sprintf("http://%s/api/v1/patient/signup/verify-email?token=%s", os.Getenv("IP_ADDRESS"), token)
	err = SendVerificationEmail(email, verificationURL)
	if err != nil {
		return "Failed to send verification email", err
	}

	return "Verification email sent successfully", nil
}

func SendVerificationEmail(email, link string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", "nuhmotp@gmail.com")
	m.SetHeader("To", email)
	m.SetHeader("Subject", "Email Verification")
	m.SetBody("text/html", fmt.Sprintf("Click the link to verify your email: <a href='%s'>Verify Email</a>", link))

	d := gomail.NewDialer("smtp.gmail.com", 587, os.Getenv("APPEMAIL"), os.Getenv("APPPASSWORD"))

	if err := d.DialAndSend(m); err != nil {
		return err
	}
	return nil
}
