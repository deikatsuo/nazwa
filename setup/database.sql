CREATE TABLE "users" (
	"user_id" serial(11) NOT NULL,
	"first_name" varchar(25) NOT NULL,
	"last_name" varchar(25),
	"username" varchar(25) NOT NULL UNIQUE,
	"password" char(100),
	"password_salt" char(10),
	"gender" bit NOT NULL,
	"register_date" TIMESTAMP NOT NULL,
	"last_active_date" TIMESTAMP,
	"balance" DECIMAL(15,2) NOT NULL DEFAULT '0',
	CONSTRAINT "users_pk" PRIMARY KEY ("user_id")
) WITH (
  OIDS=FALSE
);



CREATE TABLE "roles" (
	"user_id" int(11) NOT NULL,
	"role" char(4) NOT NULL DEFAULT 'CUST',
	CONSTRAINT "roles_pk" PRIMARY KEY ("user_id")
) WITH (
  OIDS=FALSE
);



CREATE TABLE "emails" (
	"email_id" serial(11) NOT NULL,
	"user_id" int(11) NOT NULL,
	"email" varchar(60) NOT NULL UNIQUE,
	CONSTRAINT "emails_pk" PRIMARY KEY ("email_id")
) WITH (
  OIDS=FALSE
);



CREATE TABLE "addresses" (
	"address_id" serial(11) NOT NULL,
	"user_id" int(11) NOT NULL,
	"address_one" varchar(80) NOT NULL,
	"address_two" varchar(80),
	"zip_code" int(5) NOT NULL,
	"village_id" int(11) NOT NULL,
	"sub_district_id" int(11) NOT NULL,
	"city_id" int(11) NOT NULL,
	"province_id" int(11) NOT NULL,
	"notes" varchar(255),
	CONSTRAINT "addresses_pk" PRIMARY KEY ("address_id")
) WITH (
  OIDS=FALSE
);



CREATE TABLE "orders" (
	"order_id" serial(11) NOT NULL,
	"customer_id" int(11) NOT NULL,
	"sales_id" int(11) NOT NULL,
	"surveyor_id" int(11) NOT NULL,
	"shipping_address_id" int(11) NOT NULL,
	"billing_address_id" int(11),
	"order_status" char(1) NOT NULL DEFAULT 'A',
	"order_date" DATE NOT NULL,
	"shipped_date" DATE,
	"installment" BOOLEAN NOT NULL DEFAULT '0',
	"notes" varchar(255),
	CONSTRAINT "orders_pk" PRIMARY KEY ("order_id")
) WITH (
  OIDS=FALSE
);



CREATE TABLE "provinces" (
	"province_id" serial(11) NOT NULL,
	"province" varchar(50) NOT NULL,
	CONSTRAINT "provinces_pk" PRIMARY KEY ("province_id")
) WITH (
  OIDS=FALSE
);



CREATE TABLE "cities" (
	"city_id" serial(11) NOT NULL,
	"province_id" int(11) NOT NULL,
	"city" varchar(50) NOT NULL,
	CONSTRAINT "cities_pk" PRIMARY KEY ("city_id")
) WITH (
  OIDS=FALSE
);



CREATE TABLE "sub-districts" (
	"sub_district_id" serial(11) NOT NULL,
	"cityid" int(11) NOT NULL,
	"subdistrict" varchar(50) NOT NULL,
	CONSTRAINT "sub-districts_pk" PRIMARY KEY ("sub_district_id")
) WITH (
  OIDS=FALSE
);



CREATE TABLE "villages" (
	"village_id" serial(11) NOT NULL,
	"sub_district_id" int(11) NOT NULL,
	"village" varchar(50) NOT NULL,
	CONSTRAINT "villages_pk" PRIMARY KEY ("village_id")
) WITH (
  OIDS=FALSE
);



CREATE TABLE "order_items" (
	"order_item_id" serial(11) NOT NULL,
	"order_id" int(11) NOT NULL,
	"product_id" int(11) NOT NULL,
	"quantity" int2 NOT NULL,
	CONSTRAINT "order_items_pk" PRIMARY KEY ("order_item_id")
) WITH (
  OIDS=FALSE
);



CREATE TABLE "products" (
	"product_id" serial(11) NOT NULL,
	"product" varchar(100) NOT NULL,
	"price" DECIMAL(15,2) NOT NULL,
	"notes" TEXT(255),
	CONSTRAINT "products_pk" PRIMARY KEY ("product_id")
) WITH (
  OIDS=FALSE
);



CREATE TABLE "installments" (
	"installment_id" serial(11) NOT NULL,
	"order_id" int(11) NOT NULL,
	"collector_id" int(11) NOT NULL,
	"collector_two_id" int(11),
	"recipient_id" int(11) NOT NULL,
	"recipient_two_id" int(11) NOT NULL,
	"duration" int2 NOT NULL,
	"total" DECIMAL(15,2) NOT NULL,
	"remaining" DECIMAL(15,2) NOT NULL,
	"monthly" DECIMAL(15,2) NOT NULL,
	"completed" BOOLEAN NOT NULL DEFAULT '0',
	CONSTRAINT "installments_pk" PRIMARY KEY ("installment_id")
) WITH (
  OIDS=FALSE
);



