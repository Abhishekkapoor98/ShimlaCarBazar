CREATE TABLE "users" (
  "id" bigserial PRIMARY KEY,
  "email" varchar UNIQUE NOT NULL,
  "password_hash" varchar NOT NULL,
  "role" varchar NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "cars" (
  "id" bigserial PRIMARY KEY,
  "title" varchar NOT NULL,
  "brand" varchar NOT NULL,
  "model" varchar NOT NULL,
  "condition_type" varchar NOT NULL,
  "price" bigint NOT NULL,
  "status" varchar NOT NULL,
  "image_urls" text NOT NULL DEFAULT '',
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "sell_proposals" (
  "id" bigserial PRIMARY KEY,
  "customer_id" bigint NOT NULL,
  "car_brand" varchar NOT NULL,
  "car_model" varchar NOT NULL,
  "expected_price" bigint NOT NULL,
  "phone_number" varchar NOT NULL DEFAULT '',
  "image_urls" text NOT NULL DEFAULT '',
  "status" varchar NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "rc_transfers" (
  "id" bigserial PRIMARY KEY,
  "car_id" bigint NOT NULL,
  "buyer_id" bigint NOT NULL,
  "current_status" varchar NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

ALTER TABLE "sell_proposals" ADD FOREIGN KEY ("customer_id") REFERENCES "users" ("id");
ALTER TABLE "rc_transfers" ADD FOREIGN KEY ("car_id") REFERENCES "cars" ("id");
ALTER TABLE "rc_transfers" ADD FOREIGN KEY ("buyer_id") REFERENCES "users" ("id");