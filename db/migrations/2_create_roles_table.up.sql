CREATE TABLE IF NOT EXISTS roles(
  id SERIAL PRIMARY KEY,
  name VARCHAR(50) UNIQUE NOT NULL
);

CREATE TABLE IF NOT EXISTS user_roles(
  id SERIAL PRIMARY KEY,
  role_id INT REFERENCES roles ON DELETE CASCADE,
  user_id INT REFERENCES users ON DELETE CASCADE
);

CREATE UNIQUE INDEX index_user_roles_on_role_id_and_user_id ON user_roles (role_id, user_id);
CREATE INDEX index_user_roles_on_user_id ON user_roles (user_id);
