package handler

import (
	"net/http"

	"recruit/internal/model"
	"recruit/pkg/httpx"
)

func (s *Server) registerResumeRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/resumes", s.createResume)
	mux.HandleFunc("GET /api/resumes/{id}", s.getResume)
	mux.HandleFunc("PUT /api/resumes/{id}", s.updateResume)
	mux.HandleFunc("GET /api/candidates/{id}/resume", s.getCandidateResume)
}

type resumeRequest struct {
	CandidateID string   `json:"candidate_id"`
	Summary     string   `json:"summary"`
	Skills      []string `json:"skills"`
	Experience  string   `json:"experience"`
}

func (s *Server) createResume(w http.ResponseWriter, r *http.Request) {
	var req resumeRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	res, err := s.svc.CreateResume(model.Resume{
		CandidateID: req.CandidateID, Summary: req.Summary, Skills: req.Skills, Experience: req.Experience,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, res)
}

func (s *Server) getResume(w http.ResponseWriter, r *http.Request) {
	res, err := s.svc.GetResume(r.PathValue("id"))
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, res)
}

func (s *Server) updateResume(w http.ResponseWriter, r *http.Request) {
	var req resumeRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	res, err := s.svc.UpdateResume(r.PathValue("id"), model.Resume{
		Summary: req.Summary, Skills: req.Skills, Experience: req.Experience,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, res)
}

func (s *Server) getCandidateResume(w http.ResponseWriter, r *http.Request) {
	res, err := s.svc.GetResumeByCandidate(r.PathValue("id"))
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, res)
}
