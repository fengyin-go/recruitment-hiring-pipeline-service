package handler

import (
	"net/http"

	"recruit/pkg/httpx"
)

func (s *Server) registerExportRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/export/positions", s.exportPositions)
	mux.HandleFunc("GET /api/export/candidates", s.exportCandidates)
}

func (s *Server) exportPositions(w http.ResponseWriter, r *http.Request) {
	items, err := s.svc.ExportPositions()
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, items)
}

func (s *Server) exportCandidates(w http.ResponseWriter, r *http.Request) {
	items, err := s.svc.ExportCandidates()
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, items)
}
