package util

import (
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"net/http"
	"os"
)

// CheckRequiredVariables checks if the required environment variables are set and returns the list of variables
// that are not set.
func CheckRequiredVariables(requiredVariables []string) []string {
	missingVariables := []string{}
	for _, variable := range requiredVariables {
		if value, exists := os.LookupEnv(variable); !exists || value == "" {
			missingVariables = append(missingVariables, variable)
		}
	}
	return missingVariables
}

// EncodeResource encodes a string to base64.
func EncodeResource(username string, clientSecret string) string {
	authHeaderString := fmt.Sprintf("%s:%s", username, clientSecret)
	return base64.StdEncoding.EncodeToString([]byte(authHeaderString))
}

// GetTransportSettings disables certificate validation for http requests if the insecure env variable has been set.
func GetTransportSettings() *http.Transport {
	variableName := "INSECURE_SKIP_TLS_VERIFY"
	if value, exists := os.LookupEnv(variableName); exists && value == "true" {
		return &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}
	}

	return &http.Transport{}
}
