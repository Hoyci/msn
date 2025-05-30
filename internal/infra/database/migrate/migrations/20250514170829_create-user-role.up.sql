CREATE TABLE IF NOT EXISTS roles (
  id          VARCHAR(255) PRIMARY KEY DEFAULT new_id('role'),
  name        VARCHAR(255) NOT NULL UNIQUE,
  created_at  TIMESTAMP   NOT NULL DEFAULT NOW(),
  updated_at  TIMESTAMP,
  deleted_at  TIMESTAMP
);

INSERT INTO roles (name) 
VALUES
  ('client'),
  ('professional')
