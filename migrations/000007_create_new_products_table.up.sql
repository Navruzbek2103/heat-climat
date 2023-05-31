CREATE TABLE IF NOT EXISTS new_products(
    id SERIAL NOT NULL PRIMARY KEY,
    brand_id INT REFERENCES brands(id),
    title VARCHAR(100) NOT NULL,
    new_price NUMERIC(10, 2) NOT NULL,
    old_price NUMERIC(10, 2),
    description TEXT,
    display_type VARCHAR(100),
    os_type VARCHAR(100),
    camera VARCHAR(60),
    dioganal FLOAT,
    charactestics JSON,
    rating FLOAT,
    created_at TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP WITHOUT TIME ZONE
);