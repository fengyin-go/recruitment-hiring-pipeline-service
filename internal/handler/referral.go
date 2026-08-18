package handler

import (
	"net/http"

	"recruit/internal/model"
	"recruit/pkg/httpx"
)

func (s *Server) registerReferralRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/referrals", s.createReferral)
	mux.HandleFunc("GET /api/referrals", s.listReferrals)
	mux.HandleFunc("GET /api/referrals/{id}", s.getReferral)
	mux.HandleFunc("POST /api/referrals/{id}/hired", s.markReferralHired)
}

type referralRequest struct {
	Referrer    string `json:"referrer"`
	CandidateID string `json:"candidate_id"`
	PositionID  string `json:"position_id"`
	Bonus       int64  `json:"bonus"`
}

func (s *Server) createReferral(w http.ResponseWriter, r *http.Request) {
	var req referralRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	ref, err := s.svc.CreateReferral(req.Referrer, req.CandidateID, req.PositionID, req.Bonus)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, ref)
}

func (s *Server) listReferrals(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	filter := model.ReferralFilter{
		Referrer: r.URL.Query().Get("referrer"),
		Status:   r.URL.Query().Get("status"),
	}
	items, total, err := s.svc.ListReferrals(filter, pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}

func (s *Server) getReferral(w http.ResponseWriter, r *http.Request) {
	ref, err := s.svc.GetReferral(r.PathValue("id"))
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, ref)
}

func (s *Server) markReferralHired(w http.ResponseWriter, r *http.Request) {
	ref, err := s.svc.MarkReferralHired(r.PathValue("id"))
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, ref)
}
