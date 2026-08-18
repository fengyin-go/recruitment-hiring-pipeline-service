package store

import (
	"testing"
	"time"

	"recruit/internal/model"
)

func newTestStore() *MemoryStore { return NewMemoryStore() }

func TestDepartmentCRUD(t *testing.T) {
	s := newTestStore()
	if err := s.CreateDepartment(&model.Department{ID: "d1", Name: "研发部"}); err != nil {
		t.Fatalf("create: %v", err)
	}
	if err := s.CreateDepartment(&model.Department{ID: "d2", Name: "研发部"}); err != ErrConflict {
		t.Fatalf("expect conflict, got %v", err)
	}
	got, _ := s.GetDepartment("d1")
	got.Name = "平台研发"
	if err := s.UpdateDepartment(got); err != nil {
		t.Fatalf("update: %v", err)
	}
	if err := s.DeleteDepartment("d1"); err != nil {
		t.Fatalf("delete: %v", err)
	}
	if _, err := s.GetDepartment("d1"); err != ErrNotFound {
		t.Fatalf("expect not found, got %v", err)
	}
}

func TestPositionCRUD(t *testing.T) {
	s := newTestStore()
	p := &model.Position{ID: "p1", Title: "后端工程师", DepartmentID: "d1", Headcount: 3}
	if err := s.CreatePosition(p); err != nil {
		t.Fatalf("create: %v", err)
	}
	if err := s.CreatePosition(&model.Position{ID: "p2", Title: "后端工程师", DepartmentID: "d1", Headcount: 1}); err != ErrConflict {
		t.Fatalf("expect conflict, got %v", err)
	}
	got, _ := s.GetPosition("p1")
	got.Status = model.PositionClosed
	if err := s.UpdatePosition(got); err != nil {
		t.Fatalf("update: %v", err)
	}
	if err := s.DeletePosition("p1"); err != nil {
		t.Fatalf("delete: %v", err)
	}
}

func TestCandidateCRUD(t *testing.T) {
	s := newTestStore()
	c := &model.Candidate{ID: "c1", Name: "张三", Email: "z@x.com", PositionID: "p1"}
	if err := s.CreateCandidate(c); err != nil {
		t.Fatalf("create: %v", err)
	}
	if err := s.CreateCandidate(&model.Candidate{ID: "c2", Name: "李四", Email: "z@x.com", PositionID: "p1"}); err != ErrConflict {
		t.Fatalf("expect conflict by email+position, got %v", err)
	}
	got, _ := s.GetCandidate("c1")
	got.Status = model.CandidateInterviewing
	if err := s.UpdateCandidate(got); err != nil {
		t.Fatalf("update: %v", err)
	}
	if err := s.DeleteCandidate("c1"); err != nil {
		t.Fatalf("delete: %v", err)
	}
}

func TestResumeCRUD(t *testing.T) {
	s := newTestStore()
	r := &model.Resume{ID: "r1", CandidateID: "c1", Summary: "5 年后端经验"}
	if err := s.CreateResume(r); err != nil {
		t.Fatalf("create: %v", err)
	}
	if err := s.CreateResume(&model.Resume{ID: "r2", CandidateID: "c1", Summary: "x"}); err != ErrConflict {
		t.Fatalf("expect conflict, got %v", err)
	}
	if _, err := s.GetResumeByCandidate("c1"); err != nil {
		t.Fatalf("get by candidate: %v", err)
	}
	got, _ := s.GetResume("r1")
	got.Summary = "更新"
	if err := s.UpdateResume(got); err != nil {
		t.Fatalf("update: %v", err)
	}
}

func TestInterviewCRUD(t *testing.T) {
	s := newTestStore()
	iv := &model.Interview{ID: "i1", CandidateID: "c1", PositionID: "p1", Interviewer: "王经理", ScheduledAt: time.Now()}
	if err := s.CreateInterview(iv); err != nil {
		t.Fatalf("create: %v", err)
	}
	got, _ := s.GetInterview("i1")
	got.Status = model.InterviewCompleted
	if err := s.UpdateInterview(got); err != nil {
		t.Fatalf("update: %v", err)
	}
	if n := len(s.ListInterviews()); n != 1 {
		t.Fatalf("list len = %d", n)
	}
}

