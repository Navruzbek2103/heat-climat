CREATE TABLE IF NOT EXISTS product_media(
    id SERIAL NOT NULL PRIMARY KEY,
    product_id INT NOT NULL REFERENCES products(id),
    media_link TEXT
);