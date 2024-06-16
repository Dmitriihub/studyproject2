ALTER TABLE integrations.user_emails
ADD COLUMN created_by varchar(255) NOT NULL,
ADD COLUMN created_by_uuid uuid NOT NULL,
ADD COLUMN settings json NOT NULL DEFAULT '{}'::jsonb;

ALTER TABLE integrations.user_emails
DROP COLUMN IF EXISTS domain;
