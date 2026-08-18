package handler

import (
	"net/http"

	"recruit/internal/model"
	"recruit/pkg/httpx"
)

func (s *Server) registerOfferRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/offers", s.createOffer)
	mux.HandleFunc("GET /api/offers", s.listOffers)
	mux.HandleFunc("GET /api/offers/{id}", s.getOffer)
	mux.HandleFunc("POST /api/offers/{id}/accept", s.acceptOffer)
	mux.HandleFunc("POST /api/offers/{id}/decline", s.declineOffer)
}

type offerRequest struct {
	CandidateID string `json:"candidate_id"`
	PositionID  string `json:"position_id"`
	Salary      int64  `json:"salary"`
}

func (s *Server) createOffer(w http.ResponseWriter, r *http.Request) {
	var req offerRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	o, err := s.svc.CreateOffer(req.CandidateID, req.PositionID, req.Salary)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, o)
}

func (s *Server) listOffers(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	filter := model.OfferFilter{
		CandidateID: r.URL.Query().Get("candidate_id"),
		PositionID:  r.URL.Query().Get("position_id"),
		Status:      r.URL.Query().Get("status"),
	}
	items, total, err := s.svc.ListOffers(filter, pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}

func (s *Server) getOffer(w http.ResponseWriter, r *http.Request) {
	o, err := s.svc.GetOffer(r.PathValue("id"))
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, o)
}

func (s *Server) acceptOffer(w http.ResponseWriter, r *http.Request) {
	o, err := s.svc.AcceptOffer(r.PathValue("id"))
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, o)
}

func (s *Server) declineOffer(w http.ResponseWriter, r *http.Request) {
	o, err := s.svc.DeclineOffer(r.PathValue("id"))
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, o)
}
