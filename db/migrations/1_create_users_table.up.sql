CREATE TABLE IF NOT EXISTS users(
  id SERIAL PRIMARY KEY,
  login VARCHAR(255) UNIQUE NOT NULL,
  password VARCHAR(255) NOT NULL,
  created_at TIMESTAMP DEFAULT (NOW() AT TIME ZONE 'utc')
);
