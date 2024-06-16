ALTER TABLE public.legal_entities
ALTER COLUMN phone DROP DEFAULT,
ALTER COLUMN fax DROP DEFAULT;

ALTER TABLE public.legal_entities
ALTER COLUMN phone TYPE bigint USING phone::bigint,
ALTER COLUMN fax TYPE bigint USING fax::bigint;

ALTER TABLE public.legal_entities
RENAME COLUMN mail_index TO mail_postal_code;

ALTER TABLE public.legal_entities
ALTER COLUMN phone SET DEFAULT 0,
ALTER COLUMN fax SET DEFAULT 0;
