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
	// Parse query parameters for filtering
	filterKabKota := r.URL.Query().Get("kab_kota")
	filterBidang := r.URL.Query().Get("bidang")

	// Build cache key with filters
	cacheKey := "heatmap"
	if filterKabKota != "" || filterBidang != "" {
		cacheKey += "_" + filterKabKota + "_" + filterBidang
	}

	// Check cache
	if b, ok := s.Cache.Get(cacheKey); ok {
		w.Header().Set("Content-Type", "text/csv; charset=utf-8")
		w.Header().Set("Cache-Control", "public, max-age=30")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(b)
		return
	}

	// Build WHERE clause for filters
	var filterClauses []string
	if filterKabKota != "" {
		filterClauses = append(filterClauses, "  kab_kota = '"+sanitizeSQL(filterKabKota)+"'")
	}
	if filterBidang != "" {
		filterClauses = append(filterClauses, "  bidang = '"+sanitizeSQL(filterBidang)+"'")
	}

	havingClause := ""
	if len(filterClauses) > 0 {
		havingClause = "\nHAVING\n" + strings.Join(filterClauses, "\n  AND ")
	}

	// Build SQL with optional filters
	sql := `
WITH base AS (
  SELECT
    s."group" AS status,
    l.name AS bidang,
    ia.assignee_id,
    SUBSTRING(
      COALESCE(u.display_name,'') || ' ' ||
      COALESCE(u.first_name,'')   || ' ' ||
      COALESCE(u.last_name,''),
      '(21[0-9]{2})'
    ) AS kab_kode
  FROM issues i
  JOIN projects p ON p.id = i.project_id
  JOIN workspaces w ON w.id = p.workspace_id
  JOIN states s ON s.id = i.state_id
  LEFT JOIN issue_assignees ia ON ia.issue_id = i.id AND ia.deleted_at IS NULL
  LEFT JOIN users u ON u.id = ia.assignee_id
  JOIN issue_labels il ON il.issue_id = i.id AND il.deleted_at IS NULL
  JOIN labels l ON l.id = il.label_id
  WHERE w.id = '58f6ec9b-f0ae-4e68-8f05-8f1d9ddf9cac'::uuid
    AND p.id = 'cfc12151-e169-4caf-bca9-3eb83ed588ee'::uuid
    AND i.deleted_at IS NULL
    AND s."group" != 'cancelled'
),
data AS (
  SELECT
    CASE
      WHEN assignee_id IS NULL THEN 'Belum Ditugaskan'
      WHEN kab_kode IS NULL THEN 'Kode Kab/Kota Tidak Terbaca'
      WHEN kab_kode = '2101' THEN 'Karimun'
      WHEN kab_kode = '2102' THEN 'Bintan'
      WHEN kab_kode = '2103' THEN 'Natuna'
      WHEN kab_kode = '2104' THEN 'Lingga'
      WHEN kab_kode = '2105' THEN 'Kep. Anambas'
      WHEN kab_kode = '2171' THEN 'Batam'
      WHEN kab_kode = '2172' THEN 'Tanjung Pinang'
      ELSE 'Lainnya'
    END AS kab_kota,
    bidang,
    status
  FROM base
)
SELECT
  kab_kota,
  bidang,
  COUNT(*) AS total,
  COUNT(*) FILTER (WHERE status = 'completed') AS selesai,
  (
    ROUND(
      (COUNT(*) FILTER (WHERE status = 'completed')::numeric * 100)
      / NULLIF(COUNT(*)::numeric, 0),
      2
    )
  )::float8 AS persen_selesai
FROM data
GROUP BY kab_kota, bidang` + havingClause + `
ORDER BY kab_kota, bidang;
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
	s.Cache.Set(cacheKey, b)

	// Write response
	w.Header().Set("Content-Type", "text/csv; charset=utf-8")
	w.Header().Set("Cache-Control", "public, max-age=30")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(b)
}

// GET /api/v1/issues_detail.csv
func (s *Server) IssuesDetail(w http.ResponseWriter, r *http.Request) {
	// Parse query parameters for filtering
	filterKabKota := r.URL.Query().Get("kab_kota")
	filterBidang := r.URL.Query().Get("bidang")
	filterStatus := r.URL.Query().Get("status")

	// Build cache key with filters
	cacheKey := "issues_detail"
	if filterKabKota != "" || filterBidang != "" || filterStatus != "" {
		cacheKey += "_" + filterKabKota + "_" + filterBidang + "_" + filterStatus
	}

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

	// Build email -> nama mapping from ALL directory entries (provinsi + kabkot)
	var namaCases []string
	allRows := append(dir.Provinsi, dir.Kabkot...)
	for _, row := range allRows {
		email := strings.ToLower(row.Email)
		namaCases = append(namaCases, "      WHEN LOWER(u.email) = '"+email+"' THEN '"+sanitizeSQL(row.Nama)+"'")
	}

	// Build WHERE clause for filters
	var filterClauses []string
	if filterKabKota != "" {
		filterClauses = append(filterClauses, "  kab_kota = '"+sanitizeSQL(filterKabKota)+"'")
	}
	if filterBidang != "" {
		filterClauses = append(filterClauses, "  bidang = '"+sanitizeSQL(filterBidang)+"'")
	}
	if filterStatus != "" {
		filterClauses = append(filterClauses, "  status = '"+sanitizeSQL(filterStatus)+"'")
	}

	whereClause := ""
	if len(filterClauses) > 0 {
		whereClause = "\nWHERE\n" + strings.Join(filterClauses, "\n  AND ")
	}

	// Build dynamic SQL
	sql := `
