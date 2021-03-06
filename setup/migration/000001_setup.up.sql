/* Nazwa
 *  Membuat tabel utama
 *  Authored by Deri Herdianto
 *  c 2020
 */
-- Tabel role/peran
CREATE TABLE "role" (
    "id" smallint GENERATED BY DEFAULT AS IDENTITY NOT NULL,
    "name" varchar(10) NOT NULL,
    PRIMARY KEY (id)
);

-- Buat peran
INSERT INTO "role" ("id", "name")
    VALUES (0, 'dev'), (1, 'admin'), (2, 'collector'), (3, 'driver'), (4, 'surveyor'), (5, 'sales'), (6, 'customer');

-- Susun ulang nomor sequence untuk tabel role
SELECT
    setval(pg_get_serial_sequence('role', 'id'), coalesce((
            SELECT
                max(id) + 1 FROM "role"), 1), FALSE);

-- SELECT MAX(id) FROM "role";
-- SELECT pg_get_serial_sequence('role', 'id');
-- SELECT nextval('public.role_id_seq');
-- Buat gender
CREATE TYPE "gender" AS ENUM (
    'f',
    'm'
);

-- Tabel user/pengguna
CREATE TABLE "user" (
    "id" int GENERATED ALWAYS AS IDENTITY NOT NULL,
    "first_name" varchar(25) NOT NULL,
    "last_name" varchar(25),
    "ric" varchar(16) UNIQUE NOT NULL,
    "username" varchar(25) UNIQUE,
    "avatar" varchar(50) DEFAULT 'default.jpg' NOT NULL,
    "password" char(60),
    "gender" GENDER NOT NULL,
    "occupation" varchar(25),
    "created_at" timestamp DEFAULT CURRENT_TIMESTAMP,
    "created_by" int,
    "balance" numeric(15) DEFAULT 0,
    PRIMARY KEY (id),
    FOREIGN KEY (created_by) REFERENCES "user" ("id")
);

-- Tabel KK
CREATE TABLE "family_card" (
    "id" int GENERATED ALWAYS AS IDENTITY NOT NULL,
    "number" varchar(16) UNIQUE NOT NULL,
    PRIMARY KEY (id)
);

-- Tabel relasi keluarga
CREATE TABLE "family" (
    "id" int GENERATED ALWAYS AS IDENTITY NOT NULL,
    "family_card_id" int NOT NULL,
    "user_id" int NOT NULL,
    PRIMARY KEY (id),
    FOREIGN KEY (family_card_id) REFERENCES "family_card" ("id"),
    FOREIGN KEY (user_id) REFERENCES "user" ("id")
);

-- Tabel peran user/pengguna
CREATE TABLE "user_role" (
    "id" int GENERATED ALWAYS AS IDENTITY NOT NULL,
    "role_id" int NOT NULL,
    "user_id" int UNIQUE NOT NULL,
    PRIMARY KEY (id),
    FOREIGN KEY (role_id) REFERENCES "role" ("id"),
    FOREIGN KEY (user_id) REFERENCES "user" ("id")
);

-- Tabel email
-- Pengguna boleh memiliki beberapa email
CREATE TABLE "email" (
    "id" int GENERATED ALWAYS AS IDENTITY NOT NULL,
    "user_id" int NOT NULL,
    "email" varchar(60) UNIQUE NOT NULL,
    "verified" boolean DEFAULT 'false' NOT NULL,
    PRIMARY KEY (id),
    FOREIGN KEY (user_id) REFERENCES "user" ("id")
);

-- Tabel nomor HP
-- Pengguna boleh memiliki lebih dari satu nomor HP
CREATE TABLE "phone" (
    "id" int GENERATED ALWAYS AS IDENTITY NOT NULL,
    "user_id" int NOT NULL,
    "phone" varchar(15) UNIQUE NOT NULL,
    "verified" boolean DEFAULT 'false' NOT NULL,
    PRIMARY KEY (id),
    FOREIGN KEY (user_id) REFERENCES "user" ("id")
);

