CREATE TABLE bank_accounts (
    "uuid" uuid NOT NULL DEFAULT gen_random_uuid() PRIMARY KEY,
    "legal_entity_uuid" uuid NOT NULL REFERENCES legal_entities(uuid) ON DELETE CASCADE,
    "bik" VARCHAR(20) NOT NULL,
    "bank_name" VARCHAR(255) NOT NULL,
    "address" VARCHAR(255) NOT NULL,
    "corr_account" VARCHAR(20) NOT NULL,
    "account_number" VARCHAR(20) NOT NULL,
    "currency" VARCHAR(10) NOT NULL,
    "comment" TEXT,
    "is_active" BOOLEAN NOT NULL DEFAULT FALSE,
    "created_at" timestamptz NOT NULL DEFAULT now(),
    "updated_at" timestamptz NOT NULL DEFAULT now(),
    "deleted_at" timestamptz
);

CREATE INDEX "idx_bank_accounts_legal_entity_uuid" ON bank_accounts ("legal_entity_uuid");