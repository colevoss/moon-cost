package auth

import (
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"fmt"
	"moon-cost/assert"
)

// The PasswordSalter interface wraps the SaltPassword function which
// returns salted password
type PasswordSalter interface {
	SaltPassword() string
}

// Concats a raw password and a salt and sha256 encodes them
type Sha256SaltedPassword struct {
	Password string
	Salt     string
}

func (s Sha256SaltedPassword) SaltPassword() string {
	assert.Ok(s.Password != "", "Password must be populated and not blank")
	assert.Ok(s.Salt != "", "Salt must be populated and not blank")

	combined := fmt.Sprintf("%s%s", s.Password, s.Salt)
	h := sha256.New()
	h.Write([]byte(combined))
	hashed := h.Sum(nil)

	return fmt.Sprintf("%x", hashed)
}

// Used to represent a basic string as a SaltedPassword
type BasicSaltedPassword string

func (b BasicSaltedPassword) SaltPassword() string {
	return string(b)
}

// The Salt interface wraps the
// Salt generates a string that should be used when salting a password
type Salt interface {
	Salt() string
}

// Salt implementation that returns a cryptographically secure random string
// of Length length
type RandomSalt struct {
	Length int
}

func (r RandomSalt) Salt() string {
	salt := make([]byte, r.Length)

	rand.Read(salt)

	return fmt.Sprintf("%x", string(salt))
}

func ComparePasswords(expected PasswordSalter, actual PasswordSalter) bool {
	expectedValue := expected.SaltPassword()
	actualValue := actual.SaltPassword()

	compared := subtle.ConstantTimeCompare([]byte(expectedValue), []byte(actualValue))

	return compared == 1
}
