package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/Adityoexs/contract-employee/backend/internal/model"
)

const dateLayout = "2006-01-02"

// KaryawanRepository handles DB operations for karyawan_kontrak.
type KaryawanRepository struct {
	db *sql.DB
}

// NewKaryawanRepository returns a new repository.
func NewKaryawanRepository(db *sql.DB) *KaryawanRepository {
	return &KaryawanRepository{db: db}
}

func GenerateKode(next int) string {
	return fmt.Sprintf("KAR-%03d", next)
}

// FindAll returns all karyawan records.
func (r *KaryawanRepository) FindAll() ([]model.Karyawan, error) {
	rows, err := r.db.Query(
		`SELECT id, kode, nama, tanggal_mulai, tanggal_habis, created_at, updated_at
		 FROM public.karyawan_kontrak ORDER BY id ASC`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []model.Karyawan
	for rows.Next() {
		var k model.Karyawan
		if err := rows.Scan(&k.ID, &k.Kode, &k.Nama, &k.TanggalMulai, &k.TanggalHabis, &k.CreatedAt, &k.UpdatedAt); err != nil {
			return nil, err
		}
		list = append(list, k)
	}
	return list, rows.Err()
}

// FindByID returns a single karyawan or sql.ErrNoRows.
func (r *KaryawanRepository) FindByID(id int) (*model.Karyawan, error) {
	var k model.Karyawan
	err := r.db.QueryRow(
		`SELECT id, kode, nama, tanggal_mulai, tanggal_habis, created_at, updated_at
		 FROM public.karyawan_kontrak WHERE id = $1`,
		id,
	).Scan(&k.ID, &k.Kode, &k.Nama, &k.TanggalMulai, &k.TanggalHabis, &k.CreatedAt, &k.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &k, nil
}

// KodeExists returns true if a kode is already used by another record.
func (r *KaryawanRepository) KodeExists(kode string, excludeID int) (bool, error) {
	var count int
	err := r.db.QueryRow(
		`SELECT COUNT(*) FROM public.karyawan_kontrak WHERE kode = $1 AND id <> $2`,
		strings.TrimSpace(kode), excludeID,
	).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *KaryawanRepository) generateKode(tx *sql.Tx) (string, error) {
	var lastID sql.NullInt64
	err := tx.QueryRow(`SELECT COALESCE(MAX(id), 0) FROM public.karyawan_kontrak`).Scan(&lastID)
	if err != nil {
		return "", err
	}

	next := int(lastID.Int64) + 1
	return fmt.Sprintf("KAR-%03d", next), nil
}

// Create inserts a new karyawan and returns it with generated id.
func (r *KaryawanRepository) Create(req model.KaryawanCreateRequest) (*model.Karyawan, error) {
	tm, err := time.Parse(dateLayout, req.TanggalMulai)
	if err != nil {
		return nil, errors.New("tanggal_mulai format tidak valid")
	}
	th, err := time.Parse(dateLayout, req.TanggalHabis)
	if err != nil {
		return nil, errors.New("tanggal_habis format tidak valid")
	}

	tx, err := r.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	kode, err := r.generateKode(tx)
	if err != nil {
		return nil, err
	}

	var k model.Karyawan
	err = tx.QueryRow(
		`INSERT INTO public.karyawan_kontrak (kode, nama, tanggal_mulai, tanggal_habis)
		 VALUES ($1, $2, $3, $4)
		 RETURNING id, kode, nama, tanggal_mulai, tanggal_habis, created_at, updated_at`,
		kode, strings.TrimSpace(req.Nama), tm, th,
	).Scan(&k.ID, &k.Kode, &k.Nama, &k.TanggalMulai, &k.TanggalHabis, &k.CreatedAt, &k.UpdatedAt)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return &k, nil
}

// Update modifies an existing karyawan and returns the updated record.
func (r *KaryawanRepository) Update(id int, req model.KaryawanUpdateRequest) (*model.Karyawan, error) {
	tm, err := time.Parse(dateLayout, req.TanggalMulai)
	if err != nil {
		return nil, errors.New("tanggal_mulai format tidak valid")
	}
	th, err := time.Parse(dateLayout, req.TanggalHabis)
	if err != nil {
		return nil, errors.New("tanggal_habis format tidak valid")
	}

	var k model.Karyawan
	err = r.db.QueryRow(
		`UPDATE public.karyawan_kontrak
		 SET kode=$1, nama=$2, tanggal_mulai=$3, tanggal_habis=$4, updated_at=NOW()
		 WHERE id=$5
		 RETURNING id, kode, nama, tanggal_mulai, tanggal_habis, created_at, updated_at`,
		strings.TrimSpace(req.Kode), strings.TrimSpace(req.Nama), tm, th, id,
	).Scan(&k.ID, &k.Kode, &k.Nama, &k.TanggalMulai, &k.TanggalHabis, &k.CreatedAt, &k.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &k, nil
}

// Delete removes a karyawan by id.
func (r *KaryawanRepository) Delete(id int) error {
	res, err := r.db.Exec(`DELETE FROM public.karyawan_kontrak WHERE id = $1`, id)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return sql.ErrNoRows
	}
	return nil
}

// FindMaxKodeSequence returns the highest numeric suffix among kodes matching
// the pattern "KAR-NNN", or 0 if there are none.
func (r *KaryawanRepository) FindMaxKodeSequence() (int, error) {
	var maxNum int
	err := r.db.QueryRow(
		`SELECT COALESCE(MAX(CAST(SUBSTRING(kode FROM 'KAR-(\d+)') AS INTEGER)), 0)
		 FROM public.karyawan_kontrak
		 WHERE kode ~ '^KAR-\d+$'`,
	).Scan(&maxNum)
	return maxNum, err
}

// BulkCreate inserts multiple karyawan rows and returns the number inserted.
// Kodes are pre-generated by the caller.
func (r *KaryawanRepository) BulkCreate(rows []model.ImportRow, startSeq int) (int, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return 0, err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	stmt, err := tx.Prepare(
		`INSERT INTO public.karyawan_kontrak (kode, nama, tanggal_mulai, tanggal_habis)
		 VALUES ($1, $2, $3, $4)`,
	)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	for i, row := range rows {
		kode := fmt.Sprintf("KAR-%03d", startSeq+i)
		tm, _ := time.Parse(dateLayout, row.TanggalMulai)
		th, _ := time.Parse(dateLayout, row.TanggalHabis)
		if _, err = stmt.Exec(kode, strings.TrimSpace(row.Nama), tm, th); err != nil {
			return 0, err
		}
	}

	err = tx.Commit()
	return len(rows), err
}
