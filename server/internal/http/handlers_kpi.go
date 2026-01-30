package httpx

import (
	"context"
	"net/http"
	"os"
	"time"
)

// KPIProvinsiTemplate handles KPI Provinsi using SQL template
func (s *Server) KPIProvinsiTemplate(w http.ResponseWriter, r *http.Request) {
	// Parse query parameters for filtering
	filters := map[string]string{
		"bidang":  r.URL.Query().Get("bidang"),
		"jabatan": r.URL.Query().Get("jabatan"),
	}

	// Build cache key
	cacheKey := buildCacheKey("kpi_provinsi", filters)

	// Check cache
	if b, ok := s.Cache.Get(cacheKey); ok {
		w.Header().Set("Content-Type", "text/csv; charset=utf-8")
		w.Header().Set("Cache-Control", "public, max-age=30")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(b)
		return
	}

	// Load directory from CSV
	csvPath := os.Getenv("DIRECTORY_CSV_PATH")
	if csvPath == "" {
		csvPath = "data/daftar_pengguna_serumpun.csv"
	}

	dir, err := LoadDirectoryFromCSV(csvPath)
	if err != nil {
		http.Error(w, "failed to load directory: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if len(dir.Provinsi) == 0 {
		http.Error(w, "no provinsi staff found in directory", http.StatusInternalServerError)
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
	whereClause := buildAdditionalWhere(filters)

	// Replace placeholders
	sql := buildDynamicSQL(sqlTemplate, map[string]string{
		"{{NAMA_CASES}}":    namaCases,
		"{{BIDANG_CASES}}":  bidangCases,
		"{{JABATAN_CASES}}": jabatanCases,
		"{{EMAILS}}":        emails,
		"{{WHERE_CLAUSE}}":  whereClause,
	})

	// Execute query
	ctx, cancel := context.WithTimeout(r.Context(), 25*time.Second)
	defer cancel()

	b, err := QueryToCSV(ctx, s.DB, sql)
	if err != nil {
		http.Error(w, "query failed: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Set cache
	s.Cache.Set(cacheKey, b)

	// Write response
	w.Header().Set("Content-Type", "text/csv; charset=utf-8")
	w.Header().Set("Cache-Control", "public, max-age=30")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(b)
}

// KPIKabkotTemplate handles KPI Kabkot using SQL template
func (s *Server) KPIKabkotTemplate(w http.ResponseWriter, r *http.Request) {
	// Parse query parameters for filtering
	filters := map[string]string{
		"bidang":   r.URL.Query().Get("bidang"),
		"instansi": r.URL.Query().Get("instansi"),
		"jabatan":  r.URL.Query().Get("jabatan"),
	}

	// Build cache key
	cacheKey := buildCacheKey("kpi_kabkot", filters)

	// Check cache
	if b, ok := s.Cache.Get(cacheKey); ok {
		w.Header().Set("Content-Type", "text/csv; charset=utf-8")
		w.Header().Set("Cache-Control", "public, max-age=30")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(b)
		return
	}

	// Load directory from CSV
	csvPath := os.Getenv("DIRECTORY_CSV_PATH")
	if csvPath == "" {
		csvPath = "data/daftar_pengguna_serumpun.csv"
	}

	dir, err := LoadDirectoryFromCSV(csvPath)
	if err != nil {
		http.Error(w, "failed to load directory: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if len(dir.Kabkot) == 0 {
		http.Error(w, "no ketua kabkot found in directory", http.StatusInternalServerError)
		return
	}

	// Load SQL template
	sqlTemplate, err := s.Queries.Load("kpi_kabkot.sql")
	if err != nil {
		http.Error(w, "failed to load sql template: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Build dynamic parts
	namaCases := buildNamaCases(dir.Kabkot)
	bidangCases := buildBidangCases(dir.Kabkot)
	instansiCases := buildInstansiCases(dir.Kabkot)
	jabatanCases := buildJabatanCases(dir.Kabkot)
	emails := buildEmailList(dir.Kabkot)
	whereClause := buildAdditionalWhere(filters)

	// Replace placeholders
	sql := buildDynamicSQL(sqlTemplate, map[string]string{
		"{{NAMA_CASES}}":     namaCases,
		"{{BIDANG_CASES}}":   bidangCases,
		"{{INSTANSI_CASES}}": instansiCases,
		"{{JABATAN_CASES}}":  jabatanCases,
		"{{EMAILS}}":         emails,
		"{{WHERE_CLAUSE}}":   whereClause,
	})

	// Execute query
	ctx, cancel := context.WithTimeout(r.Context(), 25*time.Second)
	defer cancel()

	b, err := QueryToCSV(ctx, s.DB, sql)
	if err != nil {
		http.Error(w, "query failed: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Set cache
	s.Cache.Set(cacheKey, b)

	// Write response
	w.Header().Set("Content-Type", "text/csv; charset=utf-8")
	w.Header().Set("Cache-Control", "public, max-age=30")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(b)
}

// HeatmapTemplate handles Heatmap using SQL template
func (s *Server) HeatmapTemplate(w http.ResponseWriter, r *http.Request) {
	// Parse query parameters for filtering
	filters := map[string]string{
		"kab_kota": r.URL.Query().Get("kab_kota"),
		"bidang":   r.URL.Query().Get("bidang"),
	}

	// Build cache key
	cacheKey := buildCacheKey("heatmap", filters)

	// Check cache
	if b, ok := s.Cache.Get(cacheKey); ok {
		w.Header().Set("Content-Type", "text/csv; charset=utf-8")
		w.Header().Set("Cache-Control", "public, max-age=30")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(b)
		return
	}

	// Load SQL template
	sqlTemplate, err := s.Queries.Load("heatmap.sql")
	if err != nil {
		http.Error(w, "failed to load sql template: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Build HAVING clause for filters
	havingClause := buildHavingClause(filters)

	// Replace placeholders
	sql := buildDynamicSQL(sqlTemplate, map[string]string{
		"{{HAVING_CLAUSE}}": havingClause,
	})

	// Execute query
	ctx, cancel := context.WithTimeout(r.Context(), 25*time.Second)
	defer cancel()

	b, err := QueryToCSV(ctx, s.DB, sql)
	if err != nil {
		http.Error(w, "query failed: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Set cache
	s.Cache.Set(cacheKey, b)

	// Write response
	w.Header().Set("Content-Type", "text/csv; charset=utf-8")
	w.Header().Set("Cache-Control", "public, max-age=30")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(b)
}
