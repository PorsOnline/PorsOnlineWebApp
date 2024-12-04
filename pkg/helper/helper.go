package helper

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/joho/godotenv"
	"gopkg.in/gomail.v2"
)

func GetConfig(key string) string {
	return os.Getenv(key)
}
func LoadEnvFile() {
	err := godotenv.Load(".env")
	if err != nil {
		panic("Error loading .env file")
	}
}

func SendEmail(to string, message string) {
	LoadEnvFile()
	GOOGLE_SMTP_PASSWORD := GetConfig("GOOGLE_SMTP_PASSWORD")
	m := gomail.NewMessage()
	m.SetHeader("From", "r.nikookolah@gmail.com")
	m.SetHeader("To", to)
	m.SetHeader("Subject", "PosOnline verification code.")
	bodyMessage := fmt.Sprintf("Now you can try to send the code %s", message)
	m.SetBody("text/plain", bodyMessage)

	// SMTP server configuration
	d := gomail.NewDialer("smtp.gmail.com", 587, "r.nikookolah@gmail.com", GOOGLE_SMTP_PASSWORD)

	// Send the email
	if err := d.DialAndSend(m); err != nil {
		fmt.Println("Failed to send email:", err)
	} else {
		fmt.Println("Email sent successfully.")
	}
}
func GetRandomCode() int {
	rand.Seed(time.Now().UnixNano())

	// Generate a random 5-digit number (10000 to 99999)
	return rand.Intn(90000) + 10000
}
