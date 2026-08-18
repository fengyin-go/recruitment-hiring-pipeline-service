package handler

import (
	"net/http"
	"testing"
)

func TestGetMissingEntityReturnsNotFound(t *testing.T) {
	s := newTestServer()
	paths := []string{
		"/api/departments/missing",
		"/api/positions/missing",
		"/api/candidates/missing",
		"/api/resumes/missing",
		"/api/interviews/missing",
		"/api/offers/missing",
		"/api/referrals/missing",
	}
	for _, path := range paths {
		rec, _ := doRequest(s, http.MethodGet, path, nil)
		if rec.Code != http.StatusNotFound {
			t.Fatalf("%s status = %d, want 404", path, rec.Code)
		}
	}
}
