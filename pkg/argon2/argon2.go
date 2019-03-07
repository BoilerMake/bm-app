package argon2

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"golang.org/x/crypto/argon2"
)

var (
	ErrInvalidEncodedHash = errors.New("hash is in an unrecognizable format")
	ErrAlgorithmMismatch  = errors.New("hash does not use argon2id")
	ErrVersionMismatch    = errors.New("hash's argon2id version differs from current")
	ErrEmptyPassword      = errors.New("password cannot be empty")
)

const (
	DefaultHashTime       = 4
	DefaultHashMemory     = 64 * 1024
	DefaultHashThreads    = 4
	DefaultHashKeyLength  = 64
	DefaultHashSaltLength = 16
)

var (
	DefaultParameters = &Parameters{
		Time:       DefaultHashTime,
		Memory:     DefaultHashMemory,
		Threads:    DefaultHashThreads,
		KeyLength:  DefaultHashKeyLength,
		SaltLength: DefaultHashSaltLength,
	}
)

type Parameters struct {
	Time       uint32 // Number of passes to run
	Memory     uint32 // Size of memory to use in KiB
	Threads    uint8  // Number of threads to use
	KeyLength  uint32 // Length of resulting hash
	SaltLength uint32 // Length of salt to generate
}

// HashPassword returns an encoded string with all the necessary parameters
// encoded, separated by "$".
func (p *Parameters) HashPassword(password string) (encoded string, err error) {
	if password == "" {
		return encoded, ErrEmptyPassword
	}

	salt, err := randomBytes(p.SaltLength)
	if err != nil {
		return encoded, err
	}

	hash := argon2.IDKey([]byte(password), salt, p.Time, p.Memory, p.Threads, p.KeyLength)

	b64salt := base64.RawStdEncoding.EncodeToString(salt)
	b64hash := base64.RawStdEncoding.EncodeToString(hash)

	encoded = fmt.Sprintf("$%s$v=%d$t=%d,m=%d,p=%d$%s$%s",
		"argon2id",
		argon2.Version,
		p.Time,
		p.Memory,
		p.Threads,
		b64salt,
		b64hash,
	)

	return encoded, err
}

// CheckPassword compares a string to a known encoded hash, using all the same
// parameters and salt as the known hash.  In cases where you might expect a
// normal function to return an err, we'll just return false.
func CheckPassword(password string, encoded string) (match bool) {
	params, salt, key, err := decode(encoded)
	if err != nil {
		// TODO maybe consider logging errors here or forgetting what's
		// written above and returning/logging those errors because if an error
		// shows up here then something's v messed up
		return false
	}

	inputHash := argon2.IDKey([]byte(password), salt, params.Time, params.Memory, params.Threads, params.KeyLength)

	if subtle.ConstantTimeCompare(key, inputHash) == 1 {
		return true
	}

	return false
}

// decode separates an encoded hash at every "$" and returns the parameters,
// salt, and hash for that encoded string.  It's pretty much the same as
// ParseParameters but also gets the salt and hash.
func decode(encoded string) (p *Parameters, salt, hash []byte, err error) {
	params := Parameters{}

	split := strings.Split(encoded, "$")
	if len(split) != 6 {
		return &params, salt, hash, ErrInvalidEncodedHash
	}

	if split[1] != "argon2id" {
		return &params, salt, hash, ErrAlgorithmMismatch
	}

	var version int
	_, err = fmt.Sscanf(split[2], "v=%d", &version)
	if err != nil {
		return &params, salt, hash, err
	}
	if version != argon2.Version {
		return &params, salt, hash, ErrVersionMismatch
	}

	_, err = fmt.Sscanf(split[3], "t=%d,m=%d,p=%d", &params.Time, &params.Memory, &params.Threads)
	if err != nil {
		return &params, salt, hash, err
	}

	salt, err = base64.RawStdEncoding.DecodeString(split[4])
	if err != nil {
		return &params, salt, hash, err
	}
	params.SaltLength = uint32(len(salt))

	hash, err = base64.RawStdEncoding.DecodeString(split[5])
	if err != nil {
		return &params, salt, hash, err
	}
	params.KeyLength = uint32(len(hash))

	return &params, salt, hash, err
}

func randomBytes(n uint32) (b []byte, err error) {
	b = make([]byte, n)
	_, err = rand.Read(b)
	return b, err
}
