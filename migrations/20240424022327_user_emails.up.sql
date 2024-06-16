CREATE TABLE user_emails (
    "uuid" uuid NOT NULL DEFAULT gen_random_uuid() PRIMARY KEY,
    "email" varchar(255) NOT NULL UNIQUE,
    "password" varchar(255) NOT NULL,
    "created_at" timestamptz NOT NULL DEFAULT now(),
    "updated_at" timestamptz NOT NULL DEFAULT now(),
    "deleted_at" timestamptz,
    "last_fetched_at" timestamptz DEFAULT NULL
);
