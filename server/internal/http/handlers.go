package httpx

import (
	"context"
	"net/http"
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

// GET /csv/kpi.csv
func (s *Server) KPI(w http.ResponseWriter, r *http.Request) {
	// SQL hardcode ID -> no args
	s.serveCSV(w, r, "kpi", "kpi.sql")
}

// GET /csv/progress_kabkot.csv
func (s *Server) ProgressKabkot(w http.ResponseWriter, r *http.Request) {
	s.serveCSV(w, r, "progress_kabkot", "progress_kabkot.sql")
}

// GET /csv/progress_bidang.csv
func (s *Server) ProgressBidang(w http.ResponseWriter, r *http.Request) {
	s.serveCSV(w, r, "progress_bidang", "progress_bidang.sql")
}

// GET /csv/heatmap.csv
func (s *Server) Heatmap(w http.ResponseWriter, r *http.Request) {
	// file name MUST match what you saved
	// if your file is named heatmap_kabkot_bidang.sql, keep it consistent
	s.serveCSV(w, r, "heatmap", "heatmap_kabkot_bidang.sql")
}

// GET /csv/issues_detail.csv
func (s *Server) IssuesDetail(w http.ResponseWriter, r *http.Request) {
	s.serveCSV(w, r, "issues_detail", "issues_detail.sql")
}
