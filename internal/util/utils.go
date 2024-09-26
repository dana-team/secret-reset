package util

import (
	"encoding/base64"
	"fmt"
	"os"
)

// CheckRequiredVariables checks if the required environment variables are set and returns the list of variables
// that are not set.
func CheckRequiredVariables(requiredVariables []string) []string {
	missingVariables := []string{}
	for _, varibale := range requiredVariables {
		if value, exists := os.LookupEnv(varibale); !exists || value == "" {
			missingVariables = append(missingVariables, value)
		}
	}
	return missingVariables
}

// EncodeResource encodes a string to base64.
func EncodeResource(username string, clientSecret string) string {
	authHeaderString := fmt.Sprintf("%s:%s", username, clientSecret)
	return base64.StdEncoding.EncodeToString([]byte(authHeaderString))
}
