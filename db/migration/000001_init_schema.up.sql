CREATE TYPE "user_status" AS ENUM (
  'ACTIVE',
  'BLOCKED'
);

CREATE TYPE "wallet_status" AS ENUM (
  'ACTIVE',
  'INACTIVE',
  'BLOCKED'
);

CREATE TYPE "bank_account_status" AS ENUM (
  'IN_VERIFICATION',
  'VERIFIED',
  'VERIFICATION_FAILED'
);

CREATE TABLE "users" (
    "id" bigserial PRIMARY KEY,
    "username" varchar NOT NULL,
    "hashed_password" varchar NOT NULL,
    "status" user_status NOT NULL,
    "full_name" varchar NOT NULL,
    "email" varchar NOT NULL,
    "password_changed_at" timestamp NOT NULL DEFAULT 'now()',
    "created_at" timestamp NOT NULL DEFAULT 'now()',
    "updated_at" timestamp NOT NULL DEFAULT 'now()'
);

CREATE TABLE "wallets" (
    "id" bigserial PRIMARY KEY,
    "name" varchar NOT NULL,
    "address" varchar NOT NULL,
    "status" wallet_status NOT NULL,
    "user_id" bigint NOT NULL,
    "bank_account_id" bigint NOT NULL,
    "balance" bigint NOT NULL,
    "currency" varchar NOT NULL,
    "created_at" timestamp NOT NULL DEFAULT 'now()',
    "updated_at" timestamp NOT NULL DEFAULT 'now()'
);

CREATE TABLE "bank_accounts" (
    "id" bigserial PRIMARY KEY,
    "account_no" varchar NOT NULL,
    "ifsc" varchar NOT NULL,
    "bank_name" varchar NOT NULL,
    "status" bank_account_status NOT NULL,
    "user_id" bigint NOT NULL,
    "currency" varchar NOT NULL,
    "created_at" timestamp NOT NULL DEFAULT 'now()',
    "updated_at" timestamp NOT NULL DEFAULT 'now()'
);

CREATE TABLE "currencies" (
    "code" varchar PRIMARY KEY,
    "fraction" bigint NOT NULL,
    "created_at" timestamp NOT NULL DEFAULT 'now()'
);

CREATE TABLE "entries" (
    "id" bigserial PRIMARY KEY,
    "wallet_id" bigint NOT NULL,
    "amount" bigint NOT NULL,
    "created_at" timestamp NOT NULL DEFAULT 'now()'
);

CREATE TABLE "transfers" (
    "id" bigserial PRIMARY KEY,
    "from_wallet_id" bigint NOT NULL,
    "to_wallet_id" bigint NOT NULL,
    "amount" bigint NOT NULL,
    "created_at" timestamp NOT NULL DEFAULT 'now()'
);

ALTER TABLE "wallets" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

ALTER TABLE "wallets" ADD FOREIGN KEY ("bank_account_id") REFERENCES "bank_accounts" ("id");

ALTER TABLE "wallets" ADD FOREIGN KEY ("currency") REFERENCES "currencies" ("code");

ALTER TABLE "bank_accounts" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

ALTER TABLE "bank_accounts" ADD FOREIGN KEY ("currency") REFERENCES "currencies" ("code");

ALTER TABLE "entries" ADD FOREIGN KEY ("wallet_id") REFERENCES "wallets" ("id");

ALTER TABLE "transfers" ADD FOREIGN KEY ("from_wallet_id") REFERENCES "wallets" ("id");

ALTER TABLE "transfers" ADD FOREIGN KEY ("to_wallet_id") REFERENCES "wallets" ("id");

CREATE INDEX ON "users" ("username");

CREATE INDEX ON "users" ("email");

CREATE INDEX ON "wallets" ("user_id");

CREATE UNIQUE INDEX ON "wallets" ("address");

CREATE UNIQUE INDEX ON "wallets" ("bank_account_id");

CREATE INDEX ON "bank_accounts" ("user_id");

CREATE UNIQUE INDEX ON "bank_accounts" ("account_no", "ifsc");

CREATE INDEX ON "entries" ("wallet_id");

CREATE INDEX ON "transfers" ("from_wallet_id");

CREATE INDEX ON "transfers" ("to_wallet_id");

CREATE INDEX ON "transfers" ("from_wallet_id", "to_wallet_id");
