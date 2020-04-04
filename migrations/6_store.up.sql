CREATE TABLE store (
    id UUID NOT NULL PRIMARY KEY DEFAULT gen_random_uuid(),
    customer_id TEXT NOT NULL,
    name TEXT NOT NULL
);