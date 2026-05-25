CREATE TABLE wallet (
  id         SERIAL PRIMARY KEY,
  user_id    INT NOT NULL UNIQUE,
  balance    DECIMAL(15,2) NOT NULL DEFAULT 0,
  updated_at TIMESTAMP NOT NULL DEFAULT NOW(),

  FOREIGN KEY (user_id) REFERENCES users(id)
);