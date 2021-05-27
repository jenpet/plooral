INSERT INTO
    security_credentials (password, active)
VALUES
    (crypt('user-protected', gen_salt('bf')), true),
    (crypt('owner-protected', gen_salt('bf')), true),
    (crypt('user-hidden', gen_salt('bf')), true),
    (crypt('owner-hidden', gen_salt('bf')), true);

INSERT INTO
    organizations (slug, title, description, hidden, protected, user_credentials_id, owner_credentials_id)
VALUES
    (
        'org-tests-regular',
        'Regular Organization',
        'Sample organization for org testing.',
        false,
        false,
        NULL,
        NULL
    ),
    (
        'org-tests-protected',
        'Protected Organization',
        'Sample organization for org testing.',
        false,
        true,
        (SELECT id FROM security_credentials WHERE password = crypt('user-protected', password)),
        (SELECT id FROM security_credentials WHERE password = crypt('owner-protected', password))
    ),
    (
        'org-tests-hidden',
        'Hidden Organization',
        'Sample organization for org testing.',
        true,
        true,
        (SELECT id FROM security_credentials WHERE password = crypt('user-hidden', password)),
        (SELECT id FROM security_credentials WHERE password = crypt('owner-hidden', password))
    );