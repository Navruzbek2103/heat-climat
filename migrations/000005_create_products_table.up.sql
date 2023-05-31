CREATE TABLE IF NOT EXISTS products(
    id SERIAL NOT NULL PRIMARY KEY,
    brand_id INT REFERENCES brands(id),
    category_id INT NOT NULL REFERENCES categories(id),
    title VARCHAR(100) NOT NULL,
    price NUMERIC(10, 2) NOT NULL,
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