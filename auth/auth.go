package auth

import (
	"context"
	"errors"
)

var (
	SignupAccountExistsError = errors.New("Account already exists")
)

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

type signupAccount struct {
	email    string
	salt     string
	password string
}

type signupUser struct {
	Firstname string
	Lastname  string
}

type SignupResult struct {
	User    User
	Account Account
}

func (s *Service) Signup(ctx context.Context, input SignupInput) (SignupResult, error) {
	salt := s.Salt.Generate()

	saltedPass := Sha256SaltedPassword{
		Password: input.Password,
		Salt:     salt,
	}

	password := saltedPass.Value()

	createAccountInput := signupAccount{
		email:    input.Email,
		salt:     salt,
		password: password,
	}

	createUserInput := signupUser{
		Firstname: input.Firstname,
		Lastname:  input.Lastname,
	}

	account, err := s.Repo.CreateAccount(ctx, createUserInput, createAccountInput)

	return account, err
}
