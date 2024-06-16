ALTER TABLE integrations.sent_emails
DROP COLUMN sender_email_uuid;

ALTER TABLE integrations.sent_emails
ADD COLUMN uuid uuid NOT NULL DEFAULT gen_random_uuid() PRIMARY KEY,
ADD COLUMN sender_email_uuid uuid NOT NULL;