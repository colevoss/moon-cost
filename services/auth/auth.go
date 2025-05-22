package auth

import (
	"context"
	"errors"
	"log/slog"
)

var (
	ErrSignupAccountExists = errors.New("account already exists")
	ErrAccountNotFound     = errors.New("account not found")
)

var DefaultSalt = RandomSalt{
	Length: 32,
}

type Service struct {
	Salt Salt
	Repo Repo
}

type SignupInput struct {
	Email     string
	Password  string
	Firstname string
	Lastname  string
}

type SignupResult struct {
	UserId    string `json:"userId"`
	AccountId string `json:"accountId"`
}

func (s *Service) SignUp(ctx context.Context, input SignupInput) (SignupResult, error) {
	exists, err := s.Repo.AccountExists(ctx, input.Email)

	if err != nil {
		return SignupResult{}, err
	}

	if exists {
		slog.WarnContext(
			ctx,
			"user with email already exists",
			slog.String("email", input.Email),
		)

		return SignupResult{}, ErrSignupAccountExists
	}

	salt := s.Salt.Salt()

	saltedPass := Sha256SaltedPassword{
		Password: input.Password,
		Salt:     salt,
	}

	password := saltedPass.SaltPassword()

	accountInput := createAccountAccount{
		email:    input.Email,
		salt:     salt,
		password: password,
	}

	userInput := createAccountUser{
		firstname: input.Firstname,
		lastname:  input.Lastname,
	}

	account, err := s.Repo.CreateAccount(ctx, createAccountInput{
		account: accountInput,
		user:    userInput,
	})

	return account, err
}

type SignInInput struct {
	Email    string
	Password string
}

type SignInResult struct {
	UserId    string `json:"userId"`
	AccountId string `json:"accountId"`
}

func (s *Service) SignIn(ctx context.Context, input SignInInput) (SignInResult, error) {
	var signin SignInResult

	account, err := s.Repo.FindAccountByEmail(ctx, input.Email)

	// we don't return here because we want to still do password comparison to keep consistent timing
	// it could be argued that the log will produce inconsistent timing as well
	if err != nil {
		slog.WarnContext(
			ctx,
			"error querying for user by email",
			slog.String("email", input.Email),
			slog.String("err", err.Error()),
		)
	}

	submittedPass := Sha256SaltedPassword{
		Password: input.Password,
		Salt:     account.Salt,
	}

	expectedPass := BasicSaltedPassword(account.Password)
	passwordsMatch := ComparePasswords(submittedPass, expectedPass)

	if !passwordsMatch {
		return signin, ErrAccountNotFound
	}

	signin.UserId = account.UserId
	signin.AccountId = account.Id

	return signin, nil
}
