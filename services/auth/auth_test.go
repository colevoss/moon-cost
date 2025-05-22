package auth

import (
	"context"
	"fmt"
	"moon-cost/moon"
	"runtime"
	"testing"
)

type testSigUpAuthRepo struct {
	NoopRepo
	exists bool

	createUserId    string
	createAccountId string

	createdPassword string
}

func (t *testSigUpAuthRepo) AccountExists(context.Context, string) (bool, error) {
	return t.exists, nil
}

func (t *testSigUpAuthRepo) CreateAccount(ctx context.Context, input createAccountInput) (SignupResult, error) {
	t.createdPassword = input.account.password

	return SignupResult{
		UserId:    t.createUserId,
		AccountId: t.createAccountId,
	}, nil
}

func TestAuthServiceSignUpNewUser(t *testing.T) {
	repo := testSigUpAuthRepo{
		exists:          false,
		createUserId:    "user-id",
		createAccountId: "account-id",
	}

	service := Service{
		Salt: StaticSalt{"salt"},
		Repo: &repo,
	}

	input := SignupInput{
		Email:     "email",
		Password:  "password",
		Firstname: "firstname",
		Lastname:  "lastname",
	}

	// sha256 of "passwordsalt"
	expectedPassword := "7a37b85c8918eac19a9089c0fa5a2ab4dce3f90528dcdeec108b23ddf3607b99"

	_, err := service.SignUp(context.Background(), input)

	if err != nil {
		t.Errorf("serivce.Signup() = _, %s. want err", err)
	}

	if repo.createdPassword != expectedPassword {
		t.Errorf("created user password = %s. want %s", repo.createdPassword, expectedPassword)
	}
}

func TestAuthServiceSignUpExistinguser(t *testing.T) {
	moon.DisableSlog(t)

	repo := testSigUpAuthRepo{
		exists: true,
	}

	service := Service{
		Salt: DefaultSalt,
		Repo: &repo,
	}

	_, err := service.SignUp(context.Background(), SignupInput{})

	if err != ErrSignupAccountExists {
		t.Errorf("service.Signup() = _, %s. want %s", err, ErrSignupAccountExists)
	}
}

type testSignInAuthRepo struct {
	NoopRepo
	exists  bool
	account Account
}

func (t *testSignInAuthRepo) FindAccountByEmail(context.Context, string) (Account, error) {
	if !t.exists {
		return Account{}, ErrAccountNotFound
	}

	return t.account, nil
}

func TestAuthServiceSignInNotFoundWhenNoUser(t *testing.T) {
	repo := testSignInAuthRepo{
		exists: false,
	}

	service := Service{
		Salt: DefaultSalt,
		Repo: &repo,
	}

	input := SignInInput{
		Email:    "email",
		Password: "password",
	}

	_, err := service.SignIn(context.Background(), input)

	if err != ErrAccountNotFound {
		t.Errorf("service.SignIn() = _, %s. want %s", err, ErrAccountNotFound)
	}
}

func TestAuthServiceSignInWrongPassword(t *testing.T) {
	salted := Sha256SaltedPassword{
		Password: "password",
		Salt:     DefaultSalt.Salt(),
	}

	account := Account{
		Email:    "email",
		Password: salted.SaltPassword(),
		Salt:     salted.Salt,
	}

	repo := testSignInAuthRepo{
		exists:  true,
		account: account,
	}

	service := Service{
		Repo: &repo,
	}

	input := SignInInput{
		Email:    "email",
		Password: "wrong-password",
	}

	_, err := service.SignIn(context.Background(), input)

	if err != ErrAccountNotFound {
		t.Errorf("service.SignIn(wrong password) = _, %s. want %s", err, ErrAccountNotFound)
	}
}

func TestAuthServiceSignInSuccess(t *testing.T) {
	_, filename, _, _ := runtime.Caller(0)

	fmt.Printf("filename: %v\n", filename)

	salted := Sha256SaltedPassword{
		Password: "password",
		Salt:     DefaultSalt.Salt(),
	}

	account := Account{
		Id:       "account-id",
		UserId:   "user-id",
		Email:    "email",
		Password: salted.SaltPassword(),
		Salt:     salted.Salt,
	}

	repo := testSignInAuthRepo{
		exists:  true,
		account: account,
	}

	service := Service{
		Repo: &repo,
	}

	input := SignInInput{
		Email:    "email",
		Password: "password",
	}

	res, err := service.SignIn(context.Background(), input)

	if err != nil {
		t.Errorf("service.SignIn(wrong password) = _, %s error. want nil", err)
	}

	if res.UserId != account.UserId {
		t.Errorf("signingResult.UserId = %s. want %s", res.UserId, account.UserId)
	}

	if res.AccountId != account.Id {
		t.Errorf("signingResult.AccountId = %s. want %s", res.AccountId, account.Id)
	}
}
