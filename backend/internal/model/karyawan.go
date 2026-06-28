package model

import "time"

// DateLayout is the canonical date format used across create/update flows.
const DateLayout = "2006-01-02"

// Karyawan represents a contract employee record.
type Karyawan struct {
	ID           int       `db:"id"            json:"id"`
	Kode         string    `db:"kode"          json:"kode"`
	Nama         string    `db:"nama"          json:"nama"`
	TanggalMulai time.Time `db:"tanggal_mulai" json:"tanggal_mulai"`
	TanggalHabis time.Time `db:"tanggal_habis" json:"tanggal_habis"`
	CreatedAt    time.Time `db:"created_at"    json:"created_at"`
	UpdatedAt    time.Time `db:"updated_at"    json:"updated_at"`
}

// KaryawanRequest is the payload for create/update operations.
type KaryawanCreateRequest struct {
	Nama         string `json:"nama"`
	TanggalMulai string `json:"tanggal_mulai"`
	TanggalHabis string `json:"tanggal_habis"`
}

// KaryawanUpdateRequest is used when updating a record.
type KaryawanUpdateRequest struct {
	Kode         string `json:"kode"`
	Nama         string `json:"nama"`
	TanggalMulai string `json:"tanggal_mulai"`
	TanggalHabis string `json:"tanggal_habis"`
}

// ImportRow represents a validated row from an Excel import file.
type ImportRow struct {
	Nama         string `json:"nama"`
	TanggalMulai string `json:"tanggal_mulai"`
	TanggalHabis string `json:"tanggal_habis"`
}

// ImportError represents a validation error for a specific row in an import file.
type ImportError struct {
	Row     int    `json:"row"`
	Field   string `json:"field"`
	Message string `json:"message"`
}
