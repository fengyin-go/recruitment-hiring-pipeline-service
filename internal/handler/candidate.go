package handler

import (
	"net/http"

	"recruit/internal/model"
	"recruit/pkg/httpx"
)

func (s *Server) registerCandidateRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/candidates", s.createCandidate)
	mux.HandleFunc("GET /api/candidates", s.listCandidates)
	mux.HandleFunc("GET /api/candidates/{id}", s.getCandidate)
	mux.HandleFunc("POST /api/candidates/{id}/reject", s.rejectCandidate)
	mux.HandleFunc("POST /api/candidates/batch-reject", s.batchReject)
}

type candidateRequest struct {
	Name       string `json:"name"`
	Email      string `json:"email"`
	Phone      string `json:"phone"`
	PositionID string `json:"position_id"`
}

func (s *Server) createCandidate(w http.ResponseWriter, r *http.Request) {
	var req candidateRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	c, err := s.svc.CreateCandidate(model.Candidate{
		Name: req.Name, Email: req.Email, Phone: req.Phone, PositionID: req.PositionID,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, c)
}

func (s *Server) listCandidates(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	filter := model.CandidateFilter{
		PositionID: r.URL.Query().Get("position_id"),
		Status:     r.URL.Query().Get("status"),
		Keyword:    r.URL.Query().Get("keyword"),
	}
	items, total, err := s.svc.ListCandidates(filter, pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}

func (s *Server) getCandidate(w http.ResponseWriter, r *http.Request) {
	c, err := s.svc.GetCandidate(r.PathValue("id"))
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, c)
}

func (s *Server) rejectCandidate(w http.ResponseWriter, r *http.Request) {
	c, err := s.svc.RejectCandidate(r.PathValue("id"))
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, c)
}

type batchRejectRequest struct {
	IDs []string `json:"ids"`
}

func (s *Server) batchReject(w http.ResponseWriter, r *http.Request) {
	var req batchRejectRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	n, err := s.svc.BatchReject(req.IDs)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, map[string]int{"rejected": n})
}
