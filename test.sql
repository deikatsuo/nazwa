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
GROUP BY u.first_name, u.last_name, u.username, u.avatar, u.gender, u.created_at, u.balance, r.name;

SELECT a.id, a.name, a.one, a.two, a.zip, p.name AS province_name, c.name AS city_name, d.name AS district_name, v.name AS village_name
FROM "address" a
JOIN "province" p ON p.id=a.province_id
JOIN "city" c ON c.id=a.city_id
JOIN "district" d ON d.id=a.district_id
JOIN "village" v ON v.id=a.village_id
WHERE user_id=1

-- Buat peran
INSERT INTO "product" ("name", "code","base_price","price","type","credit_twelve","credit_fifteen")
VALUES  ('Spring bed merah jambu','SP001','1000000','1500000','spring bed','3500000','4000000'),
        ('Liswar Jumbo','LW001','500000','1000000','liswar','2000000','2500000'),
        ('Laptop AEX','LA001','1000000','1500000','laptop','3500000','4000000'),
        ('Processor i3','PR001','1000000','1500000','processor','3500000','4000000'),
        ('Motor honda x','MO001','12000000','15000000','motor','35000000','50000000'),
        ('Mobil chevrolet s','MB001','1000000','1500000','mobil','3500000','4000000'),
        ('Kapal ferry Z','KP001','1000000','1500000','kapal','3500000','4000000'),
        ('Jet Pack Radio 2','JT001','1000000','1500000','pesawat','3500000','4000000'),
        ('Pesawat jet hover','JT002','1000000','1500000','pesawat','3500000','4000000'),
        ('Rumah box solo','RM001','1000000','1500000','rumah','3500000','4000000');

