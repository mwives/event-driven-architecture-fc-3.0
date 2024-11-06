CREATE DATABASE IF NOT EXISTS wallet;

USE wallet;

CREATE TABLE
  clients (
    id VARCHAR(255) PRIMARY KEY,
    name VARCHAR(255),
    email VARCHAR(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
  );

CREATE TABLE
  accounts (
    id VARCHAR(255) PRIMARY KEY,
    client_id VARCHAR(255),
    balance FLOAT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
  );

CREATE TABLE
  transactions (
    id VARCHAR(255) PRIMARY KEY,
    account_id_from VARCHAR(255),
    account_id_to VARCHAR(255),
    amount FLOAT,
    created_at TIMESTAMP
  );

INSERT INTO
  clients (id, name, email)
VALUES
  (
    'ac73550c-5cd0-407e-a4e4-30d4c5cdf297',
    'John Doe',
    'john@doe.com'
  );

INSERT INTO
  clients (id, name, email)
VALUES
  (
    '9257ceb9-a731-4ad1-859c-0308b7966fda',
    'Jane Doe',
    'jane@doe.com'
  );

INSERT INTO
  accounts (id, client_id, balance)
VALUES
  (
    '14c4e358-06d4-4031-b7e0-fe56662c6e97',
    'ac73550c-5cd0-407e-a4e4-30d4c5cdf297',
    1000.0
  );

INSERT INTO
  accounts (id, client_id, balance)
VALUES
  (
    'd9b93621-3b57-4755-b0e5-0bce147c7c3a',
    '9257ceb9-a731-4ad1-859c-0308b7966fda',
    1000.0
  );