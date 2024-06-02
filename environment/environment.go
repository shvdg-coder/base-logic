package environment

import (
	"encoding/base64"
	_ "github.com/joho/godotenv/autoload" // Load environment variables from a .env file
	"os"
)

// GetValueAsString retrieves an environment string value given a key.
func GetValueAsString(key string) string {
	return os.Getenv(key)
}

// GetValueAsBytes retrieves an environment []byte value given a key.
func GetValueAsBytes(key string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(GetValueAsString(key))
}
