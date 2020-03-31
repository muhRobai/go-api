DO $$ BEGIN
  CREATE EXTENSION pgcrypto;
EXCEPTION
  WHEN duplicate_object THEN null;
END $$;

CREATE TABLE customer (
    id UUID NOT NULL PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(50) NOT NULL,
    birthdate DATE NOT NULL
);