-- Tabel provinsi
CREATE TABLE "province" (
    "id" int GENERATED BY DEFAULT AS IDENTITY NOT NULL,
    "parent" int NOT NULL,
    "name" varchar(50) NOT NULL,
    "original" boolean DEFAULT 'true' NOT NULL,
    "created_by" int,
    PRIMARY KEY (id),
    FOREIGN KEY (created_by) REFERENCES "user" ("id")
);

-- Tabel kota/kabupaten
CREATE TABLE "city" (
    "id" int GENERATED BY DEFAULT AS IDENTITY NOT NULL,
    "parent" int NOT NULL,
    "name" varchar(50) NOT NULL,
    "original" boolean DEFAULT 'true' NOT NULL,
    "created_by" int,
    PRIMARY KEY (id),
    FOREIGN KEY (parent) REFERENCES "province" ("id") ON DELETE CASCADE,
    FOREIGN KEY (created_by) REFERENCES "user" ("id")
);

-- Tabel distrik/kecamatan
CREATE TABLE "district" (
    "id" int GENERATED BY DEFAULT AS IDENTITY NOT NULL,
    "parent" int NOT NULL,
    "name" varchar(50) NOT NULL,
    "original" boolean DEFAULT 'true' NOT NULL,
    "created_by" int,
    PRIMARY KEY (id),
    FOREIGN KEY (parent) REFERENCES "city" ("id") ON DELETE CASCADE,
    FOREIGN KEY (created_by) REFERENCES "user" ("id")
);

-- Tabel kelurahan/desa
CREATE TABLE "village" (
    "id" int GENERATED BY DEFAULT AS IDENTITY NOT NULL,
    "parent" int NOT NULL,
    "name" varchar(50) NOT NULL,
    "original" boolean DEFAULT 'true' NOT NULL,
    "created_by" int,
    PRIMARY KEY (id),
    FOREIGN KEY (parent) REFERENCES "district" ("id") ON DELETE CASCADE,
    FOREIGN KEY (created_by) REFERENCES "user" ("id")
);

-- Tabel zona
CREATE TABLE "zone" (
    "id" int GENERATED BY DEFAULT AS IDENTITY NOT NULL,
    "collector_id" int,
    "name" varchar(25),
    "created_by" int,
    PRIMARY KEY (id),
    FOREIGN KEY (collector_id) REFERENCES "user" ("id") ON DELETE SET NULL,
    FOREIGN KEY (created_by) REFERENCES "user" ("id") ON DELETE SET NULL
);

-- Buat zona bawaan
INSERT INTO "zone" ("name")
    VALUES ('Zona 1'), ('Zona 2'), ('Zona 3');

-- Susun serial sequence
SELECT
    setval(pg_get_serial_sequence('zone', 'id'), coalesce((
            SELECT
                max(id) + 1 FROM "zone"), 1), FALSE);

-- Tabel line/arah zona
CREATE TABLE "zone_line" (
    "id" int GENERATED ALWAYS AS IDENTITY NOT NULL,
    "code" varchar(5) UNIQUE NOT NULL,
    "name" varchar(25) NOT NULL,
    PRIMARY KEY (id)
);

-- Tabel line list
CREATE TABLE "zone_line_list" (
    "id" int GENERATED ALWAYS AS IDENTITY NOT NULL,
    "zone_line_code" varchar(5) NOT NULL,
    "number" int NOT NULL,
    PRIMARY KEY (id)
);

-- Tabel list arah yang termasuk dalam zona
CREATE TABLE "zone_list" (
    "id" int GENERATED ALWAYS AS IDENTITY NOT NULL,
    "zone_id" int NOT NULL,
    "zone_line_id" int UNIQUE NOT NULL,
    PRIMARY KEY (id),
    FOREIGN KEY (zone_id) REFERENCES "zone" ("id") ON DELETE CASCADE,
    FOREIGN KEY (zone_line_id) REFERENCES "zone_line" ("id")
);

