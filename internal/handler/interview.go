package handler

import (
	"net/http"
	"time"

	"recruit/internal/model"
	"recruit/pkg/httpx"
)

func (s *Server) registerInterviewRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/interviews", s.scheduleInterview)
	mux.HandleFunc("GET /api/interviews", s.listInterviews)
	mux.HandleFunc("GET /api/interviews/{id}", s.getInterview)
	mux.HandleFunc("POST /api/interviews/{id}/complete", s.completeInterview)
}

type scheduleInterviewRequest struct {
	CandidateID string `json:"candidate_id"`
	PositionID  string `json:"position_id"`
	Interviewer string `json:"interviewer"`
	ScheduledAt string `json:"scheduled_at"`
}

func (s *Server) scheduleInterview(w http.ResponseWriter, r *http.Request) {
	var req scheduleInterviewRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	scheduledAt, err := parseTime(req.ScheduledAt)
	if err != nil {
		httpx.BadRequest(w, "scheduled_at 时间格式非法")
		return
	}
	iv, err := s.svc.ScheduleInterview(req.CandidateID, req.PositionID, req.Interviewer, scheduledAt)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, iv)
}

func (s *Server) listInterviews(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	filter := model.InterviewFilter{
		CandidateID: r.URL.Query().Get("candidate_id"),
		PositionID:  r.URL.Query().Get("position_id"),
		Status:      r.URL.Query().Get("status"),
	}
	items, total, err := s.svc.ListInterviews(filter, pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}

func (s *Server) getInterview(w http.ResponseWriter, r *http.Request) {
	iv, err := s.svc.GetInterview(r.PathValue("id"))
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, iv)
}

type completeInterviewRequest struct {
	Passed   bool   `json:"passed"`
	Feedback string `json:"feedback"`
}

func (s *Server) completeInterview(w http.ResponseWriter, r *http.Request) {
	var req completeInterviewRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	iv, err := s.svc.CompleteInterview(r.PathValue("id"), req.Passed, req.Feedback)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, iv)
}

func parseTime(s string) (time.Time, error) {
	if t, err := time.Parse(time.RFC3339, s); err == nil {
		return t, nil
	}
	return time.Parse("2006-01-02 15:04:05", s)
}
