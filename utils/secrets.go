package utils

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// GetSecretFromKey reads a file containing key=value pairs and returns the value for the given key.
// It falls back to environment variables if the file or key is not found.
func GetSecretFromKey(fileName string, key string) string {
	filePath := fmt.Sprintf("/run/secrets/%s", fileName)
	content, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Printf("Error reading secret file %q: %v. Using environment variable fallback.\n", fileName, err)
		return os.Getenv(key)
	}

	secrets := make(map[string]string)
	scanner := bufio.NewScanner(strings.NewReader(string(content)))
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			secrets[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
		}
	}

	if val, ok := secrets[key]; ok {
		return val
	}

	fmt.Printf("Secret %q not found in %q. Using environment variable fallback.\n", key, fileName)
	return os.Getenv(key)
}

// GetSecret reads the entire content of a secret file.
// It falls back to environment variables if the file is not found.
func GetSecret(name string) string {
	filePath := fmt.Sprintf("/run/secrets/%s", name)
	content, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Printf("Secret %q not found. Using environment variable fallback.\n", name)
		return os.Getenv(name)
	}
	return strings.TrimSpace(string(content))
}
