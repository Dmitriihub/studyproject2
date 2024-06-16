ALTER TABLE public.legal_entities
DROP COLUMN IF EXISTS first_name,
DROP COLUMN IF EXISTS last_name,
DROP COLUMN IF EXISTS middle_name,
DROP COLUMN IF EXISTS document_type,
DROP COLUMN IF EXISTS document_series,
DROP COLUMN IF EXISTS document_number,
DROP COLUMN IF EXISTS document_issuer,
DROP COLUMN IF EXISTS document_issue_date,
DROP COLUMN IF EXISTS document_department_code,
DROP COLUMN IF EXISTS ogrnip,
DROP COLUMN IF EXISTS certificate_number,
DROP COLUMN IF EXISTS certificate_date;