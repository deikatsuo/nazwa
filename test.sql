SELECT u.id, u.first_name, u.last_name, p.phone, e.email, r.name AS role_name
FROM "user" u
LEFT JOIN "phone" p ON u.id = p.user_id
LEFT JOIN "email" e ON u.id = e.user_id
LEFT JOIN "user_role" ur ON ur.user_id = u.id
LEFT JOIN "role" r ON ur.role_id = r.id
WHERE r.id = 1;

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
INSERT INTO "product" ("name", "code","base_price","price","type","brand")
VALUES  ('Spring bed biru jambu','SP002','1000000','1500000','spring bed','samsung'),
        ('Liswar Jumbo2x','LW002','500000','1000000','liswar','axis'),
        ('Laptop Dell','LA002','1000000','1500000','laptop','nokia'),
        ('Processor i5','PR002','1000000','1500000','processor','xiaomi'),
        ('Motor yamaha','MO002','12000000','15000000','motor','honda'),
        ('Mobil ford','MB002','1000000','1500000','mobil','yamaha'),
        ('Kapal cargo','KP002','1000000','1500000','kapal','burik'),
        ('Jet Pack Laser','JT003','1000000','1500000','pesawat','kucuk'),
        ('Pesawat jet slam','JT004','1000000','1500000','pesawat','apple'),
        ('Pesawat jet tempur','JT005','1000000','1500000','pesawat','huawei'),
        ('Rumah lipat','RM002','1000000','1500000','rumah','acer');

SELECT
		u.id,
		u.first_name,
		u.last_name,
		u.username,
		u.avatar,
		u.gender,
		TO_CHAR(u.created_at, 'MM/DD/YYYY HH12:MI:SS AM') AS created_at,
		u.balance,
		INITCAP(r.name) AS role
		FROM "user" u
		LEFT JOIN "user_role" ur ON ur.user_id=u.id
		LEFT JOIN "role" r ON r.id=ur.role_id
        WHERE u.id > 1 AND r.id = 1
        LIMIT 10;

ALTER TABLE "product"
ALTER COLUMN "thumbnail" TYPE VARCHAR(50);

ALTER TABLE "product"
ADD COLUMN "thumbnail" VARCHAR(50);

ALTER TABLE "product_photo"
ADD COLUMN "photo" VARCHAR(50) NOT NULL UNIQUE;

INSERT INTO "product_photo" ("product_id", "photo") VALUES ('3', '01234567892.png');

INSERT INTO "order" ("customer_id", "sales_id", "surveyor_id", "shipping_address_id", "billing_address_id", "status", "credit", "notes", "first_time", "code")
VALUES ('12', '12', '7', '4', '4', 'completed', 'true', 'kredit mania', 'true', 'TRX2556555');

INSERT INTO "order_item" ("order_id", "product_id", "quantity", "notes")
VALUES ('7', '5', '5', 'contoh kredit');

SELECT
		id,
		name,
		TO_CHAR(base_price,'Rp999G999G999G999G999') AS base_price,
		TO_CHAR(price,'Rp999G999G999G999G999') AS price,
		code,
		TO_CHAR(created_at, 'MM/DD/YYYY HH12:MI:SS AM') AS created_at,
		type,
		brand
		FROM "product"
		WHERE id=$1
		LIMIT 1