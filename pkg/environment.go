package pkg

import (
	"encoding/base64"
	"fmt"
	_ "github.com/joho/godotenv/autoload" // Load environment variables from a .env file
	"os"
	"strconv"
)

// GetEnvValueAsString retrieves an environment string value given a key.
func GetEnvValueAsString(key string) string {
	return os.Getenv(key)
}

// GetEnvValueAsBoolean retrieves an environment boolean value given a key.
func GetEnvValueAsBoolean(key string) bool {
	boolValue, err := strconv.ParseBool(os.Getenv(key))
	if err != nil {
		fmt.Printf("Error while converting environment value '%s' to bool", key)
	}
	return boolValue
}

// GetEnvValueAsBytes retrieves an environment []byte value given a key.
func GetEnvValueAsBytes(key string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(GetEnvValueAsString(key))
}
