CREATE TYPE "color" AS ENUM (
  'red',
  'yellow',
  'blue'
);

CREATE TABLE "user" (
  "id" bigserial PRIMARY KEY,
  "username" varchar(255) UNIQUE NOT NULL,
  "name" varchar(64) NOT NULL,
  "hashed_password" varchar(255) NOT NULL,
  "password_changed_at" timestamp  NOT NULL DEFAULT (now()),
  "deleted_at" timestamp,
  "created_at" timestamp NOT NULL DEFAULT (now())
);

CREATE TABLE "card" (
  "id" bigserial PRIMARY KEY,
  "author" bigint NOT NULL,
  "accused" bigint NOT NULL,
  "color" color NOT NULL,
  "description" varchar(255),
  "event" bigint,
  "created_at" timestamp NOT NULL DEFAULT (now())
);

CREATE TABLE "comment" (
  "id" bigserial PRIMARY KEY,
  "message" text NOT NULL,
  "author" bigint NOT NULL,
  "card" bigint NOT NULL,
  "deleted_at" timestamp,
  "created_at" timestamp NOT NULL DEFAULT (now())
);

CREATE TABLE "event" (
  "id" bigserial PRIMARY KEY,
  "name" varchar(255) NOT NULL,
  "date" date NOT NULL,
  "deleted_at" timestamp,
  "created_at" timestamp NOT NULL DEFAULT (now())
);

CREATE INDEX ON "user" ("username");

CREATE INDEX ON "card" ("accused");

CREATE INDEX ON "comment" ("card");

ALTER TABLE "card" ADD FOREIGN KEY ("author") REFERENCES "user" ("id");

ALTER TABLE "card" ADD FOREIGN KEY ("accused") REFERENCES "user" ("id");

ALTER TABLE "card" ADD FOREIGN KEY ("event") REFERENCES "event" ("id");

ALTER TABLE "comment" ADD FOREIGN KEY ("author") REFERENCES "user" ("id");

ALTER TABLE "comment" ADD FOREIGN KEY ("card") REFERENCES "card" ("id");

CREATE TABLE "event_user" (
  "event_id" bigserial,
  "user_id" bigserial,
  PRIMARY KEY ("event_id", "user_id")
);

ALTER TABLE "event_user" ADD FOREIGN KEY ("event_id") REFERENCES "event" ("id");

ALTER TABLE "event_user" ADD FOREIGN KEY ("user_id") REFERENCES "user" ("id");

