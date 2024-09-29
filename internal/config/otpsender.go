package config

import (
	"encoding/base64"
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/go-redis/redis"
	"gopkg.in/gomail.v2"
)

// Redis client setup (ensure Redis is running)
var Rdb = redis.NewClient(&redis.Options{
	Addr:     "localhost:6379", // Redis address
	Password: "",               // no password set
	DB:       0,
})

// var ctx = context.Background()

// Function to generate a secure random token
func GenerateToken() (string, error) {
	token := make([]byte, 32)
	if _, err := rand.Read(token); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(token), nil
}
func SignUpverify(email string) (string, error) {

	// Save patient details in DB and generate token
	token, err := GenerateToken()
	if err != nil {
		return "Failed to generate token", err
	}

	// Store the token in Redis with a 15-minute expiry
	err = Rdb.Set(token, email, 15*time.Minute).Err()
	if err != nil {
		return "Failed to store token in Redis", err
	}

	// Send the email with verification link
	verificationURL := fmt.Sprintf("http://localhost:8080/api/v1/patient/signup/verify-email?token=%s", token)
	err = SendVerificationEmail(email, verificationURL)
	if err != nil {
		return "Failed to send verification email", err
	}

	return "Verification email sent successfully", nil
}

// Send email function using GoMail
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
