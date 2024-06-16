ALTER TABLE public.legal_entities
ALTER COLUMN phone DROP DEFAULT,
ALTER COLUMN fax DROP DEFAULT;

ALTER TABLE public.legal_entities
RENAME COLUMN mail_postal_code TO mail_index;

ALTER TABLE public.legal_entities
ALTER COLUMN phone TYPE varchar(20) USING phone::varchar,
ALTER COLUMN fax TYPE varchar(20) USING fax::varchar;

ALTER TABLE public.legal_entities
ALTER COLUMN phone SET DEFAULT '',
ALTER COLUMN fax SET DEFAULT '';
