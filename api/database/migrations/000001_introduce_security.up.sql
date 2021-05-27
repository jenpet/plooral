CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE IF NOT EXISTS security_credentials (
    id                          SERIAL PRIMARY KEY,
    password                    TEXT NOT NULL,
    created                     TIMESTAMP DEFAULT (NOW() AT TIME ZONE 'utc'),
    active                      BOOL DEFAULT TRUE
);