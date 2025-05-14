package auth

import (
	"context"
	"errors"
	"log/slog"
	"moon-cost/logging"
)

var (
	SignupAccountExistsError = errors.New("Account already exists")
)

var defaultSalt RandomSalt

type Service struct {
	Salt   Salt
	Repo   Repo
	Logger *slog.Logger
}

func NewService(repo Repo, logger *slog.Logger) *Service {
	return &Service{
		Repo:   repo,
		Salt:   defaultSalt,
		Logger: logging.Logger(logger, slog.String("service", "auth")),
	}
}

type Signup struct {
	Email     string
	Password  string
	Firstname string
	Lastname  string
}

type SignupResult struct {
	User    User
	Account Account
}

func (s *Service) Signup(ctx context.Context, input Signup) (SignupResult, error) {
	salt := s.Salt.Salt()

	saltedPass := Sha256SaltedPassword{
		Password: input.Password,
		Salt:     salt,
	}

	password := saltedPass.SaltPassword()

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
