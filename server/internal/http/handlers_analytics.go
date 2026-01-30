package httpx

import (
	"context"
	"net/http"
	"os"
	"time"
)

// IssuesDetailTemplate handles Issues Detail using SQL template
func (s *Server) IssuesDetailTemplate(w http.ResponseWriter, r *http.Request) {
	// Parse query parameters for filtering
	filters := map[string]string{
		"scope":    r.URL.Query().Get("scope"),
		"kab_kota": r.URL.Query().Get("kab_kota"),
		"bidang":   r.URL.Query().Get("bidang"),
		"status":   r.URL.Query().Get("status"),
	}

	// Build cache key
	cacheKey := buildCacheKey("issues_detail", filters)

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

	// Load SQL template
	sqlTemplate, err := s.Queries.Load("issues_detail.sql")
	if err != nil {
		http.Error(w, "failed to load sql template: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Build dynamic parts
	allRows := append(dir.Provinsi, dir.Kabkot...)
	namaCases := buildNamaCases(allRows)
	scopeCases := buildScopeCases(allRows)
	whereClause := buildWhereClause(filters)

	// Replace placeholders
	sql := buildDynamicSQL(sqlTemplate, map[string]string{
		"{{NAMA_CASES}}":   namaCases,
		"{{SCOPE_CASES}}":  scopeCases,
		"{{WHERE_CLAUSE}}": whereClause,
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

// TimelineTemplate handles Timeline using SQL template
func (s *Server) TimelineTemplate(w http.ResponseWriter, r *http.Request) {
	// Parse query parameters for filtering
	filters := map[string]string{
		"scope":    r.URL.Query().Get("scope"),
		"kab_kota": r.URL.Query().Get("kab_kota"),
		"bidang":   r.URL.Query().Get("bidang"),
		"status":   r.URL.Query().Get("status"),
	}

	// Build cache key
	cacheKey := buildCacheKey("timeline", filters)

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

	// Load SQL template
	sqlTemplate, err := s.Queries.Load("timeline.sql")
	if err != nil {
		http.Error(w, "failed to load sql template: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Build dynamic parts
	allRows := append(dir.Provinsi, dir.Kabkot...)
	namaCases := buildNamaCases(allRows)
	scopeCases := buildScopeCases(allRows)
	whereClause := buildWhereClause(filters)

	// Replace placeholders
	sql := buildDynamicSQL(sqlTemplate, map[string]string{
		"{{NAMA_CASES}}":   namaCases,
		"{{SCOPE_CASES}}":  scopeCases,
		"{{WHERE_CLAUSE}}": whereClause,
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

// LeaderboardTemplate handles Leaderboard using SQL template
func (s *Server) LeaderboardTemplate(w http.ResponseWriter, r *http.Request) {
	// Parse query parameters for filtering
	filters := map[string]string{
		"scope":  r.URL.Query().Get("scope"),
		"bidang": r.URL.Query().Get("bidang"),
	}

	// Build cache key
	cacheKey := buildCacheKey("leaderboard", filters)

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

	// Load SQL template
	sqlTemplate, err := s.Queries.Load("leaderboard.sql")
	if err != nil {
		http.Error(w, "failed to load sql template: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Build dynamic parts
	allRows := append(dir.Provinsi, dir.Kabkot...)
	namaCases := buildNamaCases(allRows)
	scopeCases := buildScopeCases(allRows)
	instansiCases := buildInstansiCases(allRows)
	whereClause := buildWhereClause(filters)

	// Replace placeholders
	sql := buildDynamicSQL(sqlTemplate, map[string]string{
		"{{NAMA_CASES}}":     namaCases,
		"{{SCOPE_CASES}}":    scopeCases,
		"{{INSTANSI_CASES}}": instansiCases,
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

// WorkloadTemplate handles Workload using SQL template
func (s *Server) WorkloadTemplate(w http.ResponseWriter, r *http.Request) {
	// Parse query parameters for filtering
	filters := map[string]string{
		"scope":  r.URL.Query().Get("scope"),
		"bidang": r.URL.Query().Get("bidang"),
	}

	// Build cache key
	cacheKey := buildCacheKey("workload", filters)

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

	// Load SQL template
	sqlTemplate, err := s.Queries.Load("workload.sql")
	if err != nil {
		http.Error(w, "failed to load sql template: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Build dynamic parts
	allRows := append(dir.Provinsi, dir.Kabkot...)
	namaCases := buildNamaCases(allRows)
	scopeCases := buildScopeCases(allRows)
	instansiCases := buildInstansiCases(allRows)
	whereClause := buildWhereClause(filters)

	// Replace placeholders
	sql := buildDynamicSQL(sqlTemplate, map[string]string{
		"{{NAMA_CASES}}":     namaCases,
		"{{SCOPE_CASES}}":    scopeCases,
		"{{INSTANSI_CASES}}": instansiCases,
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
