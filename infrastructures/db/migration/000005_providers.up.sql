CREATE TABLE providers (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL,
    type VARCHAR(30) NOT NULL,
    license_number VARCHAR(100) NOT NULL UNIQUE,
    status VARCHAR(20) NOT NULL DEFAULT 'PENDING',
    county VARCHAR(100),
    address TEXT,
    phone VARCHAR(20),
    email VARCHAR(255),
    contact_person VARCHAR(255),
    user_id UUID REFERENCES users(id),
    created_by UUID,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_providers_license_number ON providers(license_number);
CREATE INDEX idx_providers_status ON providers(status);
CREATE INDEX idx_providers_type ON providers(type);
CREATE INDEX idx_providers_county ON providers(county);
CREATE INDEX idx_providers_user_id ON providers(user_id);

CREATE TABLE contracts (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    provider_id UUID NOT NULL REFERENCES providers(id) ON DELETE CASCADE,
    start_date TIMESTAMPTZ NOT NULL,
    end_date TIMESTAMPTZ NOT NULL,
    terms TEXT NOT NULL DEFAULT '',
    status VARCHAR(20) NOT NULL DEFAULT 'ACTIVE',
    created_by UUID,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_contracts_provider_id ON contracts(provider_id);
CREATE INDEX idx_contracts_status ON contracts(status);

CREATE TABLE rate_cards (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    provider_id UUID NOT NULL REFERENCES providers(id) ON DELETE CASCADE,
    procedure_code VARCHAR(20) NOT NULL,
    procedure_name VARCHAR(255) NOT NULL,
    rate_amount BIGINT NOT NULL,
    effective_date TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    age_from INT NOT NULL DEFAULT 0,
    age_to INT NOT NULL DEFAULT 150,
    gender VARCHAR(10),
    relationship VARCHAR(20),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_rate_cards_provider_id ON rate_cards(provider_id);
CREATE INDEX idx_rate_cards_procedure_code ON rate_cards(procedure_code);
CREATE INDEX idx_rate_cards_provider_procedure ON rate_cards(provider_id, procedure_code);
