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
	DB            *pgxpool.Pool
	Cache         *cache.Cache
	Queries       *queries.Store
	WorkspaceName string
	ProjectName   string
	KabkotaKey    string
}

func (s *Server) serveCSV(w http.ResponseWriter, r *http.Request, cacheKey, queryFile string, args ...any) {
	// cache
	if b, ok := s.Cache.Get(cacheKey); ok {
		w.Header().Set("Content-Type", "text/csv; charset=utf-8")
		w.Header().Set("Cache-Control", "public, max-age=30")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(b)
		return
	}

	sql, err := s.Queries.Load(queryFile)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 20*time.Second)
	defer cancel()

	b, err := QueryToCSV(ctx, s.DB, sql, args...)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	s.Cache.Set(cacheKey, b)

	w.Header().Set("Content-Type", "text/csv; charset=utf-8")
	w.Header().Set("Cache-Control", "public, max-age=30")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(b)
}

func (s *Server) KPI(w http.ResponseWriter, r *http.Request) {
	s.serveCSV(w, r, "kpi", "kpi.sql", s.WorkspaceName, s.ProjectName)
}

func (s *Server) ProgressKabkot(w http.ResponseWriter, r *http.Request) {
	s.serveCSV(w, r, "progress_kabkot", "progress_kabkot.sql", s.WorkspaceName, s.ProjectName, s.KabkotaKey)
}

func (s *Server) ProgressBidang(w http.ResponseWriter, r *http.Request) {
	s.serveCSV(w, r, "progress_bidang", "progress_bidang.sql", s.WorkspaceName, s.ProjectName)
}

func (s *Server) Heatmap(w http.ResponseWriter, r *http.Request) {
	s.serveCSV(w, r, "heatmap", "heatmap_kabkot_bidang.sql", s.WorkspaceName, s.ProjectName, s.KabkotaKey)
}

func (s *Server) IssuesDetail(w http.ResponseWriter, r *http.Request) {
	s.serveCSV(w, r, "issues_detail", "issues_detail.sql", s.WorkspaceName, s.ProjectName, s.KabkotaKey)
}