CREATE TABLE "installment_histories" (
	"installment_history_id" serial(11) NOT NULL,
	"issuer_id" int(11) NOT NULL,
	"collector_id" int(11) NOT NULL,
	"recipient_id" int(11) NOT NULL,
	"issued_date" DATE NOT NULL,
	"paid" DECIMAL(15,2) NOT NULL DEFAULT '0',
	CONSTRAINT "installment_histories_pk" PRIMARY KEY ("installment_history_id")
) WITH (
  OIDS=FALSE
);




ALTER TABLE "roles" ADD CONSTRAINT "roles_fk0" FOREIGN KEY ("user_id") REFERENCES "users"("user_id");

ALTER TABLE "emails" ADD CONSTRAINT "emails_fk0" FOREIGN KEY ("user_id") REFERENCES "users"("user_id");

ALTER TABLE "addresses" ADD CONSTRAINT "addresses_fk0" FOREIGN KEY ("user_id") REFERENCES "users"("user_id");
ALTER TABLE "addresses" ADD CONSTRAINT "addresses_fk1" FOREIGN KEY ("village_id") REFERENCES "villages"("village_id");
ALTER TABLE "addresses" ADD CONSTRAINT "addresses_fk2" FOREIGN KEY ("sub_district_id") REFERENCES "sub-districts"("sub_district_id");
ALTER TABLE "addresses" ADD CONSTRAINT "addresses_fk3" FOREIGN KEY ("city_id") REFERENCES "cities"("city_id");
ALTER TABLE "addresses" ADD CONSTRAINT "addresses_fk4" FOREIGN KEY ("province_id") REFERENCES "provinces"("province_id");

ALTER TABLE "orders" ADD CONSTRAINT "orders_fk0" FOREIGN KEY ("customer_id") REFERENCES "users"("user_id");
ALTER TABLE "orders" ADD CONSTRAINT "orders_fk1" FOREIGN KEY ("sales_id") REFERENCES "users"("user_id");
ALTER TABLE "orders" ADD CONSTRAINT "orders_fk2" FOREIGN KEY ("surveyor_id") REFERENCES "users"("user_id");
ALTER TABLE "orders" ADD CONSTRAINT "orders_fk3" FOREIGN KEY ("shipping_address_id") REFERENCES "addresses"("address_id");
ALTER TABLE "orders" ADD CONSTRAINT "orders_fk4" FOREIGN KEY ("billing_address_id") REFERENCES "addresses"("address_id");


ALTER TABLE "cities" ADD CONSTRAINT "cities_fk0" FOREIGN KEY ("province_id") REFERENCES "provinces"("province_id");

ALTER TABLE "sub-districts" ADD CONSTRAINT "sub-districts_fk0" FOREIGN KEY ("cityid") REFERENCES "cities"("city_id");

ALTER TABLE "villages" ADD CONSTRAINT "villages_fk0" FOREIGN KEY ("sub_district_id") REFERENCES "sub-districts"("sub_district_id");

ALTER TABLE "order_items" ADD CONSTRAINT "order_items_fk0" FOREIGN KEY ("order_id") REFERENCES "orders"("order_id");
ALTER TABLE "order_items" ADD CONSTRAINT "order_items_fk1" FOREIGN KEY ("product_id") REFERENCES "products"("product_id");


ALTER TABLE "installments" ADD CONSTRAINT "installments_fk0" FOREIGN KEY ("order_id") REFERENCES "orders"("order_id");
ALTER TABLE "installments" ADD CONSTRAINT "installments_fk1" FOREIGN KEY ("collector_id") REFERENCES "users"("user_id");
ALTER TABLE "installments" ADD CONSTRAINT "installments_fk2" FOREIGN KEY ("collector_two_id") REFERENCES "users"("user_id");
ALTER TABLE "installments" ADD CONSTRAINT "installments_fk3" FOREIGN KEY ("recipient_id") REFERENCES "users"("user_id");
ALTER TABLE "installments" ADD CONSTRAINT "installments_fk4" FOREIGN KEY ("recipient_two_id") REFERENCES "users"("user_id");

ALTER TABLE "installment_histories" ADD CONSTRAINT "installment_histories_fk0" FOREIGN KEY ("issuer_id") REFERENCES "users"("user_id");
ALTER TABLE "installment_histories" ADD CONSTRAINT "installment_histories_fk1" FOREIGN KEY ("collector_id") REFERENCES "users"("user_id");
ALTER TABLE "installment_histories" ADD CONSTRAINT "installment_histories_fk2" FOREIGN KEY ("recipient_id") REFERENCES "users"("user_id");

