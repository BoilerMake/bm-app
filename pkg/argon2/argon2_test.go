package argon2

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"reflect"
	"testing"
)

const (
	algorithm = "argon2id"
	password  = "hunter2"
)

func TestHashPassword(t *testing.T) {
	tt := []struct {
		name      string
		input     string
		shouldErr bool
		err       error
	}{
		{
			name:      "standard",
			input:     password,
			shouldErr: false,
		},
		{
			name:      "random password",
			input:     "aosifu02j0as8dfu($#UFPS)",
			shouldErr: false,
		},
		{
			name:      "empty",
			input:     "",
			shouldErr: true,
			err:       ErrEmptyPassword,
		},
	}

	for _, tc := range tt {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			result, err := DefaultParameters.HashPassword(tc.input)
			if err != nil && !tc.shouldErr {
				t.Fatalf("decoding failed: %v", err)
			} else if err != nil && tc.shouldErr {
				if tc.err != err {
					t.Fatalf("expected error: %v, got: %v", tc.err, err)
				} else {
					return
				}
			}

			params, _, _, err := decode(result)
			if err != nil {
				t.Fatalf("decoding encoded password failed: %v", err)
			}

			if !reflect.DeepEqual(params, DefaultParameters) {
				t.Errorf("expected params to be: %v, got: %v", DefaultParameters, params)
			}
		})
	}
}

func TestHashingSamePasswords(t *testing.T) {
	// Testing that the same password results in a different hash.
	// It's possible for passwords hashes to collide, but that chance is super
	// small.  If it (somehow) happens just run the test again, and if it's
	// happening consistently then something's probably wrong.

	e1, err := DefaultParameters.HashPassword(password)
	if err != nil {
		t.Fatalf("failed while hashing password: %v", err)
	}

	e2, err := DefaultParameters.HashPassword(password)
	if err != nil {
		t.Fatalf("failed while hashing password: %v", err)
	}

	if e1 == e2 {
		t.Errorf("resulting encoded passwords are the same")
	}

	_, _, hash1, err := decode(e1)
	if err != nil {
		t.Fatalf("failed while hashing password: %v", err)
	}

	_, _, hash2, err := decode(e2)
	if err != nil {
		t.Fatalf("failed while hashing password: %v", err)
	}

	// Testing hashes specifically in case salts aren't being accounted for or
	// something along those lines.
	if bytes.Equal(hash1, hash2) {
		t.Errorf("resulting encoded hashes are the same")
	}
}

func TestCheckPassword(t *testing.T) {
	e1, err := DefaultParameters.HashPassword(password)
	if err != nil {
		t.Fatalf("failed while hashing password: %v", err)
	}

	tt := []struct {
		name     string
		password string
		expected bool
	}{
		{
			name:     "close password",
			password: "hunter1",
			expected: false,
		},
		{
			name:     "empty",
			password: "",
			expected: false,
		},
		{
			name:     "correct",
			password: password,
			expected: true,
		},
	}

	for _, tc := range tt {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			if CheckPassword(tc.password, e1) != tc.expected {
				t.Errorf("expected CheckPassword to return %v", tc.expected)
			}
		})
	}
}

func TestDecode(t *testing.T) {
	testSalt := []byte("S4LTS4LTS4LTS4LT")
	testHash := []byte("H4SH 0R K3Y 0R WH4AT3V3RH4SH 0R K3Y 0R WH4AT3V3RH4SH 0R K3Y 0R W")

	b64Salt := base64.RawStdEncoding.EncodeToString([]byte(testSalt))
	b64hash := base64.RawStdEncoding.EncodeToString([]byte(testHash))

	type decodeResponse struct {
		params Parameters
		salt   []byte
		hash   []byte
		err    error
	}

	tt := []struct {
		name      string
		input     string
		expected  decodeResponse
		shouldErr bool
	}{
		{
			name: "standard",
			input: fmt.Sprintf("$%s$v=%d$t=%d,m=%d,p=%d$%s$%s",
				algorithm,
				19,
				DefaultHashTime,
				DefaultHashMemory,
				DefaultHashThreads,
				b64Salt,
				b64hash,
			),
			expected: decodeResponse{
				params: *DefaultParameters,
				salt:   []byte(testSalt),
				hash:   []byte(testHash),
				err:    nil,
			},
			shouldErr: false,
		},
		{
			name: "wrong algorithm",
			input: fmt.Sprintf("$%s$v=%d$t=%d,m=%d,p=%d$%s$%s",
				"argonjk tho not actually",
				19,
				DefaultHashTime,
				DefaultHashMemory,
				DefaultHashThreads,
				b64Salt,
				b64hash,
			),
			expected: decodeResponse{
				err: ErrAlgorithmMismatch,
			},
			shouldErr: true,
		},
		{
			name: "wrong version",
			input: fmt.Sprintf("$%s$v=%d$t=%d,m=%d,p=%d$%s$%s",
				"argon2id",
				2,
				DefaultHashTime,
				DefaultHashMemory,
				DefaultHashThreads,
				b64Salt,
				b64hash,
			),
			expected: decodeResponse{
				err: ErrVersionMismatch,
			},
			shouldErr: true,
		},
		{
			name: "missing parameters",
			input: fmt.Sprintf("$%s$v=%d$t=%d,p=%d$%s$%s",
				"argon2id",
				19,
				DefaultHashTime,
				//DefaultHashMemory,
				DefaultHashThreads,
				b64Salt,
				b64hash,
			),
			expected: decodeResponse{
				err: errors.New("input does not match format"),
			},
			shouldErr: true,
		},
		{
			name: "missing part (aka missing $)",
			input: fmt.Sprintf("$%s$v=%d$t=%d,m=%d,p=%d$%s",
				"argon2id",
				19,
				DefaultHashTime,
				DefaultHashMemory,
				DefaultHashThreads,
				b64hash,
			),
			expected: decodeResponse{
				err: ErrInvalidEncodedHash,
			},
			shouldErr: true,
		},
		{
			name:  "empty",
			input: "",
			expected: decodeResponse{
				err: ErrInvalidEncodedHash,
			},
			shouldErr: true,
		},
	}

	for _, tc := range tt {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			params, salt, hash, err := decode(tc.input)
			if err != nil && !tc.shouldErr {
				t.Fatalf("decoding failed: %v", err)
			} else if err != nil && tc.shouldErr {
				if tc.expected.err.Error() != err.Error() {
					t.Fatalf("expected error: %v, got: %v", tc.expected.err, err)
				} else {
					return
				}
			}

			result := decodeResponse{*params, salt, hash, err}

			if !reflect.DeepEqual(result, tc.expected) {
				t.Errorf("expected: %v, got: %v", tc.expected, result)
			}
		})
	}
}
