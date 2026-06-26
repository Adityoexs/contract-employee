package model

import "time"

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
type KaryawanRequest struct {
	Kode         string `json:"kode"`
	Nama         string `json:"nama"`
	TanggalMulai string `json:"tanggal_mulai"`
	TanggalHabis string `json:"tanggal_habis"`
}
