/* Nazwa
*  Membuat tabel utama
*  Authored by Deri Herdianto
*  c 2020
*/

-- Tabel role/peran
CREATE TABLE "role" (
    "id" SMALLINT GENERATED BY DEFAULT AS IDENTITY NOT NULL,
    "name" VARCHAR(10) NOT NULL,
    PRIMARY KEY (id)
);

-- Buat peran
INSERT INTO "role" ("id", "name")
VALUES  (0, 'dev'),
        (1, 'admin'),
        (2, 'collector'),
        (3, 'driver'),
        (4, 'surveyor'),
        (5, 'sales'),
        (6, 'customer'),
        (7, 'substitute');

-- Buat gender
CREATE TYPE "gender" AS ENUM ('f', 'm');

-- Tabel user/pengguna
CREATE TABLE "user" (
    "id" INT GENERATED ALWAYS AS IDENTITY NOT NULL,
    "first_name" VARCHAR(25) NOT NULL,
    "last_name" VARCHAR(25),
    "ric" VARCHAR(16) UNIQUE NOT NULL,
    "username" VARCHAR(25) UNIQUE,
    "avatar" VARCHAR(50) DEFAULT 'default.jpg' NOT NULL,
    "password" CHAR(60),
    "gender" GENDER NOT NULL,
    "occupation" VARCHAR(25),
    "created_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    "created_by" INT,
    "balance" NUMERIC(15) DEFAULT 0,
    PRIMARY KEY (id),
    FOREIGN KEY (created_by) REFERENCES "user"("id")
);

-- Tabel KK
CREATE TABLE "family_card" (
    "id" INT GENERATED ALWAYS AS IDENTITY NOT NULL,
    "number" VARCHAR(16) UNIQUE NOT NULL,
    PRIMARY KEY (id)
);

-- Tabel relasi keluarga
CREATE TABLE "family" (
    "id" INT GENERATED ALWAYS AS IDENTITY NOT NULL,
    "family_card_id" INT NOT NULL,
    "user_id" INT NOT NULL,
    PRIMARY KEY (id),
    FOREIGN KEY (family_card_id) REFERENCES "family_card"("id"),
    FOREIGN KEY (user_id) REFERENCES "user"("id")
);

-- Tabel peran user/pengguna
CREATE TABLE "user_role" (
    "id" INT GENERATED ALWAYS AS IDENTITY NOT NULL,
    "role_id" INT NOT NULL,
    "user_id" INT UNIQUE NOT NULL,
    PRIMARY KEY (id),
    FOREIGN KEY (role_id) REFERENCES "role"("id"),
    FOREIGN KEY (user_id) REFERENCES "user"("id")
);

-- Tabel email
-- Pengguna boleh memiliki beberapa email
CREATE TABLE "email" (
    "id" INT GENERATED ALWAYS AS IDENTITY NOT NULL,
    "user_id" INT NOT NULL,
    "email" VARCHAR(60) UNIQUE NOT NULL,
    "verified" BOOLEAN DEFAULT 'false' NOT NULL,
    PRIMARY KEY (id),
    FOREIGN KEY (user_id) REFERENCES "user"("id")
);

-- Tabel nomor HP
-- Pengguna boleh memiliki lebih dari satu nomor HP
CREATE TABLE "phone" (
    "id" INT GENERATED ALWAYS AS IDENTITY NOT NULL,
    "user_id" INT NOT NULL,
    "phone" VARCHAR(15) UNIQUE NOT NULL,
    "verified" BOOLEAN DEFAULT 'false' NOT NULL,
    PRIMARY KEY (id),
    FOREIGN KEY (user_id) REFERENCES "user"("id")
);

-- Tabel provinsi
CREATE TABLE "province" (
    "id" INT GENERATED BY DEFAULT AS IDENTITY NOT NULL,
    "parent" INT NOT NULL,
    "name" VARCHAR (50) NOT NULL,
    "deault" BOOLEAN DEFAULT 'true' NOT NULL,
    PRIMARY KEY (id)
);

-- Tabel kota/kabupaten
CREATE TABLE "city" (
    "id" INT GENERATED BY DEFAULT AS IDENTITY NOT NULL,
    "parent" INT NOT NULL,
    "name" VARCHAR (50) NOT NULL,
    "deault" BOOLEAN DEFAULT 'true' NOT NULL,
    PRIMARY KEY (id),
    FOREIGN KEY (parent) REFERENCES "province"("id")
);

-- Tabel distrik/kecamatan
CREATE TABLE "district" (
    "id" INT GENERATED BY DEFAULT AS IDENTITY NOT NULL,
    "parent" INT NOT NULL,
    "name" VARCHAR(50) NOT NULL,
    "default" BOOLEAN DEFAULT 'true' NOT NULL,
    PRIMARY KEY (id),
    FOREIGN KEY (parent) REFERENCES "city"("id")
);

-- Tabel kelurahan/desa
CREATE TABLE "village" (
    "id" BIGINT GENERATED BY DEFAULT AS IDENTITY NOT NULL,
    "parent" INT NOT NULL,
    "name" VARCHAR(50) NOT NULL,
    "default" BOOLEAN DEFAULT 'true' NOT NULL,
    PRIMARY KEY (id),
    FOREIGN KEY (parent) REFERENCES "district"("id")
);

