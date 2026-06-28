package handler

import (
	"bytes"
	"fmt"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/Adityoexs/contract-employee/backend/internal/model"
	"github.com/extrame/xls"
	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
)

const (
	maxImportSize = 500 * 1024 // 500 KB
	importSheet   = "Karyawan"
	dateLayout    = "2006-01-02"
)

// ImportExcel godoc – POST /api/karyawan/import
// Accepts a multipart/form-data upload with field "file" (.xlsx or .xls, ≤ 500 KB).
// Validates all rows first; if any are invalid the whole file is rejected.
// On success all rows are inserted with auto-generated KAR-NNN kodes.
func (h *KaryawanHandler) ImportExcel(c *gin.Context) {
	fh, err := c.FormFile("file")
	if err != nil {
		respond(c, http.StatusBadRequest, "File wajib diunggah", nil, nil)
		return
	}

	if fh.Size > maxImportSize {
		respond(c, http.StatusRequestEntityTooLarge,
			fmt.Sprintf("Ukuran file melebihi batas %d KB", maxImportSize/1024), nil, nil)
		return
	}

	ext := strings.ToLower(filepath.Ext(fh.Filename))
	if ext != ".xlsx" && ext != ".xls" {
		respond(c, http.StatusBadRequest, "Format file tidak didukung, gunakan .xlsx atau .xls", nil, nil)
		return
	}

	f, err := fh.Open()
	if err != nil {
		respond(c, http.StatusInternalServerError, "Gagal membuka file", nil, nil)
		return
	}
	defer f.Close()

	buf := new(bytes.Buffer)
	if _, err := buf.ReadFrom(f); err != nil {
		respond(c, http.StatusInternalServerError, "Gagal membaca file", nil, nil)
		return
	}

	var (
		rows       []model.ImportRow
		parseErrs  []model.ImportError
		parseErr   error
	)

	if ext == ".xlsx" {
		rows, parseErrs, parseErr = parseXlsx(buf.Bytes())
	} else {
		rows, parseErrs, parseErr = parseXls(buf.Bytes())
	}

	if parseErr != nil {
		respond(c, http.StatusUnprocessableEntity, "Gagal memproses file Excel", nil, parseErr.Error())
		return
	}

	if len(parseErrs) > 0 {
		respond(c, http.StatusUnprocessableEntity, "Validasi file gagal, periksa isi file", nil, parseErrs)
		return
	}

	if len(rows) == 0 {
		respond(c, http.StatusUnprocessableEntity, "File tidak memiliki data", nil, nil)
		return
	}

	startSeq, err := h.svc.FindMaxKodeSequence()
	if err != nil {
		respond(c, http.StatusInternalServerError, "Gagal membaca data kode", nil, err.Error())
		return
	}
	startSeq++ // first new kode is max+1

	inserted, err := h.svc.BulkCreate(rows, startSeq)
	if err != nil {
		respond(c, http.StatusInternalServerError, "Gagal menyimpan data", nil, err.Error())
		return
	}

	respond(c, http.StatusOK, "Import berhasil", map[string]int{"inserted": inserted}, nil)
}

// DownloadTemplate godoc – GET /api/karyawan/import/template
// Returns an Excel (.xlsx) template with required column headers and one example row.
func (h *KaryawanHandler) DownloadTemplate(c *gin.Context) {
	xl := excelize.NewFile()
	defer xl.Close()

	sheet := importSheet
	idx, err := xl.NewSheet(sheet)
	if err != nil {
		respond(c, http.StatusInternalServerError, "Gagal membuat template", nil, err.Error())
		return
	}
	xl.SetActiveSheet(idx)
	_ = xl.DeleteSheet("Sheet1")

	headers := []string{"nama", "tanggal_mulai", "tanggal_habis"}
	for i, v := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		xl.SetCellValue(sheet, cell, v)
	}

	example := []string{"Budi Santoso", "2024-01-01", "2025-01-01"}
	for i, v := range example {
		cell, _ := excelize.CoordinatesToCellName(i+1, 2)
		xl.SetCellValue(sheet, cell, v)
	}

	buf, err := xl.WriteToBuffer()
	if err != nil {
		respond(c, http.StatusInternalServerError, "Gagal membuat template", nil, err.Error())
		return
	}

	c.Header("Content-Disposition", "attachment; filename=template_karyawan.xlsx")
	c.Data(http.StatusOK,
		"application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
		buf.Bytes())
}

// ─── parsers ──────────────────────────────────────────────────────────────────

