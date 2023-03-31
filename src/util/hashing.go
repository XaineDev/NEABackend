package util

import (
	"crypto/rand"
	"encoding/base64"
	"golang.org/x/crypto/argon2"
)

var (
	params = argonParams{
		memory:      64 * 1024,
		iterations:  3,
		parallelism: 2,
		saltLength:  16,
		keyLength:   32,
	}
)

type argonParams struct {
	memory      uint32
	iterations  uint32
	parallelism uint8
	saltLength  uint32
	keyLength   uint32
}

func HashPassowrd(password string) (hash string, salt string, err error) {
	// Generate a cryptographically secure random salt.
	salt, err = generateSalt()
	if err != nil {
		return "", "", err
	}

	// Pass the plaintext password, salt and parameters to the argon2.IDKey
	// function. This will generate a hash of the password using the Argon2id
	// variant.
	hash = EncodeToString(argon2.IDKey([]byte(password), []byte(salt), params.iterations, params.memory, params.parallelism, params.keyLength))

	return hash, salt, nil
}

func HashPasswordWithSalt(password string, salt string) (hash string, err error) {
	// Pass the plaintext password, salt and parameters to the argon2.IDKey
	// function. This will generate a hash of the password using the Argon2id
	// variant.
	hash = EncodeToString(argon2.IDKey([]byte(password), []byte(salt), params.iterations, params.memory, params.parallelism, params.keyLength))

	return hash, nil
}

func generateRandomBytes(n uint32) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func generateSalt() (string, error) {
	b, err := generateRandomBytes(params.saltLength)
	if err != nil {
		return "", err
	}
	return EncodeToString(b), nil
}

func EncodeToString(b []byte) string {
	return base64.StdEncoding.EncodeToString(b)
}
