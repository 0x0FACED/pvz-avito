CREATE SCHEMA IF NOT EXISTS avito;

CREATE TYPE avito.role_enum AS ENUM ('employee', 'moderator');
CREATE TYPE avito.city_enum AS ENUM ('Москва', 'Санкт-Петербург', 'Казань');
CREATE TYPE avito.status_enum AS ENUM ('in_progress', 'close');
CREATE TYPE avito.product_type_enum AS ENUM ('электроника', 'одежда', 'обувь');

CREATE TABLE IF NOT EXISTS avito.users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(320) NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    role avito.role_enum NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS avito.pvz (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    registration_date TIMESTAMP NOT NULL DEFAULT NOW(),
    city avito.city_enum NOT NULL
);

CREATE TABLE IF NOT EXISTS avito.receptions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    date_time TIMESTAMP NOT NULL DEFAULT NOW(),
    pvz_id UUID NOT NULL REFERENCES avito.pvz(id) ON DELETE CASCADE,
    status avito.status_enum NOT NULL,
    CONSTRAINT unique_active_reception UNIQUE (pvz_id) WHERE (status = 'in_progress')
);

CREATE TABLE IF NOT EXISTS avito.products (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    date_time TIMESTAMP NOT NULL DEFAULT NOW(),
    type avito.product_type_enum NOT NULL,
    reception_id UUID NOT NULL REFERENCES avito.receptions(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_receptions_pvz_id ON avito.receptions(pvz_id);
CREATE INDEX IF NOT EXISTS idx_receptions_status ON avito.receptions(status);
CREATE INDEX IF NOT EXISTS idx_products_reception_id ON avito.products(reception_id);
CREATE INDEX IF NOT EXISTS idx_products_date_time ON avito.products(date_time);