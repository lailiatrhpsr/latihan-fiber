-- ---------------------------------------------------------------
-- permissions baru untuk resource students
-- ---------------------------------------------------------------
INSERT INTO permissions (name, description) VALUES
    ('student:list',       'Melihat daftar seluruh mahasiswa'),
    ('student:read:any',   'Melihat data mahasiswa mana pun'),
    ('student:create',     'Mendaftarkan data mahasiswa baru'),
    ('student:update:any', 'Mengubah data mahasiswa mana pun'),
    ('student:delete',     'Menghapus data mahasiswa')
ON CONFLICT (name) DO NOTHING;

INSERT INTO role_permissions (role_name, permission_name) VALUES
    ('admin', 'student:list'),
    ('admin', 'student:read:any'),
    ('admin', 'student:create'),
    ('admin', 'student:update:any'),
    ('admin', 'student:delete'),
    ('staff', 'student:list'),
    ('staff', 'student:read:any'),
    ('staff', 'student:create')
ON CONFLICT DO NOTHING;
-- role 'user' sengaja tidak diberi permission apa pun: mereka hanya boleh
-- melihat dan mengubah data mahasiswa yang owner_id-nya adalah dirinya
-- sendiri, lewat pemeriksaan kepemilikan di service (bukan lewat permission).

-- ---------------------------------------------------------------
-- owner_id — menandai user mana yang mendaftarkan baris data tersebut.
-- Inilah dasar pemeriksaan kepemilikan pada GET/PUT/PATCH /students/:id.
-- ---------------------------------------------------------------
ALTER TABLE students ADD COLUMN IF NOT EXISTS owner_id INTEGER;

-- Tabel students sudah berisi data lama dengan owner_id masih NULL, 
-- sehingga NOT NULL akan gagal diterapkan. Karena pemilik asli data lama tidak dapat diketahui, 
-- data tersebut akan diberikan kepada admin pertama sebagai pengelola data warisan. 
-- Pengisian owner_id wajib dilakukan sebelum menambahkan constraint NOT NULL dan FOREIGN KEY

UPDATE students
SET owner_id = (SELECT id FROM users WHERE role = 'admin' ORDER BY id ASC LIMIT 1)
WHERE owner_id IS NULL;

ALTER TABLE students ALTER COLUMN owner_id SET NOT NULL;

ALTER TABLE students DROP CONSTRAINT IF EXISTS students_owner_id_fkey;
ALTER TABLE students
    ADD CONSTRAINT students_owner_id_fkey
    FOREIGN KEY (owner_id) REFERENCES users(id) ON DELETE CASCADE;

CREATE INDEX IF NOT EXISTS students_owner_id_idx ON students (owner_id);