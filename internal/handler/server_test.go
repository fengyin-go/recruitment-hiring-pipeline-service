package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"recruit/internal/config"
	"recruit/internal/model"
	"recruit/internal/service"
	"recruit/internal/store"
	"recruit/pkg/httpx"
	"recruit/pkg/logger"
)

func newTestServer() *Server {
	cfg := &config.Config{MaxPageSize: 100}
	log := logger.NewLevel(logger.LevelError)
	st := store.NewMemoryStore()
	svc := service.New(st, log, cfg)
	return NewServer(svc, log, cfg)
}

func doRequest(s *Server, method, path string, body interface{}) (*httptest.ResponseRecorder, *httpx.Response) {
	var buf bytes.Buffer
	if body != nil {
		_ = json.NewEncoder(&buf).Encode(body)
	}
	req := httptest.NewRequest(method, path, &buf)
	rec := httptest.NewRecorder()
	s.Routes().ServeHTTP(rec, req)
	resp := &httpx.Response{}
	_ = json.Unmarshal(rec.Body.Bytes(), resp)
	return rec, resp
}

func decodeData(t *testing.T, resp *httpx.Response, dst interface{}) {
	t.Helper()
	raw, _ := json.Marshal(resp.Data)
	if err := json.Unmarshal(raw, dst); err != nil {
		t.Fatalf("decode data: %v", err)
	}
}

func setupDepartmentAndPosition(t *testing.T, s *Server) *model.Position {
	t.Helper()
	_, resp := doRequest(s, http.MethodPost, "/api/departments", map[string]string{"name": "研发部"})
	var d model.Department
	decodeData(t, resp, &d)
	_, resp = doRequest(s, http.MethodPost, "/api/positions", map[string]interface{}{
		"title": "后端工程师", "department_id": d.ID, "headcount": 3, "salary_min": 20000, "salary_max": 40000,
	})
	var p model.Position
	decodeData(t, resp, &p)
	return &p
}

func TestDepartmentEndpoints(t *testing.T) {
	s := newTestServer()
	_, resp := doRequest(s, http.MethodPost, "/api/departments", map[string]string{"name": "研发部"})
	if resp.Code != 0 {
		t.Fatalf("create: %s", resp.Message)
	}
	rec, _ := doRequest(s, http.MethodPost, "/api/departments", map[string]string{"name": "研发部"})
	if rec.Code != http.StatusConflict {
		t.Fatalf("duplicate status = %d", rec.Code)
	}
	rec, _ = doRequest(s, http.MethodGet, "/api/departments/missing", nil)
	if rec.Code != http.StatusNotFound {
		t.Fatalf("missing status = %d", rec.Code)
	}
}

func TestFullRecruitmentFlowEndpoints(t *testing.T) {
	s := newTestServer()
	pos := setupDepartmentAndPosition(t, s)

	_, resp := doRequest(s, http.MethodPost, "/api/candidates", map[string]string{
		"name": "张三", "email": "z@x.com", "position_id": pos.ID,
	})
	var c model.Candidate
	decodeData(t, resp, &c)

	_, resp = doRequest(s, http.MethodPost, "/api/interviews", map[string]string{
		"candidate_id": c.ID, "position_id": pos.ID, "interviewer": "王经理", "scheduled_at": "2026-08-20T10:00:00Z",
	})
	var iv model.Interview
	decodeData(t, resp, &iv)

	rec, _ := doRequest(s, http.MethodPost, "/api/interviews/"+iv.ID+"/complete", map[string]interface{}{"passed": true, "feedback": "通过"})
	if rec.Code != http.StatusOK {
		t.Fatalf("complete status = %d", rec.Code)
	}

	_, resp = doRequest(s, http.MethodPost, "/api/offers", map[string]interface{}{
		"candidate_id": c.ID, "position_id": pos.ID, "salary": 35000,
	})
	var o model.Offer
	decodeData(t, resp, &o)

	rec, _ = doRequest(s, http.MethodPost, "/api/offers/"+o.ID+"/accept", nil)
	if rec.Code != http.StatusOK {
		t.Fatalf("accept status = %d", rec.Code)
	}

	_, resp = doRequest(s, http.MethodGet, "/api/candidates/"+c.ID, nil)
	decodeData(t, resp, &c)
	if c.Status != model.CandidateHired {
		t.Fatalf("final status = %s, want hired", c.Status)
	}
}

