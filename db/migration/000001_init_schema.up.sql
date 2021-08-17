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

CREATE TYPE "payment_request_status" AS ENUM (
  'WAITING_APPROVAL',
  'APPROVED',
  'REFUSED',
  'PAYMENT_SUCCESS',
  'PAYMENT_FAILED'
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
    "from_wallet_id" bigint    NOT NULL,
    "to_wallet_id"   bigint    NOT NULL,
    "amount"         bigint    NOT NULL,
    "created_at"     timestamp NOT NULL DEFAULT 'now()'
);

CREATE TABLE "payment_requests"
(
    "id"             bigserial PRIMARY KEY,
    "from_wallet_id" bigint                 NOT NULL,
    "to_wallet_id"   bigint                 NOT NULL,
    "amount"         bigint                 NOT NULL,
    "status"         payment_request_status NOT NULL,
    "created_at"     timestamp              NOT NULL DEFAULT 'now()'
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

ALTER TABLE "payment_requests"
    ADD FOREIGN KEY ("from_wallet_id") REFERENCES "wallets" ("id");

ALTER TABLE "payment_requests"
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

CREATE INDEX ON "payment_requests" ("from_wallet_id");

CREATE INDEX ON "payment_requests" ("to_wallet_id");

INSERT INTO currencies (code, fraction)
values ('INR', 2);

INSERT INTO users (username, hashed_password, status, full_name, email)
VALUES ('grabINRUser', '$2a$10$N.kqx5ktzyNSQcmc.XjUTOWGwiJuhpjrh9KbO0kM9U61q0tiw7aBy', 'ACTIVE',
        'My Wallet Main Acct INR', 'mywalletinr@gmail.com');

INSERT INTO bank_accounts (account_no, ifsc, bank_name, currency, user_id, status)
VALUES ('1234567890', 'HDFC000076', 'HDFC BANK', 'INR', '1', 'VERIFIED');

INSERT INTO wallets (address, status, user_id, bank_account_id, organization_wallet_id, balance, currency)
VALUES ('grabinr@my.wallet', 'ACTIVE', 1, 1, 1, 100000, 'INR');

INSERT INTO currencies (code, fraction)
values ('USD', 2);

INSERT INTO users (username, hashed_password, status, full_name, email)
VALUES ('grabUSDUser', '$2a$10$N.kqx5ktzyNSQcmc.XjUTOWGwiJuhpjrh9KbO0kM9U61q0tiw7aBy', 'ACTIVE',
        'My Wallet Main Acct USD', 'mywalletusd@gmail.com');

INSERT INTO bank_accounts (account_no, ifsc, bank_name, currency, user_id, status)
VALUES ('1234567891', 'HDFC000076', 'HDFC BANK', 'USD', '2', 'VERIFIED');

INSERT INTO wallets (address, status, user_id, bank_account_id, organization_wallet_id, balance, currency)
VALUES ('grabusd@my.wallet', 'ACTIVE', 2, 2, 2, 100000, 'USD');

INSERT INTO currencies (code, fraction)
values ('EUR', 2);

INSERT INTO users (username, hashed_password, status, full_name, email)
VALUES ('grabEURUser', '$2a$10$N.kqx5ktzyNSQcmc.XjUTOWGwiJuhpjrh9KbO0kM9U61q0tiw7aBy', 'ACTIVE',
        'My Wallet Main Acct EUR', 'mywalleteur@gmail.com');

INSERT INTO bank_accounts (account_no, ifsc, bank_name, currency, user_id, status)
VALUES ('1234567892', 'HDFC000076', 'HDFC BANK', 'EUR', '3', 'VERIFIED');

INSERT INTO wallets (address, status, user_id, bank_account_id, organization_wallet_id, balance, currency)
VALUES ('grabeur@my.wallet', 'ACTIVE', 3, 3, 3, 100000, 'EUR');