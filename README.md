# Contract Employee — Karyawan Kontrak

Aplikasi manajemen karyawan kontrak dengan:
- **Frontend:** React + TypeScript (Vite)
- **Backend:** Golang (Gin)
- **Database:** PostgreSQL

---

## Fitur

- Lihat daftar karyawan kontrak
- Tambah karyawan baru
- Edit data karyawan
- Detail karyawan
- Hapus karyawan
- Validasi input di frontend dan backend
- Response API konsisten (`message`, `data`, `errors`)

---

## Struktur Project

```
contract-employee/
├── backend/          # Golang REST API
│   ├── cmd/
│   │   └── main.go
│   ├── internal/
│   │   ├── db/
│   │   ├── handler/
│   │   ├── model/
│   │   └── repository/
│   ├── .env.example
│   └── go.mod
├── database/
│   └── schema.sql    # DDL PostgreSQL
└── frontend/         # React + TypeScript
    ├── src/
    │   ├── components/
    │   ├── pages/
    │   ├── services/
    │   └── types/
    └── package.json
```

---

## Setup Database (PostgreSQL)

1. Buat database:
   ```sql
   CREATE DATABASE contract_employee;
   ```

2. Jalankan schema:
   ```bash
   psql -U postgres -d contract_employee -f database/schema.sql
   ```

---

## Menjalankan Backend (Golang)

1. Masuk ke folder backend:
   ```bash
   cd backend
   ```

2. Salin dan sesuaikan konfigurasi:
   ```bash
   cp .env.example .env
   # Edit .env sesuai konfigurasi database Anda
   ```

3. Jalankan server:
   ```bash
   go run cmd/main.go
   ```

   Server berjalan di `http://localhost:8080`

### REST API Endpoints

| Method | Endpoint              | Keterangan           |
|--------|-----------------------|----------------------|
| GET    | /api/karyawan         | Daftar semua karyawan |
| GET    | /api/karyawan/:id     | Detail karyawan       |
| POST   | /api/karyawan         | Tambah karyawan baru  |
| PUT    | /api/karyawan/:id     | Update karyawan       |
| DELETE | /api/karyawan/:id     | Hapus karyawan        |

### Contoh Request & Response

**POST /api/karyawan**
```json
{
  "kode": "EMP001",
  "nama": "Budi Santoso",
  "tanggal_mulai": "2024-01-01",
  "tanggal_habis": "2025-01-01"
}
```

**Response sukses (201):**
```json
{
  "message": "Data berhasil ditambahkan",
  "data": {
    "id": 1,
    "kode": "EMP001",
    "nama": "Budi Santoso",
    "tanggal_mulai": "2024-01-01T00:00:00Z",
    "tanggal_habis": "2025-01-01T00:00:00Z",
    "created_at": "...",
    "updated_at": "..."
  }
}
```

**Response validasi gagal (422):**
```json
{
  "message": "Validasi gagal",
  "errors": {
    "kode": "Kode sudah digunakan, gunakan kode yang lain",
    "tanggal_habis": "Tanggal habis harus lebih besar atau sama dengan tanggal mulai"
  }
}
```

---

## Menjalankan Frontend (React)

1. Masuk ke folder frontend:
   ```bash
   cd frontend
   ```

2. Install dependencies (jika belum):
   ```bash
   npm install
   ```

3. Salin konfigurasi:
   ```bash
   cp .env.example .env
   # VITE_API_URL=http://localhost:8080/api
   ```

4. Jalankan dev server:
   ```bash
   npm run dev
   ```

   Aplikasi berjalan di `http://localhost:5173`

---

## Konfigurasi Backend (.env)

```env
PORT=8080

DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=contract_employee
DB_SSLMODE=disable
```

---

## Validasi Input

| Field           | Aturan                                                    |
|-----------------|-----------------------------------------------------------|
| `kode`          | Wajib, tidak boleh kosong/spasi, harus unik               |
| `nama`          | Wajib, tidak boleh kosong/spasi                           |
| `tanggal_mulai` | Wajib, format YYYY-MM-DD                                  |
| `tanggal_habis` | Wajib, format YYYY-MM-DD, >= tanggal_mulai                |

Validasi dilakukan di **frontend** (Zod + react-hook-form) dan **backend** (Go) sebagai sumber kebenaran utama.
