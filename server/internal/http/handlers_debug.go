package httpx

import (
	"fmt"
	"net/http"
	"os"
	"strings"
)

// DebugDirectory shows directory loading status
func (s *Server) DebugDirectory(w http.ResponseWriter, r *http.Request) {
	csvPath := os.Getenv("DIRECTORY_CSV_PATH")
	if csvPath == "" {
		csvPath = "data/daftar_pengguna_serumpun.csv"
	}

	dir, err := LoadDirectoryFromCSV(csvPath)
	if err != nil {
		http.Error(w, "failed to load directory: "+err.Error(), http.StatusInternalServerError)
		return
	}

	response := fmt.Sprintf(`Directory Loading Status:
CSV Path: %s
Provinsi Count: %d
Kabkot Count: %d
Bidang List: %v

Sample Provinsi (first 3):
`, csvPath, len(dir.Provinsi), len(dir.Kabkot), dir.BidangList)

	for i, row := range dir.Provinsi {
		if i >= 3 {
			break
		}
		response += fmt.Sprintf("  - %s (%s) - %s - %s\n", row.Nama, row.Email, row.Jabatan, row.Bidang)
	}

	response += "\nSample Kabkot (first 3):\n"
	for i, row := range dir.Kabkot {
		if i >= 3 {
			break
		}
		response += fmt.Sprintf("  - %s (%s) - %s - %s - %s\n", row.Nama, row.Email, row.Instansi, row.Jabatan, row.Bidang)
	}

	// Test building cases
	response += "\n\nTest buildNamaCases (first 100 chars):\n"
	namaCases := buildNamaCases(dir.Provinsi)
	if len(namaCases) > 100 {
		response += namaCases[:100] + "..."
	} else {
		response += namaCases
	}

	response += "\n\nTest buildEmailList (first 100 chars):\n"
	emails := buildEmailList(dir.Provinsi)
	if len(emails) > 100 {
		response += emails[:100] + "..."
	} else {
		response += emails
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(response))
}

// DebugSQL shows generated SQL for KPI Provinsi
func (s *Server) DebugSQL(w http.ResponseWriter, r *http.Request) {
	csvPath := os.Getenv("DIRECTORY_CSV_PATH")
	if csvPath == "" {
		csvPath = "data/daftar_pengguna_serumpun.csv"
	}

	dir, err := LoadDirectoryFromCSV(csvPath)
	if err != nil {
		http.Error(w, "failed to load directory: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Load SQL template
	sqlTemplate, err := s.Queries.Load("kpi_provinsi.sql")
	if err != nil {
		http.Error(w, "failed to load sql template: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Build dynamic parts
	namaCases := buildNamaCases(dir.Provinsi)
	bidangCases := buildBidangCases(dir.Provinsi)
	jabatanCases := buildJabatanCases(dir.Provinsi)
	emails := buildEmailList(dir.Provinsi)
	whereClause := buildAdditionalWhere(map[string]string{})

	// Replace placeholders
	sql := buildDynamicSQL(sqlTemplate, map[string]string{
		"{{NAMA_CASES}}":    namaCases,
		"{{BIDANG_CASES}}":  bidangCases,
		"{{JABATAN_CASES}}": jabatanCases,
		"{{EMAILS}}":        emails,
		"{{WHERE_CLAUSE}}":  whereClause,
	})

	response := fmt.Sprintf(`Generated SQL for KPI Provinsi:
Provinsi Count: %d
namaCases lines: %d
bidangCases lines: %d
jabatanCases lines: %d
emails count: %d

Full SQL:
%s
`, len(dir.Provinsi),
		len(strings.Split(namaCases, "\n")),
		len(strings.Split(bidangCases, "\n")),
		len(strings.Split(jabatanCases, "\n")),
		len(strings.Split(emails, ",")),
		sql)

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(response))
}
