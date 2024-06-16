CREATE SCHEMA IF NOT EXISTS integrations;

CREATE TABLE integrations.sent_emails (
    "sender_email_uuid" uuid NOT NULL DEFAULT gen_random_uuid() PRIMARY KEY,
    "sender_email" varchar(255) NOT NULL,
    "recipient_email" varchar(255) NOT NULL,
    "subject" varchar(255),
    "body" text,
    "created_at" timestamptz NOT NULL DEFAULT now()
);
