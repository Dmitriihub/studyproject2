CREATE TABLE bank_accounts (
    uuid UUID PRIMARY KEY,
    legal_entity_id UUID NOT NULL,
    bic VARCHAR(20),
    bank_name VARCHAR(255),
    bank_address VARCHAR(255),
    correspondent_account VARCHAR(50),
    settlement_account VARCHAR(50) NOT NULL,
    currency VARCHAR(20),
    comment TEXT,
    is_primary BOOLEAN DEFAULT false,
    created_at TIMESTAMP WITHOUT TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITHOUT TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_legal_entity FOREIGN KEY(legal_entity_id) REFERENCES legal_entities(uuid) ON DELETE CASCADE
);