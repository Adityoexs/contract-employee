package handler

import (
	"database/sql"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/Adityoexs/contract-employee/backend/internal/model"
	"github.com/Adityoexs/contract-employee/backend/internal/repository"
	"github.com/gin-gonic/gin"
)

const dateLayout = "2006-01-02"

// KaryawanHandler holds the repository dependency.
type KaryawanHandler struct {
	repo *repository.KaryawanRepository
}

// NewKaryawanHandler creates a new handler.
func NewKaryawanHandler(repo *repository.KaryawanRepository) *KaryawanHandler {
	return &KaryawanHandler{repo: repo}
}

// apiResponse is the standard JSON envelope.
type apiResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Errors  interface{} `json:"errors,omitempty"`
}

// ─── helpers ──────────────────────────────────────────────────────────────────

func respond(c *gin.Context, status int, message string, data interface{}, errs interface{}) {
	c.JSON(status, apiResponse{Message: message, Data: data, Errors: errs})
}

// validateRequest validates the incoming KaryawanRequest and returns a map of
// field -> error message. An empty map means the request is valid.
func validateRequest(req model.KaryawanRequest, excludeID int, repo *repository.KaryawanRepository) map[string]string {
	errs := make(map[string]string)

	if strings.TrimSpace(req.Kode) == "" {
		errs["kode"] = "Kode wajib diisi dan tidak boleh kosong/spasi"
	} else {
		exists, err := repo.KodeExists(req.Kode, excludeID)
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

	// Cross-field validation
	if errTm == nil && errTh == nil && !tm.IsZero() && !th.IsZero() {
		if th.Before(tm) {
			errs["tanggal_habis"] = "Tanggal habis harus lebih besar atau sama dengan tanggal mulai"
		}
	}

	return errs
}

// ─── handlers ────────────────────────────────────────────────────────────────

// GetAll godoc – GET /api/karyawan
func (h *KaryawanHandler) GetAll(c *gin.Context) {
	list, err := h.repo.FindAll()
	if err != nil {
		respond(c, http.StatusInternalServerError, "Gagal mengambil data", nil, err.Error())
		return
	}
	if list == nil {
		list = []model.Karyawan{}
	}
	respond(c, http.StatusOK, "Berhasil", list, nil)
}

// GetByID godoc – GET /api/karyawan/:id
func (h *KaryawanHandler) GetByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		respond(c, http.StatusBadRequest, "ID tidak valid", nil, nil)
		return
	}
	k, err := h.repo.FindByID(id)
	if errors.Is(err, sql.ErrNoRows) {
		respond(c, http.StatusNotFound, "Data tidak ditemukan", nil, nil)
		return
	}
	if err != nil {
		respond(c, http.StatusInternalServerError, "Gagal mengambil data", nil, err.Error())
		return
	}
	respond(c, http.StatusOK, "Berhasil", k, nil)
}

// Create godoc – POST /api/karyawan
func (h *KaryawanHandler) Create(c *gin.Context) {
	var req model.KaryawanRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respond(c, http.StatusBadRequest, "Request tidak valid", nil, err.Error())
		return
	}

	if errs := validateRequest(req, 0, h.repo); len(errs) > 0 {
		respond(c, http.StatusUnprocessableEntity, "Validasi gagal", nil, errs)
		return
	}

	k, err := h.repo.Create(req)
	if err != nil {
		respond(c, http.StatusInternalServerError, "Gagal menyimpan data", nil, err.Error())
		return
	}
	respond(c, http.StatusCreated, "Data berhasil ditambahkan", k, nil)
}

// Update godoc – PUT /api/karyawan/:id
func (h *KaryawanHandler) Update(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		respond(c, http.StatusBadRequest, "ID tidak valid", nil, nil)
		return
	}

	var req model.KaryawanRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respond(c, http.StatusBadRequest, "Request tidak valid", nil, err.Error())
		return
	}

	if errs := validateRequest(req, id, h.repo); len(errs) > 0 {
		respond(c, http.StatusUnprocessableEntity, "Validasi gagal", nil, errs)
		return
	}

	k, err := h.repo.Update(id, req)
	if errors.Is(err, sql.ErrNoRows) {
		respond(c, http.StatusNotFound, "Data tidak ditemukan", nil, nil)
		return
	}
	if err != nil {
		respond(c, http.StatusInternalServerError, "Gagal mengupdate data", nil, err.Error())
		return
	}
	respond(c, http.StatusOK, "Data berhasil diupdate", k, nil)
}

// Delete godoc – DELETE /api/karyawan/:id
func (h *KaryawanHandler) Delete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		respond(c, http.StatusBadRequest, "ID tidak valid", nil, nil)
		return
	}

	err = h.repo.Delete(id)
	if errors.Is(err, sql.ErrNoRows) {
		respond(c, http.StatusNotFound, "Data tidak ditemukan", nil, nil)
		return
	}
	if err != nil {
		respond(c, http.StatusInternalServerError, "Gagal menghapus data", nil, err.Error())
		return
	}
	respond(c, http.StatusOK, "Data berhasil dihapus", nil, nil)
}
