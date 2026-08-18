package handler

import (
	"net/http"

	"recruit/internal/model"
	"recruit/pkg/httpx"
)

func (s *Server) registerPositionRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/positions", s.createPosition)
	mux.HandleFunc("GET /api/positions", s.listPositions)
	mux.HandleFunc("GET /api/positions/{id}", s.getPosition)
	mux.HandleFunc("PUT /api/positions/{id}", s.updatePosition)
	mux.HandleFunc("DELETE /api/positions/{id}", s.deletePosition)
}

type positionRequest struct {
	Title        string `json:"title"`
	DepartmentID string `json:"department_id"`
	Description  string `json:"description"`
	Headcount    int    `json:"headcount"`
	SalaryMin    int64  `json:"salary_min"`
	SalaryMax    int64  `json:"salary_max"`
	Status       string `json:"status"`
}

func (s *Server) createPosition(w http.ResponseWriter, r *http.Request) {
	var req positionRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	p, err := s.svc.CreatePosition(model.Position{
		Title: req.Title, DepartmentID: req.DepartmentID, Description: req.Description,
		Headcount: req.Headcount, SalaryMin: req.SalaryMin, SalaryMax: req.SalaryMax,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, p)
}

func (s *Server) listPositions(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	filter := model.PositionFilter{
		DepartmentID: r.URL.Query().Get("department_id"),
		Status:       r.URL.Query().Get("status"),
		Keyword:      r.URL.Query().Get("keyword"),
	}
	items, total, err := s.svc.ListPositions(filter, pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}

func (s *Server) getPosition(w http.ResponseWriter, r *http.Request) {
	p, err := s.svc.GetPosition(r.PathValue("id"))
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, p)
}

func (s *Server) updatePosition(w http.ResponseWriter, r *http.Request) {
	var req positionRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	p, err := s.svc.UpdatePosition(r.PathValue("id"), model.Position{
		Title: req.Title, Description: req.Description, Headcount: req.Headcount,
		SalaryMin: req.SalaryMin, SalaryMax: req.SalaryMax, Status: req.Status,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, p)
}

func (s *Server) deletePosition(w http.ResponseWriter, r *http.Request) {
	if err := s.svc.DeletePosition(r.PathValue("id")); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.NoContent(w)
}
