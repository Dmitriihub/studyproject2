ALTER TABLE integrations.user_emails SET SCHEMA public;

ALTER TABLE integrations.emails SET SCHEMA public;

DROP SCHEMA integrations CASCADE;