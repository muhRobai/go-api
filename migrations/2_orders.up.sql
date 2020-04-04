CREATE TABLE orders (
    id UUID NOT NULL PRIMARY KEY DEFAULT gen_random_uuid(),
    buyer UUID NOT NULL,
    status INT NOT NULL,
    payment_type INT NOT NULL,
    payment_code JSON NOT NULL,
    effective_time TIMESTAMP NOT NULL,
    expiry_time TIMESTAMP NOT NULL,
    created_time TIMESTAMP NOT NULL,
    update_time TIMESTAMP NULL,
    delete_time TIMESTAMP NULL
);