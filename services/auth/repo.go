package auth

import (
	"context"
)

type createAccountAccount struct {
	email    string
	salt     string
	password string
}

type createAccountUser struct {
	firstname string
	lastname  string
}

type createAccountInput struct {
	account createAccountAccount
	user    createAccountUser
}

type Repo interface {
	CreateAccount(context.Context, createAccountInput) (SignupResult, error)
	FindAccountByEmail(context.Context, string) (Account, error)
	AccountExists(context.Context, string) (bool, error)
}

type NoopRepo struct{}

func (n *NoopRepo) CreateAccount(context.Context, createAccountInput) (SignupResult, error) {
	return SignupResult{}, nil
}

func (n *NoopRepo) FindAccountByEmail(context.Context, string) (Account, error) {
	return Account{}, nil
}

func (n *NoopRepo) AccountExists(context.Context, string) (bool, error) {
	return false, nil
}
