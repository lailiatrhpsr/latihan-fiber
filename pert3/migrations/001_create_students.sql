CREATE TABLE IF NOT EXISTS students (
    id SERIAL PRIMARY KEY,
    nim VARCHAR(20) NOT NULL,
    name VARCHAR(150) NOT NULL,
    grade NUMERIC(5,2) NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Indeks unik untuk NIM (case-insensitive)
CREATE UNIQUE INDEX IF NOT EXISTS students_nim_lower_key 
ON students (LOWER(nim));

-- Indeks B-Tree untuk mempercepat pencarian nama
CREATE INDEX IF NOT EXISTS students_name_idx 
ON students (name);