func TestOfferCRUD(t *testing.T) {
	s := newTestStore()
	o := &model.Offer{ID: "o1", CandidateID: "c1", PositionID: "p1", Salary: 30000}
	if err := s.CreateOffer(o); err != nil {
		t.Fatalf("create: %v", err)
	}
	if err := s.CreateOffer(&model.Offer{ID: "o2", CandidateID: "c1", PositionID: "p1", Salary: 1}); err != ErrConflict {
		t.Fatalf("expect conflict by candidate, got %v", err)
	}
	got, _ := s.GetOffer("o1")
	got.Status = model.OfferAccepted
	if err := s.UpdateOffer(got); err != nil {
		t.Fatalf("update: %v", err)
	}
}

func TestReferralCRUD(t *testing.T) {
	s := newTestStore()
	ref := &model.Referral{ID: "ref1", Referrer: "老员工王", CandidateID: "c1", PositionID: "p1", Bonus: 5000}
	if err := s.CreateReferral(ref); err != nil {
		t.Fatalf("create: %v", err)
	}
	if err := s.CreateReferral(&model.Referral{ID: "ref2", Referrer: "x", CandidateID: "c1", PositionID: "p1"}); err != ErrConflict {
		t.Fatalf("expect conflict by candidate, got %v", err)
	}
	got, _ := s.GetReferral("ref1")
	got.Status = model.ReferralHired
	if err := s.UpdateReferral(got); err != nil {
		t.Fatalf("update: %v", err)
	}
	if n := len(s.ListReferrals()); n != 1 {
		t.Fatalf("list len = %d", n)
	}
}

func TestUpdateNotFound(t *testing.T) {
	s := newTestStore()
	if err := s.UpdatePosition(&model.Position{ID: "missing"}); err != ErrNotFound {
		t.Fatalf("expect not found, got %v", err)
	}
	if err := s.UpdateCandidate(&model.Candidate{ID: "missing"}); err != ErrNotFound {
		t.Fatalf("expect not found, got %v", err)
	}
	if err := s.UpdateOffer(&model.Offer{ID: "missing"}); err != ErrNotFound {
		t.Fatalf("expect not found, got %v", err)
	}
	if err := s.UpdateReferral(&model.Referral{ID: "missing"}); err != ErrNotFound {
		t.Fatalf("expect not found, got %v", err)
	}
}

func TestUpdateConflicts(t *testing.T) {
	s := newTestStore()

	s.CreateDepartment(&model.Department{ID: "d1", Name: "研发部"})
	s.CreateDepartment(&model.Department{ID: "d2", Name: "产品部"})
	d2, _ := s.GetDepartment("d2")
	d2.Name = "研发部"
	if err := s.UpdateDepartment(d2); err != ErrConflict {
		t.Fatalf("department update conflict: got %v", err)
	}

	s.CreatePosition(&model.Position{ID: "p1", Title: "后端", DepartmentID: "d1", Headcount: 1})
	s.CreatePosition(&model.Position{ID: "p2", Title: "前端", DepartmentID: "d1", Headcount: 1})
	p2, _ := s.GetPosition("p2")
	p2.Title = "后端"
	if err := s.UpdatePosition(p2); err != ErrConflict {
		t.Fatalf("position update conflict: got %v", err)
	}
}

func TestListEmpty(t *testing.T) {
	s := newTestStore()
	if n := len(s.ListDepartments()); n != 0 {
		t.Fatalf("empty departments len = %d", n)
	}
	if n := len(s.ListPositions()); n != 0 {
		t.Fatalf("empty positions len = %d", n)
	}
	if n := len(s.ListCandidates()); n != 0 {
		t.Fatalf("empty candidates len = %d", n)
	}
	if n := len(s.ListInterviews()); n != 0 {
		t.Fatalf("empty interviews len = %d", n)
	}
	if n := len(s.ListOffers()); n != 0 {
		t.Fatalf("empty offers len = %d", n)
	}
	if n := len(s.ListReferrals()); n != 0 {
		t.Fatalf("empty referrals len = %d", n)
	}
}

func TestGetResumeNotFound(t *testing.T) {
	s := newTestStore()
	if _, err := s.GetResume("missing"); err != ErrNotFound {
		t.Fatalf("expect not found, got %v", err)
	}
	if _, err := s.GetResumeByCandidate("missing"); err != ErrNotFound {
		t.Fatalf("expect not found by candidate, got %v", err)
	}
}
