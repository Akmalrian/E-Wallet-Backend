CREATE TABLE topup_details (
  id                SERIAL PRIMARY KEY,
  transaction_id    INT NOT NULL UNIQUE,
  wallet_id         INT NOT NULL,
  payment_method_id INT NOT NULL,
  order_amount      DECIMAL(15,2) NOT NULL,
  tax_amount        DECIMAL(15,2) NOT NULL DEFAULT 0,
  delivery_fee      DECIMAL(15,2) NOT NULL DEFAULT 0,
  total_amount      DECIMAL(15,2) NOT NULL,
  created_at        TIMESTAMP NOT NULL DEFAULT NOW(),

  FOREIGN KEY (transaction_id)    REFERENCES transactions(id),
  FOREIGN KEY (wallet_id)         REFERENCES wallet(id),
  FOREIGN KEY (payment_method_id) REFERENCES payment_methods(id)
);