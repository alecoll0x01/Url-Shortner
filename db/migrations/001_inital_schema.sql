-- migrate:up

CREATE TABLE urls (
    id         BIGSERIAL    PRIMARY KEY,
    code       VARCHAR(32)  NOT NULL UNIQUE,
    original   TEXT         NOT NULL,
    clicks     BIGINT       NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    expires_at TIMESTAMPTZ
);

CREATE INDEX idx_urls_code ON urls (code);

-- migrate:down

DROP TABLE IF EXISTS urls;