package service

import (
	"testing"
	"time"

	"recruit/internal/config"
	"recruit/internal/model"
	"recruit/internal/store"
	"recruit/pkg/logger"
)

func newTestService() *Service {
	cfg := &config.Config{MaxPageSize: 100}
	log := logger.NewLevel(logger.LevelError)
	return New(store.NewMemoryStore(), log, cfg)
}

func setupPosition(t *testing.T) (*Service, *model.Department, *model.Position) {
	t.Helper()
	s := newTestService()
	dep, err := s.CreateDepartment(model.Department{Name: "研发部"})
	if err != nil {
		t.Fatalf("create department: %v", err)
	}
	pos, err := s.CreatePosition(model.Position{
		Title: "后端工程师", DepartmentID: dep.ID, Headcount: 3, SalaryMin: 20000, SalaryMax: 40000,
	})
	if err != nil {
		t.Fatalf("create position: %v", err)
	}
	return s, dep, pos
}

func TestCreateCandidateRequiresOpenPosition(t *testing.T) {
	s, _, pos := setupPosition(t)
	if _, err := s.CreateCandidate(model.Candidate{Name: "张三", Email: "z@x.com", PositionID: pos.ID}); err != nil {
		t.Fatalf("create candidate: %v", err)
	}
	// 关闭职位后不能再投递
	s.UpdatePosition(pos.ID, model.Position{Status: model.PositionClosed})
	if _, err := s.CreateCandidate(model.Candidate{Name: "李四", Email: "l@x.com", PositionID: pos.ID}); err == nil {
		t.Fatalf("expect error applying to closed position")
	}
}

func TestFullRecruitmentFlow(t *testing.T) {
	s, _, pos := setupPosition(t)

	c, err := s.CreateCandidate(model.Candidate{Name: "张三", Email: "z@x.com", PositionID: pos.ID})
	if err != nil {
		t.Fatalf("create candidate: %v", err)
	}
	if c.Status != model.CandidateApplied {
		t.Fatalf("status = %s, want applied", c.Status)
	}

	// 安排面试 -> interviewing
	iv, err := s.ScheduleInterview(c.ID, pos.ID, "王经理", time.Now().Add(24*time.Hour))
	if err != nil {
		t.Fatalf("schedule interview: %v", err)
	}
	c, _ = s.GetCandidate(c.ID)
	if c.Status != model.CandidateInterviewing {
		t.Fatalf("after schedule status = %s, want interviewing", c.Status)
	}

	// 面试通过 -> offered
	if _, err := s.CompleteInterview(iv.ID, true, "技术扎实"); err != nil {
		t.Fatalf("complete interview: %v", err)
	}
	c, _ = s.GetCandidate(c.ID)
	if c.Status != model.CandidateOffered {
		t.Fatalf("after interview status = %s, want offered", c.Status)
	}

	// 发 Offer
	o, err := s.CreateOffer(c.ID, pos.ID, 35000)
	if err != nil {
		t.Fatalf("create offer: %v", err)
	}
	if o.Status != model.OfferPending {
		t.Fatalf("offer status = %s, want pending", o.Status)
	}

	// 接受 -> hired，职位录用数 +1
	if _, err := s.AcceptOffer(o.ID); err != nil {
		t.Fatalf("accept offer: %v", err)
	}
	c, _ = s.GetCandidate(c.ID)
	if c.Status != model.CandidateHired {
		t.Fatalf("after accept status = %s, want hired", c.Status)
	}
	pos, _ = s.GetPosition(pos.ID)
	if pos.HiredCount != 1 {
		t.Fatalf("hired count = %d, want 1", pos.HiredCount)
	}
}

func TestInterviewFailRejectsCandidate(t *testing.T) {
	s, _, pos := setupPosition(t)
	c, _ := s.CreateCandidate(model.Candidate{Name: "张三", Email: "z@x.com", PositionID: pos.ID})
	iv, _ := s.ScheduleInterview(c.ID, pos.ID, "王经理", time.Now())
	s.CompleteInterview(iv.ID, false, "不符合要求")
	c, _ = s.GetCandidate(c.ID)
	if c.Status != model.CandidateRejected {
		t.Fatalf("status = %s, want rejected", c.Status)
	}
}

