package httpx

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func NewRouter(s *Server) http.Handler {
	r := chi.NewRouter()

	r.Get("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	})

	r.Route("/api/v1", func(r chi.Router) {
		// Debug endpoints
		r.Get("/debug/directory", s.DebugDirectory)
		r.Get("/debug/sql", s.DebugSQL)

		// Core KPI endpoints (using SQL templates)
		r.Get("/kpi_provinsi.csv", s.KPIProvinsiTemplate)
		r.Get("/kpi_kabkot.csv", s.KPIKabkotTemplate)

		// Supporting endpoints (using SQL templates)
		r.Get("/heatmap.csv", s.HeatmapTemplate)
		r.Get("/issues_detail.csv", s.IssuesDetailTemplate)

		// Advanced analytics endpoints (using SQL templates)
		r.Get("/timeline.csv", s.TimelineTemplate)
		r.Get("/leaderboard.csv", s.LeaderboardTemplate)
		r.Get("/workload.csv", s.WorkloadTemplate)
	})

	return r
}
