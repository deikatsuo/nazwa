SELECT u.id, u.first_name, u.last_name, p.phone, e.email, r.name AS role_name
FROM "user" u
LEFT JOIN "phone" p ON u.id = p.user_id
LEFT JOIN "email" e ON u.id = e.user_id
LEFT JOIN "user_role" ur ON ur.user_id = u.id
LEFT JOIN "role" r ON ur.role_id = r.id;

SELECT id
FROM "user"
WHERE username='rika'
UNION
SELECT user_id
FROM "phone"
WHERE phone='081282683627'
UNION
SELECT user_id
FROM "email"
WHERE email='thisisderi@gmail.com';

SELECT
    u.first_name,
    u.last_name,
    u.username,
    u.avatar,
    u.gender,
    u.created_at,
    u.balance,
    string_agg(DISTINCT p.phone, ',' ORDER BY p.phone) AS phone,
    string_agg(DISTINCT e.email, ',' ORDER BY e.email) AS email,
    r.name AS role
FROM "user" u
LEFT JOIN "phone" p ON p.user_id=u.id
LEFT JOIN "email" e ON e.user_id=u.id
LEFT JOIN "user_role" ur ON ur.user_id=u.id
LEFT JOIN "role" r ON r.id=ur.role_id
GROUP BY u.first_name, u.last_name, u.username, u.avatar, u.gender, u.created_at, u.balance, r.name