CREATE TABLE IF NOT EXISTS "user" (
    "id" INTEGER GENERATED BY DEFAULT AS IDENTITY UNIQUE PRIMARY KEY NOT NULL,
    "username" varchar(255) UNIQUE NOT NULL DEFAULT 'noname',
    "password_hash" varchar(255) NOT NULL,
    "created_at" timestamp NOT NULL DEFAULT (now()),
    "is_deleted" boolean NOT NULL DEFAULT false
);

CREATE TABLE IF NOT EXISTS "chat" (
    "id" INTEGER GENERATED BY DEFAULT AS IDENTITY UNIQUE PRIMARY KEY NOT NULL,
    "name" varchar(255) UNIQUE NOT NULL,
    "created_at" timestamp NOT NULL DEFAULT (now()),
    "is_deleted" boolean NOT NULL DEFAULT false
);

CREATE TABLE IF NOT EXISTS "message" (
    "id" INTEGER GENERATED BY DEFAULT AS IDENTITY UNIQUE PRIMARY KEY NOT NULL,
    "text" varchar(255) NOT NULL,
    "created_at" timestamp NOT NULL DEFAULT (now()),
    "is_deleted" boolean NOT NULL DEFAULT false
);

CREATE TABLE IF NOT EXISTS "users_chat" (
    "id" INTEGER GENERATED BY DEFAULT AS IDENTITY UNIQUE PRIMARY KEY NOT NULL,
    "user_id" integer NOT NULL,
    "chat_id" integer NOT NULL
);

CREATE TABLE IF NOT EXISTS "chats_messages" (
    "id" INTEGER GENERATED BY DEFAULT AS IDENTITY UNIQUE PRIMARY KEY NOT NULL,
    "users_chat_id" integer NOT NULL,
    "message_id" integer NOT NULL
);


ALTER TABLE "users_chat" ADD FOREIGN KEY ("user_id") REFERENCES "user" ("id") ON DELETE NO ACTION;

ALTER TABLE "users_chat" ADD FOREIGN KEY ("chat_id") REFERENCES "chat" ("id") ON DELETE CASCADE;

ALTER TABLE "chats_messages" ADD FOREIGN KEY ("users_chat_id") REFERENCES "users_chat" ("id") ON DELETE NO ACTION;

ALTER TABLE "chats_messages" ADD FOREIGN KEY ("message_id") REFERENCES "message" ("id") ON DELETE CASCADE;