func TestDeclineOfferRejectsCandidate(t *testing.T) {
	s, _, pos := setupPosition(t)
	c, _ := s.CreateCandidate(model.Candidate{Name: "张三", Email: "z@x.com", PositionID: pos.ID})
	iv, _ := s.ScheduleInterview(c.ID, pos.ID, "王经理", time.Now())
	s.CompleteInterview(iv.ID, true, "")
	o, _ := s.CreateOffer(c.ID, pos.ID, 35000)
	s.DeclineOffer(o.ID)
	c, _ = s.GetCandidate(c.ID)
	if c.Status != model.CandidateRejected {
		t.Fatalf("status = %s, want rejected", c.Status)
	}
}

func TestInvalidCandidateTransitions(t *testing.T) {
	s, _, pos := setupPosition(t)
	c, _ := s.CreateCandidate(model.Candidate{Name: "张三", Email: "z@x.com", PositionID: pos.ID})
	// applied 不能直接 hired
	if _, err := s.CreateOffer(c.ID, pos.ID, 30000); err == nil {
		t.Fatalf("expect error offering to applied candidate")
	}
}

func TestResumeService(t *testing.T) {
	s, _, pos := setupPosition(t)
	c, _ := s.CreateCandidate(model.Candidate{Name: "张三", Email: "z@x.com", PositionID: pos.ID})
	res, err := s.CreateResume(model.Resume{CandidateID: c.ID, Summary: "5 年后端", Skills: []string{"Go", "MySQL"}})
	if err != nil {
		t.Fatalf("create resume: %v", err)
	}
	if _, err := s.CreateResume(model.Resume{CandidateID: c.ID, Summary: "x"}); err == nil {
		t.Fatalf("expect conflict on duplicate resume")
	}
	byCandidate, err := s.GetResumeByCandidate(c.ID)
	if err != nil || byCandidate.ID != res.ID {
		t.Fatalf("get by candidate: %v", err)
	}
}

func TestDeletePositionBlocked(t *testing.T) {
	s, _, pos := setupPosition(t)
	c, _ := s.CreateCandidate(model.Candidate{Name: "张三", Email: "z@x.com", PositionID: pos.ID})
	if err := s.DeletePosition(pos.ID); err == nil {
		t.Fatalf("expect error deleting position with active candidate")
	}
	// 终态后可删
	s.RejectCandidate(c.ID)
	if err := s.DeletePosition(pos.ID); err != nil {
		t.Fatalf("delete after reject: %v", err)
	}
}

func TestDeleteDepartmentBlocked(t *testing.T) {
	s, dep, _ := setupPosition(t)
	if err := s.DeleteDepartment(dep.ID); err == nil {
		t.Fatalf("expect error deleting department with position")
	}
}

func TestCandidateStats(t *testing.T) {
	s, _, pos := setupPosition(t)
	s.CreateCandidate(model.Candidate{Name: "张三", Email: "z@x.com", PositionID: pos.ID})
	s.CreateCandidate(model.Candidate{Name: "李四", Email: "l@x.com", PositionID: pos.ID})
	stats, err := s.CandidateStats()
	if err != nil || stats.Total != 2 {
		t.Fatalf("candidate stats: %v total=%d", err, stats.Total)
	}
}

func TestRecruitmentFunnel(t *testing.T) {
	s, _, pos := setupPosition(t)
	c1, _ := s.CreateCandidate(model.Candidate{Name: "张三", Email: "z@x.com", PositionID: pos.ID})
	iv, _ := s.ScheduleInterview(c1.ID, pos.ID, "王经理", time.Now())
	s.CompleteInterview(iv.ID, true, "")
	o, _ := s.CreateOffer(c1.ID, pos.ID, 35000)
	s.AcceptOffer(o.ID)
	// 第二个候选只投递
	s.CreateCandidate(model.Candidate{Name: "李四", Email: "l@x.com", PositionID: pos.ID})

	funnel, err := s.RecruitmentFunnel()
	if err != nil {
		t.Fatalf("funnel: %v", err)
	}
	f := *funnel
	if f[0].Count != 2 || f[1].Count != 1 || f[2].Count != 1 || f[3].Count != 1 {
		t.Fatalf("funnel = %+v", f)
	}
}

