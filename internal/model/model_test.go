package model

import (
	"testing"
)

func TestPositionValidate(t *testing.T) {
	p := &Position{Title: "后端", DepartmentID: "d1", Headcount: 3, SalaryMin: 20000, SalaryMax: 40000}
	if err := p.Validate(); err != nil {
		t.Fatalf("valid: %v", err)
	}
	if p.Status != PositionOpen {
		t.Fatalf("default status = %s", p.Status)
	}
	p.SalaryMin = 50000
	if err := p.Validate(); err == nil {
		t.Fatalf("expect error for min > max")
	}
}

func TestCandidateValidate(t *testing.T) {
	c := &Candidate{Name: "张三", Email: "z@x.com", PositionID: "p1"}
	if err := c.Validate(); err != nil {
		t.Fatalf("valid: %v", err)
	}
	if c.Status != CandidateApplied {
		t.Fatalf("default status = %s", c.Status)
	}
	c.Name = ""
	if err := c.Validate(); err == nil {
		t.Fatalf("expect error for empty name")
	}
}

func TestCandidateTransitions(t *testing.T) {
	cases := []struct {
		from, to string
		want     bool
	}{
		{CandidateApplied, CandidateInterviewing, true},
		{CandidateApplied, CandidateRejected, true},
		{CandidateApplied, CandidateHired, false},
		{CandidateInterviewing, CandidateOffered, true},
		{CandidateInterviewing, CandidateRejected, true},
		{CandidateOffered, CandidateHired, true},
		{CandidateOffered, CandidateRejected, true},
		{CandidateHired, CandidateRejected, false},
		{CandidateRejected, CandidateHired, false},
	}
	for _, c := range cases {
		if got := CanTransitionCandidate(c.from, c.to); got != c.want {
			t.Fatalf("transition %s->%s = %v, want %v", c.from, c.to, got, c.want)
		}
	}
}

func TestIsTerminalCandidate(t *testing.T) {
	for _, s := range []string{CandidateHired, CandidateRejected} {
		if !IsTerminalCandidate(s) {
			t.Fatalf("%s should be terminal", s)
		}
	}
	for _, s := range []string{CandidateApplied, CandidateInterviewing, CandidateOffered} {
		if IsTerminalCandidate(s) {
			t.Fatalf("%s should not be terminal", s)
		}
	}
}

func TestInterviewTransitions(t *testing.T) {
	if !CanTransitionInterview(InterviewScheduled, InterviewCompleted) {
		t.Fatalf("scheduled->completed should be allowed")
	}
	if !CanTransitionInterview(InterviewScheduled, InterviewCanceled) {
		t.Fatalf("scheduled->canceled should be allowed")
	}
	if CanTransitionInterview(InterviewCompleted, InterviewScheduled) {
		t.Fatalf("completed->scheduled should not be allowed")
	}
}

func TestOfferTransitions(t *testing.T) {
	if !CanTransitionOffer(OfferPending, OfferAccepted) {
		t.Fatalf("pending->accepted should be allowed")
	}
	if !CanTransitionOffer(OfferPending, OfferDeclined) {
		t.Fatalf("pending->declined should be allowed")
	}
	if CanTransitionOffer(OfferAccepted, OfferDeclined) {
		t.Fatalf("accepted->declined should not be allowed")
	}
}

func TestOfferValidate(t *testing.T) {
	o := &Offer{CandidateID: "c1", PositionID: "p1", Salary: 30000}
	if err := o.Validate(); err != nil {
		t.Fatalf("valid: %v", err)
	}
	if o.Status != OfferPending {
		t.Fatalf("default status = %s", o.Status)
	}
	o.Salary = 0
	if err := o.Validate(); err == nil {
		t.Fatalf("expect error for zero salary")
	}
}

func TestFilters(t *testing.T) {
	p := &Position{Title: "后端工程师", DepartmentID: "d1", Status: PositionOpen}
	if !(PositionFilter{Status: PositionOpen}).Match(p) {
		t.Fatalf("position filter status")
	}
	if !(PositionFilter{Keyword: "后端"}).Match(p) {
		t.Fatalf("position filter keyword")
	}

	c := &Candidate{Name: "张三", Email: "z@x.com", PositionID: "p1", Status: CandidateApplied}
	if !(CandidateFilter{Status: CandidateApplied}).Match(c) {
		t.Fatalf("candidate filter status")
	}
	if !(CandidateFilter{Keyword: "张三"}).Match(c) {
		t.Fatalf("candidate filter keyword")
	}
}

