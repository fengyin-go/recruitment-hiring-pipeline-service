package handler

import (
	"net/http"

	"recruit/pkg/httpx"
)

func (s *Server) registerStatsRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/stats/candidates", s.candidateStats)
	mux.HandleFunc("GET /api/stats/funnel", s.funnel)
	mux.HandleFunc("GET /api/stats/positions", s.positionStats)
	mux.HandleFunc("GET /api/stats/interviews", s.interviewStats)
	mux.HandleFunc("GET /api/stats/departments", s.departmentStats)
	mux.HandleFunc("GET /api/stats/monthly", s.monthlyStats)
	mux.HandleFunc("GET /api/stats/referrals", s.referralStats)
}

func (s *Server) monthlyStats(w http.ResponseWriter, r *http.Request) {
	stats, err := s.svc.MonthlyStats()
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, stats)
}

func (s *Server) referralStats(w http.ResponseWriter, r *http.Request) {
	stats, err := s.svc.ReferralStats()
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, stats)
}

func (s *Server) candidateStats(w http.ResponseWriter, r *http.Request) {
	stats, err := s.svc.CandidateStats()
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, stats)
}

func (s *Server) funnel(w http.ResponseWriter, r *http.Request) {
	funnel, err := s.svc.RecruitmentFunnel()
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, funnel)
}

func (s *Server) positionStats(w http.ResponseWriter, r *http.Request) {
	stats, err := s.svc.ListPositionStats()
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, stats)
}

func (s *Server) interviewStats(w http.ResponseWriter, r *http.Request) {
	stats, err := s.svc.InterviewStats()
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, stats)
}

func (s *Server) departmentStats(w http.ResponseWriter, r *http.Request) {
	stats, err := s.svc.ListDepartmentStats()
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, stats)
}