func TestPositionAndDepartmentStats(t *testing.T) {
	s, dep, pos := setupPosition(t)
	s.CreateCandidate(model.Candidate{Name: "张三", Email: "z@x.com", PositionID: pos.ID})

	ps, err := s.ListPositionStats()
	if err != nil || len(ps) != 1 {
		t.Fatalf("position stats: %v", err)
	}
	if ps[0].Candidates != 1 {
		t.Fatalf("position candidates = %d", ps[0].Candidates)
	}

	ds, err := s.ListDepartmentStats()
	if err != nil || len(ds) != 1 {
		t.Fatalf("department stats: %v", err)
	}
	if ds[0].Positions != 1 || ds[0].Candidates != 1 {
		t.Fatalf("department stats = %+v", ds[0])
	}
	_ = dep
}

func TestInterviewStats(t *testing.T) {
	s, _, pos := setupPosition(t)
	c1, _ := s.CreateCandidate(model.Candidate{Name: "张三", Email: "z@x.com", PositionID: pos.ID})
	iv, _ := s.ScheduleInterview(c1.ID, pos.ID, "王经理", time.Now())
	s.CompleteInterview(iv.ID, true, "")
	stats, err := s.InterviewStats()
	if err != nil {
		t.Fatalf("interview stats: %v", err)
	}
	if stats.Completed != 1 || stats.Passed != 1 || stats.PassRate != 1.0 {
		t.Fatalf("interview stats = %+v", stats)
	}
}

func TestListPagination(t *testing.T) {
	s, _, pos := setupPosition(t)
	for i := 0; i < 5; i++ {
		s.CreateCandidate(model.Candidate{Name: "候选" + string(rune('A'+i)), Email: string(rune('a'+i)) + "@x.com", PositionID: pos.ID})
	}
	items, total, _ := s.ListCandidates(model.CandidateFilter{}, 1, 2)
	if total != 5 || len(items) != 2 {
		t.Fatalf("page1 total=%d len=%d", total, len(items))
	}
	items, total, _ = s.ListCandidates(model.CandidateFilter{}, 3, 2)
	if total != 5 || len(items) != 1 {
		t.Fatalf("page3 total=%d len=%d", total, len(items))
	}
}

func TestReferralFlow(t *testing.T) {
	s, _, pos := setupPosition(t)
	c, _ := s.CreateCandidate(model.Candidate{Name: "张三", Email: "z@x.com", PositionID: pos.ID})

	ref, err := s.CreateReferral("老员工王", c.ID, pos.ID, 5000)
	if err != nil {
		t.Fatalf("create referral: %v", err)
	}
	if ref.Status != model.ReferralPending {
		t.Fatalf("status = %s, want pending", ref.Status)
	}
	// 重复内推同一候选人 -> 冲突
	if _, err := s.CreateReferral("另一个人", c.ID, pos.ID, 3000); err == nil {
		t.Fatalf("expect conflict on duplicate referral")
	}
	// 内推成功
	ref, err = s.MarkReferralHired(ref.ID)
	if err != nil {
		t.Fatalf("mark hired: %v", err)
	}
	if ref.Status != model.ReferralHired || ref.HiredAt == nil {
		t.Fatalf("referral after hired = %+v", ref)
	}
	stats, _ := s.ReferralStats()
	if stats.Total != 1 || stats.Hired != 1 || stats.TotalBonus != 5000 {
		t.Fatalf("referral stats = %+v", stats)
	}
}

