package auth

import (
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"fmt"
)

type SaltedPassword interface {
	Value() string
}

// Concats a raw password and a salt and sha256 encodes them
type Sha256SaltedPassword struct {
	Password string
	Salt     string
}

func (s Sha256SaltedPassword) Value() string {
	combined := s.Password + s.Salt
	h := sha256.New()
	h.Write([]byte(combined))
	hashed := h.Sum(nil)

	return fmt.Sprintf("%x", hashed)
}

// Used to represent a basic string as a SaltedPassword
type BasicSaltedPassword string

func (b BasicSaltedPassword) Value() string {
	return string(b)
}

type Salt interface {
	Generate() string
}

type RandomSalt struct {
	Length int
}

func (r RandomSalt) Generate() string {
	salt := make([]byte, r.Length)

	rand.Read(salt)

	return fmt.Sprintf("%x", string(salt))
}

func comparePasswords(expected SaltedPassword, actual SaltedPassword) bool {
	expectedValue := expected.Value()
	actualValue := actual.Value()

	compared := subtle.ConstantTimeCompare([]byte(expectedValue), []byte(actualValue))

	return compared == 1
}
