CREATE TABLE IF NOT EXISTS brands(
    id SERIAL NOT NULL PRIMARY KEY,
    brand_name VARCHAR(100) NOT NULL, 
    logo TEXT
);