func TestBatchReject(t *testing.T) {
	s, _, pos := setupPosition(t)
	var ids []string
	for i := 0; i < 3; i++ {
		c, _ := s.CreateCandidate(model.Candidate{Name: "候选" + string(rune('A'+i)), Email: string(rune('a'+i)) + "@x.com", PositionID: pos.ID})
		ids = append(ids, c.ID)
	}
	n, err := s.BatchReject(ids)
	if err != nil || n != 3 {
		t.Fatalf("batch reject: %v n=%d", err, n)
	}
}

func TestMonthlyAndExport(t *testing.T) {
	s, _, pos := setupPosition(t)
	s.CreateCandidate(model.Candidate{Name: "张三", Email: "z@x.com", PositionID: pos.ID})

	monthly, err := s.MonthlyStats()
	if err != nil || len(monthly) != 1 {
		t.Fatalf("monthly stats: %v len=%d", err, len(monthly))
	}
	posExport, err := s.ExportPositions()
	if err != nil || len(posExport) != 1 {
		t.Fatalf("export positions: %v len=%d", err, len(posExport))
	}
	if posExport[0].DepartmentName != "研发部" {
		t.Fatalf("export position = %+v", posExport[0])
	}
	candExport, err := s.ExportCandidates()
	if err != nil || len(candExport) != 1 {
		t.Fatalf("export candidates: %v len=%d", err, len(candExport))
	}
	if candExport[0].PositionTitle != "后端工程师" {
		t.Fatalf("export candidate = %+v", candExport[0])
	}
}

func TestCandidateEmailConflict(t *testing.T) {
	s, _, pos := setupPosition(t)
	s.CreateCandidate(model.Candidate{Name: "张三", Email: "z@x.com", PositionID: pos.ID})
	if _, err := s.CreateCandidate(model.Candidate{Name: "李四", Email: "z@x.com", PositionID: pos.ID}); err == nil {
		t.Fatalf("expect conflict on same email+position")
	}
}

func TestScheduleInterviewPositionMismatch(t *testing.T) {
	s, _, pos := setupPosition(t)
	c, _ := s.CreateCandidate(model.Candidate{Name: "张三", Email: "z@x.com", PositionID: pos.ID})
	if _, err := s.ScheduleInterview(c.ID, "wrong-position", "王经理", time.Now()); err == nil {
		t.Fatalf("expect error for position mismatch")
	}
}

func TestInterviewInvalidTransition(t *testing.T) {
	s, _, pos := setupPosition(t)
	c, _ := s.CreateCandidate(model.Candidate{Name: "张三", Email: "z@x.com", PositionID: pos.ID})
	iv, _ := s.ScheduleInterview(c.ID, pos.ID, "王经理", time.Now())
	if _, err := s.CompleteInterview(iv.ID, true, ""); err != nil {
		t.Fatalf("complete: %v", err)
	}
	// 已完成面试不能再完成
	if _, err := s.CompleteInterview(iv.ID, true, ""); err == nil {
		t.Fatalf("expect error re-completing interview")
	}
}

func TestCreateOfferInvalidSalary(t *testing.T) {
	s, _, pos := setupPosition(t)
	c, _ := s.CreateCandidate(model.Candidate{Name: "张三", Email: "z@x.com", PositionID: pos.ID})
	iv, _ := s.ScheduleInterview(c.ID, pos.ID, "王经理", time.Now())
	s.CompleteInterview(iv.ID, true, "")
	if _, err := s.CreateOffer(c.ID, pos.ID, 0); err == nil {
		t.Fatalf("expect error for zero salary")
	}
}

func TestListFilters(t *testing.T) {
	s, _, pos := setupPosition(t)
	c, _ := s.CreateCandidate(model.Candidate{Name: "张三", Email: "z@x.com", PositionID: pos.ID})
	s.ScheduleInterview(c.ID, pos.ID, "王经理", time.Now())

	// 候选人按状态过滤
	items, total, _ := s.ListCandidates(model.CandidateFilter{Status: model.CandidateInterviewing}, 1, 10)
	if total != 1 || len(items) != 1 {
		t.Fatalf("candidate filter: total=%d len=%d", total, len(items))
	}
	// 面试按职位过滤
	ivs, total, _ := s.ListInterviews(model.InterviewFilter{PositionID: pos.ID}, 1, 10)
	if total != 1 || len(ivs) != 1 {
		t.Fatalf("interview filter: total=%d len=%d", total, len(ivs))
	}
}

