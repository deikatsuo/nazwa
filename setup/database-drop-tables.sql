ALTER TABLE "roles" DROP CONSTRAINT IF EXISTS "roles_fk0";

ALTER TABLE "emails" DROP CONSTRAINT IF EXISTS "emails_fk0";

ALTER TABLE "addresses" DROP CONSTRAINT IF EXISTS "addresses_fk0";

ALTER TABLE "addresses" DROP CONSTRAINT IF EXISTS "addresses_fk1";

ALTER TABLE "addresses" DROP CONSTRAINT IF EXISTS "addresses_fk2";

ALTER TABLE "addresses" DROP CONSTRAINT IF EXISTS "addresses_fk3";

ALTER TABLE "addresses" DROP CONSTRAINT IF EXISTS "addresses_fk4";

ALTER TABLE "orders" DROP CONSTRAINT IF EXISTS "orders_fk0";

ALTER TABLE "orders" DROP CONSTRAINT IF EXISTS "orders_fk1";

ALTER TABLE "orders" DROP CONSTRAINT IF EXISTS "orders_fk2";

ALTER TABLE "orders" DROP CONSTRAINT IF EXISTS "orders_fk3";

ALTER TABLE "orders" DROP CONSTRAINT IF EXISTS "orders_fk4";

ALTER TABLE "cities" DROP CONSTRAINT IF EXISTS "cities_fk0";

ALTER TABLE "sub-districts" DROP CONSTRAINT IF EXISTS "sub-districts_fk0";

ALTER TABLE "villages" DROP CONSTRAINT IF EXISTS "villages_fk0";

ALTER TABLE "order_items" DROP CONSTRAINT IF EXISTS "order_items_fk0";

ALTER TABLE "order_items" DROP CONSTRAINT IF EXISTS "order_items_fk1";

ALTER TABLE "installments" DROP CONSTRAINT IF EXISTS "installments_fk0";

ALTER TABLE "installments" DROP CONSTRAINT IF EXISTS "installments_fk1";

ALTER TABLE "installments" DROP CONSTRAINT IF EXISTS "installments_fk2";

ALTER TABLE "installments" DROP CONSTRAINT IF EXISTS "installments_fk3";

ALTER TABLE "installments" DROP CONSTRAINT IF EXISTS "installments_fk4";

ALTER TABLE "installment_histories" DROP CONSTRAINT IF EXISTS "installment_histories_fk0";

ALTER TABLE "installment_histories" DROP CONSTRAINT IF EXISTS "installment_histories_fk1";

ALTER TABLE "installment_histories" DROP CONSTRAINT IF EXISTS "installment_histories_fk2";

DROP TABLE IF EXISTS "users";

DROP TABLE IF EXISTS "roles";

DROP TABLE IF EXISTS "emails";

DROP TABLE IF EXISTS "addresses";

DROP TABLE IF EXISTS "orders";

DROP TABLE IF EXISTS "provinces";

DROP TABLE IF EXISTS "cities";

DROP TABLE IF EXISTS "sub-districts";

DROP TABLE IF EXISTS "villages";

DROP TABLE IF EXISTS "order_items";

DROP TABLE IF EXISTS "products";

DROP TABLE IF EXISTS "installments";

DROP TABLE IF EXISTS "installment_histories";

