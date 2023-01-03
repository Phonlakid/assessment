-- Sequence and defined type
CREATE SEQUENCE IF NOT EXISTS expenses;

-- Table Definition
CREATE TABLE IF NOT EXISTS expenses (
    id SERIAL PRIMARY KEY,
    title TEXT,
    amount FLOAT,
    note TEXT,
    tags TEXT[]
);