-- Issues Detail: Get ketua_bidang nama from CSV
WITH 
base AS (
  SELECT
    i.id AS issue_id,
    i.name AS issue_title,
    s."group" AS status,
    i.start_date,
    i.target_date,

    ia.assignee_id,
    u.email AS assignee_email,

    -- Get nama from CSV, fallback to users table
    CASE
` + strings.Join(namaCases, "\n") + `
      ELSE COALESCE(
        u.display_name,
        NULLIF(TRIM(COALESCE(u.first_name,'') || ' ' || COALESCE(u.last_name,'')), ''),
        '-'
      )
    END AS assignee_name,

    -- ekstraksi kode kab/kota
    SUBSTRING(
      COALESCE(u.display_name,'') || ' ' ||
      COALESCE(u.first_name,'')   || ' ' ||
      COALESCE(u.last_name,''),
      '(21[0-9]{2})'
    ) AS kab_kode,

    l.name AS bidang
  FROM issues i
  JOIN projects p ON p.id = i.project_id
  JOIN workspaces w ON w.id = p.workspace_id
  JOIN states s ON s.id = i.state_id
  LEFT JOIN issue_assignees ia ON ia.issue_id = i.id AND ia.deleted_at IS NULL
  LEFT JOIN users u ON u.id = ia.assignee_id
  LEFT JOIN issue_labels il ON il.issue_id = i.id AND il.deleted_at IS NULL
  LEFT JOIN labels l ON l.id = il.label_id
  WHERE w.id = '58f6ec9b-f0ae-4e68-8f05-8f1d9ddf9cac'::uuid
    AND p.id = 'cfc12151-e169-4caf-bca9-3eb83ed588ee'::uuid
    AND i.deleted_at IS NULL
    AND s."group" != 'cancelled'
),
latest_comment AS (
  SELECT DISTINCT ON (ic.issue_id)
    ic.issue_id,
    ic.comment_stripped AS last_comment,
    ic.created_at AS comment_time
  FROM issue_comments ic
  ORDER BY ic.issue_id, ic.created_at DESC
),
final AS (
  SELECT
    b.issue_title,
    CASE
      WHEN b.assignee_id IS NULL THEN 'Belum Ditugaskan'
      WHEN b.kab_kode IS NULL THEN 'Kode Kab/Kota Tidak Terbaca'
      WHEN b.kab_kode = '2101' THEN 'Karimun'
      WHEN b.kab_kode = '2102' THEN 'Bintan'
      WHEN b.kab_kode = '2103' THEN 'Natuna'
      WHEN b.kab_kode = '2104' THEN 'Lingga'
      WHEN b.kab_kode = '2105' THEN 'Kep. Anambas'
      WHEN b.kab_kode = '2171' THEN 'Batam'
      WHEN b.kab_kode = '2172' THEN 'Tanjung Pinang'
      ELSE 'Lainnya'
    END AS kab_kota,
    COALESCE(b.bidang, '-') AS bidang,
    b.status,
    b.assignee_name AS ketua_bidang,
    COALESCE(b.assignee_email, '-') AS email_ketua_bidang,
    b.start_date,
    b.target_date,
    lc.last_comment,
    lc.comment_time
  FROM base b
  LEFT JOIN latest_comment lc ON lc.issue_id = b.issue_id
)
SELECT * FROM final` + whereClause + `
ORDER BY issue_title;
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
	s.Cache.Set(cacheKey, b)

	// Write response
	w.Header().Set("Content-Type", "text/csv; charset=utf-8")
	w.Header().Set("Cache-Control", "public, max-age=30")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(b)
}

