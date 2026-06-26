-- Schema for contract-employee application
-- Database: PostgreSQL

CREATE TABLE IF NOT EXISTS karyawan_kontrak (
    id          SERIAL PRIMARY KEY,
    kode        VARCHAR(20)  NOT NULL UNIQUE,
    nama        VARCHAR(100) NOT NULL,
    tanggal_mulai DATE        NOT NULL,
    tanggal_habis DATE        NOT NULL,
    created_at  TIMESTAMP    NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMP    NOT NULL DEFAULT NOW(),
    CONSTRAINT chk_tanggal CHECK (tanggal_habis >= tanggal_mulai)
);

-- Sample data
INSERT INTO karyawan_kontrak (kode, nama, tanggal_mulai, tanggal_habis) VALUES
    ('EMP001', 'Budi Santoso',  '2024-01-01', '2025-01-01'),
    ('EMP002', 'Ani Rahayu',    '2024-03-01', '2025-03-01'),
    ('EMP003', 'Citra Dewi',    '2024-06-01', '2025-06-01')
ON CONFLICT (kode) DO NOTHING;
