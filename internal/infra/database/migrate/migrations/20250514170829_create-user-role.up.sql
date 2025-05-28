CREATE TABLE IF NOT EXISTS user_roles (
  id          VARCHAR(255) PRIMARY KEY DEFAULT new_id('user_role'),
  name        VARCHAR(255) NOT NULL UNIQUE,
  created_at  TIMESTAMP   NOT NULL DEFAULT NOW(),
  updated_at  TIMESTAMP,
  deleted_at  TIMESTAMP
);

INSERT INTO user_roles (name) 
VALUES
  ('client'),
  ('professional')
