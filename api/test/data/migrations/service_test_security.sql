INSERT INTO
    security_passwords (password, active)
VALUES
    (crypt('active-password', gen_salt('bf')), true),
    (crypt('inactive-password', gen_salt('bf')), false);