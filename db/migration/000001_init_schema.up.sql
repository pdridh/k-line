CREATE TABLE "menu_items" (
  "id" serial PRIMARY KEY,
  "name" text UNIQUE NOT NULL,
  "description" text,
  "price" float NOT NULL,
  "requires_ticket" bool NOT NULL,
  "created_at" timestamp DEFAULT (now())
);

CREATE TYPE user_type AS ENUM (
  'waiter',
  'kitchen',
  'admin'
);

CREATE TABLE "users" (
  "id" uuid PRIMARY KEY DEFAULT (gen_random_uuid()),
  "email" text UNIQUE NOT NULL,
  "name" text NOT NULL,
  "type" user_type NOT NULL,
  "password" text NOT NULL,
  "created_at" timestamp DEFAULT (now())
);

CREATE INDEX ON "menu_items" ("id");

CREATE INDEX ON "menu_items" ("name");

CREATE INDEX ON "users" ("id");

CREATE INDEX ON "users" ("email");


