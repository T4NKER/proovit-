CREATE TABLE transactions (
  transactionID TEXT PRIMARY KEY,

  amount REAL NOT NULL,

  spent INTEGER CHECK(spent IN (0, 1)) NOT NULL,

  createdAt DATETIME NOT NULL
);