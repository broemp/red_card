CREATE TABLE "session" (
  "id" uuid PRIMARY KEY,
  "user_id" bigint NOT NULL,
  "refresh_token" varchar NOT NULL,
  "user_agent" varchar NOT NULL,
  "client_ip" varchar NOT NULL,
  "is_blocked" boolean  NOT NULL DEFAULT false,
  "expires_at" timestamp NOT NULL,
  "created_at" timestamp NOT NULL DEFAULT (now())
);

ALTER TABLE "session" ADD FOREIGN KEY ("user_id") REFERENCES "user" ("id");
