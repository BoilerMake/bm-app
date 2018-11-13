// Package env handles loading and reading environment variables from a file.
package env

import (
	"bufio"
	"fmt"
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

	lineCount := 0
	s := bufio.NewScanner(r)
	for s.Scan() {
		lineCount++
		line := s.Text()

		// TODO maybe use other string method instead of split, you're rejoining anyway...
		// Ignore commented and blank lines
		if len(line) != 0 && line[0] != '#' {
			split := strings.Split(line, "=")

			// Make sure there's actually a value assigned
			if len(split) < 2 {
				return parsed, fmt.Errorf("error parsing .env at line %d", lineCount)
			}

			// Ignore placeholder assignments (like "PORT=")
			if split[1] != "" {
				key := strings.TrimSpace(split[0])
				value := strings.TrimSpace(strings.Join(split[1:], "="))

				parsed[key] = value
			}
		}
	}

	if err = s.Err(); err != nil {
		return parsed, err
	}

	return parsed, err
}
