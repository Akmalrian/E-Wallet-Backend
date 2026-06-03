CREATE TABLE IF NOT EXISTS transactions (
  id         SERIAL PRIMARY KEY,
  user_id    INT NOT NULL,
  type       VARCHAR(20) NOT NULL,
  amount     DECIMAL(15,2) NOT NULL,
  status     VARCHAR(20) NOT NULL DEFAULT 'pending',
  created_at TIMESTAMP NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMP,

  FOREIGN KEY (user_id) REFERENCES users(id)
);