-- Tabel address/alamat
CREATE TABLE "address" (
    "id" int GENERATED ALWAYS AS IDENTITY NOT NULL,
    "user_id" int NOT NULL,
    "name" varchar(25) NOT NULL,
    "description" varchar(100),
    "one" varchar(80) NOT NULL,
    "two" varchar(80),
    "zip" varchar(5),
    "province_id" int NOT NULL,
    "city_id" int NOT NULL,
    "district_id" int NOT NULL,
    "village_id" bigint NOT NULL,
    PRIMARY KEY (id),
    FOREIGN KEY (user_id) REFERENCES "user" ("id"),
    FOREIGN KEY (province_id) REFERENCES "province" ("id") ON DELETE SET NULL,
    FOREIGN KEY (city_id) REFERENCES "city" ("id") ON DELETE SET NULL,
    FOREIGN KEY (district_id) REFERENCES "district" ("id") ON DELETE SET NULL,
    FOREIGN KEY (village_id) REFERENCES "village" ("id") ON DELETE SET NULL
);

-- Tabel barang/produk
CREATE TABLE "product" (
    "id" int GENERATED ALWAYS AS IDENTITY NOT NULL,
    "created_by" int NOT NULL,
    "name" varchar(100) NOT NULL,
    "slug" varchar(100) UNIQUE NOT NULL,
    "stock" int NOT NULL DEFAULT 0,
    "code" varchar(10) UNIQUE NOT NULL,
    "base_price" numeric(15) DEFAULT 0,
    "price" numeric(15) DEFAULT 0,
    "type" varchar(25),
    "brand" varchar(25),
    "thumbnail" varchar(50),
    "created_at" timestamp DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (id),
    FOREIGN KEY (created_by) REFERENCES "user" ("id")
);

-- Tabel harga kredit barang/produk
CREATE TABLE "product_credit_price" (
    "id" int GENERATED ALWAYS AS IDENTITY NOT NULL,
    "product_id" int NOT NULL,
    "duration" smallint NOT NULL,
    "price" numeric(15) DEFAULT 0,
    PRIMARY KEY (id),
    FOREIGN KEY (product_id) REFERENCES "product" ("id")
);

-- Tabel photo produk
CREATE TABLE "product_photo" (
    "id" int GENERATED ALWAYS AS IDENTITY NOT NULL,
    "product_id" int NOT NULL,
    "photo" varchar(50) NOT NULL UNIQUE,
    PRIMARY KEY (id),
    FOREIGN KEY (product_id) REFERENCES "product" ("id")
);

-- Tabel order/penjualan
CREATE TABLE "order" (
    "id" int GENERATED ALWAYS AS IDENTITY NOT NULL,
    "customer_id" int NOT NULL,
    "sales_id" int,
    "surveyor_id" int,
    "driver_id" int,
    "shipping_address_id" int NOT NULL,
    "billing_address_id" int NOT NULL,
    "code" varchar(50) NOT NULL UNIQUE,
    "status" varchar(25) NOT NULL DEFAULT 'pending',
    "credit" boolean DEFAULT 'false' NOT NULL,
    "notes" varchar(100),
    "order_date" date NOT NULL DEFAULT CURRENT_DATE,
    "shipping_date" date NOT NULL DEFAULT CURRENT_DATE,
    "created_at" timestamp DEFAULT CURRENT_TIMESTAMP,
    "created_by" int NOT NULL,
    "deposit" numeric(15) DEFAULT 0,
    "price_total" numeric(15) NOT NULL,
    "base_price_total" numeric(15) NOT NULL,
    PRIMARY KEY (id),
    FOREIGN KEY (customer_id) REFERENCES "user" ("id"),
    FOREIGN KEY (sales_id) REFERENCES "user" ("id"),
    FOREIGN KEY (surveyor_id) REFERENCES "user" ("id"),
    FOREIGN KEY (driver_id) REFERENCES "user" ("id"),
    FOREIGN KEY (shipping_address_id) REFERENCES "address" ("id"),
    FOREIGN KEY (billing_address_id) REFERENCES "address" ("id"),
    FOREIGN KEY (created_by) REFERENCES "user" ("id")
);

