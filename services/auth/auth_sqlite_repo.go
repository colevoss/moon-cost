package auth

import (
	"context"
	"database/sql"
	"fmt"
)

type SQLiteRepo struct {
	db *sql.DB
}

func (s *SQLiteRepo) CreateAccount(ctx context.Context, userInput signupUser, accountInput signupAccount) (SignupResult, error) {
	signupResult := SignupResult{}

	tx, err := s.db.BeginTx(ctx, nil)

	if err != nil {
		return signupResult, err
	}

	defer tx.Rollback()

	createUserRes, err := tx.ExecContext(
		ctx,
		`INSERT INTO users (firstname, lastname) VALUES (?, ?)`,
		userInput.Firstname,
		userInput.Lastname,
	)

	if err != nil {
		return signupResult, err
	}

	userId, err := createUserRes.LastInsertId()

	if err != nil {
		return signupResult, err
	}

	createAccountRes, err := tx.ExecContext(
		ctx,
		`INSERT INTO accounts (email, password, salt, active, userId)
    VALUES (?, ?, ?, ?, ?)`,
		accountInput.email,
		accountInput.password,
		accountInput.salt,
		1,
		userId,
	)

	if err != nil {
		return signupResult, err
	}

	// TODO: Get account
	accountId, err := createAccountRes.LastInsertId()

	if err != nil {
		return signupResult, err
	}

	fmt.Printf("Account Id %d", accountId)

	signupResult.User, err = s.getUserTx(ctx, tx, int(userId))

	if err != nil {
		return signupResult, err
	}

	return signupResult, nil
}

const getUserQuery = `SELECT id, firstname, lastname FROM users WHERE id = ?`

func (s *SQLiteRepo) getUserTx(ctx context.Context, tx *sql.Tx, userId int) (User, error) {
	var user User

	err := tx.QueryRowContext(
		ctx,
		getUserQuery,
		userId,
	).Scan(
		&user.Id,
		&user.Firstname,
		&user.Lastname,
	)

	return user, err
}
