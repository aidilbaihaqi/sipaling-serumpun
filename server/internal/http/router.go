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
		// Core KPI endpoints
		r.Get("/kpi_provinsi.csv", s.KPIProvinsi)
		r.Get("/kpi_kabkot.csv", s.KPIKabkot)

		// Supporting endpoints
		r.Get("/heatmap.csv", s.Heatmap)
		r.Get("/issues_detail.csv", s.IssuesDetail)
	})

	return r
}
