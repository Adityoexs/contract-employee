package service

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/Adityoexs/contract-employee/backend/internal/model"
	"github.com/Adityoexs/contract-employee/backend/internal/repository"
)

const dateLayout = model.DateLayout

// ValidationError carries field-level validation messages returned by the
// service layer so that handlers can translate them into 422 responses.
type ValidationError struct {
	Fields map[string]string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation error: %v", e.Fields)
}

// KaryawanService defines the business operations for the karyawan domain.
type KaryawanService interface {
	GetAll() ([]model.Karyawan, error)
	GetByID(id int) (*model.Karyawan, error)
	Create(req model.KaryawanCreateRequest) (*model.Karyawan, error)
	Update(id int, req model.KaryawanUpdateRequest) (*model.Karyawan, error)
	Delete(id int) error
	FindMaxKodeSequence() (int, error)
	BulkCreate(rows []model.ImportRow, startSeq int) (int, error)
}

type karyawanService struct {
	repo *repository.KaryawanRepository
}

// NewKaryawanService returns a KaryawanService backed by the given repository.
func NewKaryawanService(repo *repository.KaryawanRepository) KaryawanService {
	return &karyawanService{repo: repo}
}

// ─── read operations ──────────────────────────────────────────────────────────

func (s *karyawanService) GetAll() ([]model.Karyawan, error) {
	return s.repo.FindAll()
}

func (s *karyawanService) GetByID(id int) (*model.Karyawan, error) {
	return s.repo.FindByID(id)
}

// ─── write operations ─────────────────────────────────────────────────────────

// Create validates the request, generates a kode, and persists the record.
func (s *karyawanService) Create(req model.KaryawanCreateRequest) (*model.Karyawan, error) {
	if errs := validateCreateRequest(req); len(errs) > 0 {
		return nil, &ValidationError{Fields: errs}
	}
	return s.repo.Create(req)
}

// Update validates the request (including kode uniqueness) and persists changes.
func (s *karyawanService) Update(id int, req model.KaryawanUpdateRequest) (*model.Karyawan, error) {
	if errs := s.validateUpdateRequest(req, id); len(errs) > 0 {
		return nil, &ValidationError{Fields: errs}
	}
	return s.repo.Update(id, req)
}

// Delete removes the record with the given id.
func (s *karyawanService) Delete(id int) error {
	return s.repo.Delete(id)
}

// FindMaxKodeSequence returns the highest numeric suffix among KAR-NNN kodes.
func (s *karyawanService) FindMaxKodeSequence() (int, error) {
	return s.repo.FindMaxKodeSequence()
}

// BulkCreate inserts multiple rows with pre-assigned kodes starting at startSeq.
func (s *karyawanService) BulkCreate(rows []model.ImportRow, startSeq int) (int, error) {
	return s.repo.BulkCreate(rows, startSeq)
}

// ─── validation helpers ───────────────────────────────────────────────────────

func validateCreateRequest(req model.KaryawanCreateRequest) map[string]string {
	errs := make(map[string]string)

	if strings.TrimSpace(req.Nama) == "" {
		errs["nama"] = "Nama wajib diisi dan tidak boleh kosong/spasi"
	}

	var tm, th time.Time
	var errTm, errTh error

	if strings.TrimSpace(req.TanggalMulai) == "" {
		errs["tanggal_mulai"] = "Tanggal mulai wajib diisi"
	} else {
		tm, errTm = time.Parse(dateLayout, req.TanggalMulai)
		if errTm != nil {
			errs["tanggal_mulai"] = "Format tanggal mulai tidak valid (gunakan YYYY-MM-DD)"
		}
	}

	if strings.TrimSpace(req.TanggalHabis) == "" {
		errs["tanggal_habis"] = "Tanggal habis wajib diisi"
	} else {
		th, errTh = time.Parse(dateLayout, req.TanggalHabis)
		if errTh != nil {
			errs["tanggal_habis"] = "Format tanggal habis tidak valid (gunakan YYYY-MM-DD)"
		}
	}

	if errTm == nil && errTh == nil && !tm.IsZero() && !th.IsZero() {
		if th.Before(tm) {
			errs["tanggal_habis"] = "Tanggal habis harus lebih besar atau sama dengan tanggal mulai"
		}
	}

	return errs
}

func (s *karyawanService) validateUpdateRequest(req model.KaryawanUpdateRequest, excludeID int) map[string]string {
	errs := make(map[string]string)

	if strings.TrimSpace(req.Kode) == "" {
		errs["kode"] = "Kode wajib diisi dan tidak boleh kosong/spasi"
	} else {
		exists, err := s.repo.KodeExists(req.Kode, excludeID)
		if err != nil {
			errs["kode"] = "Gagal memvalidasi kode"
		} else if exists {
			errs["kode"] = "Kode sudah digunakan, gunakan kode yang lain"
		}
	}

	if strings.TrimSpace(req.Nama) == "" {
		errs["nama"] = "Nama wajib diisi dan tidak boleh kosong/spasi"
	}

	var tm, th time.Time
	var errTm, errTh error

	if strings.TrimSpace(req.TanggalMulai) == "" {
		errs["tanggal_mulai"] = "Tanggal mulai wajib diisi"
	} else {
		tm, errTm = time.Parse(dateLayout, req.TanggalMulai)
		if errTm != nil {
			errs["tanggal_mulai"] = "Format tanggal mulai tidak valid (gunakan YYYY-MM-DD)"
		}
	}

	if strings.TrimSpace(req.TanggalHabis) == "" {
		errs["tanggal_habis"] = "Tanggal habis wajib diisi"
	} else {
		th, errTh = time.Parse(dateLayout, req.TanggalHabis)
		if errTh != nil {
			errs["tanggal_habis"] = "Format tanggal habis tidak valid (gunakan YYYY-MM-DD)"
		}
	}

	if errTm == nil && errTh == nil && !tm.IsZero() && !th.IsZero() {
		if th.Before(tm) {
			errs["tanggal_habis"] = "Tanggal habis harus lebih besar atau sama dengan tanggal mulai"
		}
	}

	return errs
}

// IsNotFound returns true if err is sql.ErrNoRows.
func IsNotFound(err error) bool {
	return errors.Is(err, sql.ErrNoRows)
}
