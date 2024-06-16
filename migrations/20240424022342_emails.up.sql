CREATE TABLE emails (
    "uuid" UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    "user_uuid" uuid NOT NULL REFERENCES user_emails (uuid) ON DELETE CASCADE,
    "email" varchar(255) NOT NULL,
    "sender" varchar(255),
    "subject" varchar(255),
    "body_text" text,
    "received_at" timestamptz,
    "body_hash" varchar(255) UNIQUE NOT NULL,
    "body" bytea,
    "created_at" timestamptz NOT NULL DEFAULT now(),
    "updated_at" timestamptz NOT NULL DEFAULT now(),
    "deleted_at" timestamptz
);

CREATE INDEX emails_body_hash ON emails(body_hash);