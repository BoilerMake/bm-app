// Package env handles loading and reading environment variables from a file.
package env

import (
	"bufio"
	"io"
	"os"

	"strings"
)

// Load reads in the file called .env and sets all given environement variables.
// If overwrite is set to true, Load will use the values present in .env instead of the current values.
func Load(overwrite bool) (err error) {
	// env file is ALWAYS is .env
	f, err := os.Open(".env")
	if err != nil {
		return err
	}

	parsedEnv, err := parseFile(f)
	if err != nil {
		return err
	}

	// Set variables
	for key, value := range parsedEnv {
		if os.Getenv("key") == "" || overwrite {
			os.Setenv(key, value)
		}
	}

	return err
}

// parseFile will read an io.reader line by line to make a string string map of environment key values.
// It will ignore lines that start with a #.
func parseFile(r io.Reader) (parsed map[string]string, err error) {
	parsed = make(map[string]string)

	s := bufio.NewScanner(r)
	for s.Scan() {
		line := s.Text()

		// Ignore commented and blank lines
		if line[0] != '#' && len(line) != 0 {
			split := strings.Split(line, "=")

			key := strings.TrimSpace(split[0])
			value := strings.TrimSpace(split[1])

			parsed[key] = value
		}
	}

	if err = s.Err(); err != nil {
		return parsed, err
	}

	return parsed, err
}