func parseXlsx(data []byte) ([]model.ImportRow, []model.ImportError, error) {
	xl, err := excelize.OpenReader(bytes.NewReader(data))
	if err != nil {
		return nil, nil, fmt.Errorf("file tidak dapat dibaca sebagai xlsx: %w", err)
	}
	defer xl.Close()

	sheets := xl.GetSheetList()
	if len(sheets) == 0 {
		return nil, nil, fmt.Errorf("file tidak memiliki sheet")
	}

	allRows, err := xl.GetRows(sheets[0])
	if err != nil {
		return nil, nil, err
	}

	rows, errs := validateImportRows(allRows)
	return rows, errs, nil
}

func parseXls(data []byte) ([]model.ImportRow, []model.ImportError, error) {
	wb, err := xls.OpenReader(bytes.NewReader(data), "utf-8")
	if err != nil {
		return nil, nil, fmt.Errorf("file tidak dapat dibaca sebagai xls: %w", err)
	}

	sheet := wb.GetSheet(0)
	if sheet == nil {
		return nil, nil, fmt.Errorf("file tidak memiliki sheet")
	}

	var allRows [][]string
	for rowIdx := 0; rowIdx <= int(sheet.MaxRow); rowIdx++ {
		row := sheet.Row(rowIdx)
		cols := make([]string, row.LastCol())
		for colIdx := 0; colIdx < row.LastCol(); colIdx++ {
			cols[colIdx] = row.Col(colIdx)
		}
		allRows = append(allRows, cols)
	}

	rows, errs := validateImportRows(allRows)
	return rows, errs, nil
}

// validateImportRows locates the header row, then validates every data row.
// Completely empty rows are silently skipped.
// If any data row is invalid, its errors are collected and the whole batch is rejected.
func validateImportRows(allRows [][]string) ([]model.ImportRow, []model.ImportError) {
	var (
		rows   []model.ImportRow
		errors []model.ImportError
	)

	if len(allRows) == 0 {
		return rows, errors
	}

	// Find header row.
	headerIdx := -1
	namaCol, mulaiCol, habisCol := -1, -1, -1
	for ri, row := range allRows {
		for ci, cell := range row {
			switch strings.ToLower(strings.TrimSpace(cell)) {
			case "nama":
				namaCol = ci
			case "tanggal_mulai":
				mulaiCol = ci
			case "tanggal_habis":
				habisCol = ci
			}
		}
		if namaCol >= 0 && mulaiCol >= 0 && habisCol >= 0 {
			headerIdx = ri
			break
		}
	}

	if headerIdx < 0 {
		errors = append(errors, model.ImportError{
			Row:     1,
			Field:   "header",
			Message: "Kolom wajib tidak ditemukan: nama, tanggal_mulai, tanggal_habis",
		})
		return rows, errors
	}

	getCell := func(row []string, col int) string {
		if col < len(row) {
			return strings.TrimSpace(row[col])
		}
		return ""
	}

	for ri, row := range allRows[headerIdx+1:] {
		rowNum := headerIdx + ri + 2 // 1-based Excel row number

		nama := getCell(row, namaCol)
		mulai := getCell(row, mulaiCol)
		habis := getCell(row, habisCol)

		// Skip truly blank rows.
		if nama == "" && mulai == "" && habis == "" {
			continue
		}

		rowOK := true

		if nama == "" {
			errors = append(errors, model.ImportError{Row: rowNum, Field: "nama", Message: "Nama wajib diisi"})
			rowOK = false
		}

		var tm, th time.Time
		var errTm, errTh error

		if mulai == "" {
			errors = append(errors, model.ImportError{Row: rowNum, Field: "tanggal_mulai", Message: "Tanggal mulai wajib diisi"})
			rowOK = false
		} else {
			tm, errTm = time.Parse(dateLayout, mulai)
			if errTm != nil {
				errors = append(errors, model.ImportError{Row: rowNum, Field: "tanggal_mulai", Message: "Format tanggal mulai tidak valid (gunakan YYYY-MM-DD)"})
				rowOK = false
			}
		}

		if habis == "" {
			errors = append(errors, model.ImportError{Row: rowNum, Field: "tanggal_habis", Message: "Tanggal habis wajib diisi"})
			rowOK = false
		} else {
			th, errTh = time.Parse(dateLayout, habis)
			if errTh != nil {
				errors = append(errors, model.ImportError{Row: rowNum, Field: "tanggal_habis", Message: "Format tanggal habis tidak valid (gunakan YYYY-MM-DD)"})
				rowOK = false
			}
		}

		if rowOK && errTm == nil && errTh == nil && th.Before(tm) {
			errors = append(errors, model.ImportError{Row: rowNum, Field: "tanggal_habis", Message: "Tanggal habis harus lebih besar atau sama dengan tanggal mulai"})
			rowOK = false
		}

		if rowOK {
			rows = append(rows, model.ImportRow{Nama: nama, TanggalMulai: mulai, TanggalHabis: habis})
		}
	}

	return rows, errors
}
