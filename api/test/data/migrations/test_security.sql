INSERT INTO
    security_credentials (id, password, active)
VALUES
    (1000, crypt('active-password', gen_salt('bf')), true),
    (1001, crypt('inactive-password', gen_salt('bf')), false);