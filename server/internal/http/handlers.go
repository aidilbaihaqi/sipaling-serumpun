package httpx

import (
	"context"
	"net/http"
	"os"
	"strings"
	"time"

	"serumpun-data-api/internal/cache"
	"serumpun-data-api/internal/queries"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Server struct {
	DB      *pgxpool.Pool
	Cache   *cache.Cache
	Queries *queries.Store
}

// serveCSV loads SQL from file, runs query (NO REQUIRED ARGS), returns CSV.
// args... is optional: if your SQL uses $1/$2, you can pass them; if not, call without args.
func (s *Server) serveCSV(w http.ResponseWriter, r *http.Request, cacheKey, queryFile string, args ...any) {
	// cache hit
	if b, ok := s.Cache.Get(cacheKey); ok {
		w.Header().Set("Content-Type", "text/csv; charset=utf-8")
		w.Header().Set("Cache-Control", "public, max-age=30")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(b)
		return
	}

	// load SQL
	sql, err := s.Queries.Load(queryFile)
	if err != nil {
		http.Error(w, "failed to load sql: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// run query -> csv
	ctx, cancel := context.WithTimeout(r.Context(), 25*time.Second)
	defer cancel()

	b, err := QueryToCSV(ctx, s.DB, sql, args...)
	if err != nil {
		http.Error(w, "query failed: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// set cache
	s.Cache.Set(cacheKey, b)

	// write response
	w.Header().Set("Content-Type", "text/csv; charset=utf-8")
	w.Header().Set("Cache-Control", "public, max-age=30")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(b)
}

// =========================
// CSV Endpoints
// =========================

// GET /api/v1/kpi.csv
func (s *Server) KPI(w http.ResponseWriter, r *http.Request) {
	// SQL hardcode ID -> no args
	s.serveCSV(w, r, "kpi", "kpi.sql")
}

// GET /api/v1/progress_kabkot.csv
func (s *Server) ProgressKabkot(w http.ResponseWriter, r *http.Request) {
	s.serveCSV(w, r, "progress_kabkot", "progress_kabkot.sql")
}

// GET /api/v1/progress_bidang.csv
func (s *Server) ProgressBidang(w http.ResponseWriter, r *http.Request) {
	s.serveCSV(w, r, "progress_bidang", "progress_bidang.sql")
}

// GET /api/v1/heatmap.csv
func (s *Server) Heatmap(w http.ResponseWriter, r *http.Request) {
	// file name MUST match what you saved
	// if your file is named heatmap_kabkot_bidang.sql, keep it consistent
	s.serveCSV(w, r, "heatmap", "heatmap_kabkot_bidang.sql")
}

// GET /api/v1/issues_detail.csv
func (s *Server) IssuesDetail(w http.ResponseWriter, r *http.Request) {
	s.serveCSV(w, r, "issues_detail", "issues_detail.sql")
}

// GET /api/v1/kpi_provinsi.csv
func (s *Server) KPIProvinsi(w http.ResponseWriter, r *http.Request) {
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

	// Build dynamic SQL with email-bidang-jabatan mapping
	// Only include provinsi staff (Ketua + Anggota Bidang, exclude Sekretariat)
	var emailCases, bidangCases, jabatanCases []string
	var emails []string

	for _, row := range dir.Provinsi {
		email := strings.ToLower(row.Email)
		emails = append(emails, "'"+email+"'")
		emailCases = append(emailCases, "      WHEN LOWER(u.email) = '"+email+"' THEN TRUE")
		bidangCases = append(bidangCases, "      WHEN LOWER(u.email) = '"+email+"' THEN '"+row.Bidang+"'")
		jabatanCases = append(jabatanCases, "      WHEN LOWER(u.email) = '"+email+"' THEN '"+row.Jabatan+"'")
	}

	if len(emails) == 0 {
		http.Error(w, "no provinsi staff found in directory", http.StatusInternalServerError)
		return
	}

	// Build dynamic SQL
	sql := `
-- KPI Provinsi: Dynamic from CSV
WITH 
directory AS (
  SELECT 
    u.id AS user_id,
    LOWER(u.email) AS email,
    COALESCE(
      u.display_name,
      NULLIF(TRIM(COALESCE(u.first_name,'') || ' ' || COALESCE(u.last_name,'')), ''),
      '-'
    ) AS nama,
    'BPS Provinsi Kepulauan Riau' AS instansi,
    CASE
` + strings.Join(jabatanCases, "\n") + `
      ELSE 'Anggota'
    END AS jabatan,
    CASE
` + strings.Join(bidangCases, "\n") + `
      ELSE NULL
    END AS bidang
  FROM users u
  WHERE LOWER(u.email) IN (` + strings.Join(emails, ", ") + `)
),

issue_agg AS (
  SELECT
    ia.assignee_id AS user_id,
    l.name AS bidang,
    COUNT(*) FILTER (WHERE s."group" = 'backlog') AS backlog,
    COUNT(*) FILTER (WHERE s."group" = 'unstarted') AS todo,
    COUNT(*) FILTER (WHERE s."group" IN ('started', 'triage')) AS in_progress,
    COUNT(*) FILTER (WHERE s."group" = 'completed') AS done
  FROM issues i
  JOIN projects p ON p.id = i.project_id
  JOIN workspaces w ON w.id = p.workspace_id
  JOIN states s ON s.id = i.state_id
  JOIN issue_assignees ia ON ia.issue_id = i.id AND ia.deleted_at IS NULL
  JOIN issue_labels il ON il.issue_id = i.id AND il.deleted_at IS NULL
  JOIN labels l ON l.id = il.label_id
  WHERE w.id = '58f6ec9b-f0ae-4e68-8f05-8f1d9ddf9cac'
    AND p.id = 'cfc12151-e169-4caf-bca9-3eb83ed588ee'
    AND i.deleted_at IS NULL
    AND s."group" != 'cancelled'
  GROUP BY ia.assignee_id, l.name
)

SELECT
  d.nama,
  d.email,
  d.bidang,
  d.instansi,
  d.jabatan,
  COALESCE(ia.backlog, 0) AS backlog,
  COALESCE(ia.todo, 0) AS todo,
  COALESCE(ia.in_progress, 0) AS in_progress,
  COALESCE(ia.done, 0) AS done,
  CASE 
    WHEN COALESCE(ia.backlog, 0) + COALESCE(ia.todo, 0) + COALESCE(ia.in_progress, 0) + COALESCE(ia.done, 0) = 0 
    THEN 0
    ELSE ROUND(
      (COALESCE(ia.done, 0)::numeric * 100) 
      / NULLIF(
        (COALESCE(ia.backlog, 0) + COALESCE(ia.todo, 0) + COALESCE(ia.in_progress, 0) + COALESCE(ia.done, 0))::numeric, 
        0
      ),
      2
    )
  END::float8 AS percent
FROM directory d
LEFT JOIN issue_agg ia ON ia.user_id = d.user_id AND ia.bidang = d.bidang
WHERE d.bidang IS NOT NULL
ORDER BY d.bidang, d.jabatan DESC, d.nama;
`

	// Execute query
	ctx, cancel := context.WithTimeout(r.Context(), 25*time.Second)
	defer cancel()

	b, err := QueryToCSV(ctx, s.DB, sql)
	if err != nil {
		http.Error(w, "query failed: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Set cache
	s.Cache.Set("kpi_provinsi", b)

	// Write response
	w.Header().Set("Content-Type", "text/csv; charset=utf-8")
	w.Header().Set("Cache-Control", "public, max-age=30")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(b)
}

// GET /api/v1/kpi_kabkot.csv
func (s *Server) KPIKabkot(w http.ResponseWriter, r *http.Request) {
	s.serveCSV(w, r, "kpi_kabkot", "kpi_kabkot.sql")
}
