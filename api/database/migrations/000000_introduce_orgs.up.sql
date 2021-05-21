CREATE TABLE IF NOT EXISTS organizations
(
    id                          SERIAL PRIMARY KEY,
    slug                        VARCHAR(200) NOT NULL,
    title                       TEXT NOT NULL,
    description                 TEXT NULL,
    hidden                      BOOLEAN DEFAULT FALSE,
    protected                   BOOLEAN DEFAULT FALSE,
    tags                        JSONB DEFAULT '[]',
    created                     TIMESTAMP DEFAULT (NOW() AT TIME ZONE 'utc'),
    modified                    TIMESTAMP NULL,
    UNIQUE                      (slug)
);