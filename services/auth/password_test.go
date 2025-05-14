package auth

import (
	"testing"
)

var testSha256SaltedPassString = "7a37b85c8918eac19a9089c0fa5a2ab4dce3f90528dcdeec108b23ddf3607b99"

var testSha256SaltedPass = Sha256SaltedPassword{
	Password: "password",
	Salt:     "salt",
}

func TestComparePassword(t *testing.T) {
	tests := []struct {
		expected PasswordSalter
		actual   PasswordSalter
		matches  bool
	}{
		{expected: BasicSaltedPassword("password"), actual: BasicSaltedPassword("password"), matches: true},
		{expected: BasicSaltedPassword("password"), actual: BasicSaltedPassword("notpassword"), matches: false},
		{expected: testSha256SaltedPass, actual: BasicSaltedPassword(testSha256SaltedPassString), matches: true},
	}

	for _, test := range tests {
		matches := ComparePasswords(test.expected, test.actual)

		if matches != test.matches {
			t.Errorf("comparePasswords(%s, %s) == %t. want %t", test.expected, test.actual, matches, test.matches)
		}
	}
}

func TestSha256SaltedPassword(t *testing.T) {
	value := testSha256SaltedPass.SaltPassword()

	if value != testSha256SaltedPassString {
		t.Errorf("Sha256SaltedPassword.Value() == %s. want %s", value, testSha256SaltedPassString)
	}
}
