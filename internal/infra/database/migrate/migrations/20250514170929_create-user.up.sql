CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE IF NOT EXISTS "users" (
  "id"             VARCHAR(255) PRIMARY KEY,
  "name"           VARCHAR(255) NOT NULL,
  "email"          VARCHAR(255) UNIQUE NOT NULL,
  "password"       VARCHAR(255) NOT NULL,
  "avatar_url"     TEXT,
  "user_role_id"   VARCHAR(255) REFERENCES user_roles(id) ON DELETE SET NULL,
  "subcategory_id" VARCHAR(255) REFERENCES subcategories(id) ON DELETE SET NULL,
  "created_at"     TIMESTAMP     NOT NULL DEFAULT NOW(),
  "updated_at"     TIMESTAMP,
  "deleted_at"     TIMESTAMP
);
