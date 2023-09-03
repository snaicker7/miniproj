

CREATE TABLE IF NOT EXISTS "users" (
		"id" INT PRIMARY KEY NOT NULL,
		"name" TEXT,
		"mobile" VARCHAR(10) UNIQUE,
		"longitude" DECIMAL(9,6),
		"latitude" DECIMAL(9,6),
		"created_at" TIMESTAMP DEFAULT (now()),
		"updated_at" TIMESTAMP 
		);



        CREATE TABLE IF NOT EXISTS "products" (
		"product_id" SERIAL  PRIMARY KEY,
		"product_name" varchar(255),
		"product_description" TEXT,
		"product_images" TEXT[],
		"compressed_product_images" TEXT[],
		"product_price" DECIMAL(10, 2),
		"created_at" TIMESTAMP,
		"updated_at" TIMESTAMP
		);
