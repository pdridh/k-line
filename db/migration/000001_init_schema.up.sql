CREATE TYPE "table_status" AS ENUM (
  'available',
  'occupied',
  'closed'
);

CREATE TYPE "user_type" AS ENUM (
  'admin',
  'register',
  'kitchen',
  'waiter',
  'driver'
);

CREATE TYPE "order_status" AS ENUM (
  'ongoing',
  'completed',
  'cancelled'
);

CREATE TYPE "order_type" AS ENUM (
  'dining',
  'delivery',
  'takeaway'
);

CREATE TYPE "order_item_status" AS ENUM (
  'pending',
  'preparing',
  'ready',
  'served',
  'delivered',
  'cancelled'
);

CREATE TYPE "delivery_status" AS ENUM (
  'pending',
  'preparing',
  'ready',
  'dispatched',
  'delivered',
  'failed',
  'cancelled'
);

CREATE TYPE "takeaway_status" AS ENUM (
  'pending',
  'preparing',
  'picked_up',
  'cancelled',
  'no_show'
);

CREATE TABLE "menu_items" (
  "id" serial PRIMARY KEY,
  "name" text UNIQUE NOT NULL,
  "description" text,
  "price" float NOT NULL,
  "requires_ticket" bool NOT NULL,
  "created_at" timestamp DEFAULT (now())
);

CREATE TABLE "tables" (
  "id" text PRIMARY KEY,
  "capacity" smallint NOT NULL,
  "status" table_status NOT NULL DEFAULT 'available',
  "notes" text
);

CREATE TABLE "users" (
  "id" uuid PRIMARY KEY DEFAULT (gen_random_uuid()),
  "email" text UNIQUE NOT NULL,
  "name" text NOT NULL,
  "type" user_type NOT NULL,
  "password" text NOT NULL,
  "created_at" timestamp DEFAULT (now())
);

CREATE TABLE "orders" (
  "id" uuid PRIMARY KEY DEFAULT (gen_random_uuid()),
  "type" order_type NOT NULL DEFAULT 'dining',
  "employee_id" uuid NOT NULL,
  "status" order_status NOT NULL DEFAULT 'ongoing',
  "table_id" text,
  "created_at" timestamp DEFAULT (now()),
  "completed_at" timestamp
);

CREATE TABLE "order_items" (
  "id" bigserial PRIMARY KEY,
  "order_id" uuid NOT NULL,
  "item_id" int NOT NULL,
  "quantity" int NOT NULL,
  "notes" text,
  "status" order_item_status NOT NULL DEFAULT 'pending',
  "added_at" timestamp DEFAULT (now())
);

CREATE TABLE "delivery_details" (
  "order_id" uuid PRIMARY KEY,
  "address" text NOT NULL,
  "contact" text NOT NULL,
  "driver_id" uuid NOT NULL,
  "status" delivery_status NOT NULL DEFAULT 'pending',
  "dispatched_at" timestamp,
  "delivered_at" timestamp
);

CREATE TABLE "takeaway_details" (
  "order_id" uuid PRIMARY KEY,
  "contact" text NOT NULL,
  "status" takeaway_status NOT NULL DEFAULT 'pending',
  "notes" text,
  "picked_at" timestamp
);

CREATE INDEX ON "menu_items" ("name");

CREATE INDEX ON "tables" ("status");

CREATE INDEX ON "users" ("email");

CREATE INDEX ON "orders" ("table_id");

CREATE INDEX ON "orders" ("status");

CREATE INDEX ON "order_items" ("order_id");

CREATE INDEX ON "order_items" ("status");

CREATE INDEX ON "delivery_details" ("driver_id");

CREATE INDEX ON "delivery_details" ("status");

CREATE INDEX ON "takeaway_details" ("contact");

CREATE INDEX ON "takeaway_details" ("status");

ALTER TABLE "orders" ADD FOREIGN KEY ("table_id") REFERENCES "tables" ("id") ON DELETE CASCADE;

ALTER TABLE "delivery_details" ADD FOREIGN KEY ("order_id") REFERENCES "orders" ("id") ON DELETE CASCADE;

ALTER TABLE "takeaway_details" ADD FOREIGN KEY ("order_id") REFERENCES "orders" ("id") ON DELETE CASCADE;

ALTER TABLE "orders" ADD FOREIGN KEY ("employee_id") REFERENCES "users" ("id") ON DELETE CASCADE;

ALTER TABLE "order_items" ADD FOREIGN KEY ("order_id") REFERENCES "orders" ("id") ON DELETE CASCADE;

ALTER TABLE "order_items" ADD FOREIGN KEY ("item_id") REFERENCES "menu_items" ("id") ON DELETE CASCADE;

ALTER TABLE "delivery_details" ADD FOREIGN KEY ("driver_id") REFERENCES "users" ("id") ON DELETE CASCADE;
