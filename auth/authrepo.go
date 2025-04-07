package auth

import "context"

type Repo interface {
	CreateAccount(context.Context, signupUser, signupAccount) (SignupResult, error)
}
