CREATE TABLE IF NOT EXISTS new_product_media(
    id SERIAL NOT NULL PRIMARY KEY,
    new_product_id INT NOT NULL REFERENCES new_products(id),
    media_link TEXT
);