-- Tabel anggota keluarga penanggung jawab user
CREATE TABLE "order_user_substitute" (
    "id" int GENERATED ALWAYS AS IDENTITY NOT NULL,
    "ric" varchar(16),
    "first_name" varchar(25) NOT NULL,
    "last_name" varchar(25),
    "gender" GENDER NOT NULL,
    "substitute_to" int NOT NULL,
    "created_at" timestamp DEFAULT CURRENT_TIMESTAMP,
    "created_by" int,
    PRIMARY KEY (id),
    FOREIGN KEY (substitute_to) REFERENCES "order" ("id") ON DELETE CASCADE,
    FOREIGN KEY (created_by) REFERENCES "user" ("id")
);

-- Tabel order item
CREATE TABLE "order_item" (
    "id" int GENERATED ALWAYS AS IDENTITY NOT NULL,
    "order_id" int NOT NULL,
    "product_id" int NOT NULL,
    "quantity" int NOT NULL DEFAULT 1,
    "notes" varchar(100),
    "base_price" numeric(15) DEFAULT 0,
    "price" numeric(15) DEFAULT 0,
    "discount" numeric(15) DEFAULT 0,
    PRIMARY KEY (id),
    FOREIGN KEY (order_id) REFERENCES "order" ("id") ON DELETE CASCADE,
    FOREIGN KEY (product_id) REFERENCES "product" ("id")
);

-- Tabel detail credit
CREATE TABLE "order_credit_detail" (
    "id" int GENERATED ALWAYS AS IDENTITY NOT NULL,
    "order_id" int NOT NULL UNIQUE,
    "zone_line_id" int,
    "credit_code" varchar(25) NOT NULL,
    "monthly" int NOT NULL,
    "duration" smallint NOT NULL,
    "due" smallint NOT NULL,
    "total" numeric(15) NOT NULL,
    "remaining" numeric(15) NOT NULL,
    "lucky_discount" numeric(15) DEFAULT 0,
    "last_paid" date,
    "active" bool NOT NULL DEFAULT 'true',
    "done" bool NOT NULL DEFAULT 'false',
    PRIMARY KEY (id),
    FOREIGN KEY (order_id) REFERENCES "order" ("id") ON DELETE CASCADE,
    FOREIGN KEY (zone_line_id) REFERENCES "zone_line" ("id") ON DELETE SET NULL
);

-- Tabel kredit bulanan
CREATE TABLE "order_monthly_credit" (
    "id" int GENERATED ALWAYS AS IDENTITY NOT NULL,
    "order_id" int NOT NULL,
    "code" varchar(20) NOT NULL UNIQUE,
    "nth" smallint NOT NULL,
    "due_date" date NOT NULL,
    "print_date" date,
    "promise" date,
    "paid" numeric(15) DEFAULT 0 NOT NULL,
    "notes" varchar(100),
    "position" char(3) DEFAULT 'in' NOT NULL,
    "printed" boolean DEFAULT 'false' NOT NULL,
    "done" boolean DEFAULT 'false' NOT NULL,
    PRIMARY KEY (id),
    FOREIGN KEY (order_id) REFERENCES "order" ("id") ON DELETE CASCADE
);

-- Tabel log kredit bulanan
CREATE TABLE "log_order_monthly_credit" (
    "id" int GENERATED ALWAYS AS IDENTITY NOT NULL,
    "order_monthly_credit_id" int,
    "money_in" numeric(15) NOT NULL,
    "created_at" timestamp DEFAULT CURRENT_TIMESTAMP,
    "created_by" int,
    "collected_by" int,
    PRIMARY KEY (id),
    FOREIGN KEY (order_monthly_credit_id) REFERENCES "order_monthly_credit" ("id") ON DELETE CASCADE,
    FOREIGN KEY (created_by) REFERENCES "user" ("id") ON DELETE SET NULL,
    FOREIGN KEY (collected_by) REFERENCES "user" ("id") ON DELETE SET NULL
);

