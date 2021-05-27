ALTER TABLE
    organizations
ADD COLUMN
    user_credentials_id        INTEGER REFERENCES security_credentials(id) NULL,
ADD COLUMN
    owner_credentials_id       INTEGER REFERENCES security_credentials(id) NULL;