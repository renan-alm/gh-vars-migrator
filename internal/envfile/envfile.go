// Package envfile provides a lightweight .env file parser that loads
// key-value pairs into the process environment. Variables already set
// in the environment are never overwritten.
package envfile

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"
)

// Load reads a .env file and sets any variables that are not already
// present in the environment. It silently returns nil when the file
// does not exist so callers don't need to guard with os.Stat first.
func Load(path string) error {
	f, err := os.Open(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil // missing .env file is not an error
		}
		return fmt.Errorf("opening env file: %w", err)
	}
	defer f.Close() //nolint:errcheck // best-effort close on read-only file

	scanner := bufio.NewScanner(f)
	lineNum := 0

	for scanner.Scan() {
		lineNum++
		line := strings.TrimSpace(scanner.Text())

		// Skip blank lines and comments.
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Strip optional "export " prefix.
		line = strings.TrimPrefix(line, "export ")

		key, value, err := parseLine(line)
		if err != nil {
			return fmt.Errorf("env file line %d: %w", lineNum, err)
		}

		// Only set variables that are not already in the environment so
		// real env vars and CLI flags always take precedence.
		if _, exists := os.LookupEnv(key); !exists {
			if err := os.Setenv(key, value); err != nil {
				return fmt.Errorf("setting env var %s: %w", key, err)
			}
		}
	}

	return scanner.Err()
}

// parseLine splits a "KEY=VALUE" line and returns the unquoted key and
// value. It supports unquoted, single-quoted, and double-quoted values.
func parseLine(line string) (string, string, error) {
	parts := strings.SplitN(line, "=", 2)
	if len(parts) != 2 {
		return "", "", fmt.Errorf("expected KEY=VALUE, got %q", line)
	}

	key := strings.TrimSpace(parts[0])
	value := strings.TrimSpace(parts[1])

	if key == "" {
		return "", "", fmt.Errorf("empty key in %q", line)
	}

	// Remove surrounding quotes if present.
	if len(value) >= 2 {
		first, last := value[0], value[len(value)-1]
		if (first == '"' && last == '"') || (first == '\'' && last == '\'') {
			value = value[1 : len(value)-1]
		}
	}

	return key, value, nil
}