func TestResumeEndpoints(t *testing.T) {
	s := newTestServer()
	pos := setupDepartmentAndPosition(t, s)
	_, resp := doRequest(s, http.MethodPost, "/api/candidates", map[string]string{
		"name": "张三", "email": "z@x.com", "position_id": pos.ID,
	})
	var c model.Candidate
	decodeData(t, resp, &c)

	rec, _ := doRequest(s, http.MethodPost, "/api/resumes", map[string]interface{}{
		"candidate_id": c.ID, "summary": "5 年后端", "skills": []string{"Go"},
	})
	if rec.Code != http.StatusCreated {
		t.Fatalf("resume status = %d", rec.Code)
	}
	rec, _ = doRequest(s, http.MethodGet, "/api/candidates/"+c.ID+"/resume", nil)
	if rec.Code != http.StatusOK {
		t.Fatalf("get resume status = %d", rec.Code)
	}
}

func TestStatsEndpoints(t *testing.T) {
	s := newTestServer()
	setupDepartmentAndPosition(t, s)
	for _, path := range []string{
		"/api/stats/candidates", "/api/stats/funnel", "/api/stats/positions",
		"/api/stats/interviews", "/api/stats/departments",
	} {
		rec, _ := doRequest(s, http.MethodGet, path, nil)
		if rec.Code != http.StatusOK {
			t.Fatalf("%s status = %d", path, rec.Code)
		}
	}
}

func TestMalformedJSONAndNotFound(t *testing.T) {
	s := newTestServer()
	req := httptest.NewRequest(http.MethodPost, "/api/departments", strings.NewReader("{bad"))
	rec := httptest.NewRecorder()
	s.Routes().ServeHTTP(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("malformed status = %d", rec.Code)
	}
	req = httptest.NewRequest(http.MethodGet, "/api/no-such", nil)
	rec = httptest.NewRecorder()
	s.Routes().ServeHTTP(rec, req)
	if rec.Code != http.StatusNotFound {
		t.Fatalf("not found status = %d", rec.Code)
	}
}

func TestReferralEndpoints(t *testing.T) {
	s := newTestServer()
	pos := setupDepartmentAndPosition(t, s)
	_, resp := doRequest(s, http.MethodPost, "/api/candidates", map[string]string{
		"name": "张三", "email": "z@x.com", "position_id": pos.ID,
	})
	var c model.Candidate
	decodeData(t, resp, &c)

	rec, _ := doRequest(s, http.MethodPost, "/api/referrals", map[string]interface{}{
		"referrer": "老员工王", "candidate_id": c.ID, "position_id": pos.ID, "bonus": 5000,
	})
	if rec.Code != http.StatusCreated {
		t.Fatalf("referral status = %d", rec.Code)
	}
	rec, _ = doRequest(s, http.MethodGet, "/api/stats/referrals", nil)
	if rec.Code != http.StatusOK {
		t.Fatalf("referral stats status = %d", rec.Code)
	}
	rec, _ = doRequest(s, http.MethodGet, "/api/stats/monthly", nil)
	if rec.Code != http.StatusOK {
		t.Fatalf("monthly status = %d", rec.Code)
	}
}

