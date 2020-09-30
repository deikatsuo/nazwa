SELECT u.first_name, u.last_name, p.phone, e.email, r.name AS role_name
FROM "user" u
JOIN "phone" p ON u.id = p.user_id
JOIN "email" e ON u.id = e.user_id
JOIN "user_role" ur ON u.id = ur.user_id
JOIN "role" r ON ur.role_id = r.id;

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
WHERE email='thisisderi@gmail.com'