func TestReferralValidate(t *testing.T) {
	r := &Referral{Referrer: "老员工王", CandidateID: "c1", PositionID: "p1", Bonus: 5000}
	if err := r.Validate(); err != nil {
		t.Fatalf("valid: %v", err)
	}
	if r.Status != ReferralPending {
		t.Fatalf("default status = %s", r.Status)
	}
	r.Bonus = -1
	if err := r.Validate(); err == nil {
		t.Fatalf("expect error for negative bonus")
	}
}

func TestReferralFilter(t *testing.T) {
	r := &Referral{Referrer: "老员工王", CandidateID: "c1", PositionID: "p1", Status: ReferralHired}
	if !(ReferralFilter{Status: ReferralHired}).Match(r) {
		t.Fatalf("expect match by status")
	}
	if !(ReferralFilter{Referrer: "老员工王"}).Match(r) {
		t.Fatalf("expect match by referrer")
	}
	if (ReferralFilter{Referrer: "别人"}).Match(r) {
		t.Fatalf("expect no match")
	}
}

func TestDepartmentValidate(t *testing.T) {
	if err := (&Department{Name: "研发部"}).Validate(); err != nil {
		t.Fatalf("valid: %v", err)
	}
	if err := (&Department{}).Validate(); err == nil {
		t.Fatalf("expect error for empty name")
	}
}

func TestResumeValidate(t *testing.T) {
	r := &Resume{CandidateID: "c1", Summary: "5 年后端"}
	if err := r.Validate(); err != nil {
		t.Fatalf("valid: %v", err)
	}
	r.Summary = ""
	if err := r.Validate(); err == nil {
		t.Fatalf("expect error for empty summary")
	}
}

func TestInterviewValidate(t *testing.T) {
	iv := &Interview{CandidateID: "c1", PositionID: "p1", Interviewer: "王经理"}
	if err := iv.Validate(); err != nil {
		t.Fatalf("valid: %v", err)
	}
	if iv.Status != InterviewScheduled {
		t.Fatalf("default status = %s", iv.Status)
	}
	iv.Interviewer = ""
	if err := iv.Validate(); err == nil {
		t.Fatalf("expect error for empty interviewer")
	}
}

func TestInterviewAndOfferFilter(t *testing.T) {
	iv := &Interview{CandidateID: "c1", PositionID: "p1", Status: InterviewCompleted}
	if !(InterviewFilter{Status: InterviewCompleted}).Match(iv) {
		t.Fatalf("interview filter status")
	}
	if (InterviewFilter{PositionID: "p2"}).Match(iv) {
		t.Fatalf("interview filter position no match")
	}

	o := &Offer{CandidateID: "c1", PositionID: "p1", Status: OfferAccepted}
	if !(OfferFilter{Status: OfferAccepted}).Match(o) {
		t.Fatalf("offer filter status")
	}
}

func TestPositionFilterDepartment(t *testing.T) {
	p := &Position{Title: "后端", DepartmentID: "d1", Status: PositionOpen}
	if !(PositionFilter{DepartmentID: "d1"}).Match(p) {
		t.Fatalf("expect match by department")
	}
	if (PositionFilter{DepartmentID: "d2"}).Match(p) {
		t.Fatalf("expect no match by department")
	}
}

func TestCandidateFilterPosition(t *testing.T) {
	c := &Candidate{Name: "张三", Email: "z@x.com", PositionID: "p1", Status: CandidateApplied}
	if !(CandidateFilter{PositionID: "p1"}).Match(c) {
		t.Fatalf("expect match by position")
	}
	if (CandidateFilter{PositionID: "p2"}).Match(c) {
		t.Fatalf("expect no match by position")
	}
}

func TestValidationErrorIsValidation(t *testing.T) {
	err := NewValidationError("field", "消息")
	if !IsValidationError(err) {
		t.Fatalf("expect validation error")
	}
	if err.Error() != "field: 消息" {
		t.Fatalf("error string = %s", err.Error())
	}
}

func TestPositionHeadcountValidation(t *testing.T) {
	p := &Position{Title: "后端", DepartmentID: "d1", Headcount: 0}
	if err := p.Validate(); err == nil {
		t.Fatalf("expect error for zero headcount")
	}
}

func TestOfferAndInterviewDefaultStatus(t *testing.T) {
	o := &Offer{CandidateID: "c1", PositionID: "p1", Salary: 100}
	_ = o.Validate()
	if o.Status != OfferPending {
		t.Fatalf("offer default status = %s", o.Status)
	}
	iv := &Interview{CandidateID: "c1", PositionID: "p1", Interviewer: "x"}
	_ = iv.Validate()
	if iv.Status != InterviewScheduled {
		t.Fatalf("interview default status = %s", iv.Status)
	}
}
