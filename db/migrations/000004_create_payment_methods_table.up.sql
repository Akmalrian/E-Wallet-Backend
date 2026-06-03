CREATE TABLE IF NOT EXISTS payment_methods (
  id           SERIAL PRIMARY KEY,
  payment_name VARCHAR(50) NOT NULL UNIQUE
);