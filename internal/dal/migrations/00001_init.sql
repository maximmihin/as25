-- +goose Up
-- +goose StatementBegin

-- CREATE TYPE city AS ENUM (
--     'Msc',
--     'Spb',
--     'Kzn'
-- );

CREATE TYPE city AS ENUM (
    'Москва',
    'Казань',
    'СанктПетербург'
);

CREATE TABLE pvz
(
    id UUID PRIMARY KEY,
    located_at city not null,
    created_at timestamp not null -- todo date without time??
);



CREATE TYPE reception_progress AS ENUM (
    'in_progress',
    'close'
);

CREATE TABLE receptions
(
    id UUID PRIMARY KEY,
    pvz_id uuid not null ,
    status reception_progress not null,
    created_at timestamp not null,

    CONSTRAINT receptions_fk_pvz_id
        FOREIGN KEY (pvz_id)
            REFERENCES pvz
);

CREATE UNIQUE INDEX receptions_idx_one_in_progress_per_pvz
    ON receptions (pvz_id)
    WHERE status = 'in_progress';



CREATE TYPE products_type AS ENUM (
    'электроника',
    'одежда',
    'обувь'
);

CREATE TABLE products
(
    id UUID PRIMARY KEY,
    created_at timestamp not null,
    type products_type not null ,
    reception_id uuid not null,

    CONSTRAINT products_fk_reception_id
        FOREIGN KEY (reception_id)
            REFERENCES receptions
);

CREATE INDEX products_idx_queue_lifo ON products (created_at DESC);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE pvz CASCADE;
DROP TABLE receptions CASCADE;
DROP TABLE products CASCADE;

DROP TYPE city CASCADE;
DROP TYPE reception_progress CASCADE;
DROP TYPE products_type CASCADE;
-- +goose StatementEnd
