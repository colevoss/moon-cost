package auth

import (
	"context"
	"database/sql"
	"log/slog"
	"moon-cost/moon"
)

type SQLiteRepo struct {
	DB *sql.DB
	ID moon.IDGenerator
}

func (s *SQLiteRepo) CreateAccount(ctx context.Context, input createAccountInput) (SignupResult, error) {
	var signupResult SignupResult

	tx, err := s.DB.BeginTx(ctx, nil)

	if err != nil {
		return signupResult, err
	}

	defer tx.Rollback()

	userId, err := s.createUser(ctx, tx, input.user)

	if err != nil {
		return signupResult, err
	}

	signupResult.UserId = userId.String()

	accountId, err := s.createAccount(ctx, tx, input.account, userId)

	if err != nil {
		return signupResult, err
	}

	signupResult.AccountId = accountId.String()

	if err := tx.Commit(); err != nil {
		slog.ErrorContext(
			ctx,
			"error committing create account transaction",
			slog.String("err", err.Error()),
		)
		return signupResult, err
	}

	return signupResult, nil
}

func (s *SQLiteRepo) AccountExists(ctx context.Context, email string) (bool, error) {
	_, err := s.findAccountByEmail(ctx, email)

	if err == ErrAccountNotFound {
		return false, nil
	}

	return true, err
}

func (s *SQLiteRepo) FindAccountByEmail(ctx context.Context, email string) (Account, error) {
	account, err := s.findAccountByEmail(ctx, email)

	return account, err
}

func (s *SQLiteRepo) findAccountByEmail(ctx context.Context, email string) (Account, error) {
	slog.DebugContext(
		ctx,
		"querying for account by email",
		slog.String("email", email),
	)

	row := s.DB.QueryRowContext(
		ctx,
		AuthQueries.SQLite("findAccountByEmail"),
		email,
	)

	var account Account

	err := row.Scan(
		&account.Id,
		&account.Email,
		&account.Password,
		&account.Salt,
		&account.Active,
		&account.UserId,
	)

	if err == sql.ErrNoRows {
		return account, ErrAccountNotFound
	}

	return account, err
}

func (s *SQLiteRepo) createUser(ctx context.Context, e moon.Execer, input createAccountUser) (moon.UUID, error) {
	id, err := s.ID.ID()

	if err != nil {
		slog.ErrorContext(ctx, "error creating uuid", slog.String("err", err.Error()))
		return id, err
	}

	slog.DebugContext(
		ctx,
		"creating user",
		slog.String("id", id.String()),
	)

	_, err = e.ExecContext(
		ctx,
		AuthQueries.SQLite("createUser"),
		id,
		input.firstname,
		input.lastname,
	)

	if err != nil {
		slog.DebugContext(
			ctx,
			"error creating user in sqlite",
			slog.String("err", err.Error()),
			slog.String("id", id.String()),
		)
	}

	slog.InfoContext(
		ctx,
		"created user in sqlite",
		slog.String("id", id.String()),
	)

	return id, err
}

func (s *SQLiteRepo) createAccount(ctx context.Context, e moon.Execer, input createAccountAccount, userId moon.UUID) (moon.UUID, error) {
	id, err := s.ID.ID()

	if err != nil {
		slog.ErrorContext(ctx, "error creating uuid", slog.String("err", err.Error()))
		return id, err
	}

	slog.DebugContext(
		ctx,
		"creating account",
		slog.String("id", id.String()),
	)

	_, err = e.ExecContext(
		ctx,
		AuthQueries.SQLite("createAccount"),
		id,
		input.email,
		input.password,
		input.salt,
		moon.DBTrue,
		userId,
	)

	if err != nil {
		slog.ErrorContext(
			ctx,
			"error creating account in sqlite",
			slog.String("id", id.String()),
			slog.String("error", err.Error()),
		)
	}

	slog.InfoContext(
		ctx,
		"account created in sqlite",
		slog.String("id", id.String()),
	)

	return id, err
}
