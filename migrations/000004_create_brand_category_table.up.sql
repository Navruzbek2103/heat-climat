CREATE TABLE IF NOT EXISTS brand_category(
    brand_id INT NOT NULL REFERENCES brands(id),
    category_id INT REFERENCES categories(id),
    product_count INT
);