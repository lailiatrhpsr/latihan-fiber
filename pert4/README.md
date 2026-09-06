# REST API Students — PostgreSQL (Fiber)

Tugas Mandiri: REST API CRUD untuk data mahasiswa (`students`) menggunakan [Go Fiber v2](https://gofiber.io/) dan PostgreSQL sebagai database (via `pgx/v5`).

## Tech Stack

- **Bahasa**: Go
- **Framework**: [gofiber/fiber v2](https://github.com/gofiber/fiber)
- **Database**: PostgreSQL
- **Driver DB**: [jackc/pgx/v5](https://github.com/jackc/pgx) (dengan connection pool `pgxpool`)
- **Env loader**: [joho/godotenv](https://github.com/joho/godotenv)

## Struktur Folder

Direstrukturisasi pada Pertemuan 4 mengikuti struktur baku Clean Architecture
mata kuliah ini. Perilaku HTTP tidak berubah dari Pertemuan 3 — yang berubah
hanya letak kodenya (lihat `Modul_4_-_Clean_Architecture.docx` untuk
penjelasan lengkap tiap layer).

```
pert4/
├── app/
│   ├── model/                     # Layer 1 — Entities
│   │   ├── student.go             #   struct Student, request DTO, ListQuery
│   │   └── web_response.go        #   WebResponse & Meta (bentuk response seragam)
│   ├── repository/                # Layer 3 — Interface Adapters (gateway)
│   │   └── student_repository.go  #   interface + implementasi query ke tabel students
│   └── service/                   # Layer 2 (rules) + Layer 3 (controller)
│       ├── student_rules.go       #   business rules MURNI, tidak impor Fiber sama sekali
│       ├── student_rules_test.go  #   unit test untuk business rules
│       └── student_service.go     #   menerima fiber.Ctx, orkestrasi request
├── config/                        # Layer 4 — Frameworks & Drivers
│   ├── app.go                     #   rakit instance Fiber + middleware + route
│   ├── env.go                     #   load & baca environment variable (.env)
│   └── logger.go                  #   logger terstruktur (slog + rotasi file lumberjack)
├── database/                      # Layer 4
│   └── postgres.go                #   inisialisasi PostgreSQL connection pool (pgxpool)
├── helper/                        # Layer 3 — Interface Adapters (presenter)
│   ├── response.go                #   pembungkus response sukses/gagal
│   └── request.go                 #   baca query string & parameter :id
├── middleware/                    # Layer 4
│   └── middleware.go              #   RequireJSON, request logger, pendaftaran middleware global
├── route/                         # Layer 4
│   └── route.go                   #   pemetaan URL ke method service (tanpa logika)
├── logs/                          # output log (tidak di-commit)
├── migrations/
│   └── 001_create_students.sql    # DDL: buat tabel students + index
├── main.go                        # Entry point: hanya perakitan (repository -> service -> app)
├── .env / .env.example            # Konfigurasi environment
├── README.md                      # Dokumentasi utama proyek & petunjuk instalasi
└── go.mod / go.sum                # Dependency Go (di root project, bukan di sini)
```

## Persiapan & Menjalankan

### 1. Siapkan database PostgreSQL

Buat database sesuai `.env` (default: `praktikum_backend`), lalu jalankan migration:

```bash
psql -U postgres -d praktikum_backend -f migrations/001_create_students.sql
```

Ini akan membuat tabel `students` beserta:
- Unique index case-insensitive pada `nim`
- Index biasa pada `name` untuk mempercepat pencarian

### 2. Konfigurasi environment

Salin `.env.example` menjadi `.env`, lalu sesuaikan:

| Variable | Default | Keterangan |
|---|---|---|
| `APP_PORT` | `3000` | Port server HTTP |
| `DB_HOST` | `localhost` | Host PostgreSQL |
| `DB_PORT` | `5432` | Port PostgreSQL |
| `DB_USER` | `postgres` | Username DB |
| `DB_PASSWORD` | — | Password DB |
| `DB_NAME` | `praktikum_backend` | Nama database |
| `DB_SSLMODE` | `disable` | Mode SSL koneksi ke Postgres |
| `DB_MAX_CONNS` | `10` | Maksimum koneksi dalam connection pool |

### 3. Jalankan server

Dari root project (folder yang berisi `go.mod`):

```bash
go run ./pert3
```

Server akan jalan di `http://localhost:<APP_PORT>` (default `http://localhost:3000`).

> **Catatan:** Go adalah bahasa *compiled* — setiap ada perubahan kode, server harus **dihentikan lalu dijalankan ulang** (`go run` lagi, atau `go build` + eksekusi ulang binary-nya) agar perubahan ter-apply.

## Format Response

Semua response (sukses maupun gagal) dibungkus dalam struktur seragam (`WebResponse`):

```json
{
  "message": "string",
  "data": { },
  "meta": { },
  "errors": [ ]
}
```

- `data` hanya muncul untuk response sukses yang membawa data.
- `meta` hanya muncul di endpoint list (berisi info pagination).
- `errors` hanya muncul saat validasi gagal (berisi daftar pesan error per field).

## Middleware

- `requestid` — menambahkan Request ID unik ke tiap request (dipakai di access log).
- `logger` — mencatat log akses (`waktu | request id | method | path | status | latency`).
- `cors` — mengizinkan request lintas origin.
- `requireJSON` — middleware khusus grup `/api/v1/students`, mewajibkan header `Content-Type: application/json` untuk method `POST`, `PUT`, `PATCH`. Kalau tidak sesuai → `415 Unsupported Media Type`.

## Endpoints

### Health Check

**`GET /api/v1/health`**

Mengecek server hidup **dan** koneksi ke database (`pool.Ping`).

| Kondisi | Status |
|---|---|
| Server & DB normal | `200 OK` |
| Koneksi ke database gagal/putus | `503 Service Unavailable` |

Contoh sukses:
```json
{ "message": "server dan basis data berjalan", "data": { "database": "connected" } }
```

---

### List Students

**`GET /api/v1/students`**

Mendukung query parameter berikut (semua opsional):

| Param | Default | Keterangan |
|---|---|---|
| `page` | `1` | Halaman ke berapa |
| `limit` | `10` | Jumlah data per halaman |
| `search` | — | Filter nama mahasiswa (case-insensitive, `ILIKE`) |
| `sort` | `id` | Kolom pengurutan. Whitelist: `id`, `name`, `grade` |
| `order` | `asc` | Arah urutan: `asc` atau `desc` |
| `is_active` | — | Filter status aktif (`true`/`false`) |

Response sukses (`200 OK`) menyertakan `meta` untuk pagination:
```json
{
  "message": "daftar mahasiswa berhasil diambil",
  "data": [ { "id": 1, "nim": "123", "name": "Budi", "grade": 80, "is_active": true, "created_at": "..." } ],
  "meta": { "page": 1, "limit": 10, "total_data": 1, "total_pages": 1 }
}
```

---

### Get Student by ID

**`GET /api/v1/students/:id`**

| Kondisi | Status |
|---|---|
| Data ditemukan | `200 OK` |
| `id` bukan angka | `400 Bad Request` |
| Data tidak ada di database | `404 Not Found` |
| Error database lain (query gagal, tabel bermasalah, dll.) | `500 Internal Server Error` |

---

### Create Student

**`POST /api/v1/students`**

Body:
```json
{ "nim": "12345", "name": "Budi Santoso", "grade": 85.5, "is_active": true }
```

| Kondisi | Status |
|---|---|
| Berhasil dibuat | `201 Created` (menyertakan header `Location`) |
| `Content-Type` bukan `application/json` | `415 Unsupported Media Type` |
| Body tidak valid (gagal parse JSON) | `400 Bad Request` |
| `nim`/`name` kosong atau `grade` di luar 0–100 | `422 Unprocessable Entity` |
| `nim` sudah terdaftar | `409 Conflict` |
| Error database lain | `500 Internal Server Error` |

---

### Replace Student (full update)

**`PUT /api/v1/students/:id`**

Body sama seperti Create — **semua field wajib diisi ulang** (replace total).

| Kondisi | Status |
|---|---|
| Berhasil diperbarui | `200 OK` |
| `id` bukan angka | `400 Bad Request` |
| Body tidak valid | `400 Bad Request` |
| Validasi gagal | `422 Unprocessable Entity` |
| Data tidak ditemukan | `404 Not Found` |
| `nim` bentrok dengan mahasiswa lain | `409 Conflict` |
| Error database lain | `500 Internal Server Error` |

---

### Patch Student (partial update)

**`PATCH /api/v1/students/:id`**

Body hanya berisi field yang ingin diubah (semua opsional):
```json
{ "grade": 90 }
```

Field yang tidak dikirim tidak akan diubah. Alurnya: ambil data lama dulu (`FindByID`) → replace field yang dikirim → simpan.

| Kondisi | Status |
|---|---|
| Berhasil diperbarui | `200 OK` |
| `id` bukan angka | `400 Bad Request` |
| Body tidak valid | `400 Bad Request` |
| Data tidak ditemukan | `404 Not Found` |
| `grade` di luar 0–100 | `422 Unprocessable Entity` |
| `nim` bentrok dengan mahasiswa lain | `409 Conflict` |
| Error database lain | `500 Internal Server Error` |

---

### Delete Student

**`DELETE /api/v1/students/:id`**

| Kondisi | Status |
|---|---|
| Berhasil dihapus | `204 No Content` (tanpa body) |
| `id` bukan angka | `400 Bad Request` |
| Data tidak ditemukan | `404 Not Found` |
| Error database lain | `500 Internal Server Error` |

---

### 404 Fallback

Request ke path yang tidak terdaftar akan mendapat:
```json
{ "message": "endpoint tidak ditemukan" }
```
dengan status `404 Not Found`.

## Filosofi Penanganan Error (500 vs 503)

Aplikasi ini sengaja membedakan dua level kegagalan:

- **503 (`/health` saja)** — dipakai khusus untuk mengecek *ketersediaan* dependency (koneksi database). Sifatnya dianggap sementara/transient: begitu database kembali normal, request berikutnya akan berhasil tanpa perubahan kode apa pun.
- **500 (semua endpoint data)** — dipakai untuk semua kegagalan pemrosesan request di sisi server (query salah, tabel/kolom tidak sesuai, dll). Pesan yang dikirim ke client **selalu digeneralisir** ("terjadi kesalahan pada server") agar detail internal (nama tabel, struktur query, jenis database) tidak bocor ke client — detail asli errornya cukup ada di log server untuk kebutuhan debugging developer.

## Pengujian Manual

Testing 500 (query database bermasalah):
```sql
ALTER TABLE students RENAME TO students_backup;  -- lalu hit endpoint /students
ALTER TABLE students_backup RENAME TO students;  -- kembalikan
```

Testing 503 (database down):
```bash
sudo service postgresql stop   # lalu hit /api/v1/health
sudo service postgresql start  # nyalakan lagi
```