package auth

import "context"

type signupAccount struct {
	email    string
	salt     string
	password string
}

type signupUser struct {
	Firstname string
	Lastname  string
}

type Repo interface {
	CreateAccount(context.Context, signupUser, signupAccount) (SignupResult, error)
}

type NoopRepo struct{}

func (n *NoopRepo) CreateAccount(context.Context, signupUser, signupAccount) (SignupResult, error) {
	return SignupResult{}, nil
}
