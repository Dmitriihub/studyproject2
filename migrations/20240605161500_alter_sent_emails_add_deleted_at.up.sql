ALTER TABLE
    "integrations"."emails_sent"
ADD
    COLUMN "deleted_at" timestamptz;