// Package env handles loading and reading environment variables from a file.
package env

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
)

var (
	ErrNoDotEnv = errors.New("no .env file found")
)

// Load reads in the file called .env and sets all given environement variables.
// If overwrite is set to true, Load will use the values present in .env instead of the current values.
func Load(overwrite bool) (err error) {
	// env file will ALWAYS be .env
	f, err := os.Open(".env")
	if err != nil {
		if err == os.ErrNotExist {
			return ErrNoDotEnv
		}

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

		// Ignore commented and blank lines
		if len(line) != 0 && line[0] != '#' {
			equalsPos := strings.Index(line, "=")

			// Make sure there's actually a value assigned
			if equalsPos == -1 {
				return parsed, fmt.Errorf(".env expected variable assignment at line %d", lineCount)
			}

			// Ignore placeholder assignments, e.g. "PORT="
			if line[equalsPos+1:] != "" {
				key := strings.TrimSpace(line[0:equalsPos])
				value := strings.TrimSpace(line[equalsPos+1:])

				parsed[key] = value
			}
		}
	}

	if err = s.Err(); err != nil {
		return parsed, err
	}

	return parsed, err
}
