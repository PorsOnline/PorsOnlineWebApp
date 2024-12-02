package helper

import (
	"os"

	"github.com/joho/godotenv"
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
