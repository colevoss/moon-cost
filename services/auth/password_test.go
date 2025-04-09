package auth

import (
	"moon-cost/moontest"
	"testing"
)

var testSha256SaltedPassString = "7a37b85c8918eac19a9089c0fa5a2ab4dce3f90528dcdeec108b23ddf3607b99"

var testSha256SaltedPass = Sha256SaltedPassword{
	Password: "password",
	Salt:     "salt",
}

func TestComparePassword(t *testing.T) {
	tests := []struct {
		expected SaltedPassword
		actual   SaltedPassword
		matches  bool
	}{
		{expected: BasicSaltedPassword("password"), actual: BasicSaltedPassword("password"), matches: true},
		{expected: BasicSaltedPassword("password"), actual: BasicSaltedPassword("notpassword"), matches: false},
		{expected: testSha256SaltedPass, actual: BasicSaltedPassword(testSha256SaltedPassString), matches: true},
	}

	for _, test := range tests {
		matches := comparePasswords(test.expected, test.actual)

		moontest.Assert(t, matches == test.matches, "Expected %v to match %v", test.expected, test.actual)
	}
}

func TestSha256SaltedPassword(t *testing.T) {
	value := testSha256SaltedPass.Value()

	moontest.Assert(t, testSha256SaltedPassString == value, "Expected %s. Got %s", testSha256SaltedPassString, value)
}