-- Tabel address/alamat
CREATE TABLE "address" (
    "id" INT GENERATED ALWAYS AS IDENTITY NOT NULL,
    "user_id" INT NOT NULL,
    "name" VARCHAR(25) NOT NULL,
    "description" VARCHAR(100),
    "one" VARCHAR(80) NOT NULL,
    "two" VARCHAR(80),
    "zip" VARCHAR(5),
    "province_id" INT NOT NULL,
    "city_id" INT NOT NULL,
    "district_id" INT NOT NULL,
    "village_id" BIGINT NOT NULL,
    PRIMARY KEY (id),
    FOREIGN KEY (user_id) REFERENCES "user"("id"),
    FOREIGN KEY (province_id) REFERENCES "province"("id"),
    FOREIGN KEY (city_id) REFERENCES "city"("id"),
    FOREIGN KEY (district_id) REFERENCES "district"("id"),
    FOREIGN KEY (village_id) REFERENCES "village"("id")
);

-- Tabel barang/produk
CREATE TABLE "product" (
    "id" INT GENERATED ALWAYS AS IDENTITY NOT NULL,
    "created_by" INT NOT NULL,
    "name" VARCHAR(100) NOT NULL,
    "code" VARCHAR(10) UNIQUE NOT NULL,
    "base_price" NUMERIC(15) DEFAULT 0,
    "price" NUMERIC(15) DEFAULT 0,
    "type" VARCHAR(25),
    "brand" VARCHAR(25),
    "thumbnail" VARCHAR(50),
    "created_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (id),
    FOREIGN KEY (created_by) REFERENCES "user"("id")
);

-- Tabel harga kredit barang/produk
CREATE TABLE "product_credit_price" (
    "id" INT GENERATED ALWAYS AS IDENTITY NOT NULL,
    "product_id" INT NOT NULL,
    "duration" SMALLINT NOT NULL,
    "price" NUMERIC(15) DEFAULT 0,
    PRIMARY KEY (id),
    FOREIGN KEY (product_id) REFERENCES "product"("id")
);

-- Tabel photo produk
CREATE TABLE "product_photo" (
    "id" INT GENERATED ALWAYS AS IDENTITY NOT NULL,
    "product_id" INT NOT NULL,
    "photo" VARCHAR(50) NOT NULL UNIQUE,
    PRIMARY KEY (id),
    FOREIGN KEY (product_id) REFERENCES "product"("id")
);

-- Tabel order/penjualan
CREATE TABLE "order" (
    "id" INT GENERATED ALWAYS AS IDENTITY NOT NULL,
    "customer_id" INT NOT NULL,
    "sales_id" INT,
    "surveyor_id" INT,
    "collector_id" INT,
    "driver_id" INT,
    "shipping_address_id" INT NOT NULL,
    "billing_address_id" INT,
    "code" VARCHAR(50) NOT NULL UNIQUE,
    "status" VARCHAR(25) NOT NULL DEFAULT 'pending',
    "credit" BOOLEAN DEFAULT 'false' NOT NULL,
    "notes" VARCHAR(100),
    "order_date" DATE NOT NULL DEFAULT CURRENT_DATE,
    "shipping_date" DATE NOT NULL DEFAULT CURRENT_DATE,
    "created_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    "created_by" INT NOT NULL,
    "deposit" NUMERIC(15) DEFAULT 0,
    "price_total" NUMERIC(15) NOT NULL,
    "base_price_total" NUMERIC(15) NOT NULL,
    PRIMARY KEY (id),
    FOREIGN KEY (customer_id) REFERENCES "user"("id"),
    FOREIGN KEY (sales_id) REFERENCES "user"("id"),
    FOREIGN KEY (surveyor_id) REFERENCES "user"("id"),
    FOREIGN KEY (driver_id) REFERENCES "user"("id"),
    FOREIGN KEY (shipping_address_id) REFERENCES "address"("id"),
    FOREIGN KEY (billing_address_id) REFERENCES "address"("id"),
    FOREIGN KEY (created_by) REFERENCES "user"("id")
);

-- Tabel anggota keluarga penanggung jawab user
CREATE TABLE "order_user_substitute" (
    "id" INT GENERATED ALWAYS AS IDENTITY NOT NULL,
    "ric" VARCHAR(16),
    "first_name" VARCHAR(25) NOT NULL,
    "last_name" VARCHAR(25),
    "gender" GENDER NOT NULL,
    "substitute_to" INT NOT NULL,
    "created_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    "created_by" INT,
    PRIMARY KEY (id),
    FOREIGN KEY (substitute_to) REFERENCES "order"("id") ON DELETE CASCADE,
    FOREIGN KEY (created_by) REFERENCES "user"("id")
);

-- Tabel order item
CREATE TABLE "order_item" (
    "id" INT GENERATED ALWAYS AS IDENTITY NOT NULL,
    "order_id" INT NOT NULL,
    "product_id" INT NOT NULL,
    "quantity" INT NOT NULL DEFAULT 1,
    "notes" VARCHAR(100),
    "base_price" NUMERIC(15) DEFAULT 0,
    "price" NUMERIC(15) DEFAULT 0,
    "discount" NUMERIC(15) DEFAULT 0,
    PRIMARY KEY (id),
    FOREIGN KEY (order_id) REFERENCES "order"("id") ON DELETE CASCADE,
    FOREIGN KEY (product_id) REFERENCES "product"("id")
);

-- Tabel detail credit
CREATE TABLE "order_credit_detail" (
    "id" INT GENERATED ALWAYS AS IDENTITY NOT NULL,
    "order_id" INT NOT NULL UNIQUE,
    "monthly" INT NOT NULL,
    "duration" SMALLINT NOT NULL,
    "due" SMALLINT NOT NULL,
    "remaining"  NUMERIC(15) NOT NULL,
    "lucky_discount" NUMERIC(15) DEFAULT 0,
    PRIMARY KEY (id),
    FOREIGN KEY (order_id) REFERENCES "order"("id") ON DELETE CASCADE
);