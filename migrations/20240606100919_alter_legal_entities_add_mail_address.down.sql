ALTER TABLE public.legal_entities
DROP COLUMN IF EXISTS mail_index,
DROP COLUMN IF EXISTS mail_country,
DROP COLUMN IF EXISTS mail_region,
DROP COLUMN IF EXISTS mail_city,
DROP COLUMN IF EXISTS mail_street,
DROP COLUMN IF EXISTS mail_house,
DROP COLUMN IF EXISTS mail_apartment,
DROP COLUMN IF EXISTS mail_comments,
DROP COLUMN IF EXISTS mail_address_same_as_legal;