func TestRejectCandidateDirectly(t *testing.T) {
	s, _, pos := setupPosition(t)
	c, _ := s.CreateCandidate(model.Candidate{Name: "张三", Email: "z@x.com", PositionID: pos.ID})
	c, err := s.RejectCandidate(c.ID)
	if err != nil {
		t.Fatalf("reject: %v", err)
	}
	if c.Status != model.CandidateRejected || c.FinishedAt == nil {
		t.Fatalf("rejected candidate = %+v", c)
	}
	// 已终态不能再淘汰
	if _, err := s.RejectCandidate(c.ID); err == nil {
		t.Fatalf("expect error re-rejecting")
	}
}

func TestPositionUpdateValidation(t *testing.T) {
	s, _, pos := setupPosition(t)
	// 薪资下限高于上限 -> 校验失败
	if _, err := s.UpdatePosition(pos.ID, model.Position{SalaryMin: 50000, SalaryMax: 30000}); err == nil {
		t.Fatalf("expect error for min > max salary")
	}
	// 正常更新 headcount
	p, err := s.UpdatePosition(pos.ID, model.Position{Headcount: 5})
	if err != nil || p.Headcount != 5 {
		t.Fatalf("update headcount: %v", err)
	}
}

func TestOfferListAndGet(t *testing.T) {
	s, _, pos := setupPosition(t)
	c, _ := s.CreateCandidate(model.Candidate{Name: "张三", Email: "z@x.com", PositionID: pos.ID})
	iv, _ := s.ScheduleInterview(c.ID, pos.ID, "王经理", time.Now())
	s.CompleteInterview(iv.ID, true, "")
	o, _ := s.CreateOffer(c.ID, pos.ID, 35000)

	offers, total, _ := s.ListOffers(model.OfferFilter{PositionID: pos.ID}, 1, 10)
	if total != 1 || len(offers) != 1 {
		t.Fatalf("list offers: total=%d len=%d", total, len(offers))
	}
	got, err := s.GetOffer(o.ID)
	if err != nil || got.ID != o.ID {
		t.Fatalf("get offer: %v", err)
	}
}

func TestResumeUpdateAndList(t *testing.T) {
	s, _, pos := setupPosition(t)
	c, _ := s.CreateCandidate(model.Candidate{Name: "张三", Email: "z@x.com", PositionID: pos.ID})
	res, _ := s.CreateResume(model.Resume{CandidateID: c.ID, Summary: "5 年后端", Skills: []string{"Go"}})
	res, err := s.UpdateResume(res.ID, model.Resume{Summary: "6 年后端", Skills: []string{"Go", "K8s"}})
	if err != nil {
		t.Fatalf("update resume: %v", err)
	}
	if res.Summary != "6 年后端" || len(res.Skills) != 2 {
		t.Fatalf("updated resume = %+v", res)
	}
}

func TestGetNotFound(t *testing.T) {
	s, _, _ := setupPosition(t)
	if _, err := s.GetDepartment("missing"); err == nil {
		t.Fatalf("expect not found for department")
	}
	if _, err := s.GetPosition("missing"); err == nil {
		t.Fatalf("expect not found for position")
	}
	if _, err := s.GetCandidate("missing"); err == nil {
		t.Fatalf("expect not found for candidate")
	}
	if _, err := s.GetResume("missing"); err == nil {
		t.Fatalf("expect not found for resume")
	}
	if _, err := s.GetInterview("missing"); err == nil {
		t.Fatalf("expect not found for interview")
	}
	if _, err := s.GetOffer("missing"); err == nil {
		t.Fatalf("expect not found for offer")
	}
	if _, err := s.GetReferral("missing"); err == nil {
		t.Fatalf("expect not found for referral")
	}
}
