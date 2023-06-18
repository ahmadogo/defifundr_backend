-- SQL dump generated using DBML (dbml-lang.org)

-- Database: PostgreSQL

-- Generated at: 2023-06-11T09:41:45.519Z

CREATE TABLE
    "users" (
        "username" varchar PRIMARY KEY,
        "hashed_password" varchar NOT NULL,
        "avatar" varchar NOT NULL,
        "email" varchar UNIQUE NOT NULL,
        "is_email_verified" bool NOT NULL DEFAULT false,
        "password_changed_at" timestamptz NOT NULL DEFAULT '0001-01-01',
        "balance" varchar NOT NULL,
        "isBiomatric" bool NOT NULL DEFAULT false,
        "address" varchar NOT NULL,
        "file_path" varchar NOT NULL,
        "secret_code" varchar NOT NULL,
        "is_used" bool NOT NULL DEFAULT false,
        "created_at" timestamptz NOT NULL DEFAULT (now()),
        "expired_at" timestamptz NOT NULL DEFAULT (now() + interval '15 minutes')
    );

CREATE TABLE
    "campaigns" (
        "id" bigserial PRIMARY KEY,
        "image" varchar NOT NULL,
        "campaign_name" varchar NOT NULL
    );

CREATE TABLE
    "donations" (
        "id" bigserial PRIMARY KEY,
        "owner" varchar NOT NULL,
        "title" varchar NOT NULL,
        "campaign_type" varchar NOT NULL,
        "description" varchar NOT NULL,
        "goal" bigint NOT NULL,
        "deadline" timestamptz NOT NULL,
        "image" varchar NOT NULL,
        "created_at" timestamptz NOT NULL DEFAULT (now())
    );

CREATE TABLE
    "user_session" (
        "id" uuid PRIMARY KEY,
        "username" varchar NOT NULL,
        "refresh_token" varchar NOT NULL,
        "user_agent" varchar NOT NULL,
        "client_ip" varchar NOT NULL,
        "is_blocked" boolean NOT NULL DEFAULT false,
        "expires_at" timestamptz NOT NULL,
        "created_at" timestamptz NOT NULL DEFAULT (now())
    );

ALTER TABLE "donations"
ADD
    FOREIGN KEY ("owner") REFERENCES "users" ("username");