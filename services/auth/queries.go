package auth

import "moon-cost/moon"

var AuthQueries = moon.Queries{
	"createUser": {
		SQLite: "INSERT INTO users (id, firstname, lastname) VALUES (?, ?, ?)",
	},

	"createAccount": {
		SQLite: `
    INSERT INTO accounts (id, email, password, salt, active, userId)
    VALUES (?, ?, ?, ?, ?, ?)
    `,
	},

	"getUser": {
		SQLite: `SELECT id, firstname, lastname FROM users WHERE id = ?`,
	},

	"findAccountByEmail": {
		SQLite: `
    SELECT
      id,
      email,
      password,
      salt,
      active,
      userId
    FROM accounts WHERE email = ?
    `,
	},
}
