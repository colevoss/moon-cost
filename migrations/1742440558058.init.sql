CREATE TABLE IF NOT EXISTS organizations (
  id INTEGER PRIMARY KEY,
  name TEXT NOT NULL,
  createdAt INTEGER NOT NULL,
  updatedAt INTEGER NOT NULL
);

CREATE TABLE IF NOT EXISTS locations (
  id INTEGETER PRIMARY KEY,

  address TEXT,
  city TEXT,
  state TEXT,

  organizationId INTEGER NOT NULL,
  FOREIGN KEY(organizationId) REFERENCES organizations(id)
);

CREATE TABLE IF NOT EXISTS users (
  id INTEGER PRIMARY KEY,

  firstname TEXT NOT NULL,
  lastname TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS accounts (
  id INTEGER PRIMARY KEY,

  email TEXT NOT NULL UNIQUE,
  password TEXT NOT NULL,
  salt TEXT NOT NULL,

  active INTEGER CHECK (active IN (0, 1)),

  userId INTEGER NOT NULL,
  FOREIGN KEY(userId) REFERENCES users(id)
);
