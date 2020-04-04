CREATE TABLE product (
    id UUID NOT NULL PRIMARY KEY DEFAULT gen_random_uuid(),
    store_id UUID  NOT NULL,
    product_name TEXT NOT NULL,
    stock   INT NOT NULL,
    amount TEXT NOT NULL,
    created_time TIMESTAMP NOT NULL
);