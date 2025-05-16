CREATE TABLE IF NOT EXISTS organizations (
  id TEXT PRIMARY KEY,
  name TEXT NOT NULL,
  createdAt INTEGER NOT NULL,
  updatedAt INTEGER NOT NULL
);

CREATE TABLE IF NOT EXISTS locations (
  id TEXT PRIMARY KEY,

  address TEXT,
  city TEXT,
  state TEXT,

  organizationId TEXT NOT NULL,
  FOREIGN KEY(organizationId) REFERENCES organizations(id)
);

CREATE TABLE IF NOT EXISTS users (
  id TEXT PRIMARY KEY,

  firstname TEXT NOT NULL,
  lastname TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS accounts (
  id TEXT PRIMARY KEY,

  email TEXT NOT NULL UNIQUE,
  password TEXT NOT NULL,
  salt TEXT NOT NULL,

  active INTEGER CHECK (active IN (0, 1)),

  userId TEXT NOT NULL,
  FOREIGN KEY(userId) REFERENCES users(id)
);
