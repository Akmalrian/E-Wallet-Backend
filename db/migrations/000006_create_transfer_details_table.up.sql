CREATE TABLE IF NOT EXISTS transfer_details (
  id                 SERIAL PRIMARY KEY,
  transaction_id     INT NOT NULL UNIQUE,
  sender_wallet_id   INT NOT NULL,
  receiver_wallet_id INT NOT NULL,
  amount             DECIMAL(15,2) NOT NULL,
  notes              TEXT,
  created_at         TIMESTAMP NOT NULL DEFAULT NOW(),

  FOREIGN KEY (transaction_id)     REFERENCES transactions(id),
  FOREIGN KEY (sender_wallet_id)   REFERENCES wallet(id),
  FOREIGN KEY (receiver_wallet_id) REFERENCES wallet(id)
);