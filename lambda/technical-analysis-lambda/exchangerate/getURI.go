package exchangerate

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func getURI(name string) string {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
		value := os.Getenv(name)
		return value
	}

	// now you can use os.Getenv ...
	value := os.Getenv(name)
	return value
}
