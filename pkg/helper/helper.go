package helper

import (
	"fmt"
	"os"

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

func SendEmail(to string) {
	m := gomail.NewMessage()
	m.SetHeader("From", "r.nikookolah@gmail.com")
	m.SetHeader("To", to)
	m.SetHeader("Subject", "PosOnline verification code.")
	m.SetBody("text/plain", "Now you can try to send code .")

	// SMTP server configuration
	d := gomail.NewDialer("smtp.gmail.com", 587, "ixjd oort uizn pdee", to)

	// Send the email
	if err := d.DialAndSend(m); err != nil {
		fmt.Println("Not send.")
	}

}
