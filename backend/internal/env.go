package internal

import (
	"bufio"
	"os"
	"strings"
)

// LoadEnv loads environment variables from a .env file manually
func LoadEnv(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err // Return if file not found
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// Ignore empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Split by first '=' to get key=value
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		// Remove surrounding quotes (if present)
		value = strings.Trim(value, `"'`)

		// Set env variable
		_ = os.Setenv(key, value)
	}

	return scanner.Err()
}
