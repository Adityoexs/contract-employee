package handler

import (
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/Adityoexs/contract-employee/backend/internal/model"
	"github.com/Adityoexs/contract-employee/backend/internal/service"
	"github.com/gin-gonic/gin"
)

// KaryawanHandler handles HTTP requests for the karyawan domain.
type KaryawanHandler struct {
	svc service.KaryawanService
}

// NewKaryawanHandler creates a new handler backed by the given service.
func NewKaryawanHandler(svc service.KaryawanService) *KaryawanHandler {
	return &KaryawanHandler{svc: svc}
}

// apiResponse is the standard JSON envelope.
type apiResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Errors  interface{} `json:"errors,omitempty"`
}

func respond(c *gin.Context, status int, message string, data interface{}, errs interface{}) {
	c.JSON(status, apiResponse{Message: message, Data: data, Errors: errs})
}

// ─── handlers ────────────────────────────────────────────────────────────────

// GetAll godoc – GET /api/karyawan
func (h *KaryawanHandler) GetAll(c *gin.Context) {
	log.Println("GetAll hit")
	list, err := h.svc.GetAll()
	if err != nil {
		log.Printf("GetAll error: %#v\n", err)
		respond(c, http.StatusInternalServerError, "Gagal mengambil data", nil, err.Error())
		return
	}
	log.Printf("GetAll success: %d rows\n", len(list))
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
	k, err := h.svc.GetByID(id)
	if service.IsNotFound(err) {
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
	var req model.KaryawanCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respond(c, http.StatusBadRequest, "Request tidak valid", nil, err.Error())
		return
	}

	k, err := h.svc.Create(req)
	if err != nil {
		var ve *service.ValidationError
		if errors.As(err, &ve) {
			respond(c, http.StatusUnprocessableEntity, "Validasi gagal", nil, ve.Fields)
			return
		}
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

	var req model.KaryawanUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respond(c, http.StatusBadRequest, "Request tidak valid", nil, err.Error())
		return
	}

	k, err := h.svc.Update(id, req)
	if err != nil {
		var ve *service.ValidationError
		if errors.As(err, &ve) {
			respond(c, http.StatusUnprocessableEntity, "Validasi gagal", nil, ve.Fields)
			return
		}
		if service.IsNotFound(err) {
			respond(c, http.StatusNotFound, "Data tidak ditemukan", nil, nil)
			return
		}
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

	err = h.svc.Delete(id)
	if service.IsNotFound(err) {
		respond(c, http.StatusNotFound, "Data tidak ditemukan", nil, nil)
		return
	}
	if err != nil {
		respond(c, http.StatusInternalServerError, "Gagal menghapus data", nil, err.Error())
		return
	}
	respond(c, http.StatusOK, "Data berhasil dihapus", nil, nil)
}

