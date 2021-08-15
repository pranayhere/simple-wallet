CREATE TYPE "user_status" AS ENUM (
  'ACTIVE',
  'BLOCKED'
);

CREATE TYPE "wallet_status" AS ENUM (
  'ACTIVE',
  'INACTIVE'
);

CREATE TYPE "bank_account_status" AS ENUM (
  'IN_VERIFICATION',
  'VERIFIED',
  'VERIFICATION_FAILED'
);

CREATE TABLE "users"
(
    "id"                  bigserial PRIMARY KEY,
    "username"            varchar     NOT NULL,
    "hashed_password"     varchar     NOT NULL,
    "status"              user_status NOT NULL,
    "full_name"           varchar     NOT NULL,
    "email"               varchar     NOT NULL,
    "password_changed_at" timestamp   NOT NULL DEFAULT 'now()',
    "created_at"          timestamp   NOT NULL DEFAULT 'now()',
    "updated_at"          timestamp   NOT NULL DEFAULT 'now()'
);

CREATE TABLE "wallets"
(
    "id"                     bigserial PRIMARY KEY,
    "address"                varchar       NOT NULL,
    "status"                 wallet_status NOT NULL,
    "user_id"                bigint        NOT NULL,
    "bank_account_id"        bigint        NOT NULL,
    "organization_wallet_id" bigint        NOT NULL,
    "balance"                bigint        NOT NULL,
    "currency"               varchar       NOT NULL,
    "created_at"             timestamp     NOT NULL DEFAULT 'now()',
    "updated_at"             timestamp     NOT NULL DEFAULT 'now()'
);

CREATE TABLE "bank_accounts"
(
    "id"         bigserial PRIMARY KEY,
    "account_no" varchar             NOT NULL,
    "ifsc"       varchar             NOT NULL,
    "bank_name"  varchar             NOT NULL,
    "status"     bank_account_status NOT NULL,
    "user_id"    bigint              NOT NULL,
    "currency"   varchar             NOT NULL,
    "created_at" timestamp           NOT NULL DEFAULT 'now()',
    "updated_at" timestamp           NOT NULL DEFAULT 'now()'
);

CREATE TABLE "currencies"
(
    "code"       varchar PRIMARY KEY,
    "fraction"   bigint    NOT NULL,
    "created_at" timestamp NOT NULL DEFAULT 'now()'
);

CREATE TABLE "entries"
(
    "id"          bigserial PRIMARY KEY,
    "wallet_id"   bigint    NOT NULL,
    "amount"      bigint    NOT NULL,
    "transfer_id" bigint    NOT NULL,
    "created_at"  timestamp NOT NULL DEFAULT 'now()'
);

CREATE TABLE "transfers"
(
    "id"             bigserial PRIMARY KEY,
    "from_wallet_id" bigint        NOT NULL,
    "to_wallet_id"   bigint        NOT NULL,
    "amount"         bigint        NOT NULL,
    "created_at"     timestamp     NOT NULL DEFAULT 'now()'
);

ALTER TABLE "wallets"
    ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

ALTER TABLE "wallets"
    ADD FOREIGN KEY ("bank_account_id") REFERENCES "bank_accounts" ("id");

ALTER TABLE "wallets"
    ADD FOREIGN KEY ("organization_wallet_id") REFERENCES "wallets" ("id");

ALTER TABLE "wallets"
    ADD FOREIGN KEY ("currency") REFERENCES "currencies" ("code");

ALTER TABLE "bank_accounts"
    ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

ALTER TABLE "bank_accounts"
    ADD FOREIGN KEY ("currency") REFERENCES "currencies" ("code");

ALTER TABLE "entries"
    ADD FOREIGN KEY ("wallet_id") REFERENCES "wallets" ("id");

ALTER TABLE "transfers"
    ADD FOREIGN KEY ("from_wallet_id") REFERENCES "wallets" ("id");

ALTER TABLE "transfers"
    ADD FOREIGN KEY ("to_wallet_id") REFERENCES "wallets" ("id");

CREATE UNIQUE INDEX ON "users" ("username");

CREATE UNIQUE INDEX ON "users" ("email");

CREATE INDEX ON "wallets" ("user_id");

CREATE UNIQUE INDEX ON "wallets" ("address");

CREATE UNIQUE INDEX ON "wallets" ("bank_account_id");

CREATE INDEX ON "bank_accounts" ("user_id");

CREATE UNIQUE INDEX ON "bank_accounts" ("account_no", "ifsc");

CREATE INDEX ON "entries" ("wallet_id");

CREATE INDEX ON "transfers" ("from_wallet_id");

CREATE INDEX ON "transfers" ("to_wallet_id");

CREATE INDEX ON "transfers" ("from_wallet_id", "to_wallet_id");

INSERT INTO currencies (code, fraction)
values ('INR', 2);

INSERT INTO users (username, hashed_password, status, full_name, email)
VALUES ('myWalletINRUser', '$2a$10$VQPlcZxroJ1QF3nI8M7XsedQfOBlg.BIh4M70P3cECrVpE7jVxpca', 'ACTIVE',
        'My Wallet Main Acct INR', 'mywalletinr@gmail.com');

INSERT INTO bank_accounts (account_no, ifsc, bank_name, currency, user_id, status)
VALUES ('1234567890', 'HDFC000076', 'HDFC BANK', 'INR', '1', 'VERIFIED');

INSERT INTO wallets (address, status, user_id, bank_account_id, organization_wallet_id, balance, currency)
VALUES ('mywalletinr@my.wallet', 'ACTIVE', 1, 1, 1, 0, 'INR');

INSERT INTO currencies (code, fraction)
values ('USD', 2);

INSERT INTO users (username, hashed_password, status, full_name, email)
VALUES ('myWalletUSDUser', '$2a$10$VQPlcZxroJ1QF3nI8M7XsedQfOBlg.BIh4M70P3cECrVpE7jVxpca', 'ACTIVE',
        'My Wallet Main Acct USD', 'mywalletusd@gmail.com');

INSERT INTO bank_accounts (account_no, ifsc, bank_name, currency, user_id, status)
VALUES ('1234567891', 'HDFC000076', 'HDFC BANK', 'USD', '2', 'VERIFIED');

INSERT INTO wallets (address, status, user_id, bank_account_id, organization_wallet_id, balance, currency)
VALUES ('mywalletusd@my.wallet', 'ACTIVE', 2, 2, 2, 0, 'USD');

INSERT INTO currencies (code, fraction)
values ('EUR', 2);

INSERT INTO users (username, hashed_password, status, full_name, email)
VALUES ('myWalletEURUser', '$2a$10$VQPlcZxroJ1QF3nI8M7XsedQfOBlg.BIh4M70P3cECrVpE7jVxpca', 'ACTIVE',
        'My Wallet Main Acct EUR', 'mywalleteur@gmail.com');

INSERT INTO bank_accounts (account_no, ifsc, bank_name, currency, user_id, status)
VALUES ('1234567892', 'HDFC000076', 'HDFC BANK', 'EUR', '3', 'VERIFIED');

INSERT INTO wallets (address, status, user_id, bank_account_id, organization_wallet_id, balance, currency)
VALUES ('mywalleteur@my.wallet', 'ACTIVE', 3, 3, 3, 0, 'EUR');