func TestExportEndpoints(t *testing.T) {
	s := newTestServer()
	setupDepartmentAndPosition(t, s)
	rec, _ := doRequest(s, http.MethodGet, "/api/export/positions", nil)
	if rec.Code != http.StatusOK {
		t.Fatalf("export positions status = %d", rec.Code)
	}
	rec, _ = doRequest(s, http.MethodGet, "/api/export/candidates", nil)
	if rec.Code != http.StatusOK {
		t.Fatalf("export candidates status = %d", rec.Code)
	}
}

func TestBatchRejectEndpoint(t *testing.T) {
	s := newTestServer()
	pos := setupDepartmentAndPosition(t, s)
	var ids []string
	for i := 0; i < 2; i++ {
		_, resp := doRequest(s, http.MethodPost, "/api/candidates", map[string]string{
			"name": "候选" + string(rune('A'+i)), "email": string(rune('a'+i)) + "@x.com", "position_id": pos.ID,
		})
		var c model.Candidate
		decodeData(t, resp, &c)
		ids = append(ids, c.ID)
	}
	rec, _ := doRequest(s, http.MethodPost, "/api/candidates/batch-reject", map[string]interface{}{"ids": ids})
	if rec.Code != http.StatusOK {
		t.Fatalf("batch reject status = %d", rec.Code)
	}
}

func TestResumeUpdateEndpoint(t *testing.T) {
	s := newTestServer()
	pos := setupDepartmentAndPosition(t, s)
	_, resp := doRequest(s, http.MethodPost, "/api/candidates", map[string]string{
		"name": "张三", "email": "z@x.com", "position_id": pos.ID,
	})
	var c model.Candidate
	decodeData(t, resp, &c)
	_, resp = doRequest(s, http.MethodPost, "/api/resumes", map[string]interface{}{
		"candidate_id": c.ID, "summary": "5 年后端",
	})
	var res model.Resume
	decodeData(t, resp, &res)
	rec, _ := doRequest(s, http.MethodPut, "/api/resumes/"+res.ID, map[string]interface{}{"summary": "6 年后端"})
	if rec.Code != http.StatusOK {
		t.Fatalf("update resume status = %d", rec.Code)
	}
}

func TestPositionEndpoints(t *testing.T) {
	s := newTestServer()
	pos := setupDepartmentAndPosition(t, s)
	rec, _ := doRequest(s, http.MethodGet, "/api/positions/"+pos.ID, nil)
	if rec.Code != http.StatusOK {
		t.Fatalf("get position status = %d", rec.Code)
	}
	// 列表 + 状态过滤
	rec, _ = doRequest(s, http.MethodGet, "/api/positions?status=open", nil)
	if rec.Code != http.StatusOK {
		t.Fatalf("list positions status = %d", rec.Code)
	}
}

func TestPositionUpdateEndpoint(t *testing.T) {
	s := newTestServer()
	pos := setupDepartmentAndPosition(t, s)
	rec, _ := doRequest(s, http.MethodPut, "/api/positions/"+pos.ID, map[string]interface{}{"headcount": 5, "status": "closed"})
	if rec.Code != http.StatusOK {
		t.Fatalf("update position status = %d", rec.Code)
	}
	_, resp := doRequest(s, http.MethodGet, "/api/positions/"+pos.ID, nil)
	var p model.Position
	decodeData(t, resp, &p)
	if p.Status != model.PositionClosed || p.Headcount != 5 {
		t.Fatalf("updated position = %+v", p)
	}
}

func TestRejectCandidateEndpoint(t *testing.T) {
	s := newTestServer()
	pos := setupDepartmentAndPosition(t, s)
	_, resp := doRequest(s, http.MethodPost, "/api/candidates", map[string]string{
		"name": "张三", "email": "z@x.com", "position_id": pos.ID,
	})
	var c model.Candidate
	decodeData(t, resp, &c)
	rec, _ := doRequest(s, http.MethodPost, "/api/candidates/"+c.ID+"/reject", nil)
	if rec.Code != http.StatusOK {
		t.Fatalf("reject status = %d", rec.Code)
	}
}