// GET /api/v1/kpi_provinsi.csv
func (s *Server) KPIProvinsi(w http.ResponseWriter, r *http.Request) {
	// Parse query parameters for filtering
	filterBidang := r.URL.Query().Get("bidang")
	filterJabatan := r.URL.Query().Get("jabatan")

	// Build cache key with filters
	cacheKey := "kpi_provinsi"
	if filterBidang != "" || filterJabatan != "" {
		cacheKey += "_" + filterBidang + "_" + filterJabatan
	}

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

	// Build dynamic SQL with email-bidang-jabatan mapping
	// Only include provinsi staff (Ketua + Anggota Bidang, exclude Sekretariat)
	var emailCases, bidangCases, jabatanCases, namaCases []string
	var emails []string

	for _, row := range dir.Provinsi {
		email := strings.ToLower(row.Email)
		emails = append(emails, "'"+email+"'")
		emailCases = append(emailCases, "      WHEN LOWER(u.email) = '"+email+"' THEN TRUE")
		namaCases = append(namaCases, "      WHEN LOWER(u.email) = '"+email+"' THEN '"+sanitizeSQL(row.Nama)+"'")
		bidangCases = append(bidangCases, "      WHEN LOWER(u.email) = '"+email+"' THEN '"+row.Bidang+"'")
		jabatanCases = append(jabatanCases, "      WHEN LOWER(u.email) = '"+email+"' THEN '"+row.Jabatan+"'")
	}

	if len(emails) == 0 {
		http.Error(w, "no provinsi staff found in directory", http.StatusInternalServerError)
		return
	}

	// Build WHERE clause for filters
	var filterClauses []string
	if filterBidang != "" {
		filterClauses = append(filterClauses, "  d.bidang = '"+sanitizeSQL(filterBidang)+"'")
	}
	if filterJabatan != "" {
		filterClauses = append(filterClauses, "  d.jabatan = '"+sanitizeSQL(filterJabatan)+"'")
	}

	additionalWhere := ""
	if len(filterClauses) > 0 {
		additionalWhere = "\n  AND " + strings.Join(filterClauses, "\n  AND ")
	}

	// Build dynamic SQL
	sql := `
-- KPI Provinsi: Dynamic from CSV
WITH 
directory AS (
  SELECT 
    u.id AS user_id,
    LOWER(u.email) AS email,
    CASE
` + strings.Join(namaCases, "\n") + `
      ELSE COALESCE(u.display_name, u.first_name || ' ' || u.last_name, '-')
    END AS nama,
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
WHERE d.bidang IS NOT NULL` + additionalWhere + `
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
	s.Cache.Set(cacheKey, b)

	// Write response
	w.Header().Set("Content-Type", "text/csv; charset=utf-8")
	w.Header().Set("Cache-Control", "public, max-age=30")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(b)
}

// GET /api/v1/kpi_kabkot.csv
func (s *Server) KPIKabkot(w http.ResponseWriter, r *http.Request) {
	// Parse query parameters for filtering
	filterBidang := r.URL.Query().Get("bidang")
	filterInstansi := r.URL.Query().Get("instansi")
	filterJabatan := r.URL.Query().Get("jabatan")

	// Build cache key with filters
	cacheKey := "kpi_kabkot"
	if filterBidang != "" || filterInstansi != "" || filterJabatan != "" {
		cacheKey += "_" + filterBidang + "_" + filterInstansi + "_" + filterJabatan
	}

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

	// Build dynamic SQL with email-bidang-instansi-jabatan mapping
	// Include ALL Ketua: Kepala, Ketua Pelaksana, Ketua Sekretariat, Ketua Bidang
	var emailCases, bidangCases, instansiCases, jabatanCases, namaCases []string
	var emails []string

	for _, row := range dir.Kabkot {
		email := strings.ToLower(row.Email)
		emails = append(emails, "'"+email+"'")
		emailCases = append(emailCases, "      WHEN LOWER(u.email) = '"+email+"' THEN TRUE")
		namaCases = append(namaCases, "      WHEN LOWER(u.email) = '"+email+"' THEN '"+sanitizeSQL(row.Nama)+"'")
		bidangCases = append(bidangCases, "      WHEN LOWER(u.email) = '"+email+"' THEN '"+row.Bidang+"'")
		instansiCases = append(instansiCases, "      WHEN LOWER(u.email) = '"+email+"' THEN '"+row.Instansi+"'")
		jabatanCases = append(jabatanCases, "      WHEN LOWER(u.email) = '"+email+"' THEN '"+row.Jabatan+"'")
	}

	if len(emails) == 0 {
		http.Error(w, "no ketua kabkot found in directory", http.StatusInternalServerError)
		return
	}

	// Build WHERE clause for filters
	var filterClauses []string
	if filterBidang != "" {
		filterClauses = append(filterClauses, "  d.bidang = '"+sanitizeSQL(filterBidang)+"'")
	}
	if filterInstansi != "" {
		filterClauses = append(filterClauses, "  d.instansi = '"+sanitizeSQL(filterInstansi)+"'")
	}
	if filterJabatan != "" {
		filterClauses = append(filterClauses, "  d.jabatan = '"+sanitizeSQL(filterJabatan)+"'")
	}

	additionalWhere := ""
	if len(filterClauses) > 0 {
		additionalWhere = "\n  AND " + strings.Join(filterClauses, "\n  AND ")
	}

	// Build dynamic SQL
	sql := `
-- KPI Kabkot: Dynamic from CSV (All Ketua)
WITH 
directory AS (
  SELECT 
    u.id AS user_id,
    LOWER(u.email) AS email,
    CASE
` + strings.Join(namaCases, "\n") + `
      ELSE COALESCE(u.display_name, u.first_name || ' ' || u.last_name, '-')
    END AS nama,
    CASE
` + strings.Join(instansiCases, "\n") + `
      ELSE 'BPS Kabupaten/Kota'
    END AS instansi,
    CASE
` + strings.Join(jabatanCases, "\n") + `
      ELSE 'Ketua'
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
WHERE d.bidang IS NOT NULL` + additionalWhere + `
ORDER BY d.instansi, d.jabatan, d.bidang;
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
	s.Cache.Set(cacheKey, b)

	// Write response
	w.Header().Set("Content-Type", "text/csv; charset=utf-8")
	w.Header().Set("Cache-Control", "public, max-age=30")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(b)
}

// sanitizeSQL escapes single quotes in SQL string literals to prevent SQL injection
func sanitizeSQL(s string) string {
	return strings.ReplaceAll(s, "'", "''")
}
